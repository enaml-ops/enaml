package diff

import (
	"os"

	"github.com/xchapter7x/enaml"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Boshrelease", func() {
	var (
		err     error
		release *boshRelease
	)
	Context("Redis BOSH release 12", func() {
		BeforeEach(func() {
			f, err := os.Open("./fixtures/redis-boshrelease-12.tgz")
			Expect(err).NotTo(HaveOccurred())
			release = newBoshRelease()
			err = release.readBoshRelease(f)
		})
		It("should not have errored", func() {
			Expect(err).NotTo(HaveOccurred())
		})
		Context("Release manifest", func() {
			var rm enaml.ReleaseManifest
			BeforeEach(func() {
				rm = release.ReleaseManifest
			})
			It("should read each of the attributes from the manifest", func() {
				Expect(rm.Version).To(Equal("12"))
				Expect(rm.Name).To(Equal("redis"))
			})
			It("should read each job manifest", func() {
				Expect(rm.Jobs).To(HaveLen(2))
				Expect(rm.Jobs[0].Name).To(Equal("acceptance-tests"))
				Expect(rm.Jobs[1].Name).To(Equal("redis"))
			})
		})
		Context("Release jobs", func() {
			var (
				redisJob, testJob enaml.JobManifest
			)
			BeforeEach(func() {
				Expect(release.JobManifests).To(HaveLen(2))
				redisJob = release.JobManifests["redis"]
				testJob = release.JobManifests["acceptance-tests"]
			})
			It("contains the redis job", func() {
				Expect(redisJob).ToNot(BeNil())
				Expect(redisJob.Name).To(Equal("redis"))
				Expect(redisJob.Properties).To(HaveLen(7))
				Expect(redisJob.Properties).To(HaveKey("redis.port"))
				Expect(redisJob.Properties).To(HaveKey("redis.password"))
				Expect(redisJob.Properties).To(HaveKey("redis.master"))
				Expect(redisJob.Properties).To(HaveKey("consul.service.name"))
				Expect(redisJob.Properties).To(HaveKey("health.interval"))
				Expect(redisJob.Properties).To(HaveKey("health.disk.critical"))
				Expect(redisJob.Properties).To(HaveKey("health.disk.warning"))
			})
			It("contains the acceptance-tests job", func() {
				Expect(testJob).ToNot(BeNil())
				Expect(testJob.Name).To(Equal("acceptance-tests"))
				Expect(testJob.Properties).To(HaveLen(4))
				Expect(testJob.Properties).To(HaveKey("redis.port"))
				Expect(testJob.Properties).To(HaveKey("redis.password"))
				Expect(testJob.Properties).To(HaveKey("redis.master"))
				Expect(testJob.Properties).To(HaveKey("redis.slave"))
			})
		})
	})
})
