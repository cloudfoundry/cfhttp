package unix_transport

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUnixTransport(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UnixTransport Suite")
}
