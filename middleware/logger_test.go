package middleware_test

import (
	"net/http"

	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagertest"

	"net/http/httptest"

	"code.cloudfoundry.org/cfhttp/middleware"
	"code.cloudfoundry.org/cfhttp/middleware/middlewarefakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Logger", func() {

	var (
		dummyRequest       *http.Request
		err                error
		fakeResponseWriter *httptest.ResponseRecorder
		fakeHandler        *middlewarefakes.FakeHandler
		logger             *lagertest.TestLogger
		routePrefix        string
	)

	const fakePassword = "fakePassword"

	BeforeEach(func() {
		routePrefix = "/v0"
		dummyRequest, err = http.NewRequest("GET", "/v0/backends", nil)
		Expect(err).NotTo(HaveOccurred())
		dummyRequest.Header.Add("Authorization", fakePassword)

		fakeResponseWriter = httptest.NewRecorder()
		fakeHandler = new(middlewarefakes.FakeHandler)

		logger = lagertest.NewTestLogger("backup-download-test")
		logger.RegisterSink(lager.NewWriterSink(GinkgoWriter, lager.INFO))
	})

	It("logs requests that are prefixed with routePrefix", func() {
		loggerMiddleware := middleware.Logging{
			Logger:      logger,
			RoutePrefix: routePrefix,
		}
		loggerHandler := loggerMiddleware.Wrap(fakeHandler)

		loggerHandler.ServeHTTP(fakeResponseWriter, dummyRequest)

		logContents := logger.Buffer().Contents()
		Expect(logContents).To(ContainSubstring("request"))
		Expect(logContents).To(ContainSubstring("response"))
	})

	It("does not log requests that aren't prefixed with the routePrefix", func() {
		loggerMiddleware := middleware.Logging{
			Logger:      logger,
			RoutePrefix: "/v1",
		}
		loggerHandler := loggerMiddleware.Wrap(fakeHandler)

		loggerHandler.ServeHTTP(fakeResponseWriter, dummyRequest)

		logContents := logger.Buffer().Contents()
		Expect(logContents).To(BeEmpty())
	})

	It("does not log credentials", func() {
		loggerMiddleware := middleware.Logging{
			Logger:      logger,
			RoutePrefix: routePrefix,
		}
		loggerHandler := loggerMiddleware.Wrap(fakeHandler)

		loggerHandler.ServeHTTP(fakeResponseWriter, dummyRequest)

		logContents := logger.Buffer().Contents()
		Expect(logContents).ToNot(ContainSubstring(fakePassword))
	})

	It("calls next handler", func() {
		loggerMiddleware := middleware.Logging{
			Logger:      logger,
			RoutePrefix: routePrefix,
		}
		loggerHandler := loggerMiddleware.Wrap(fakeHandler)

		loggerHandler.ServeHTTP(fakeResponseWriter, dummyRequest)

		Expect(fakeHandler.ServeHTTPCallCount()).To(Equal(1))
		arg0, arg1 := fakeHandler.ServeHTTPArgsForCall(0)
		Expect(arg0).ToNot(BeNil())
		Expect(arg1).To(Equal(dummyRequest))
	})
})
