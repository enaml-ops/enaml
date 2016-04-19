package diff_test

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xchapter7x/enaml/diff"
)

var _ = Describe("jobdiff", func() {
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
