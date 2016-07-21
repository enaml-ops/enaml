package enamlbosh_test

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"

	"github.com/enaml-ops/enaml"
	. "github.com/enaml-ops/enaml/enamlbosh"
	"github.com/enaml-ops/enaml/enamlbosh/enamlboshfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("given *Client", func() {
	var boshclient *Client
	var server *ghttp.Server

	var (
		userControl = "my-user"
		passControl = "my-pass"
		controlTask = BoshTask{
			ID:          1180,
			State:       "processing",
			Description: "run errand acceptance_tests from deployment cf-warden",
			Timestamp:   1447033291,
			User:        "admin",
		}
	)

	Describe("basic auth tests", func() {
		BeforeEach(func() {
			server = ghttp.NewTLSServer()

			u, _ := url.Parse(server.URL())
			host, port, _ := net.SplitHostPort(u.Host)
			host = u.Scheme + "://" + host
			portInt, _ := strconv.Atoi(port)
			const skipSSLVerify = true
			boshclient = NewClientBasic(userControl, passControl, host, portInt, skipSSLVerify)
		})

		AfterEach(func() {
			server.Close()
		})

		Describe("given it is initialized with a valid bosh target", func() {

			Describe("GetTask", func() {
				Context("when called", func() {
					BeforeEach(func() {
						server.AppendHandlers(
							ghttp.CombineHandlers(
								ghttp.VerifyBasicAuth(userControl, passControl),
								ghttp.RespondWithJSONEncoded(http.StatusOK, controlTask),
							))
					})
					It("should return task info when called with a valid task ID", func() {
						task, err := boshclient.GetTask(controlTask.ID, new(enamlboshfakes.FakeHttpClientDoer))
						Ω(err).ShouldNot(HaveOccurred())
						Ω(task).Should(Equal(controlTask))
					})

					It("should return an error when called WITHOUT a valid task ID", func() {
						_, err := boshclient.GetTask(0, new(enamlboshfakes.FakeHttpClientDoer))
						Ω(err).Should(HaveOccurred())
					})
				})
			})

			Describe("PostRemoteRelease", func() {
				Context("when calling its PostRemoteRelease method with a valid url and sha", func() {
					BeforeEach(func() {
						server.AppendHandlers(
							ghttp.CombineHandlers(
								ghttp.VerifyBasicAuth(userControl, passControl),
								ghttp.RespondWithJSONEncoded(http.StatusOK, controlTask),
							))
					})
					It("then it should return valid task info for the targetted bosh", func() {
						bt, err := boshclient.PostRemoteRelease(enaml.Release{
							URL:  "https://bosh.io/d/github.com/cloudfoundry/cf-release?v=237",
							SHA1: "8996122278b03b6ba21ec673812d2075c37f1097",
						}, new(enamlboshfakes.FakeHttpClientDoer))
						Ω(err).ShouldNot(HaveOccurred())
						Ω(bt).ShouldNot(BeNil())
					})
				})

				Context("when calling its PostRemoteRelease method WITHOUT a valid url and sha", func() {
					BeforeEach(func() {
						server.AppendHandlers(
							ghttp.CombineHandlers(
								ghttp.VerifyBasicAuth(userControl, passControl),
								ghttp.RespondWithJSONEncoded(http.StatusOK, struct{}{}),
							))
					})
					It("then it should return an error as we only currently support remote release", func() {
						_, err := boshclient.PostRemoteRelease(enaml.Release{}, new(enamlboshfakes.FakeHttpClientDoer))
						Ω(err).Should(HaveOccurred())
					})
				})
			})

			Describe("PostRemoteStemcell", func() {
				Context("when calling its PostRemoteStemcell method with a valid url and sha", func() {
					BeforeEach(func() {
						server.AppendHandlers(
							ghttp.CombineHandlers(
								ghttp.VerifyBasicAuth(userControl, passControl),
								ghttp.RespondWithJSONEncoded(http.StatusOK, controlTask),
							))
					})

					It("then it should return valid task info for the targetted bosh", func() {
						bt, err := boshclient.PostRemoteStemcell(enaml.Stemcell{
							URL:  "https://bosh.io/d/stemcells/bosh-aws-xen-hvm-ubuntu-trusty-go_agent?v=3232.4",
							SHA1: "a57ef43974387441b4e8f79e8bb74834",
						}, new(enamlboshfakes.FakeHttpClientDoer))
						Ω(err).ShouldNot(HaveOccurred())
						Ω(bt).ShouldNot(BeNil())
					})
				})

				Context("when calling its PostRemoteStemcell method WITHOUT a valid url and sha", func() {
					BeforeEach(func() {
						server.AppendHandlers(
							ghttp.CombineHandlers(
								ghttp.VerifyBasicAuth(userControl, passControl),
								ghttp.RespondWithJSONEncoded(http.StatusOK, struct{}{}),
							))
					})
					It("then it should return an error as we only currently support remote stemcells", func() {
						_, err := boshclient.PostRemoteStemcell(enaml.Stemcell{}, new(enamlboshfakes.FakeHttpClientDoer))
						Ω(err).Should(HaveOccurred())
					})
				})
			})

			Describe("PostDeployment", func() {
				Context("when calling its PostDeployment method with a valid deployment", func() {
					BeforeEach(func() {
						server.AppendHandlers(
							ghttp.CombineHandlers(
								ghttp.VerifyBasicAuth(userControl, passControl),
								ghttp.RespondWithJSONEncoded(http.StatusOK, controlTask),
							))
					})

					It("then it should return valid task info for the targetted bosh", func() {
						bt, err := boshclient.PostDeployment(enaml.DeploymentManifest{}, new(enamlboshfakes.FakeHttpClientDoer))
						Ω(err).ShouldNot(HaveOccurred())
						Ω(bt).Should(Equal(controlTask))
					})
				})

				Context("when calling its GetCloudConfig method", func() {
					BeforeEach(func() {
						body, _ := ioutil.ReadFile("fixtures/getcloudconfig.yml")
						server.AppendHandlers(
							ghttp.CombineHandlers(
								ghttp.VerifyBasicAuth(userControl, passControl),
								ghttp.RespondWith(http.StatusOK, body),
							))
					})

					It("then we should be given a valid cloudconfigmanifest", func() {
						ccm, err := boshclient.GetCloudConfig(new(enamlboshfakes.FakeHttpClientDoer))
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
				Context("when calling its GetInfo method", func() {
					var bi *BoshInfo
					var err error
					BeforeEach(func() {
						body, _ := ioutil.ReadFile("fixtures/getinfo.json")
						server.AppendHandlers(
							ghttp.CombineHandlers(
								ghttp.VerifyBasicAuth(userControl, passControl),
								ghttp.RespondWith(http.StatusOK, body),
							))
						bi, err = boshclient.GetInfo(new(enamlboshfakes.FakeHttpClientDoer))
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

			Describe("GetStemcells", func() {
				Describe("given a call to a bosh with stemcells already uploaded", func() {
					var controlStemcells = []DeployedStemcell{
						{
							Name:    "bosh-warden-boshlite-ubuntu-trusty-go_agent",
							OS:      "ubuntu-trusty",
							Version: "3126",
						},
					}
					BeforeEach(func() {
						server.AppendHandlers(
							ghttp.CombineHandlers(
								ghttp.VerifyBasicAuth(userControl, passControl),
								ghttp.RespondWithJSONEncoded(http.StatusOK, controlStemcells),
							))
					})

					Context("when calling GetStemcells on that bosh", func() {
						It("should return an array containing those stemcells' metadata", func() {
							sl, err := boshclient.GetStemcells(new(enamlboshfakes.FakeHttpClientDoer))
							Ω(len(sl)).Should(Equal(1))
							Ω(err).ShouldNot(HaveOccurred())
							Ω(sl[0].Name).Should(Equal("bosh-warden-boshlite-ubuntu-trusty-go_agent"))
							Ω(sl[0].Version).Should(Equal("3126"))
							Ω(sl[0].OS).Should(Equal("ubuntu-trusty"))
						})
					})
				})

				Describe("given a call to a bosh with no stemcells available", func() {
					BeforeEach(func() {
						server.AppendHandlers(
							ghttp.CombineHandlers(
								ghttp.VerifyBasicAuth(userControl, passControl),
								ghttp.RespondWithJSONEncoded(http.StatusOK, []DeployedStemcell{}),
							))
					})

					Context("when calling GetStemcells on that bosh", func() {
						It("should return an empty array", func() {
							sl, err := boshclient.GetStemcells(new(enamlboshfakes.FakeHttpClientDoer))
							Ω(sl).Should(BeEmpty())
							Ω(err).ShouldNot(HaveOccurred())
						})
					})
				})
			})

			Describe("CheckRemoteStemcell", func() {
				Describe("given a stemcell that has not been uploaded", func() {
					var se bool
					var err error

					Context("when called using a stemcell configured with a Name, OS, and Version (but empty response from bosh)", func() {
						BeforeEach(func() {
							server.AppendHandlers(
								ghttp.CombineHandlers(
									ghttp.VerifyBasicAuth(userControl, passControl),
									ghttp.RespondWithJSONEncoded(http.StatusOK, []DeployedStemcell{}),
								))
						})

						It("then it should return a false and a no errors", func() {
							doer := new(enamlboshfakes.FakeHttpClientDoer)
							se, err = boshclient.CheckRemoteStemcell(enaml.Stemcell{
								Name:    "bosh-warden-boshlite-ubuntu-trusty-go_agent",
								OS:      "ubuntu-trusty",
								Version: "3126",
							}, doer)
							Ω(err).ShouldNot(HaveOccurred())
							Ω(se).Should(BeFalse())
						})
					})

					Context("when called using a stemcell configured with a Name, OS, and Version (but no match in bosh result set)", func() {
						var controlStemcells = []DeployedStemcell{
							{
								Name:    "bosh-warden-boshlite-ubuntu-trusty-go_agent",
								OS:      "ubuntu-trusty",
								Version: "3126",
							},
						}
						BeforeEach(func() {
							server.AppendHandlers(
								ghttp.CombineHandlers(
									ghttp.VerifyBasicAuth(userControl, passControl),
									ghttp.RespondWithJSONEncoded(http.StatusOK, controlStemcells),
								))
						})

						It("then it should return a false and a no errors", func() {
							doer := new(enamlboshfakes.FakeHttpClientDoer)
							se, err = boshclient.CheckRemoteStemcell(enaml.Stemcell{
								Name:    "no-matching-name",
								OS:      "no-matching-os",
								Version: "no-version",
							}, doer)
							Ω(err).ShouldNot(HaveOccurred())
							Ω(se).Should(BeFalse())
						})
					})
				})

				Describe("given a stemcell that already has been uploaded", func() {
					Context("when called", func() {
						var controlStemcells = []DeployedStemcell{
							{
								Name:    "bosh-warden-boshlite-ubuntu-trusty-go_agent",
								OS:      "ubuntu-trusty",
								Version: "3126",
							},
						}
						BeforeEach(func() {
							server.AppendHandlers(
								ghttp.CombineHandlers(
									ghttp.VerifyBasicAuth(userControl, passControl),
									ghttp.RespondWithJSONEncoded(http.StatusOK, controlStemcells),
								))
						})

						It("then it should return a true and no errors (when called with name, os, and version)", func() {
							doer := new(enamlboshfakes.FakeHttpClientDoer)
							se, err := boshclient.CheckRemoteStemcell(enaml.Stemcell{
								Name:    "bosh-warden-boshlite-ubuntu-trusty-go_agent",
								OS:      "ubuntu-trusty",
								Version: "3126",
							}, doer)
							Ω(err).ShouldNot(HaveOccurred())
							Ω(se).Should(BeTrue())
						})
						It("then it should return true and no errors (when called with name and version only)", func() {
							doer := new(enamlboshfakes.FakeHttpClientDoer)
							se, err := boshclient.CheckRemoteStemcell(enaml.Stemcell{
								Name:    "bosh-warden-boshlite-ubuntu-trusty-go_agent",
								Version: "3126",
							}, doer)
							Ω(err).ShouldNot(HaveOccurred())
							Ω(se).Should(BeTrue())
						})
						It("then it should return true and no errors (when called with OS and version only)", func() {
							doer := new(enamlboshfakes.FakeHttpClientDoer)
							se, err := boshclient.CheckRemoteStemcell(enaml.Stemcell{
								OS:      "ubuntu-trusty",
								Version: "3126",
							}, doer)
							Ω(err).ShouldNot(HaveOccurred())
							Ω(se).Should(BeTrue())
						})

						It("then it should return an error (when called without a version)", func() {
							doer := new(enamlboshfakes.FakeHttpClientDoer)
							se, err := boshclient.CheckRemoteStemcell(enaml.Stemcell{
								Name: "bosh-warden-boshlite-ubuntu-trusty-go_agent",
								OS:   "ubuntu-trusty",
							}, doer)
							Ω(err).Should(HaveOccurred())
							Ω(se).Should(BeFalse())
						})
					})
				})
			})

			Describe("NewCloudConfigRequest", func() {
				Context("when called with a valid config file", func() {
					It("then we should be able to generate a basic auth request", func() {
						req, err := boshclient.NewCloudConfigRequest(enaml.CloudConfigManifest{})
						Ω(err).ShouldNot(HaveOccurred())
						u, p, ok := req.BasicAuth()
						Ω(u).Should(Equal(userControl))
						Ω(p).Should(Equal(passControl))
						Ω(ok).Should(BeTrue())
					})
				})
			})
		})
	})
})
