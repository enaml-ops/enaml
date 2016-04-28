package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestEnaml(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Enaml Suite")
}
