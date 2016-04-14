package enaml_test

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xchapter7x/enaml"
	"github.com/xchapter7x/enaml/iaas/vsphere"
)

var _ = Describe("DeploymentManifest Primatives", func() {
	var fakeString = "fake-string"
	var fakeInt = 12
	var fakeBool = true
	var testMarshalledYaml = func(fixturePath string, object interface{}) {
		Context("When it is marshalled to yaml", func() {
			It("should yield a valid yaml string representation", func() {
				fixture, errFixture := ioutil.ReadFile(fixturePath)
				if errFixture != nil {
					panic(fmt.Sprintf("error reading fixture file: %s", errFixture))
				}
				ymlBytes, errYaml := yaml.Marshal(object)
				if errYaml != nil {
					panic(fmt.Sprintf("error marshalling object to yaml: %s", errYaml))
				}
				Ω(ymlBytes).Should(Equal(fixture))
			})
		})
	}

	Describe("Given a Instance", func() {
		Context("when used to generate a cloud_config manifest", func() {
			testMarshalledYaml("./fixtures/instancegroups.yml", struct {
				InstanceGroups []InstanceGroup `yaml:"instance_groups"`
			}{
				InstanceGroups: []InstanceGroup{
					InstanceGroup{
						Name:      fakeString,
						Instances: fakeInt,
						VMType:    fakeString,
						Stemcell:  fakeString,
						AZs:       []string{fakeString},
						Networks: []map[string]interface{}{
							map[string]interface{}{
								fakeString: fakeString,
							},
						},
						Jobs: []InstanceJob{
							InstanceJob{Name: fakeString, Release: fakeString},
						},
					},
				},
			})
		})

		Context("when used to generate a bosh-lite manifest", func() {
			testMarshalledYaml("./fixtures/instancegroups-lite.yml", struct {
				InstanceGroups []InstanceGroup `yaml:"instance_groups"`
			}{
				InstanceGroups: []InstanceGroup{
					InstanceGroup{
						Name:           fakeString,
						Instances:      fakeInt,
						ResourcePool:   fakeString,
						PersistentDisk: fakeInt,
						Networks: []map[string]interface{}{
							map[string]interface{}{
								fakeString: fakeString,
							},
						},
						Jobs: []InstanceJob{
							InstanceJob{Name: fakeString, Release: fakeString},
						},
					},
				},
			})
		})
	})

	XDescribe("Given a DeploymentManifest", func() {
		testMarshalledYaml("./fixtures/deploymentmanifest.yml", nil)
	})

	Describe("Given a Release", func() {
		Context("when using standard release blocks", func() {
			testMarshalledYaml("./fixtures/release.yml", struct{ Releases []Release }{
				Releases: []Release{
					Release{Name: fakeString, Version: fakeString},
				},
			})
		})

		Context("when using custom bosh-init release blocks", func() {
			testMarshalledYaml("./fixtures/releasecustom.yml", struct{ Releases []Release }{
				Releases: []Release{
					Release{Name: fakeString, URL: fakeString, SHA1: fakeString},
					Release{Name: fakeString, URL: fakeString},
				},
			})
		})
	})

	XDescribe("Given a VIPNetwork", func() {
		testMarshalledYaml("./fixtures/vipnetwork.yml", nil)
	})

	XDescribe("Given a DynamicNetwork", func() {
		testMarshalledYaml("./fixtures/dynamicnetwork.yml", nil)
	})

	XDescribe("Given a Subnet", func() {
		testMarshalledYaml("./fixtures/subnet.yml", nil)
	})

	Describe("Given a ManualNetwork", func() {
		testMarshalledYaml("./fixtures/manualnetwork.yml", struct{ Networks []ManualNetwork }{
			Networks: []ManualNetwork{
				ManualNetwork{
					Name: fakeString,
					Type: fakeString,
					Subnets: []Subnet{
						Subnet{
							Range:   fakeString,
							DNS:     fakeString,
							Gateway: fakeString,
							CloudProperties: CloudProperties{
								"name": fakeString,
							},
						},
					},
				},
			},
		})
	})

	Describe("Given a ResourcePool", func() {
		testMarshalledYaml("./fixtures/resourcepool.yml", struct {
			ResourcePools []ResourcePool `yaml:"resource_pools"`
		}{ResourcePools: []ResourcePool{
			ResourcePool{
				Name:     fakeString,
				Network:  fakeString,
				Stemcell: Stemcell{Name: fakeString, Version: fakeString},
				CloudProperties: CloudProperties{
					"instance_type":     fakeString,
					"availability_zone": fakeString,
				},
			},
		}})
	})

	Describe("Given a Stemcell", func() {
		Context("when using a standard stemcell format", func() {
			testMarshalledYaml("./fixtures/stemcell.yml", struct{ Stemcell Stemcell }{
				Stemcell: Stemcell{Name: fakeString, Version: fakeString},
			})
		})

		Context("when using a BOSH 2.0 stemcell format", func() {
			testMarshalledYaml("./fixtures/stemcell-bosh20.yml", struct{ Stemcell Stemcell }{
				Stemcell: Stemcell{Alias: fakeString, OS: fakeString, Version: fakeString},
			})
		})

		Context("when using a custom bosh-init stemcell format", func() {
			testMarshalledYaml("./fixtures/stemcellcustom.yml", struct{ Stemcell Stemcell }{
				Stemcell: Stemcell{URL: fakeString, SHA1: fakeString},
			})
		})
	})

	Describe("Given a DiskPool", func() {
		testMarshalledYaml("./fixtures/diskpool.yml", struct {
			DiskPools []DiskPool `yaml:"disk_pools"`
		}{
			DiskPools: []DiskPool{
				DiskPool{
					Name:     fakeString,
					DiskSize: fakeInt,
					CloudProperties: CloudProperties{
						"type": fakeString,
					},
				},
			}})
	})

	Describe("Given a Compilation", func() {
		testMarshalledYaml("./fixtures/compilation.yml", struct{ Compilation Compilation }{
			Compilation: Compilation{
				Workers:             fakeInt,
				Network:             fakeString,
				ReuseCompilationVMs: fakeBool,
				CloudProperties: CloudProperties{
					"instance_type":     fakeString,
					"availability_zone": fakeString,
				},
			},
		})
	})

	Describe("Given a Update", func() {
		testMarshalledYaml("./fixtures/update.yml", struct{ Update Update }{
			Update: Update{
				Canaries:        fakeInt,
				MaxInFlight:     fakeInt,
				CanaryWatchTime: fakeString,
				UpdateWatchTime: fakeString,
			},
		})
	})

	Describe("Given Jobs", func() {
		testMarshalledYaml("./fixtures/job.yml", []Job{
			Job{
				Name:      fakeString,
				Instances: fakeInt,
				Templates: []Template{
					Template{Name: fakeString, Release: fakeString},
				},
				PersistentDisk: fakeString,
				ResourcePool:   fakeString,
				Networks: []Network{
					Network{Name: fakeString},
				},
			},
			Job{
				Name:      fakeString,
				Instances: fakeInt,
				Templates: []Template{
					Template{Name: fakeString, Release: fakeString},
				},
				PersistentDisk: fakeString,
				ResourcePool:   fakeString,
				Networks: []Network{
					Network{Name: fakeString},
				},
			},
		})
	})

	Describe("Given a Network", func() {
		testMarshalledYaml("./fixtures/network.yml", struct{ Networks []Network }{
			Networks: []Network{
				Network{Name: fakeString},
			},
		})
	})

	Describe("Given a Template", func() {
		testMarshalledYaml("./fixtures/template.yml", struct {
			Templates []Template `yaml:",flow"`
		}{
			Templates: []Template{
				Template{Name: fakeString, Release: fakeString},
			},
		})
	})

	Describe("Given a CloudProvider", func() {
		testMarshalledYaml("./fixtures/cloudprovider.yml", struct {
			CloudProvider CloudProvider `yaml:"cloud_provider"`
		}{
			CloudProvider: CloudProvider{
				Template: Template{Name: fakeString, Release: fakeString},
				MBus:     fakeString,
				Properties: vsphere.CloudProviderProperties{
					VCenter: &vsphere.VCenter{
						Address:  fakeString,
						User:     fakeString,
						Password: fakeString,
						DataCenters: []vsphere.DataCenter{
							vsphere.DataCenter{
								Name:                       fakeString,
								VMFolder:                   fakeString,
								TemplateFolder:             fakeString,
								DatastorePattern:           fakeString,
								PersistentDatastorePattern: fakeString,
								DiskPath:                   fakeString,
								Clusters:                   []string{fakeString},
							},
						},
					},
					Agent:     vsphere.Agent{"mbus": fakeString},
					Blobstore: vsphere.Blobstore{"provider": fakeString, "path": fakeString},
					NTP:       vsphere.NTP{fakeString, fakeString},
				},
			},
		})
	})
})
