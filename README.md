#Enaml
## Because (EN)ough with the y(AML) already

[![wercker status](https://app.wercker.com/status/869354741507e6bced0b8199b26d4b5c/s/master "wercker status")](https://app.wercker.com/project/bykey/869354741507e6bced0b8199b26d4b5c)

### Intent
- deployment manifests as testable code
- so no one has to write another bosh deployment manifest in yaml again.

### Sample

**below is a repo showing how one might leverage the enaml primatives and
helpers**

[ENAML - CONCOURSE](https://github.com/xchapter7x/enaml-concourse-sample)


### how to use enaml as a cli
```
#install it using go get
$ go get github.com/xchapter7x/enaml/cmd/enaml

#create golang structs for job properties from a release
$ enaml generate-jobs https://bosh.io/d/github.com/concourse/concourse?v=1.1.0
```



### how your deployment could look
```golang

package concourse

import (
	"github.com/xchapter7x/enaml"
	"github.com/xchapter7x/enaml-concourse-sample/releasejobs"
)

var (
	DefaultName   = "concourse"
	DirectorUUID  = "asdfasdfasdf"
	StemcellAlias = "trusty"
)

func main() {
	enaml.Paint(NewDeployment())
}

type Deployment struct {
	enaml.Deployment
	Manifest *enaml.DeploymentManifest
}

func NewDeployment() (d Deployment) {
	d = Deployment{}
	d.Manifest = new(enaml.DeploymentManifest)
	d.Manifest.SetName(DefaultName)
	d.Manifest.SetDirectorUUID(DirectorUUID)
	d.Manifest.AddReleaseByName("concourse")
	d.Manifest.AddReleaseByName("garden-linux")
	d.Manifest.AddStemcellByName("ubuntu-trusty", StemcellAlias)
	web := enaml.NewInstanceGroup("web", 1, "web", StemcellAlias)
	web.AddAZ("z1")
	web.AddNetwork(enaml.Network{"name": "private"})
	atc := enaml.NewInstanceJob("atc", "concourse", releasejobs.Atc{
		ExternalUrl:        "something",
		BasicAuthUsername:  "user",
		BasicAuthPassword:  "password",
		PostgresqlDatabase: "&atc_db atc",
	})
	tsa := enaml.NewInstanceJob("tsa", "concourse", releasejobs.Tsa{})
	web.AddJob(atc)
	web.AddJob(tsa)
	db := enaml.NewInstanceGroup("db", 1, "database", StemcellAlias)
	worker := enaml.NewInstanceGroup("worker", 1, "worker", StemcellAlias)
	d.Manifest.AddInstanceGroup(web)
	d.Manifest.AddInstanceGroup(db)
	d.Manifest.AddInstanceGroup(worker)
	return
}

func (s Deployment) GetDeployment() enaml.DeploymentManifest {
	return *s.Manifest
}
```

### Building enaml

Enaml uses [Glide](https://github.com/Masterminds/glide) to manage vendored Go
dependencies. Glide is a tool for managing the vendor directory within a Go
package. As such, Golang 1.6+ is recommended.

1. If you haven't done so already, install glide and configure your GOPATH.
2. Open a terminal to the cloned enaml repo directory and run `glide install`
3. Run the enaml tests `go test $(glide novendor)`
4. Build the enaml executable `go build -o $GOPATH/bin/enaml cmd/enaml/main.go`
