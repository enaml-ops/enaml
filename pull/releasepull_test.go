package pull_test

import (
	"io"
	"os"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xchapter7x/enaml/pull"
)

var _ = Describe("given Release object", func() {
	Describe("given a Pull method", func() {
		var (
			release  *Release
			filename string
			err      error
		)

		Context("when called on a valid release in the cache", func() {
			var (
				releaseName       = "concourse?v=1.1.0"
				controlReleaseURL = "https://bosh.io/d/github.com/concourse/" + releaseName
				controlCacheDir   = "fixtures"
			)

			BeforeEach(func() {
				release = NewRelease(controlCacheDir)
				filename, err = release.Pull(controlReleaseURL)
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("then it should return a valid filename", func() {
				_, err = os.Stat(filename)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(filename).ShouldNot(BeEmpty())
				Ω(filename).Should(Equal(path.Join(controlCacheDir, releaseName)))
			})
		})

		Context("when called on an existing local release", func() {
			var (
				releaseName        = "concourse?v=1.1.0"
				controlReleaseFile = "fixtures/" + releaseName
				controlCacheDir    = "shouldnotbeused"
			)

			BeforeEach(func() {
				release = NewRelease(controlCacheDir)
				filename, err = release.Pull(controlReleaseFile)
			})

			It("should not have errored", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})
			It("returns the same local file", func() {
				Ω(filename).Should(Equal(controlReleaseFile))
			})
		})

		Context("when called on a local release that does not exist", func() {
			BeforeEach(func() {
				release = NewRelease("ignored")
				filename, err = release.Pull("fixtures/foobar?v=1.0")
			})

			It("should have errored", func() {
				Ω(err).Should(MatchError("Could not pull fixtures/foobar?v=1.0. The file doesn't exist or isn't a valid http(s) URL"))
			})
		})
	})

	Describe("given a Read method", func() {
		var (
			release *Release
			reader  io.ReadCloser
			err     error
		)

		Context("when called on a valid release in the cache", func() {
			var (
				releaseName       = "concourse?v=1.1.0"
				controlReleaseURL = "https://bosh.io/d/github.com/concourse/" + releaseName
				controlCacheDir   = "fixtures"
			)

			BeforeEach(func() {
				release = NewRelease(controlCacheDir)
				reader, err = release.Read(controlReleaseURL)
			})
			AfterEach(func() {
				reader.Close()
			})

			It("then it should return no errors", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})
			It("then it should return a stream to the package", func() {
				Ω(reader).ShouldNot(BeNil())
			})
		})
	})
})
