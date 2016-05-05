package enaml_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xchapter7x/enaml"
)

var _ = Describe("DeploymentManifest", func() {
	Describe("given base setters", func() {
		Context("when called", func() {
			var DefaultName = "testdeploy"
			var controlDeploymentManifest = new(DeploymentManifest)
			var testDeploymentManifest = new(DeploymentManifest)
			BeforeEach(func() {
				controlDeploymentManifest = &DeploymentManifest{
					Name: DefaultName,
					Releases: []Release{
						NewFooRelease(DefaultName, DefaultName),
					},
					Networks: []DeploymentNetwork{
						NewFooNetwork(DefaultName),
					},
					ResourcePools: []ResourcePool{
						NewFooResource(DefaultName, DefaultName),
					},
					DiskPools: []DiskPool{
						NewFooDiskPool("db", 10240),
					},
				}
			})
			It("then it should set the values in the given DeploymentManifest", func() {
				testDeploymentManifest.SetName(DefaultName)
				testDeploymentManifest.AddRelease(NewFooRelease(DefaultName, DefaultName))
				testDeploymentManifest.AddNetwork(NewFooNetwork(DefaultName))
				testDeploymentManifest.AddResourcePool(NewFooResource(DefaultName, DefaultName))
				testDeploymentManifest.AddDiskPool(NewFooDiskPool("db", 10240))
				Ω(*testDeploymentManifest).Should(Equal(*controlDeploymentManifest))
			})
		})
	})
})

func NewFooDiskPool(name string, size int) DiskPool {
	return DiskPool{
		Name:     name,
		DiskSize: size,
	}
}

func NewFooNetwork(networkName string) DeploymentNetwork {
	return &ManualNetwork{
		Name: networkName,
		Type: "manual",
		Subnets: []Subnet{
			Subnet{
				Range:   "10.0.0.0/24",
				DNS:     []string{"10.0.0.2"},
				Gateway: "10.0.0.1",
				CloudProperties: map[string]string{
					"name": "NETWORK_NAME",
				},
			},
		},
	}
}

func NewFooRelease(version, sha string) Release {
	return Release{
		Name: "concourse",
		URL:  "https://bosh.io/d/github.com/concourse/concourse?v=" + version,
		SHA1: sha,
	}
}

func NewFooResource(resourceName, networkName string) ResourcePool {
	return ResourcePool{
		Name:    resourceName,
		Network: networkName,
		CloudProperties: map[string]int{
			"cpu":  2,
			"ram":  4096,
			"disk": 10240,
		},
	}
}
