package middlewares

import (
	"net/http"
	"runtime/debug"

	"github.com/tigapilarmandiri/perkakas/common/http_response"

	"github.com/go-chi/chi/v5/middleware"
)

func Recover(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {
				logEntry := middleware.GetLogEntry(r)
				if logEntry != nil {
					logEntry.Panic(rvr, debug.Stack())
				} else {
					middleware.PrintPrettyStack(rvr)
				}

				http_response.SendForbiddenResponse(w, nil)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
