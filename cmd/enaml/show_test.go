package main

import (
	"bufio"
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/xchapter7x/enaml/pull"
)

var _ = Describe("Show", func() {
	Describe("Given an All func", func() {
		var (
			err         error
			buf         bytes.Buffer
			releaseFile string
		)
		JustBeforeEach(func() {
			s := &show{
				release:     "../../fixtures/" + releaseFile,
				releaseRepo: pull.Release{CacheDir: ".cache"},
			}
			buf.Reset()
			w := bufio.NewWriter(&buf)
			err = s.All(w)
			w.Flush()
		})
		Context("Redis BOSH release 12", func() {
			BeforeEach(func() {
				releaseFile = "redis-boshrelease-12.tgz"
			})
			It("does not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("writes all jobs and properties", func() {
				s := string(buf.Bytes())
				Expect(s).To(ContainSubstring("redis"))
				Expect(s).To(ContainSubstring("acceptance-tests"))
				Expect(s).To(ContainSubstring("redis.master"))
			})
		})
		Context("Xip Pivnet release 2.0.0", func() {
			BeforeEach(func() {
				releaseFile = "p-xip-2.0.0.pivotal"
			})
			It("does not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("writes all jobs and properties", func() {
				s := string(buf.Bytes())
				Expect(s).To(ContainSubstring("Release: xip"))
				Expect(s).To(ContainSubstring("Job:     xip"))
				Expect(s).To(ContainSubstring("------------------------------------------------------"))
				Expect(s).To(ContainSubstring("xip.named_conf"))
				Expect(s).To(ContainSubstring("  Description: The contents of named.conf (PowerDNS's BIND backend's configuration file)"))
				Expect(s).To(ContainSubstring("  Default: launch=pipe"))
			})
		})
	})
})
