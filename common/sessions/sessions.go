package sessions

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/tigapilarmandiri/perkakas/common/middlewares/authorization"
	"github.com/tigapilarmandiri/perkakas/common/rds"
	"github.com/tigapilarmandiri/perkakas/common/util"
	"github.com/tigapilarmandiri/perkakas/configs"
)

var genKey = func(key string) string {
	return fmt.Sprintf("s-%s", key)
}

func Store(ctx context.Context, token, remoteAddr, userAgent string) (err error) {
	var claim authorization.Claims
	claim, err = claimAuthToken(token)
	if err != nil {
		return
	}

	exp := time.Now().Add(24 * time.Hour * 7)
	if configs.Config.IsProduction() {
		exp = time.Now().Add(24 * time.Hour)
	}

	/* Remove old session first */
	err = Delete(ctx, claim.UserUUID)
	if err != nil {
		return
	}

	data := map[string]interface{}{
		"token":      token,
		"remoteAddr": remoteAddr,
		"userAgent":  userAgent,
	}

	var payload []byte
	payload, err = json.Marshal(data)
	if err != nil {
		return
	}

	/* Set new sessions */
	err = rds.GetClient().Set(ctx, genKey(claim.UserUUID), string(payload), time.Until(exp)).Err()

	return
}

func IsExist(ctx context.Context, userID string) (ok bool, err error) {

	var val int64
	val, err = rds.GetClient().Exists(ctx, genKey(userID)).Result()
	if err != nil {
		return
	}

	ok = val == 1

	return
}

func Delete(ctx context.Context, userID string) (err error) {
	err = rds.GetClient().Del(ctx, genKey(userID)).Err()
	return
}

func Clear(ctx context.Context, exclude ...string) (err error) {
	cmd := rds.GetClient().Keys(ctx, "*")

	if cmd.Err() != nil {
		err = cmd.Err()
		return
	}

	isExcluded := func(key string) bool {
		for _, v := range exclude {
			if strings.HasSuffix(key, v) {
				return true
			}
		}

		return false
	}

	for _, item := range cmd.Val() {
		if isExcluded(item) {
			continue
		}

		if strings.HasPrefix(item, "s-") {
			err = rds.GetClient().Del(ctx, item).Err()
			if err != nil {
				continue
			}
		}
	}

	return
}

func GetAll(ctx context.Context) (sessions []map[string]interface{}, err error) {
	cmd := rds.GetClient().Keys(ctx, "*")

	if cmd.Err() != nil {
		err = cmd.Err()
		return
	}

	for _, item := range cmd.Val() {
		if strings.HasPrefix(item, "s-") {
			payload := rds.GetClient().Get(ctx, item).Val()

			var sessionInfo map[string]interface{}
			if err = json.Unmarshal([]byte(payload), &sessionInfo); err != nil {
				return
			}

			if claim, err := claimAuthToken(sessionInfo["token"].(string)); err == nil {
				sessions = append(sessions, map[string]interface{}{
					"user_info": map[string]interface{}{
						"user_id":          claim.UserUUID,
						"username":         claim.UserName,
						"name":             claim.Name,
						"kepolisian_uuid":  claim.KepolisianUUID,
						"kepolisian_level": claim.KepolisianLevel,
					},
					"remote_addr": sessionInfo["remoteAddr"],
					"user_agent":  sessionInfo["userAgent"],
				})
			}

		}
	}

	return
}

func claimAuthToken(authToken string) (claim authorization.Claims, err error) {
	var token *jwt.Token
	token, err = jwt.Parse(authToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method :%v", token.Header["alg"])
		}

		return []byte(configs.Config.JWT.SecretKey), nil
	})

	if err != nil {
		util.Log.Error().Msg(err.Error())
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok || token.Valid) {
		err = fmt.Errorf("token is invalid")

		util.Log.Error().Msg(err.Error())
		return
	}

	var b []byte
	b, err = json.Marshal(claims)
	if err != nil {
		util.Log.Error().Msg(err.Error())
		return
	}

	err = json.Unmarshal(b, &claim)
	if err != nil {
		util.Log.Error().Msg(err.Error())
		return
	}

	return
}
