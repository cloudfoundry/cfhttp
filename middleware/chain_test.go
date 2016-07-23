package middleware_test

import (
	. "code.cloudfoundry.org/cfhttp/middleware"

	"code.cloudfoundry.org/cfhttp/middleware/middlewarefakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Chain", func() {

	It("wraps the calling handler with each middleware, in reverse order", func() {
		toBeWrapped := new(middlewarefakes.FakeHandler)
		fakeMiddleware1 := new(middlewarefakes.FakeMiddleware)
		fakeHandler1 := new(middlewarefakes.FakeHandler)
		fakeMiddleware2 := new(middlewarefakes.FakeMiddleware)
		fakeHandler2 := new(middlewarefakes.FakeHandler)
		fakeMiddleware3 := new(middlewarefakes.FakeMiddleware)
		fakeHandler3 := new(middlewarefakes.FakeHandler)

		fakeMiddleware1.WrapReturns(fakeHandler1)
		fakeMiddleware2.WrapReturns(fakeHandler2)
		fakeMiddleware3.WrapReturns(fakeHandler3)

		wrappedHandler := Chain{
			fakeMiddleware1,
			fakeMiddleware2,
			fakeMiddleware3,
		}.Wrap(toBeWrapped)

		Expect(fakeMiddleware3.WrapCallCount()).To(Equal(1))
		Expect(fakeMiddleware3.WrapArgsForCall(0)).To(BeIdenticalTo(toBeWrapped))
		Expect(fakeMiddleware2.WrapCallCount()).To(Equal(1))
		Expect(fakeMiddleware2.WrapArgsForCall(0)).To(BeIdenticalTo(fakeHandler3))
		Expect(fakeMiddleware1.WrapCallCount()).To(Equal(1))
		Expect(fakeMiddleware1.WrapArgsForCall(0)).To(BeIdenticalTo(fakeHandler2))

		Expect(wrappedHandler).To(Equal(fakeHandler1))
		Expect(wrappedHandler).To(BeIdenticalTo(fakeHandler1))
	})
})
