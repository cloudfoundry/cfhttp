package middleware_test

import (
	"net/http"
	"net/http/httptest"

	"code.cloudfoundry.org/cfhttp/middleware"
	"code.cloudfoundry.org/cfhttp/middleware/middlewarefakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HttpsEnforcer", func() {
	var (
		request           *http.Request
		writer            *httptest.ResponseRecorder
		fakeHandler       *middlewarefakes.FakeHandler
		wrappedMiddleware http.Handler
		forceHttps        bool
	)

	BeforeEach(func() {
		forceHttps = true
	})

	JustBeforeEach(func() {
		fakeHandler = new(middlewarefakes.FakeHandler)
		writer = httptest.NewRecorder()
		enforcer := middleware.HTTPSEnforcer{
			ForceHTTPS: forceHttps,
		}

		wrappedMiddleware = enforcer.Wrap(fakeHandler)
	})

	Context("With https header", func() {
		BeforeEach(func() {
			request, _ = http.NewRequest("GET", "https://localhost/foo/bar", nil)
			request.Header.Set("X-Forwarded-Proto", "https")
		})

		It("calls next middleware", func() {
			wrappedMiddleware.ServeHTTP(writer, request)

			Expect(fakeHandler.ServeHTTPCallCount()).To(Equal(1))
		})
	})

	Context("Without https header", func() {
		BeforeEach(func() {
			request, _ = http.NewRequest("GET", "http://localhost/foo/bar", nil)
			request.Header.Set("X-Forwarded-Proto", "http")
		})

		It("does not call next middleware", func() {
			wrappedMiddleware.ServeHTTP(writer, request)

			Expect(fakeHandler.ServeHTTPCallCount()).To(BeZero())
		})

		It("redirects to https", func() {
			wrappedMiddleware.ServeHTTP(writer, request)

			Expect(writer.Code).To(Equal(http.StatusFound))
			Expect(writer.HeaderMap.Get("Location")).To(Equal("https://localhost/foo/bar"))
		})

		Context("when ForceHttps is false", func() {
			BeforeEach(func() {
				forceHttps = false
			})

			It("calls the next middleware", func() {
				wrappedMiddleware.ServeHTTP(writer, request)

				Expect(fakeHandler.ServeHTTPCallCount()).To(Equal(1))
			})
		})
	})
})
