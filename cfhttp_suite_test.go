package cfhttp_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCFHTTP(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CFHTTP Suite")
}
