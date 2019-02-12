package pkg

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/template"
)

type ReleaseType string
type GCPAuthKind string
type BackendKind string

const (
	Major    ReleaseType = "major"
	Minor    ReleaseType = "minor"
	Patch    ReleaseType = "patch"
	Rebuild  ReleaseType = "rebuild"
	Green                = "\033[0;38;5;2m"
	Yellow               = "\033[0;38;5;11m"
	Red                  = "\033[0;38;5;9m"
	ColorOff             = "\033[0m"
	Success              = "Success!\n"
)

const (
	DataStore BackendKind = "datastore"
	Bolt      BackendKind = "bolt"
	MongoDB   BackendKind = "mongo"
	Dynamo    BackendKind = "dynamo"
	ApiKey    GCPAuthKind = "apiKey"
	CredsFile GCPAuthKind = "credsFile"
	ChartsAPI             = "/api/charts"
)

func GetPartsFromSemVer(semver string) ([]int, error) {
	parts := make([]int, 0)
	for _, part := range strings.Split(semver, ".") {
		p, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		parts = append(parts, p)
	}
	return parts, nil
}

func colorPrint(color, msg string) string {
	return fmt.Sprintf("%s%s%s", color, msg, ColorOff)
}

func green(msg string) string {
	return colorPrint(Green, msg)
}

func yellow(msg string) string {
	return colorPrint(Yellow, msg)
}

func red(msg string) string {
	return colorPrint(Red, msg)
}

func printError(msg string) {
	msg = "❌ ERROR: " + msg
	fmt.Println(red(msg))
}

func printWarning(msg string) {
	msg = "⚠ WARNING: " + msg
	fmt.Println(yellow(msg))
}

func printSkipping() {
	fmt.Println(yellow("⚠ SKIPPING: Action ShouldExecute returned false"))
}

func checkErr(err error, exitOnError, done bool) {
	if err != nil {
		if exitOnError {
			printError(err.Error())
			os.Exit(2)
		} else {
			printWarning(err.Error())
		}
	}
	if done {
		fmt.Println(green("DONE!"))
	}
}

type fmtData struct {
	OldVersion string
	NewVersion string
	App        string
}

func (f fmtData) Format(tmpl string) (string, error) {
	var buf bytes.Buffer
	t := template.Must(template.New("letter").Parse(tmpl))
	err := t.Execute(&buf, f)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func newFormatter(r Release) fmtData {
	return fmtData{
		NewVersion: r.newVersion.SemVer(),
		App:        r.App(),
	}
}
