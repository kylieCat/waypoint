package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/kylie-a/requests"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"path/filepath"
	"k8s.io/helm/pkg/renderutil"
	"k8s.io/helm/pkg/repo"
	helm_env "k8s.io/helm/pkg/helm/environment"
	"k8s.io/helm/pkg/helm/helmpath"
	"k8s.io/client-go/util/homedir"
	"bytes"
	"k8s.io/helm/pkg/getter"
	"sync"
	"io/ioutil"
)

type Step struct {
	StartMesg     func(r Release) string
	step          func(r Release) error
	exitOnErr     bool
	ShouldExecute func() bool
}

func (s Step) Execute(r Release) {
	fmt.Printf(s.StartMesg(r))
	if err := s.step(r); err != nil {
		if s.exitOnErr {
			printError(err.Error())
			os.Exit(2)
		} else {
			printWarning(err.Error())
			return
		}
	}
	fmt.Printf("%s\n", DONE)
}

var steps = []Step{
	{
		StartMesg: func(r Release) string {
			return fmt.Sprintf("Deleting previous image %s...", r.deploy.TaggedImageName(r.prevVersion.SemVer()))
		},
		step: func(r Release) error {
			var list []types.ImageSummary
			var err error

			listOpts := r.getDockerListOpts(r.deploy.TaggedImageName(r.prevVersion.SemVer()))
			if list, err = r.docker.ImageList(context.Background(), listOpts); err != nil {
				return err
			}

			var img types.ImageSummary
			if len(list) == 1 {
				img = list[0]
			} else {
				return errors.New(fmt.Sprintf("%d images found; skipping", len(list)))
			}
			_, err = r.docker.ImageRemove(context.Background(), img.ID, types.ImageRemoveOptions{})
			return err
		},
		exitOnErr: false,
	},
	{
		StartMesg: func(r Release) string{
			return fmt.Sprintf("Removing previous Helm chart  %s...", r.deploy.GetHelmDeleteURL(r.prevVersion.SemVer()))
		},
		step: func(r Release) error {
			var resp *requests.Response
			var err error

			if resp, err = requests.Delete(r.deploy.GetHelmDeleteURL(r.prevVersion.SemVer()), requests.WithBasicAuth("c3JlOmlhc2hvZGplZkJlY0pvZDA=")); err != nil {
				return err
			}
			if resp.Code != 200 {
				return errors.New("error deleting previous chart: " + resp.Content())
			}
			return nil
		},
		exitOnErr: false,
	},
	{
		StartMesg: func(r Release) string{
			return fmt.Sprintf("Building image %s...", r.deploy.TaggedImageName(r.newVersion.SemVer()))
		},
		step: func(r Release) error {
			var ctx io.Reader
			var err error
			var body types.ImageBuildResponse
			//var data []byte

			if ctx, err = r.deploy.GetContext(); err != nil {
				return err
			}
			opts := types.ImageBuildOptions{
				Tags: []string{r.deploy.TaggedImageName(r.newVersion.SemVer())},
			}
			if body, err = r.docker.ImageBuild(context.Background(), ctx, opts); err != nil {
				return err
			}
			defer body.Body.Close()
			if _, err = ioutil.ReadAll(body.Body); err != nil {
				return err
			}
			return nil
		},
		exitOnErr: true,
	},
	{
		StartMesg: func(r Release) string{
			return fmt.Sprintf("Pushing image %s...", r.deploy.TaggedImageName(r.newVersion.SemVer()))
		},
		step: func(r Release) error {
			var err error
			var body io.ReadCloser

			opts := types.ImagePushOptions{RegistryAuth: r.deploy.GetDockerAuth()}
			ref := r.deploy.TaggedImageName(r.newVersion.SemVer())
			if body, err = r.docker.ImagePush(context.Background(), ref, opts); err != nil {
				return err
			}
			defer body.Close()
			if _, err = ioutil.ReadAll(body); err != nil {
				// {"errorDetail":{"message":"unauthorized: You don't have the needed permissions to perform this operation, and you may have invalid credentials. To authenticate your request, follow the steps in: https://cloud.google.com/container-registry/docs/advanced-authentication"},"error":"unauthorized: You don't have the needed permissions to perform this operation, and you may have invalid credentials. To authenticate your request, follow the steps in: https://cloud.google.com/container-registry/docs/advanced-authentication"}
				return err
			}
			return nil
		},
		exitOnErr: true,
	},
	{
		StartMesg: func(r Release) string{
			return fmt.Sprintf("Creating Helm chart %s:%s...", r.deploy.App, r.newVersion.SemVer())
		},
		step: func(r Release) error {
			// helm package --version ${version} manifests/${app} -d ./manifests/ &> /tmp/build
			var chart *chart.Chart
			var err error
			path := r.deploy.GetHelmChartDir()

			if chart, err = chartutil.LoadDir(path); err != nil {
				return err
			}

			chart.Metadata.Version = r.newVersion.SemVer()

			if filepath.Base(path) != chart.Metadata.Name {
				return fmt.Errorf("directory name (%s) and Chart.yaml name (%s) must match", filepath.Base(path), chart.Metadata.Name)
			}
			if reqs, err := chartutil.LoadRequirements(chart); err == nil {
				if err := renderutil.CheckDependencies(chart, reqs); err != nil {
					return err
				}
			} else {
				if err != chartutil.ErrRequirementsNotFound {
					return err
				}
			}

			var dest string
			if r.deploy.GetHelmDestDir() == "." {
				// Save to the current working directory.
				dest, err = os.Getwd()
				if err != nil {
					return err
				}
			} else {
				// Otherwise save to set destination
				dest = r.deploy.GetHelmDestDir()
			}


			if _, err = chartutil.Save(chart, dest); err != nil {
				return fmt.Errorf("failed to save: %s", err)
			}

			helmSettings := helm_env.EnvSettings{Home: helmpath.Home(filepath.Join(homedir.HomeDir(), ".helm"))}
			if r.deploy.SaveHelmLocal() {

				lr := helmSettings.Home.LocalRepository()
				if err := repo.AddChartToLocalRepo(chart, lr); err != nil {
					return err
				}
			}

			return err
		},
		exitOnErr: true,
	},
	{
		StartMesg: func(r Release) string{
			return fmt.Sprintf("Uploading Helm chart to %s...", r.deploy.GetHelmPostURL())
		},
		step: func(r Release) error {
			var resp *requests.Response
			var err error

			body := r.deploy.GetHelmPackage(r.newVersion.SemVer())
			if resp, err = requests.Post(r.deploy.GetHelmPostURL(), bytes.NewBuffer(body), requests.WithBasicAuth("c3JlOmlhc2hvZGplZkJlY0pvZDA=")); err != nil {
				return err
			}
			if resp.Code <= 199 || resp.Code >= 300 {
				return errors.New("error adding chart: " + resp.Content())
			}
			return nil
		},
		exitOnErr: true,
	},
	{
		StartMesg: func(r Release) string{
			return "Updating Helm index file..."
		},
		step: func(r Release) error {
			//helm repo index manifests/${app} --url ${repo_url}
			path, err := filepath.Abs(r.deploy.GetHelmChartDir())
			if err != nil {
				return err
			}
			return index(path, r.deploy.GetHelmPostURL(), "")
		},
		exitOnErr: false,
	},
	{
		StartMesg: func(r Release) string{
			return "Updating Helm chart repos..."
		},
		step: func(r Release) error {
			// helm repo update
			settings := helm_env.EnvSettings{Home: helmpath.Home(filepath.Join(homedir.HomeDir(), ".helm"))}
			f, err := repo.LoadRepositoriesFile(settings.Home.RepositoryFile())
			if err != nil {
				return err
			}

			if len(f.Repositories) == 0 {
				return errors.New("no repositories")
			}

			var repos []*repo.ChartRepository
			for _, cfg := range f.Repositories {
				r, err := repo.NewChartRepository(cfg, getter.All(settings))
				if err != nil {
					return err
				}
				repos = append(repos, r)
			}

			updateCharts(repos, ioutil.Discard, settings.Home)
			return nil
		},
		exitOnErr: false,
	},
	{
		StartMesg: func(r Release) string {
			return "Deploying ${GREEN}${APP}:${VERSION}${COLOR_OFF} to ${YELLOW}${CLUSTER}-${TARGET}${COLOR_OFF}..."
		},
		step: func(r Release) error {
			//kubectl config use-context ${CLUSTER}-${TARGET} > /dev/null
			//kubectl config set-context $(kubectl config current-context) --namespace=sre > /dev/null
			//helm upgrade --install ${APP} prd/${APP} \
			//--values=manifests/${APP}/${TARGET}.yaml \
			//--set image.tag=${VERSION} \
			//--set secret.key=$(cat .${TARGET}_key | base64) \
			//--set onboarding.gitHash=$(git rev-parse --short=10 HEAD) \
			//--set onboarding.gitBranch=$(git symbolic-ref --short HEAD)
			return nil
		},
		exitOnErr: true,
	},
}

