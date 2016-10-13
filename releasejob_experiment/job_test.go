package releasejob_experiment_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	. "github.com/enaml-ops/enaml/releasejob_experiment"
	. "github.com/enaml-ops/enaml/releasejob_experiment/releasejob_experimentfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("release job", func() {
	Describe("Job Runner", func() {
		Context("when calling run w/ the build option", func() {
			Context("and given job metadata ", func() {
				var controlTargetDir string
				var err error
				var tmpfile *os.File

				BeforeEach(func() {
					tmpfile, _ = ioutil.TempFile("", "example")
					controlTargetDir, _ = ioutil.TempDir("", "tmpTarget")
					jobRunner := NewJobRunner([]string{"alksdf", "build", controlTargetDir})
					boshJobFake := new(FakeBoshJob)
					boshJobFake.MetaReturns(BoshJobMeta{
						Version:  "0.0.0",
						Name:     "fake-locator",
						Packages: []string{"jq", "gemfire", "something", "else"},
						JobProperties: []JobProperty{
							JobProperty{
								Name:        "external_dependencies.router.system_domain",
								Description: "System domain",
								EnvVar:      "EXTERNAL_DEPENDENCIES_ROUTER_SYSTEM_DOMAIN",
							},
							JobProperty{
								Name:        "gemfire.locator.addresses",
								Description: "List of GemFire Locator addresses of the form X.X.X.X",
								EnvVar:      "GEMFIRE_LOCATOR_ADDRESSES",
							},
							JobProperty{
								Name:        "gemfire.locator.port",
								Description: "Port the Locator will listen on",
								EnvVar:      "GEMFIRE_LOCATOR_PORT",
								Default:     "55221",
							},
						},
						PIDFile: tmpfile.Name(),
					})
					err = jobRunner.Run(boshJobFake)
					Ω(err).ShouldNot(HaveOccurred(), "we should be able to call run successfully without error")
				})

				AfterEach(func() {
					os.Remove(tmpfile.Name())
					os.RemoveAll(controlTargetDir)
				})

				It("should create a spec file with defined properties", func() {
					jobBinary := path.Join(controlTargetDir, "fake-locator", "spec")
					r, _ := ioutil.ReadFile(jobBinary)
					fmt.Println(string(r))
					Ω(exists(jobBinary)).Should(BeTrue(), "check if job's spec check returns true")
					Ω(string(r)).ShouldNot(BeEmpty(), "file contents should not be empty")
				})

				It("should create a monit file to start the job binary", func() {
					jobBinary := path.Join(controlTargetDir, "fake-locator", "monit")
					r, _ := ioutil.ReadFile(jobBinary)
					Ω(exists(jobBinary)).Should(BeTrue(), "check if job's monit script check returns true")
					Ω(string(r)).ShouldNot(BeEmpty(), "file contents should not be empty")
				})

				It("should create a job binary in the templates directory", func() {
					jobBinary := path.Join(controlTargetDir, "fake-locator", "templates", "fake-locator_ctl")
					Ω(exists(jobBinary)).Should(BeTrue(), "check if job binary exists check returns true")
				})

				It("should create a job and templates dir in the given jobs output path", func() {
					templatesDir := path.Join(controlTargetDir, "fake-locator", "templates")
					Ω(exists(templatesDir)).Should(BeTrue(), "check if target directory exists check returns true")
				})
			})
		})
	})
})

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}
	return false
}
