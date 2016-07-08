package generators_test

import (
	"github.com/enaml-ops/enaml"
	. "github.com/enaml-ops/enaml/generators"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Record", func() {
	Describe("StructName", func() {
		Context("when property is uaa.login.policy.timeout and no other properties", func() {
			var v Record
			properties := []string{"uaa.login.policy.timeout"}
			propertyName := "uaa.login.policy.timeout"
			BeforeEach(func() {
				v = CreateNewRecord(propertyName, enaml.JobManifestProperty{})
			})
			It("should return UaaJob when index is 0", func() {
				structName := v.StructName(0, "uaa", properties)
				Ω(structName).Should(Equal("UaaJob"))
			})
			It("should return Uaa when index is 1", func() {
				structName := v.StructName(1, "uaa", properties)
				Ω(structName).Should(Equal("Uaa"))
			})
			It("should return Login when index is 2", func() {
				structName := v.StructName(2, "uaa", properties)
				Ω(structName).Should(Equal("Login"))
			})
			It("should return Policy when index is 3", func() {
				structName := v.StructName(3, "uaa", properties)
				Ω(structName).Should(Equal("Policy"))
			})
		})
		Context("when property is uaa.login.policy.timeout and other properties have same name", func() {
			var v Record
			properties := []string{"uaa.login.policy.timeout", "uaa.jwt.policy.timeout"}
			BeforeEach(func() {
				v = CreateNewRecord("uaa.login.policy.timeout", enaml.JobManifestProperty{})
			})
			It("should return UaaJob when index is 0", func() {
				structName := v.StructName(0, "uaa", properties)
				Ω(structName).Should(Equal("UaaJob"))
			})
			It("should return Uaa when index is 1", func() {
				structName := v.StructName(1, "uaa", properties)
				Ω(structName).Should(Equal("Uaa"))
			})
			It("should return Login when index is 2", func() {
				structName := v.StructName(2, "uaa", properties)
				Ω(structName).Should(Equal("Login"))
			})
			It("should return LoginPolicy when index is 3", func() {
				structName := v.StructName(3, "uaa", properties)
				Ω(structName).Should(Equal("LoginPolicy"))
			})
		})
		Context("when property is uaa.login.client_secret and other properties have same name prefix", func() {
			var v Record
			properties := []string{"uaa.login.client_secret", "login.signups_enabled"}
			BeforeEach(func() {
				v = CreateNewRecord("uaa.login.client_secret", enaml.JobManifestProperty{})
			})
			It("should return UaaJob when index is 0", func() {
				structName := v.StructName(0, "uaa", properties)
				Ω(structName).Should(Equal("UaaJob"))
			})
			It("should return Uaa when index is 1", func() {
				structName := v.StructName(1, "uaa", properties)
				Ω(structName).Should(Equal("Uaa"))
			})
			It("should return Login when index is 2", func() {
				structName := v.StructName(2, "uaa", properties)
				Ω(structName).Should(Equal("UaaLogin"))
			})
		})
	})
	Describe("FindAllParentsOfSameNamedElement", func() {
		Context("when property is uaa.login.policy.timeout and no other properties", func() {
			var v Record
			properties := []string{"uaa.login.policy.timeout"}
			propertyName := "uaa.login.policy.timeout"
			BeforeEach(func() {
				v = CreateNewRecord(propertyName, enaml.JobManifestProperty{})
			})
			It("parentNames of policy should equal login", func() {
				parentNames := v.FindAllParentsOfSameNamedElement("policy", properties)
				Ω(len(parentNames)).Should(Equal(1))
				Ω(parentNames).Should(ConsistOf("login"))
			})
		})
		Context("when property is uaa.login.policy.timeout and other properties have same name", func() {
			var v Record
			properties := []string{"uaa.login.policy.timeout", "uaa.jwt.policy.timeout"}
			BeforeEach(func() {
				v = CreateNewRecord("uaa.login.policy.timeout", enaml.JobManifestProperty{})
			})
			It("parentNames of policy should equal login and jwt", func() {
				parentNames := v.FindAllParentsOfSameNamedElement("policy", properties)
				Ω(len(parentNames)).Should(Equal(2))
				Ω(parentNames).Should(ConsistOf("login", "jwt"))
			})
		})
		Context("when property is uaa.login.client_secret and other properties have same name prefix", func() {
			var v Record
			properties := []string{"uaa.login.client_secret", "login.signups_enabled"}
			BeforeEach(func() {
				v = CreateNewRecord("uaa.login.client_secret", enaml.JobManifestProperty{})
			})
			It("parentNames of login should equal uaa and ''", func() {
				parentNames := v.FindAllParentsOfSameNamedElement("login", properties)
				Ω(len(parentNames)).Should(Equal(2))
				Ω(parentNames).Should(ConsistOf("uaa", ""))
			})
		})
	})
	Describe("TypeName", func() {
		Context("when property is uaa.login.policy.timeout and no other properties", func() {
			var v Record
			properties := []string{"uaa.login.policy.timeout"}
			propertyName := "uaa.login.policy.timeout"
			BeforeEach(func() {
				v = CreateNewRecord(propertyName, enaml.JobManifestProperty{})
			})
			It("should return *Login when index is 1", func() {
				typeName := v.TypeName(1, properties)
				Ω(typeName).Should(Equal("*Login"))
			})
			It("should return *Policy when index is 2", func() {
				typeName := v.TypeName(2, properties)
				Ω(typeName).Should(Equal("*Policy"))
			})
			It("should return interface{} when index is 3", func() {
				typeName := v.TypeName(3, properties)
				Ω(typeName).Should(Equal("interface{}"))
			})
		})
		Context("when property is uaa.login.policy.timeout and other properties have same name", func() {
			var v Record
			properties := []string{"uaa.login.policy.timeout", "uaa.jwt.policy.timeout"}
			BeforeEach(func() {
				v = CreateNewRecord("uaa.login.policy.timeout", enaml.JobManifestProperty{})
			})
			It("should return *Login when index is 1", func() {
				typeName := v.TypeName(1, properties)
				Ω(typeName).Should(Equal("*Login"))
			})
			It("should return *LoginPolicy when index is 2", func() {
				typeName := v.TypeName(2, properties)
				Ω(typeName).Should(Equal("*LoginPolicy"))
			})
			It("should return interface{} when index is 3", func() {
				typeName := v.TypeName(3, properties)
				Ω(typeName).Should(Equal("interface{}"))
			})
		})
		Context("when property is uaa.login.client_secret and other properties have same name prefix", func() {
			var v Record
			properties := []string{"uaa.login.client_secret", "login.signups_enabled"}
			BeforeEach(func() {
				v = CreateNewRecord("uaa.login.client_secret", enaml.JobManifestProperty{})
			})
			It("should return *UaaLogin when index is 1", func() {
				typeName := v.TypeName(1, properties)
				Ω(typeName).Should(Equal("*UaaLogin"))
			})
			It("should return interface{} when index is 2", func() {
				typeName := v.TypeName(2, properties)
				Ω(typeName).Should(Equal("interface{}"))
			})
		})
	})
})
