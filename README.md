#Enaml
## Because (EN)ough with the y(AML) already

[![wercker status](https://app.wercker.com/status/598e34535cfd9cb173a1bdb633c9005b/s/master "wercker status")](https://app.wercker.com/project/bykey/598e34535cfd9cb173a1bdb633c9005b)

### Intent
- deployment manifests as testable code
- so no one has to write another bosh deployment manifest in yaml again.

### Sample

**below is a repo showing how one might leverage the enaml primatives and
helpers**

**github.com/xchapter7x/enaml-concourse-sample***


### how your deployment will look
```golang

package main

import (
	"github.com/xchapter7x/enaml"
	"github.com/xchapter7x/standard-components/diskpools"
	"github.com/xchapter7x/standard-components/networks"
	"github.com/xchapter7x/standard-components/releases"
	"github.com/xchapter7x/standard-components/resourcepools"
	"github.com/xchapter7x/standard-components/stemcells"
)

func main() {
	enaml.Paint(NewDeployment())
}

var (
	DefaultName            = "concourse"
	ConcourseVersion       = os.Getenv("CONCOURSE_VERSION")
	ConcourseSHA           = os.Getenv("CONCOURSE_SHA1")
	GardenVersion          = os.Getenv("GARDEN_VERSION")
	GardenSHA              = os.Getenv("GARDEN_SHA1")
	VSphereCPIVersion      = os.Getenv("CPI_VERSION")
	VSphereCPISHA          = os.Getenv("CPI_SHA1")
	VSphereStemcellVersion = os.Getenv("STEMCELL_VERSION")
	VSphereStemcellSHA     = os.Getenv("STEMCELL_SHA1")
)

type Deployment struct {
	enaml.Deployment
	enaml.DeploymentManifest
}

func NewDeployment() (deployment Deployment) {
	deployment = Deployment{}
	deployment.DeploymentManifest = enaml.DeploymentManifest{
		Name: DefaultName,
	}
	deployment.DeploymentManifest.Releases = []enaml.Release{
		releases.NewConcourse(ConcourseVersion, ConcourseSHA),
		releases.NewGarden(GardenVersion, GardenSHA),
	}
	deployment.DeploymentManifest.Networks = []enaml.DeploymentNetwork{
		networks.NewFooNetworkExternal(DefaultName),
	}
	deployment.DeploymentManifest.ResourcePools = []enaml.ResourcePool{
		resourcepools.NewSmallResource(DefaultName, DefaultName),
	}
	deployment.DiskPools = []enaml.DiskPool{
		diskpools.NewDiskPool("db", 10240),
	}
	return
}

func (s Deployment) VSphere() enaml.DeploymentManifest {
	s.Releases = append(s.Releases, releases.NewBoshVSphereCPI(VSphereCPIVersion, VSphereCPISHA))
	for i := range s.ResourcePools {
		s.ResourcePools[i].Stemcell = stemcells.NewUbuntuTrusty(VSphereStemcellVersion, VSphereStemcellSHA)
	}
	return s.DeploymentManifest
}

func (s Deployment) AWS() enaml.DeploymentManifest {
	panic("un-implemented iaas")
}

func (s Deployment) Azure() enaml.DeploymentManifest {
	panic("un-implemented iaas")
}

func (s Deployment) OpenStack() enaml.DeploymentManifest {
	panic("un-implemented iaas")
}
```
