package enaml_test

import (
	"io/ioutil"

	. "github.com/enaml-ops/enaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CloudConfigManifest Primatives", func() {
	Describe("given a NewCloudConfigManifest", func() {
		Context("when called with a []byte representation of the cloud config manifest", func() {
			It("then it should initialize the manifest object with the given bytes", func() {
				b, _ := ioutil.ReadFile("./fixtures/cloudconfig.yml")
				ccm := NewCloudConfigManifest(b)
				Ω(ccm.AZs[0].Name).Should(Equal("us-east-1c"))
				Ω(ccm.VMTypes[0].Name).Should(Equal("small"))
				Ω(ccm.DiskTypes[0].Name).Should(Equal("small"))
			})
		})
	})
})
