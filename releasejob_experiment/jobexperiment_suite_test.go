package releasejob_experiment_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestJobexperiment(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Jobexperiment Suite")
}
