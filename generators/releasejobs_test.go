package generators_test

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xchapter7x/enaml/generators"
)

var _ = Describe("given GenerateReleaseJobsPackage function", func() {
	Context("when called on a valid release", func() {
		var controlReleaseURL = "https://bosh.io/d/github.com/concourse/concourse?v=1.1.0"
		var controlCacheDir = "./fixtures"
		var controlOutputDir, _ = ioutil.TempDir("", "releasejobs")
		var controlNumberOfJobs = 6

		BeforeEach(func() {
			GenerateReleaseJobsPackage(controlReleaseURL, controlCacheDir, controlOutputDir)
		})

		AfterEach(func() {
			os.RemoveAll(controlOutputDir)
		})

		It("then it should create the release's job package", func() {
			_, err := ioutil.ReadDir(controlOutputDir)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("then it should create the release's job structs", func() {
			files, _ := ioutil.ReadDir(controlOutputDir)
			Ω(len(files)).Should(Equal(controlNumberOfJobs))
		})
	})
})
