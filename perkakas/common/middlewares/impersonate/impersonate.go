package impersonate

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/tigapilarmandiri/perkakas/common/http_response"
	"github.com/tigapilarmandiri/perkakas/common/middlewares/authorization"
	"github.com/tigapilarmandiri/perkakas/common/rds"
	"github.com/tigapilarmandiri/perkakas/common/util"
)

var errUnauthorized = errors.New("you're not authorized")

var genKey = func(key string) string {
	return fmt.Sprintf("impersonate-%s", key)
}

func Store(user string, token string) {
	err := rds.GetClient().Set(context.Background(), genKey(user), token, time.Minute*15).Err()
	if err != nil {
		util.Log.Error().Msg(err.Error())
	}
}

func Clear(user string) {
	err := rds.GetClient().Del(context.Background(), genKey(user)).Err()
	if err != nil {
		util.Log.Error().Msg(err.Error())
	}
}

func IsExist(user string) bool {
	val := rds.GetClient().Exists(context.Background(), genKey(user)).Val()

	return val == 1
}

func Get(user string) string {
	cmd := rds.GetClient().Get(context.Background(), genKey(user))
	if cmd.Err() != nil {
		util.Log.Error().Msg(cmd.Err().Error())
	}

	return cmd.String()
}

func ImpersonateCheck() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			claims, ok := ctx.Value(util.ContextKey(util.ContextClaims)).(authorization.Claims)
			if !ok {
				util.Log.Error().Msg(errUnauthorized.Error())

				http_response.SendForbiddenResponse(w, errUnauthorized)
				return
			}

			if !claims.Impersonate {
				http_response.SendForbiddenResponse(w, errors.New("you're not in impersonate mode"))
				return
			}

			if claims.Impersonate {
				if ok := IsExist(claims.UserName); !ok {
					err := errors.New("failed to get impersonate token")
					util.Log.Error().Msg(err.Error())

					http_response.SendForbiddenResponse(w, err)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
