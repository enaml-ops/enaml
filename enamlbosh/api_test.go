package enamlbosh_test

import (
	"net/http"
	"os"

	"github.com/enaml-ops/enaml"
	. "github.com/enaml-ops/enaml/enamlbosh"
	"github.com/enaml-ops/enaml/enamlbosh/enamlboshfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given *Client", func() {
	var boshclient *Client
	Describe("given it is initialized with a valid bosh target", func() {
		var (
			userControl = "my-user"
			passControl = "my-pass"
			hostControl = "1.2.3.4"
			portControl = 25555
		)
		BeforeEach(func() {
			boshclient = NewClient(userControl, passControl, hostControl, portControl)
		})
		Describe("GetTask", func() {
			Context("when calling its GetTask method with a valid taskid", func() {
				var bt BoshTask
				var err error
				BeforeEach(func() {
					doer := new(enamlboshfakes.FakeHttpClientDoer)
					body, _ := os.Open("fixtures/deployment_task.json")
					doer.DoReturns(&http.Response{
						Body: body,
					}, nil)
					bt, err = boshclient.GetTask(1180, doer)
				})

				It("then it should return valid task info for the targetted bosh", func() {
					Ω(err).ShouldNot(HaveOccurred())
					Ω(bt).ShouldNot(BeNil())
				})
			})

			Context("when calling its GetTask method WITHOUT a valid taskid", func() {
				var err error
				BeforeEach(func() {
					doer := new(enamlboshfakes.FakeHttpClientDoer)
					body, _ := os.Open("fixtures/deployment_task.json")
					doer.DoReturns(&http.Response{
						Body: body,
					}, nil)
					_, err = boshclient.GetTask(0, doer)
				})
				It("then is should return an error", func() {
					Ω(err).Should(HaveOccurred())
				})
			})
		})

		Describe("PostRemoteRelease", func() {
			Context("when calling its PostRemoteRelease method with a valid url and sha", func() {
				var bt BoshTask
				var err error
				BeforeEach(func() {
					doer := new(enamlboshfakes.FakeHttpClientDoer)
					body, _ := os.Open("fixtures/deployment_task.json")
					doer.DoReturns(&http.Response{
						Body: body,
					}, nil)
					bt, err = boshclient.PostRemoteRelease(enaml.Release{
						URL:  "https://bosh.io/d/github.com/cloudfoundry/cf-release?v=237",
						SHA1: "8996122278b03b6ba21ec673812d2075c37f1097",
					}, doer)
				})

				It("then it should return valid task info for the targetted bosh", func() {
					Ω(err).ShouldNot(HaveOccurred())
					Ω(bt).ShouldNot(BeNil())
				})
			})

			Context("when calling its PostRemoteRelease method WITHOUT a valid url and sha", func() {
				var err error
				BeforeEach(func() {
					doer := new(enamlboshfakes.FakeHttpClientDoer)
					_, err = boshclient.PostRemoteRelease(enaml.Release{}, doer)
				})
				It("then is should return an error as we only currently support remote release", func() {
					Ω(err).Should(HaveOccurred())
				})
			})
		})

		Describe("PostRemoteStemcell", func() {
			Context("when calling its PostRemoteStemcell method with a valid url and sha", func() {
				var bt BoshTask
				var err error
				BeforeEach(func() {
					doer := new(enamlboshfakes.FakeHttpClientDoer)
					body, _ := os.Open("fixtures/deployment_task.json")
					doer.DoReturns(&http.Response{
						Body: body,
					}, nil)
					bt, err = boshclient.PostRemoteStemcell(enaml.Stemcell{
						URL:  "https://bosh.io/d/stemcells/bosh-aws-xen-hvm-ubuntu-trusty-go_agent?v=3232.4",
						SHA1: "a57ef43974387441b4e8f79e8bb74834",
					}, doer)
				})

				It("then it should return valid task info for the targetted bosh", func() {
					Ω(err).ShouldNot(HaveOccurred())
					Ω(bt).ShouldNot(BeNil())
				})
			})

			Context("when calling its PostRemoteStemcell method WITHOUT a valid url and sha", func() {
				var err error
				BeforeEach(func() {
					doer := new(enamlboshfakes.FakeHttpClientDoer)
					_, err = boshclient.PostRemoteStemcell(enaml.Stemcell{}, doer)
				})
				It("then is should return an error as we only currently support remote stemcells", func() {
					Ω(err).Should(HaveOccurred())
				})
			})
		})
		Describe("PostDeployment", func() {
			Context("when calling its PostDeployment method with a valid doer and deployment", func() {
				var bt BoshTask
				var err error
				BeforeEach(func() {
					doer := new(enamlboshfakes.FakeHttpClientDoer)
					body, _ := os.Open("fixtures/deployment_task.json")
					doer.DoReturns(&http.Response{
						Body: body,
					}, nil)
					bt, err = boshclient.PostDeployment(enaml.DeploymentManifest{}, doer)
				})

				It("then it should return valid task info for the targetted bosh", func() {
					Ω(err).ShouldNot(HaveOccurred())
					Ω(bt).ShouldNot(BeNil())
				})
			})

			Context("what calling its GetCloudConfig method w/ a valid httpclientdoer", func() {
				var ccm *enaml.CloudConfigManifest
				var err error
				BeforeEach(func() {
					doer := new(enamlboshfakes.FakeHttpClientDoer)
					body, _ := os.Open("fixtures/getcloudconfig.yml")
					doer.DoReturns(&http.Response{
						Body: body,
					}, nil)
					ccm, err = boshclient.GetCloudConfig(doer)
				})
				It("then we should be given a valid cloudconfigmanifest", func() {
					Ω(err).ShouldNot(HaveOccurred())
					Ω(len(ccm.AZs)).Should(Equal(1))
					Ω(len(ccm.VMTypes)).Should(Equal(2))
					Ω(len(ccm.DiskTypes)).Should(Equal(3))
					Ω(len(ccm.Networks)).Should(Equal(2))
					Ω(ccm.Compilation).ShouldNot(BeNil())
				})
			})
		})

		Describe("GetInfo", func() {
			Context("when calling its GetInfo method with a valid doer", func() {
				var bi *BoshInfo
				var err error
				BeforeEach(func() {
					doer := new(enamlboshfakes.FakeHttpClientDoer)
					body, _ := os.Open("fixtures/getinfo.json")
					doer.DoReturns(&http.Response{
						Body: body,
					}, nil)
					bi, err = boshclient.GetInfo(doer)
				})
				It("then it should return valid info for the targetted bosh", func() {
					Ω(err).ShouldNot(HaveOccurred())
					Ω(bi).ShouldNot(BeNil())
				})
				It("then it should have a valid bosh name", func() {
					Ω(bi.Name).Should(Equal("my-bosh"))
				})
				It("then it should have a valid bosh guid", func() {
					Ω(bi.UUID).Should(Equal("ebecbaf0-70ce-4324-a1ea-8ea27073fc3b"))
				})
				It("then it should have a valid bosh version", func() {
					Ω(bi.Version).Should(Equal("1.3232.2.0 (00000000)"))
				})
				It("then it should have a valid bosh user", func() {
					Ω(bi.User).Should(Equal(""))
				})
				It("then it should have a valid bosh cpi", func() {
					Ω(bi.CPI).Should(Equal("aws_cpi"))
				})
				It("then it should have a valid bosh features", func() {
					Ω(bi.Features).ShouldNot(BeNil())
				})
			})
		})

		Describe("NewCloudConfigRequest", func() {
			Context("when calling its NewCloudConfigRequest method w/ a valid config file", func() {
				var req *http.Request
				BeforeEach(func() {
					req, _ = boshclient.NewCloudConfigRequest(enaml.CloudConfigManifest{})
				})
				It("then we should be able to generate a basic auth request", func() {
					u, p, ok := req.BasicAuth()
					Ω(u).Should(Equal(userControl))
					Ω(p).Should(Equal(passControl))
					Ω(ok).Should(BeTrue())
				})
			})
		})
	})
})
