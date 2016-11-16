package enaml_test

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"

	. "github.com/enaml-ops/enaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeploymentManifest", func() {

	Describe("given a NewDeploymentManifest", func() {
		Context("when called with a []byte representation of the cloud config manifest", func() {
			It("then it should initialize the manifest object with the given bytes", func() {
				b, _ := ioutil.ReadFile("./fixtures/concourse.yml")
				dm := NewDeploymentManifest(b)
				Ω(dm.Name).Should(Equal("concourse"))
				Ω(dm.DirectorUUID).Should(Equal("REPLACE_ME"))
				Ω(len(dm.Releases)).Should(Equal(2))
				Ω(len(dm.Stemcells)).Should(Equal(1))
				Ω(len(dm.InstanceGroups)).Should(Equal(3))
				Ω(dm.Update).ShouldNot(BeNil())
			})
		})
	})

	Describe("given a NewDeploymentManifestFromFile", func() {
		Context("when called with a []byte representation of the cloud config manifest", func() {
			It("then it should initialize the manifest object with the given bytes", func() {
				f, err := os.Open("./fixtures/concourse.yml")
				Ω(err).ShouldNot(HaveOccurred())
				dm := NewDeploymentManifestFromFile(f)
				Ω(dm.Name).Should(Equal("concourse"))
				Ω(dm.DirectorUUID).Should(Equal("REPLACE_ME"))
				Ω(len(dm.Releases)).Should(Equal(2))
				Ω(len(dm.Stemcells)).Should(Equal(1))
				Ω(len(dm.InstanceGroups)).Should(Equal(3))
				Ω(dm.Update).ShouldNot(BeNil())
			})
		})
	})

	Describe("given a deployment manifest", func() {
		Context("when adding a stemcell by name", func() {
			dm := new(DeploymentManifest)
			const name = "stemcell-name"
			const alias = "alias"
			dm.AddStemcellByName(name, alias)

			It("should have a stemcell with the name specified and the OS empty", func() {
				Ω(len(dm.Stemcells)).Should(Equal(1))
				Ω(dm.Stemcells[0].Alias).Should(Equal(alias))

				Ω(dm.Stemcells[0].Name).Should(Equal(name))
				Ω(dm.Stemcells[0].OS).Should(BeEmpty())
			})
		})

		Context("when adding a stemcell by OS", func() {
			dm := new(DeploymentManifest)
			const os = "ubuntu-trusty"
			const alias = "alias"
			dm.AddStemcellByOS(os, alias)

			It("should have a stemcell with the OS specified and the name empty", func() {
				Ω(len(dm.Stemcells)).Should(Equal(1))
				Ω(dm.Stemcells[0].Alias).Should(Equal(alias))

				Ω(dm.Stemcells[0].OS).Should(Equal(os))
				Ω(dm.Stemcells[0].Name).Should(BeEmpty())
			})
		})
	})

	Describe("given AddRemoteStemcell", func() {
		Context("when called with valid remote stemcell values", func() {
			var dm *DeploymentManifest
			const controlOS = "ubuntu-trusty"
			const controlVer = "1.2"
			const controlURL = "http://hi.com"
			const controlSHA = "alkshdglkashdg9243"
			const controlAlias = "stuf"

			BeforeEach(func() {
				dm = new(DeploymentManifest)
				dm.AddRemoteStemcell(controlOS, controlAlias, controlVer, controlURL, controlSHA)
			})

			It("then it should properly add a remote stemcell record", func() {
				// you can specify name OR os, make sure we didn't do both
				Ω(dm.Stemcells[0].OS).Should(Equal(controlOS))
				Ω(dm.Stemcells[0].Name).Should(BeEmpty())

				Ω(dm.Stemcells[0].Alias).Should(Equal(controlAlias))
				Ω(dm.Stemcells[0].URL).Should(Equal(controlURL))
				Ω(dm.Stemcells[0].SHA1).Should(Equal(controlSHA))
				Ω(dm.Stemcells[0].Version).Should(Equal(controlVer))
			})
		})
	})

	Describe("given AddRemoteRelease", func() {
		Context("when called with valid remote release values", func() {
			var dm *DeploymentManifest
			var controlName = "name"
			var controlVer = "1.2"
			var controlURL = "http://hi.com"
			var controlSHA = "alkshdglkashdg9243"
			BeforeEach(func() {
				dm = new(DeploymentManifest)
				dm.AddRemoteRelease(controlName, controlVer, controlURL, controlSHA)
			})
			It("then it should properly add a remote release record", func() {
				Ω(dm.Releases[0].Name).Should(Equal(controlName))
				Ω(dm.Releases[0].URL).Should(Equal(controlURL))
				Ω(dm.Releases[0].SHA1).Should(Equal(controlSHA))
				Ω(dm.Releases[0].Version).Should(Equal(controlVer))
			})
		})
	})

	Describe("given manifest tags", func() {
		var dm *DeploymentManifest

		BeforeEach(func() {
			dm = new(DeploymentManifest)
		})

		It("has no tags by default", func() {
			Ω(dm.Tags).Should(BeEmpty())
			Ω(dm.Tag("foo")).Should(BeEmpty())
		})

		It("can add tags", func() {
			dm.AddTag("foo", "bar")
			Ω(dm.Tags).Should(HaveLen(1))
			Ω(dm.Tags).Should(HaveKeyWithValue("foo", "bar"))
		})

		It("can overwrite tags", func() {
			dm.AddTag("foo", "bar")
			dm.AddTag("foo", "baz")
			Ω(dm.Tags).Should(HaveLen(1))
			Ω(dm.Tags).Should(HaveKeyWithValue("foo", "baz"))
		})

		It("can remove tags", func() {
			dm.AddTag("foo", "bar")
			Ω(dm.Tags).Should(HaveLen(1))
			dm.RemoveTag("foo")
			Ω(dm.Tags).Should(BeEmpty())
			Ω(dm.Tag("foo")).Should(BeEmpty())
		})

		It("can retrieve tags", func() {
			dm.AddTag("foo", "bar")
			Ω(dm.Tag("foo")).Should(Equal("bar"))
			Ω(dm.Tag("baz")).Should(BeEmpty())
		})

		It("marshals to the correct YAML", func() {
			dm.AddTag("foo", "bar")
			dm.AddTag("baz", "quux")
			b, err := yaml.Marshal(dm.Tags)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(b).Should(MatchYAML(`{foo: bar, baz: quux}`))
		})
	})

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
