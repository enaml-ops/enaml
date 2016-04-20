package diff_test

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xchapter7x/enaml/diff"
)

var _ = Describe("jobdiff", func() {

	Describe("given Diff struct ", func() {
		var diff *Diff
		BeforeEach(func() {
			diff = NewDiff("./fixtures")
		})
		Describe("given ReleaseDiff method", func() {
			Context("when calling ReleaseDiff with releases that have changed", func() {
				It("then there should be a set of changes", func() {
					diffset, err := diff.ReleaseDiff("fixtures/releaseA", "fixtures/releaseB")
					Ω(err).ShouldNot(HaveOccurred())
					Ω(diffset).ShouldNot(BeEmpty())
				})
			})

			Context("when calling ReleaseDiff with releases that have changed", func() {
				It("then there should not be any changes", func() {
					diffset, err := diff.ReleaseDiff("fixtures/releaseA", "fixtures/releaseA")
					Ω(err).ShouldNot(HaveOccurred())
					Ω(diffset).Should(BeEmpty())
				})
			})
		})
		Describe("given JobDiffBetweenReleases method", func() {
			Context("when called with a job which doesnt exist", func() {
				It("then it should return an error stating it can't find the jobname", func() {
					var diffset []string
					var err error

					Ω(func() {
						diffset, err = diff.JobDiffBetweenReleases("i-do-not-exist", "http://url.com/concourse?v=1.1.0", "http://url.com/concourse?v=1.1.0")
					}).ShouldNot(Panic())
					Ω(diffset).Should(BeEmpty())
					Ω(err).Should(HaveOccurred())
				})
			})
			Context("when calling JobDiffBetweenReleases on 2 releases with unchanged properties", func() {
				It("then it should return an empty changeset", func() {
					diff, err := diff.JobDiffBetweenReleases("atc", "http://url.com/concourse?v=1.1.0", "http://url.com/concourse?v=1.1.0")
					Ω(diff).Should(BeEmpty())
					Ω(err).ShouldNot(HaveOccurred())
				})
			})

			Context("when calling JobDiffBetweenReleases on 2 releases with different properties on the given job", func() {
				It("then it should return the diff set", func() {
					diff, err := diff.JobDiffBetweenReleases("atc", "http://url.com/concourse?v=1.0.1", "http://url.com/concourse?v=1.1.0")
					Ω(diff).ShouldNot(BeEmpty())
					Ω(err).ShouldNot(HaveOccurred())
				})
			})
		})
	})

	Describe("Given a JobPropertiesDiff func", func() {
		Context("When both yaml sets are the same", func() {
			var jobDiff []string
			BeforeEach(func() {
				var yamlA, _ = ioutil.ReadFile("./fixtures/jobV1.yml")
				jobDiff = JobPropertiesDiff(yamlA, yamlA)
			})

			It("should return no diff records", func() {
				Ω(len(jobDiff)).Should(Equal(0))
			})
		})
		Context("When both yaml sets are different", func() {
			var jobDiff []string
			BeforeEach(func() {
				var yamlA, _ = ioutil.ReadFile("./fixtures/jobV1.yml")
				var yamlB, _ = ioutil.ReadFile("./fixtures/jobV2.yml")
				jobDiff = JobPropertiesDiff(yamlA, yamlB)
			})

			It("should return a list of the differences", func() {
				Ω(len(jobDiff)).Should(BeNumerically(">", 0))
			})
		})
	})
})
