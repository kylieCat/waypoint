# waypoint

Wrapper for building Docker images, packaging them in to Helm charts and deploying to Kubernetes.

## Features   

 - Keeps track of the semantic version and increments the correct part of it. Helm requires semantic versioning which was a pain to manage manually.
 - Configuration is minimal and can be done at the repo level or system wide. 
 - Cleans up previous Helm releases and Docker images to help keep your system clean(er :-D)

## Planned

 - Deployments! Currently `waypoint` will build an image and create a Helm chart but doesn't do `helm upgrade/install` yet. 
 - Reduce amount of config needed. Add a `defaults` block to the conf file that will store defaults to be used unless overridden.
 - Plugable build steps. Currently the build steps are hard coded to match what I needed them to do but it should be easy to make them configurable.
 - Conditional steps. Add conditions to steps so that some steps can be skipped (don't remove previous release, only build the imsge, etc).
 - Improve the data that gets saved. Currently a pretty bare bones amount of data is persisted about a release by increasing that data it can be made more useful for analysis.
 - Version rollback. If a release fails automatically roll the version back so the latest version is still accurate.
 
## Usage

  - Add an application with an initial version (the default is `0.1.0`):
  
  ```
  waypoint new my-app [--initial 1.2.3]
  ```
  
  - Get the latest version
    
  ```
  wapoint latest my-app
  0.1.0
  ```
  
  - Increment an applications version by semantic version part:
  
  ```
  waypoint bump my-app --patch
  0.1.1
  
  waypoint bump my-app --minor
  0.2.0
  
  waypoint bump my-app --major
  1.0.0
  ```
  
  - Build a release:
  
  ```
  waypoint release --target stage --minor
  Deleting previous image gcr.io/my-project/my-app:0.2.0...DONE!
  Removing previous Helm chart https://kubernetes-charts.storage.googleapis.com/api/charts/my-app/0.2.0...DONE!
  Building image gcr.io/my-project/my-app:1.0.0...DONE!
  Pushing image gcr.io/my-project/my-app:1.0.0...DONE!
  Creating Helm chart my-app:1.0.0...DONE!
  Uploading Helm chart to https://kubernetes-charts.storage.googleapis.com/api/charts...DONE!
  Updating Helm index file...DONE!
  Updating Helm chart repos...DONE!
  ```
