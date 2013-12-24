package publichost_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPublichost(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Publichost Suite")
}
