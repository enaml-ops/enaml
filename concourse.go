package enaml

import (
	"github.com/concourse/atc"
)

const (
	BoshDeploymentResourceName = "bosh-deployment"
	BoshIOResourceName         = "bosh-io-release"
	GitResourceName            = "git"
	GithubResourceName         = "github-release"
)

type ConcourseSource map[string]interface{}
type ConcourseParams ConcourseSource

// ConcoursePipeline an object representing a concourse pipeline yaml
type ConcoursePipeline struct {
	defaultPlatform        string
	defaultImageType       string
	defaultImageRepository string
	atc.Config
}

// GetDefaultTaskImageResource - convenience helper to output default object for
// task images
func (s *ConcoursePipeline) GetDefaultTaskImageResource() atc.ImageResource {
	return atc.ImageResource{
		Type: s.defaultImageType,
		Source: atc.Source{
			"repository": s.defaultImageRepository,
		},
	}
}

// GetDefaultPlatform - platfor getter
func (s *ConcoursePipeline) GetDefaultPlatform() string {
	return s.defaultPlatform
}

//AddRawJob helper to add a job to the pipeline manifest
func (s *ConcoursePipeline) AddRawJob(job atc.JobConfig) {
	s.Jobs = append(s.Jobs, job)
}

//AddGroup helper to add a group to the pipeline manifest
func (s *ConcoursePipeline) AddGroup(name string, jobs ...string) {
	s.Groups = append(s.Groups, atc.GroupConfig{
		Name: name,
		Jobs: jobs,
	})
}

//SetImageDefaults setter
func (s *ConcoursePipeline) SetImageDefaults(platform, imagetype, repo string) {
	s.defaultPlatform = platform
	s.defaultImageType = imagetype
	s.defaultImageRepository = repo
}

//GetResourceByName convenience method to find and return a resource by name
func (s *ConcoursePipeline) GetResourceByName(name string) *atc.ResourceConfig {
	for i, v := range s.Resources {
		if v.Name == name {
			return &s.Resources[i]
		}
	}
	return nil
}

//AddRawResource helper to add a resource to the pipeline manifest
func (s *ConcoursePipeline) AddRawResource(rawResource atc.ResourceConfig) {
	s.Resources = append(s.Resources, rawResource)
}

//AddResource helper to add a resource to the pipeline manifest
func (s *ConcoursePipeline) AddResource(name string, typename string, source map[string]interface{}) {
	s.Resources = append(s.Resources, atc.ResourceConfig{
		Name:   name,
		Type:   typename,
		Source: source,
	})
}

//AddGithubResource github specific resource add
func (s *ConcoursePipeline) AddGithubResource(name string, source map[string]interface{}) {
	s.AddResource(name, GithubResourceName, source)
}

//AddBoshIOResource bosh io specific resource add
func (s *ConcoursePipeline) AddBoshIOResource(name string, source map[string]interface{}) {
	s.AddResource(name, BoshIOResourceName, source)
}

//AddBoshDeploymentResource bosh deployment resource add
func (s *ConcoursePipeline) AddBoshDeploymentResource(name string, source map[string]interface{}) {
	s.AddResource(name, BoshDeploymentResourceName, source)
}

//AddGitResource git specific resource add
func (s *ConcoursePipeline) AddGitResource(name string, source map[string]interface{}) {
	s.AddResource(name, GitResourceName, source)
}
