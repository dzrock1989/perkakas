package authentication

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/tigapilarmandiri/perkakas/common/http_response"
	"github.com/tigapilarmandiri/perkakas/common/middlewares/authorization"
	"github.com/tigapilarmandiri/perkakas/common/middlewares/impersonate"
	"github.com/tigapilarmandiri/perkakas/common/sessions"
	"github.com/tigapilarmandiri/perkakas/common/util"
	"github.com/tigapilarmandiri/perkakas/configs"

	"github.com/golang-jwt/jwt/v4"
)

var (
	errDateOrTokenEmpty = errors.New("token or date is empty")

	errTokenNotValid = errors.New("token not valid")

	errDateNotValid   = errors.New("date not valid")
	errDateIsNotEpoch = errors.New("date is not epoch")
	errDateExpired    = errors.New("date is expired")
	errHmacNotValid   = errors.New("hmac not valid")
)

// Authentication is for validate the user
// if you want to protect it with JWT set isThereAJwt to true
// if no, set isThereAJwt to false
func Authentication(isThereAJwt bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authToken := r.Header.Get("Authorization")
			date := r.Header.Get("Dates")

			if isAuthorizationOrDateEmpty(authToken, date) {
				util.Log.Error().Msg(errDateOrTokenEmpty.Error())
				http_response.SendForbiddenResponse(w, errDateOrTokenEmpty)
				return
			}

			hmac_jwtToken, err := getHmacDate_JwtToken(authToken)
			if err != nil {
				util.Log.Error().Msg(err.Error())
				http_response.SendForbiddenResponse(w, err)
				return
			}

			authToken, err = validateDate(date, hmac_jwtToken)
			if err != nil {
				util.Log.Error().Msg(err.Error())
				http_response.SendForbiddenResponse(w, err)
				return
			}

			if !isThereAJwt {
				next.ServeHTTP(w, r)
				return
			}

			token, err := jwt.Parse(authToken, func(token *jwt.Token) (any, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method :%v", token.Header["alg"])
				}

				return []byte(configs.Config.JWT.SecretKey), nil
			})
			if err != nil {
				util.Log.Error().Msg(err.Error())
				// sometimes it will error token is expired
				http_response.SendForbiddenResponse(w, err)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				http_response.SendForbiddenResponse(w, errTokenNotValid)
				return
			}
			b, err := json.Marshal(claims)
			if err != nil {
				util.Log.Error().Msg(err.Error())
				http_response.SendForbiddenResponse(w, nil)
				return
			}
			ctx := context.WithValue(r.Context(), util.ContextKey(util.ContextClaimsBytes), b)

			var userInfo authorization.Claims

			err = json.Unmarshal(b, &userInfo)
			if err != nil {
				util.Log.Error().Msg(err.Error())
				http_response.SendForbiddenResponse(w, nil)
				return
			}

			if configs.Config.Env == "local" {
				ctx = context.WithValue(ctx, util.ContextKey(util.ContextClaims), userInfo)
				r = r.WithContext(ctx)

				next.ServeHTTP(w, r)
				return
			}

			// Check the sessions
			ok, err = sessions.IsExist(r.Context(), userInfo.UserUUID)
			if err != nil {
				util.Log.Error().Msg(err.Error())
				http_response.SendForbiddenResponse(w, errors.New("failed to get session"))
				return
			}

			if !ok && !userInfo.Impersonate {
				http_response.SendRedirectResponse(w, errors.New("session not found"))
				return
			}

			if userInfo.Impersonate {
				if ok := impersonate.IsExist(userInfo.UserName); !ok {
					err = errors.New("impersonate token not found or expired")

					util.Log.Error().Msg(err.Error())
					http_response.SendForbiddenResponse(w, err)
					return
				}
			}

			ctx = context.WithValue(ctx, util.ContextKey(util.ContextClaims), userInfo)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func isAuthorizationOrDateEmpty(authorization, date string) bool {
	return authorization == "" || date == ""
}

func getHmacDate_JwtToken(token string) (string, error) {
	arrToken := strings.Split(token, "Bearer ")
	if len(arrToken) != 2 {
		return "", errTokenNotValid
	}

	return arrToken[1], nil
}

func validateDate(date, hmac_jwtToken string) (string, error) {
	arrToken := strings.SplitN(hmac_jwtToken, "_", 2)

	if len(arrToken) != 2 {
		return "", errDateNotValid
	}

	// get hmac from header date
	sig := hmac.New(sha256.New, []byte(configs.Config.JWT.DateKey))
	sig.Write([]byte(date))

	hmac_date := hex.EncodeToString(sig.Sum(nil))

	// compare hmac and authentication
	if hmac_date != arrToken[0] {
		return "", errHmacNotValid
	}

	dif := time.Hour * 24 * 7 // 7 days
	if configs.Config.IsProduction() {
		dif = time.Hour * 5
	}

	epoch, err := strconv.Atoi(date)
	if err != nil {
		util.Log.Error().Msg(err.Error())
		return "", errDateIsNotEpoch
	}

	epochTime := time.UnixMilli(int64(epoch))
	since := time.Since(epochTime)
	if since < 0 {
		since *= -1
	}

	if since > dif {
		return "", errDateExpired
	}

	return arrToken[1], nil
}
