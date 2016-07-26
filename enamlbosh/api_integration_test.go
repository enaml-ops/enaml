package enamlbosh_test

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	. "github.com/enaml-ops/enaml/enamlbosh"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("given bosh api", func() {
	testClient("admin", "admin", "https://ec2-52-41-100-248.us-west-2.compute.amazonaws.com", 25555, true)
	testClient("admin", "eadxz7dh1d4e3bhgx518", "https://ec2-52-40-154-174.us-west-2.compute.amazonaws.com", 25555, true)
})

func testClient(user, pass, host string, port int, sslIgnore bool) {
	Describe(fmt.Sprintf("when client initialized for host %s", host), func() {

		Context("when creating bosh client", func() {
			var client *Client
			var err error
			client, err = NewClient(user, pass, host, port, sslIgnore)
			It("should have returned a non-nil client and no error", func() {
				Ω(err).ShouldNot(HaveOccurred())
				Ω(client).ShouldNot(BeNil())
			})
		})

		Context("when getting cloud config", func() {
			var client *Client
			var err error
			var cloudConfig *enaml.CloudConfigManifest
			client, err = NewClient(user, pass, host, port, sslIgnore)
			It("should have returned a non-nil client and no error", func() {
				Ω(err).ShouldNot(HaveOccurred())
				Ω(client).ShouldNot(BeNil())
				cloudConfig, err = client.GetCloudConfig()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(cloudConfig).ShouldNot(BeNil())
			})
		})

	})

}
