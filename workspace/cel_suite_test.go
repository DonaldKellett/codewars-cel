package cel_test

import (
	"testing"

	"github.com/codewars/cel"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var userSolution []byte

var _ = BeforeSuite(func() {
	prog, err := cel.ReadUserSolution()
	Expect(err).NotTo(HaveOccurred())
	userSolution = prog
})

func TestCel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cel Suite")
}