func updateCharts(repos []*repo.ChartRepository, out io.Writer, home helmpath.Home) {
	fmt.Fprintln(out, "Hang tight while we grab the latest from your chart repositories...")
	var wg sync.WaitGroup
	for _, re := range repos {
		wg.Add(1)
		go func(re *repo.ChartRepository) {
			defer wg.Done()
			if re.Config.Name == "local" {
				fmt.Fprintf(out, "...Skip %s chart repository\n", re.Config.Name)
				return
			}
			err := re.DownloadIndexFile(home.Cache())
			if err != nil {
				fmt.Fprintf(out, "...Unable to get an update from the %q chart repository (%s):\n\t%s\n", re.Config.Name, re.Config.URL, err)
			} else {
				fmt.Fprintf(out, "...Successfully got an update from the %q chart repository\n", re.Config.Name)
			}
		}(re)
	}
	wg.Wait()
	fmt.Fprintln(out, "Update Complete. ⎈ Happy Helming!⎈ ")
}

func index(dir, url, mergeTo string) error {
	out := filepath.Join(dir, "index.yaml")

	i, err := repo.IndexDirectory(dir, url)
	if err != nil {
		return err
	}
	if mergeTo != "" {
		// if index.yaml is missing then create an empty one to merge into
		var i2 *repo.IndexFile
		if _, err := os.Stat(mergeTo); os.IsNotExist(err) {
			i2 = repo.NewIndexFile()
			i2.WriteFile(mergeTo, 0644)
		} else {
			i2, err = repo.LoadIndexFile(mergeTo)
			if err != nil {
				return fmt.Errorf("merge failed: %s", err)
			}
		}
		i.Merge(i2)
	}
	i.SortEntries()
	return i.WriteFile(out, 0644)
}
