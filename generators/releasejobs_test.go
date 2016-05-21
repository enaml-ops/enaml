package generators_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/enaml-ops/enaml/generators"
	"github.com/xchapter7x/lo"
)

var _ = Describe("releasejobs", func() {
	Describe("given a ProcessFile method on ReleaseJobsGenerator", func() {
		var pwd, _ = os.Getwd()
		var tmpDirName = "_test"
		var tmpPath = path.Join(pwd, tmpDirName)
		var controlOutputDir string
		var controlRelease = "./fixtures/concourse?v=1.1.0"
		var controlPackageDir = "fixtures/_concourse_1.1.0"
		var controlpackages, _ = ioutil.ReadDir(controlPackageDir)
		var namelistFromPackages = func(pkg []os.FileInfo) (list []string) {
			for _, v := range pkg {
				list = append(list, v.Name())
			}
			return
		}
		var gen *ReleaseJobsGenerator

		BeforeEach(func() {
			os.MkdirAll(tmpPath, 0700)
			controlOutputDir, _ = ioutil.TempDir(tmpPath, "releasejobs")
		})

		AfterEach(func() {
			lo.G.Debug("removing: ", controlOutputDir)
			os.RemoveAll(tmpPath)
		})

		Context("when called on a release", func() {
			BeforeEach(func() {
				gen = &ReleaseJobsGenerator{
					OutputDir: controlOutputDir,
				}
				gen.ProcessFile(controlRelease)
			})
			It("then it should create a package for each job", func() {
				packages, err := ioutil.ReadDir(controlOutputDir)
				lo.G.Debug(controlOutputDir)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(namelistFromPackages(packages)).Should(Equal(namelistFromPackages(controlpackages)))
			})

			It("then it should create a struct for the job properties", func() {
				packages, _ := ioutil.ReadDir(controlOutputDir)

				for _, pkg := range packages {
					jobs, err := ioutil.ReadDir(path.Join(controlOutputDir, pkg.Name()))
					controlJobs, _ := ioutil.ReadDir(path.Join(controlPackageDir, pkg.Name()))
					Ω(err).ShouldNot(HaveOccurred())
					Ω(namelistFromPackages(jobs)).Should(Equal(namelistFromPackages(controlJobs)))
				}
			})

			It("then it should create a struct for each of the nested property objects", func() {

				packages, _ := ioutil.ReadDir(controlOutputDir)

				for _, pkg := range packages {
					packagedir := path.Join(controlOutputDir, pkg.Name())
					jobs, _ := ioutil.ReadDir(packagedir)

					for _, job := range jobs {
						_, err := ioutil.ReadFile(path.Join(packagedir, job.Name()))
						_, errControl := ioutil.ReadFile(path.Join(controlPackageDir, pkg.Name(), job.Name()))
						Ω(err).ShouldNot(HaveOccurred())
						Ω(errControl).ShouldNot(HaveOccurred())
					}
				}
			})

			It("then it should create a struct with the correct elements and format for each nested property", func() {
				var errBuffer bytes.Buffer
				packages, _ := ioutil.ReadDir(controlOutputDir)

				for _, pkg := range packages {

					lo.G.Debug(controlOutputDir)
					lo.G.Debug(pkg.Name())
					packagedir := path.Join(controlOutputDir, pkg.Name())
					lo.G.Debug("package", packagedir)
					cmd := exec.Command("go", "build", "./"+path.Join(tmpDirName, path.Base(controlOutputDir), pkg.Name()))
					cmd.Stderr = &errBuffer
					out, err := cmd.Output()
					lo.G.Debug("out: ", out)
					Ω(errBuffer.String()).Should(BeEmpty())
					Ω(err).ShouldNot(HaveOccurred())
				}
			})
		})
	})

	Describe("given GenerateReleaseJobsPackage function", func() {
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
})
