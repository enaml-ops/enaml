package pkg_test

import (
	"io"
	"os"

	. "github.com/xchapter7x/enaml/pkg"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Walker", func() {
	Describe("Given a walk func", func() {
		var (
			walker Walker
			err    error
		)
		Context("redis-boshrelease-1.tgz", func() {
			var file FileEntry
			BeforeEach(func() {
				walker = NewWalker(readFixtureFile("redis-boshrelease-1.tgz"))
			})
			JustBeforeEach(func() {
				err = walker.Walk()
			})
			Context("When matching on release.MF", func() {
				BeforeEach(func() {
					walker.OnMatch("release.MF", func(f FileEntry) error {
						file = f
						return nil
					})
				})
				It("Calls back with release.MF", func() {
					Expect(file.FileName).To(Equal("./release.MF"))
				})
			})
			Context("When matching jobs dir", func() {
				BeforeEach(func() {
					walker.OnMatch("/jobs/", func(f FileEntry) error {
						file = f
						return nil
					})
				})
				It("has no errors", func() {
					Expect(err).NotTo(HaveOccurred())
				})
				It("Calls back with redis.tgz", func() {
					Expect(file.FileName).To(Equal("./jobs/redis.tgz"))
				})
			})
		})
		Context("p-redis-1.5.0.pivotal", func() {
			var file FileEntry
			BeforeEach(func() {
				walker = NewWalker(readFixtureFile("p-redis-1.5.0.pivotal"))
			})
			JustBeforeEach(func() {
				err = walker.Walk()
			})
			Context("When matching releases dir", func() {
				BeforeEach(func() {
					walker.OnMatch("/releases/", func(f FileEntry) error {
						file = f
						return nil
					})
				})
				It("has no errors", func() {
					Expect(err).NotTo(HaveOccurred())
				})
				It("Calls back with redis-boshrelease-12.tgz", func() {
					Expect(file.FileName).To(Equal("./releases/redis-boshrelease-12.tgz"))
				})
			})
		})
	})
})

// Create fixture command lines, cwd == root of dir to zip/tar:
// COPYFILE_DISABLE=1 tar czfv /tmp/redis-boshrelease-1.tgz --exclude=".DS_Store" .
// zip -vr /tmp/p-redis-1.5.0.pivotal . -x "*.DS_Store"
func readFixtureFile(filename string) io.Reader {
	f, err := os.Open("./fixtures/" + filename)
	Expect(err).NotTo(HaveOccurred())
	return f
}
