package enamlbosh_test

import (
	"net/http"

	"github.com/enaml-ops/enaml"
	. "github.com/enaml-ops/enaml/enamlbosh"
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
