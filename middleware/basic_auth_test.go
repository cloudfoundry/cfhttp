package middleware_test

import (
	. "code.cloudfoundry.org/cfhttp/middleware"

	"net/http"
	"net/http/httptest"

	"code.cloudfoundry.org/cfhttp/middleware/middlewarefakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BasicAuth", func() {
	var (
		fakeResponseWriter *httptest.ResponseRecorder
		fakeHandler        *middlewarefakes.FakeHandler
		m                  BasicAuth
		request            *http.Request
		err                error
	)

	BeforeEach(func() {
		fakeResponseWriter = httptest.NewRecorder()

		fakeHandler = new(middlewarefakes.FakeHandler)

		m = BasicAuth{
			Username: "foo",
			Password: "bar",
		}

		request, err = http.NewRequest("GET", "https://example.com", nil)
		Expect(err).NotTo(HaveOccurred())
	})

	It("complains when no basic auth is provided", func() {
		wrappedHandler := m.Wrap(fakeHandler)

		wrappedHandler.ServeHTTP(fakeResponseWriter, request)
		Expect(fakeHandler.ServeHTTPCallCount()).To(Equal(0))

		Expect(fakeResponseWriter.Code).To(Equal(http.StatusUnauthorized))
		Expect(fakeResponseWriter.Header().Get("WWW-Authenticate")).To(ContainSubstring("Basic realm"))
	})

	It("complains when incorrect credentials are provided", func() {
		request.SetBasicAuth("foo", "unicorn")

		wrappedHandler := m.Wrap(fakeHandler)

		wrappedHandler.ServeHTTP(fakeResponseWriter, request)
		Expect(fakeHandler.ServeHTTPCallCount()).To(Equal(0))

		Expect(fakeResponseWriter.Code).To(Equal(http.StatusUnauthorized))
		Expect(fakeResponseWriter.Header().Get("WWW-Authenticate")).To(ContainSubstring("Basic realm"))
	})

	It("passes on the handle when given the correct credentials", func() {
		request.SetBasicAuth("foo", "bar")

		wrappedHandler := m.Wrap(fakeHandler)

		wrappedHandler.ServeHTTP(fakeResponseWriter, request)

		Expect(fakeHandler.ServeHTTPCallCount()).To(Equal(1))

		rw, req := fakeHandler.ServeHTTPArgsForCall(0)
		Expect(rw).To(BeIdenticalTo(fakeResponseWriter))
		Expect(req).To(BeIdenticalTo(request))
	})
})
