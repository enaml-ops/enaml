package enamlbosh_test

import (
	"fmt"
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
		Context("what calling its GetCloudConfig method w/ a valid httpclientdoer", func() {
			var ccm *enaml.CloudConfigManifest
			var err error
			BeforeEach(func() {
				doer := new(enamlboshfakes.FakeHttpClientDoer)
				body, _ := os.Open("fixtures/getcloudconfig.yml")
				doer.DoReturns(&http.Response{
					Body: body,
				}, nil)
				fmt.Println("not using this: ", doer)
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
