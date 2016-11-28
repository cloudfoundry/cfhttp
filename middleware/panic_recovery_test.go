package middleware_test

import (
	. "code.cloudfoundry.org/cfhttp/middleware"

	"net/http"
	"net/http/httptest"

	"code.cloudfoundry.org/cfhttp/middleware/middlewarefakes"
	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagertest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PanicRecovery", func() {
	var (
		fakeResponseWriter *httptest.ResponseRecorder
		fakeHandler        *middlewarefakes.FakeHandler
		logger             *lagertest.TestLogger
		m                  PanicRecovery
		request            *http.Request
		err                error
	)

	BeforeEach(func() {
		fakeResponseWriter = httptest.NewRecorder()

		fakeHandler = new(middlewarefakes.FakeHandler)

		logger = lagertest.NewTestLogger("panic recovery")

		m = PanicRecovery{
			Logger: logger,
		}

		request, err = http.NewRequest("GET", "https://example.com", nil)
		Expect(err).NotTo(HaveOccurred())
	})

	It("forwards requests, without logging", func() {
		handler := m.Wrap(fakeHandler)

		Expect(func() {
			handler.ServeHTTP(fakeResponseWriter, request)
		}).NotTo(Panic())

		capturedResponseWriter, capturedRequest := fakeHandler.ServeHTTPArgsForCall(0)
		Expect(capturedResponseWriter).To(BeIdenticalTo(fakeResponseWriter))
		Expect(capturedRequest).To(BeIdenticalTo(request))
		Expect(logger.Logs()).To(BeEmpty())
	})

	It("catches panics, serving a 500", func() {
		fakeHandler.ServeHTTPStub = func(rw http.ResponseWriter, req *http.Request) {
			panic("foobar")
		}

		handler := m.Wrap(fakeHandler)

		Expect(func() {
			handler.ServeHTTP(fakeResponseWriter, request)
		}).NotTo(Panic())

		Expect(fakeResponseWriter.Code).To(Equal(http.StatusInternalServerError))

		Expect(logger.Logs()[0].LogLevel).To(Equal(lager.ERROR))
	})
})
