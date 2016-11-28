package middleware

import (
	"net/http"

	"code.cloudfoundry.org/lager"
)

type PanicRecovery struct {
	Logger lager.Logger
}

func (p PanicRecovery) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		defer func() {
			if panicInfo := recover(); panicInfo != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				p.Logger.Error("Panic while serving request", nil, lager.Data{
					"panicInfo": panicInfo,
				})
			}
		}()
		next.ServeHTTP(rw, req)
	})
}
