package authorization

import (
	"net/http"

	"github.com/tigapilarmandiri/perkakas/common/http_response"
	"github.com/tigapilarmandiri/perkakas/common/util"
)

func SuperAdminOnly() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claim, ok := r.Context().Value(util.ContextKey(util.ContextClaims)).(Claims)
			if !ok {
				http_response.SendForbiddenResponse(w, errUnauthorized)
				return
			}

			if !claim.IsSuperadmin {
				http_response.SendForbiddenResponse(w, errUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
