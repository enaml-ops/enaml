package generators_test

import (
	"io/ioutil"
	"os"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xchapter7x/enaml/generators"
)

var _ = Describe("releasejobs", func() {
	Describe("given a ProcessFile method on ReleaseJobsGenerator", func() {
		var controlOutputDir, _ = ioutil.TempDir("", "releasejobs")
		var controlRelease = "./fixtures/concourse?v=1.1.0"
		var controlPackageDir = "fixtures/concourse_1.1.0"
		var controlpackages, _ = ioutil.ReadDir(controlPackageDir)
		var namelistFromPackages = func(pkg []os.FileInfo) (list []string) {
			for _, v := range pkg {
				list = append(list, v.Name())
			}
			return
		}
		var gen *ReleaseJobsGenerator

		Context("when called on a release", func() {
			BeforeEach(func() {
				gen = &ReleaseJobsGenerator{
					OutputDir: controlOutputDir,
				}
				gen.ProcessFile(controlRelease)
			})
			AfterEach(func() {
				os.RemoveAll(controlOutputDir)
			})
			It("then it should create a package for each job", func() {
				packages, err := ioutil.ReadDir(controlOutputDir)
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

			XIt("then it should create a struct with the correct elements and format for each nested property", func() {

				packages, _ := ioutil.ReadDir(controlOutputDir)

				for _, pkg := range packages {
					packagedir := path.Join(controlOutputDir, pkg.Name())
					jobs, _ := ioutil.ReadDir(packagedir)

					for _, job := range jobs {
						jobBytes, _ := ioutil.ReadFile(path.Join(packagedir, job.Name()))
						controlBytes, _ := ioutil.ReadFile(path.Join(controlPackageDir, pkg.Name(), job.Name()))
						Ω(jobBytes).Should(Equal(controlBytes))
					}
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
