package enaml

import (
	"io/ioutil"
	"os"

	"github.com/xchapter7x/lo"
	"gopkg.in/yaml.v2"
)

func (s *DeploymentManifest) Bytes() (b []byte) {
	var err error

	if b, err = yaml.Marshal(s); err != nil {
		lo.G.Error("error grabbing deployment manifest bytes: ", err)
		b = nil
	}
	return
}

//NewDeploymentManifestFromFile - will read any implementor of os.File and
//initialize a deployment manifest from its bytes.
//this can be used to read a file or os.Stdin
func NewDeploymentManifestFromFile(f *os.File) *DeploymentManifest {
	var b []byte
	fi, _ := f.Stat()

	if fi.Size() > 0 {
		b, _ = ioutil.ReadAll(f)
	}
	return NewDeploymentManifest(b)
}

//NewDeploymentManifest - deployment manifest constructor
func NewDeploymentManifest(b []byte) *DeploymentManifest {
	dm := new(DeploymentManifest)
	yaml.Unmarshal(b, dm)
	return dm
}

func (s *DeploymentManifest) GetDeployment() DeploymentManifest {
	return *s
}

func (s *DeploymentManifest) SetDirectorUUID(d string) (err error) {
	s.DirectorUUID = d
	return
}

func (s *DeploymentManifest) SetName(n string) (err error) {
	s.Name = n
	return
}

func (s *DeploymentManifest) AddRelease(r Release) (err error) {
	s.Releases = append(s.Releases, r)
	return
}

//AddRemoteRelease - adds a remote release to the manifest. Url should not
//contain version information
func (s *DeploymentManifest) AddRemoteRelease(releaseName, ver, url, sha1 string) (err error) {
	s.Releases = append(s.Releases, Release{
		Name:    releaseName,
		URL:     url,
		SHA1:    sha1,
		Version: ver,
	})
	return
}

func (s *DeploymentManifest) AddReleaseByName(releaseName string) (err error) {
	s.Releases = append(s.Releases, Release{Name: releaseName, Version: "latest"})
	return
}

func (s *DeploymentManifest) AddNetwork(n DeploymentNetwork) (err error) {
	s.Networks = append(s.Networks, n)
	return
}

func (s *DeploymentManifest) AddResourcePool(r ResourcePool) (err error) {
	s.ResourcePools = append(s.ResourcePools, r)
	return
}

func (s *DeploymentManifest) AddStemcell(stemcell Stemcell) (err error) {
	s.Stemcells = append(s.Stemcells, stemcell)
	return
}

func (s *DeploymentManifest) AddStemcellByName(name, alias string) {
	s.Stemcells = append(s.Stemcells, Stemcell{
		Name:    name,
		Alias:   alias,
		Version: "latest",
	})
}

func (s *DeploymentManifest) AddStemcellByOS(os, alias string) {
	s.Stemcells = append(s.Stemcells, Stemcell{
		OS:      os,
		Alias:   alias,
		Version: "latest",
	})
}

// AddRemoteStemcell adds a remote stemcell to the manifest.
// The URL should not contain version information.
func (s *DeploymentManifest) AddRemoteStemcell(os, alias, ver, url, sha1 string) {
	s.Stemcells = append(s.Stemcells, Stemcell{
		OS:      os,
		Alias:   alias,
		Version: ver,
		URL:     url,
		SHA1:    sha1,
	})
}

func (s *DeploymentManifest) AddDiskPool(d DiskPool) (err error) {
	s.DiskPools = append(s.DiskPools, d)
	return
}

func (s *DeploymentManifest) SetCompilation(c Compilation) (err error) {
	s.Compilation = &c
	return
}

func (s *DeploymentManifest) SetUpdate(u Update) (err error) {
	s.Update = u
	return
}

func (s *DeploymentManifest) AddInstanceGroup(i *InstanceGroup) (err error) {
	s.InstanceGroups = append(s.InstanceGroups, i)
	return
}

func (s *DeploymentManifest) GetInstanceGroupByName(name string) (ig *InstanceGroup) {
	for _, i := range s.InstanceGroups {
		if i.Name == name {
			ig = i
			break
		}
	}
	return
}

func (s *DeploymentManifest) GetJobByName(name string) (job *Job) {
	for _, j := range s.Jobs {
		if j.Name == name {
			job = &j
			break
		}
	}
	return
}

func (s *DeploymentManifest) AddJob(j Job) (err error) {
	s.Jobs = append(s.Jobs, j)
	return
}

func (s *DeploymentManifest) AddProperty(k string, val interface{}) (err error) {
	if s.Properties == nil {
		s.Properties = make(map[string]interface{})
	}
	s.Properties[k] = val
	return
}

func (s *DeploymentManifest) SetCloudProvider(c CloudProvider) (err error) {
	s.CloudProvider = &c
	return
}
