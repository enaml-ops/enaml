package main

import (
	"os"

	"github.com/enaml-ops/enaml/releasejob_experiment"
	"github.com/xchapter7x/lo"
)

var Version = "0.0.0"

func main() {
	err := releasejob_experiment.NewJobRunner([]string{}).Run(new(TestBoshJob))

	if err != nil {
		lo.G.Error("error: %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}

type TestBoshJob struct{}

func (s *TestBoshJob) Meta() releasejob_experiment.BoshJobMeta {
	return releasejob_experiment.BoshJobMeta{
		Version:  Version,
		Name:     "locator",
		Packages: []string{"gemfire", "jdk8", "gemfire-utils", "jq"},
		JobProperties: []releasejob_experiment.JobProperty{
			releasejob_experiment.JobProperty{
				Name:        "external_dependencies.router.system_domain",
				Description: "System domain",
				EnvVar:      "EXTERNAL_DEPENDENCIES_ROUTER_SYSTEM_DOMAIN",
			},
			releasejob_experiment.JobProperty{
				Name:        "gemfire.locator.addresses",
				Description: "List of GemFire Locator addresses of the form X.X.X.X",
				EnvVar:      "GEMFIRE_LOCATOR_ADDRESSES",
			},
			releasejob_experiment.JobProperty{
				Name:        "gemfire.locator.port",
				Description: "Port the Locator will listen on",
				EnvVar:      "GEMFIRE_LOCATOR_PORT",
				Default:     "55221",
			},
			releasejob_experiment.JobProperty{
				Name:        "gemfire.locator.rest_port",
				Description: "Port the Locator will listen on for REST API",
				EnvVar:      "GEMFIRE_LOCATOR_REST_PORT",
				Default:     "8080",
			},
			releasejob_experiment.JobProperty{
				Name:        "gemfire.locator.vm_memory",
				Description: "RAM allocated to the locator VM in MB",
				EnvVar:      "GEMFIRE_LOCATOR_VM_MEMORY",
			},
			releasejob_experiment.JobProperty{
				Name:        "gemfire.cluster_topology.number_of_locators",
				Description: "Current topology",
				EnvVar:      "GEMFIRE_CLUSTER_TOPOLOGY_NUMBER_OF_LOCATORS",
				Default:     "2",
			},
			releasejob_experiment.JobProperty{
				Name:        "gemfire.cluster_topology.min_number_of_locators",
				Description: "min number of locators which should be present",
				EnvVar:      "GEMFIRE_CLUSTER_TOPOLOGY_MIN_NUMBER_OF_LOCATORS",
				Default:     "2",
			},
		},
	}
}

func (s *TestBoshJob) Start() error {
	return nil
}

func (s *TestBoshJob) Stop() error {
	return nil
}
