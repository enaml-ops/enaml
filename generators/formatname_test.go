package generators_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xchapter7x/enaml/generators"
)

var _ = Describe("name formatters", func() {
	Describe("ConvertToCamelCase", func() {
		Context("when given a name with formatting issues", func() {
			It("should camel case the name as output", func() {
				badname := "hello-there_badname"
				controlname := "HelloThereBadname"
				Ω(ConvertToCamelCase(badname)).Should(Equal(controlname))
			})
		})
	})
	Describe("FormatName", func() {
		Context("when given a name with formatting issues", func() {
			It("should properly format the name as output", func() {
				badname := "hello-there_badname"
				controlname := "HelloThereBadname"
				Ω(FormatName(badname)).Should(Equal(controlname))
			})
		})
	})
})
