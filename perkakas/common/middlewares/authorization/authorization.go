package authorization

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/tigapilarmandiri/perkakas/common/http_response"
	"github.com/tigapilarmandiri/perkakas/common/rds"
	"github.com/tigapilarmandiri/perkakas/common/util"
	"github.com/tigapilarmandiri/perkakas/configs"
	"go.elastic.co/apm/v2"
	"gorm.io/gorm"
)

var errUnauthorized = errors.New("you're not authorized")

type GetRedis func(context.Context, string) ([]byte, error)

var defaultGetRedis GetRedis = func(ctx context.Context, key string) ([]byte, error) {
	return rds.GetClient().Get(ctx, configs.Config.Redis.RedisAuthKey).Bytes()
}

// Authorization is to validate request from chi
func Authorization(f GetRedis) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if configs.Config.Env == "local" {
				next.ServeHTTP(w, r)
				return
			}
			ctx := r.Context()
			claims, ok := ctx.Value(util.ContextKey(util.ContextClaims)).(Claims)
			if !ok {
				util.Log.Error().Msg(errUnauthorized.Error())
				http_response.SendForbiddenResponse(w, errUnauthorized)
				return
			}

			if claims.IsSuperadmin {
				next.ServeHTTP(w, r)
				return
			}

			if f == nil {
				f = defaultGetRedis
			}

			b, err := f(ctx, configs.Config.Redis.RedisAuthKey)
			if err != nil {
				util.Log.Error().Msg(err.Error())
				http_response.SendForbiddenResponse(w, errUnauthorized)
				return
			}

			var permissions Permission
			err = json.Unmarshal(b, &permissions)
			if err != nil {
				util.Log.Error().Msg(err.Error())
				http_response.SendForbiddenResponse(w, errUnauthorized)
				return
			}

			if len(claims.Roles) == 0 {
				http_response.SendForbiddenResponse(w, errUnauthorized)
				return
			}

			if claims.IsSuperadmin && configs.Config.IsUseELK {
				trx := apm.TransactionFromContext(ctx)
				defer trx.End()

				trx.Context.SetLabel("username", claims.UserName)
			}

			path := r.URL.Path

			for _, v := range keys {
				if s := chi.URLParam(r, v); s != "" {
					path = strings.ReplaceAll(path, s, "*")
				}
			}

			path = strings.ReplaceAll(path, "/", "_")

			for _, v := range claims.Roles {
				permission, ok := permissions[v.Uuid]
				if !ok {
					continue
				}

				userPermission := strings.ToUpper(permission[path])
				switch r.Method {
				case "POST":
					if strings.Contains(userPermission, "C") {
						next.ServeHTTP(w, r)
						return
					}
				case "GET":
					if strings.Contains(userPermission, "R") {
						next.ServeHTTP(w, r)
						return
					}
				case "PATCH":
					if strings.Contains(userPermission, "U") {
						next.ServeHTTP(w, r)
						return
					}
				case "DELETE":
					if strings.Contains(userPermission, "D") {
						next.ServeHTTP(w, r)
						return
					}
				}
			}

			var roles []string
			for _, v := range claims.Roles {
				roles = append(roles, v.Name)
			}
			sRoles := strings.Join(roles, ", ")

			util.Log.Error().Msg(fmt.Sprintf("permission not permitted or not set: [%s] -> %s : %s", sRoles, r.Method, path))
			http_response.SendForbiddenResponse(w, errUnauthorized)
		}
		return http.HandlerFunc(fn)
	}
}

var (
	stmtQueryAuth                       *sql.Stmt
	stmtListPermittedWilayahId          *sql.Stmt
	stmtListPermittedWilayahIdParent    *sql.Stmt
	stmtListPermittedWilayahIdWithJenis *sql.Stmt

	stmtQueryAuthKepolisian            *sql.Stmt
	stmtListPermittedKepolisians       *sql.Stmt
	stmtListPermittedKepolisiansParent *sql.Stmt
)

// InitPreparedStatements is for ini sql.Stmt
// call it in application start
func InitPreparedStatements(gormDB *gorm.DB) error {
	var err error
	var db *sql.DB
	db, err = gormDB.DB()
	if err != nil {
		util.Log.Error().Msg(err.Error())
		return errUnauthorized
	}

	stmtQueryAuth, err = db.Prepare(`WITH RECURSIVE tree_view AS (
    SELECT id
    FROM wilayahs
    WHERE id = any ($1)
		and deleted_at is null
		and active = true

UNION ALL

    SELECT parent.id
    FROM wilayahs parent
    JOIN tree_view tv
      ON parent.parent_id = tv.id
		where parent.deleted_at is null
		and parent.active = true
)

SELECT count(1) from tree_view where id = $2`)
	if err != nil {
		return err
	}

	stmtListPermittedWilayahId, err = db.Prepare(`WITH RECURSIVE tree_view AS (
    SELECT id
    FROM wilayahs
    WHERE id = any ($1)
		and deleted_at is null
		and active = true

UNION ALL

    SELECT parent.id
    FROM wilayahs parent
    JOIN tree_view tv
      ON parent.parent_id = tv.id
		where parent.deleted_at is null
		and parent.active = true
)

SELECT id from tree_view`)
	if err != nil {
		return err
	}

	stmtListPermittedWilayahIdParent, err = db.Prepare(`WITH RECURSIVE tree_view AS (
    SELECT id,parent_id
    FROM wilayahs
    WHERE id = any ($1)
		and deleted_at is null
		and active = true

UNION ALL

    SELECT parent.id, parent.parent_id
    FROM tree_view tv
    JOIN wilayahs parent
      ON parent.id = tv.parent_id
		where parent.deleted_at is null
		and parent.active = true
)

SELECT id from tree_view`)
	if err != nil {
		return err
	}

	stmtListPermittedWilayahIdWithJenis, err = db.Prepare(`WITH RECURSIVE tree_view AS (
    SELECT id,
		jenis
    FROM wilayahs
    WHERE id = any ($1)
		and deleted_at is null
		and active = true

UNION ALL

    SELECT parent.id,
		parent.jenis
    FROM wilayahs parent
    JOIN tree_view tv
      ON parent.parent_id = tv.id
		where parent.deleted_at is null
		and parent.active = true
)

SELECT id from tree_view where jenis = $2`)

	return err
}

func InitPreparedStatementKepolisians(gormDB *gorm.DB) error {
	var err error
	var db *sql.DB
	db, err = gormDB.DB()
	if err != nil {
		util.Log.Error().Msg(err.Error())
		return errUnauthorized
	}

	stmtQueryAuthKepolisian, err = db.Prepare(`WITH RECURSIVE tree_view AS (
    SELECT id
    FROM kepolisians
    WHERE id = $1
		and deleted_at is null
		and active = true

UNION ALL

    SELECT parent.id
    FROM kepolisians parent
    JOIN tree_view tv
      ON parent.parent_id = tv.id
		where parent.deleted_at is null
		and parent.active = true
)

SELECT count(1) from tree_view where id = $2`)
	if err != nil {
		return err
	}

	stmtListPermittedKepolisians, err = db.Prepare(`WITH RECURSIVE tree_view AS (
    SELECT id
    FROM kepolisians
    WHERE id = $1
		and deleted_at is null
		and active = true

UNION ALL

    SELECT parent.id
    FROM kepolisians parent
    JOIN tree_view tv
      ON parent.parent_id = tv.id
		where parent.deleted_at is null
		and parent.active = true
)

SELECT id from tree_view`)
	if err != nil {
		return err
	}

	stmtListPermittedKepolisiansParent, err = db.Prepare(`WITH RECURSIVE tree_view AS (
	SELECT ID
		,
		parent_id
	FROM
		kepolisians
	WHERE
		ID = $1
		AND deleted_at IS NULL
		AND active = TRUE UNION ALL
	SELECT
		parent.ID,
		parent.parent_id
	FROM
		tree_view tv
		JOIN kepolisians parent ON tv.parent_id = parent.ID
	WHERE
		parent.deleted_at IS NULL
		AND parent.active = TRUE
	) SELECT
	id
FROM
	tree_view`)

	return err
}

var errUserWilayahEmpty = errors.New("user wilayah is empty")

// QueryAuthorization will validate user to access data
// userWilayahUuid is wilayah_uuid from JWT
// userWillAccessUuid is from what user will access data, its coming from user input
func QueryAuthorization(ctx context.Context, wilayahs []Wilayah, userWillAccessUuid string) error {
	if len(wilayahs) == 0 {
		return errUserWilayahEmpty
	}

	_, err := uuid.Parse(userWillAccessUuid)
	if err != nil {
		util.Log.Error().Msg(err.Error())
		return errUnauthorized
	}

	_, err = uuid.Parse(wilayahs[0].Uuid)
	if err != nil {
		util.Log.Error().Msg(err.Error())
		return errUnauthorized
	}

	var b strings.Builder

	b.WriteString("{" + wilayahs[0].Uuid)

	if len(wilayahs) > 1 {
		for _, v := range wilayahs[1:] {
			_, err = uuid.Parse(v.Uuid)
			if err != nil {
				util.Log.Error().Msg(err.Error())
				return errUnauthorized
			}
			b.WriteString("," + v.Uuid)
		}
	}

	b.WriteString("}")

	var count int64

	err = stmtQueryAuth.QueryRowContext(ctx, b.String(), userWillAccessUuid).
		Scan(&count)
	if err != nil {
		util.Log.Error().Msg(err.Error())
		return errUnauthorized
	}

	if count == 0 {
		return errUnauthorized
	}

	return nil
}

// GetListPermitedWilayahId is for get list permitted wilayah id based on user wilayah
// don't use this function if user level is super admin or level kepolisian is MABES
func GetListPermitedWilayahId(ctx context.Context, wilayahs []Wilayah, jenis string) ([]string, error) {
	if len(wilayahs) == 0 {
		return nil, errUserWilayahEmpty
	}

	_, err := uuid.Parse(wilayahs[0].Uuid)
	if err != nil {
		util.Log.Error().Msg(err.Error())
		return nil, err
	}

	var b strings.Builder

	b.WriteString("{" + wilayahs[0].Uuid)

	if len(wilayahs) > 1 {
		for _, v := range wilayahs[1:] {
			_, err = uuid.Parse(v.Uuid)
			if err != nil {
				util.Log.Error().Msg(err.Error())
				return nil, err
			}
			b.WriteString("," + v.Uuid)
		}
	}

	b.WriteString("}")

	args := []any{b.String()}

	var rows *sql.Rows
	if jenis == "" {
		rows, err = stmtListPermittedWilayahId.QueryContext(ctx, args...)
		if err != nil {
			util.Log.Error().Msg(err.Error())
			return nil, err
		}
	} else {
		args = append(args, jenis)
		rows, err = stmtListPermittedWilayahIdWithJenis.QueryContext(ctx, args...)
		if err != nil {
			util.Log.Error().Msg(err.Error())
			return nil, err
		}
	}
	defer rows.Close()

	var results []string
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			util.Log.Error().Msg(err.Error())
			return nil, err
		}
		results = append(results, id)
	}

	err = rows.Err()
	if err != nil {
		util.Log.Error().Msg(err.Error())
		return nil, err
	}

	if len(results) == 0 {
		return nil, errors.New("You can't access this data/s")
	}

	return results, nil
}

func GetListPermitedWilayahIdParent(ctx context.Context, wilayahs []Wilayah) ([]string, error) {
	if len(wilayahs) == 0 {
		return nil, errUserWilayahEmpty
	}

	_, err := uuid.Parse(wilayahs[0].Uuid)
	if err != nil {
		util.Log.Error().Msg(err.Error())
		return nil, err
	}

	var b strings.Builder

	b.WriteString("{" + wilayahs[0].Uuid)

	if len(wilayahs) > 1 {
		for _, v := range wilayahs[1:] {
			_, err = uuid.Parse(v.Uuid)
			if err != nil {
				util.Log.Error().Msg(err.Error())
				return nil, err
			}
			b.WriteString("," + v.Uuid)
		}
	}

	b.WriteString("}")

	var rows *sql.Rows
	rows, err = stmtListPermittedWilayahIdParent.QueryContext(ctx, b.String())
	if err != nil {
		util.Log.Error().Msg(err.Error())
		return nil, err
	}
	defer rows.Close()

	var results []string
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			util.Log.Error().Msg(err.Error())
			return nil, err
		}
		results = append(results, id)
	}

	err = rows.Err()
	if err != nil {
		util.Log.Error().Msg(err.Error())
		return nil, err
	}

	if len(results) == 0 {
		return nil, errors.New("You can't access this data/s")
	}

	return results, nil
}

func QueryAuthorizationKepolisian(ctx context.Context, kepolisianId uuid.UUID, userWillAccessUuid uuid.UUID) error {
	if kepolisianId.String() == "" || userWillAccessUuid.String() == "" {
		return errUnauthorized
	}

	var (
		err   error
		count int64
	)

	err = stmtQueryAuthKepolisian.QueryRowContext(ctx, kepolisianId.String(), userWillAccessUuid.String()).
		Scan(&count)
	if err != nil {
		util.Log.Error().Msg(err.Error())
		return errUnauthorized
	}

	if count == 0 {
		return errUnauthorized
	}

	return nil
}

func GetListPermitedKepolisianId(ctx context.Context, kepolisianId uuid.UUID) ([]string, error) {
	if kepolisianId.String() == "" {
		return nil, errUnauthorized
	}

	var (
		err  error
		rows *sql.Rows
	)

	rows, err = stmtListPermittedKepolisians.QueryContext(ctx, kepolisianId.String())
	if err != nil {
		util.Log.Error().Msg(err.Error())
		return nil, errUnauthorized
	}
	defer rows.Close()

	var results []string

	for rows.Next() {
		var each string
		err = rows.Scan(&each)
		if err != nil {
			util.Log.Error().Msg(err.Error())
			return nil, errUnauthorized
		}
		results = append(results, each)
	}

	err = rows.Err()
	if err != nil {
		util.Log.Error().Msg(err.Error())
		return nil, errUnauthorized
	}

	if len(results) == 0 {
		return nil, errUnauthorized
	}

	return results, nil
}

func GetListPermitedKepolisianIdParent(ctx context.Context, kepolisianId uuid.UUID) ([]string, error) {
	if kepolisianId.String() == "" {
		return nil, errUnauthorized
	}

	var (
		err  error
		rows *sql.Rows
	)

	rows, err = stmtListPermittedKepolisiansParent.QueryContext(ctx, kepolisianId.String())
	if err != nil {
		util.Log.Error().Msg(err.Error())
		return nil, errUnauthorized
	}
	defer rows.Close()

	var results []string

	for rows.Next() {
		var each string
		err = rows.Scan(&each)
		if err != nil {
			util.Log.Error().Msg(err.Error())
			return nil, errUnauthorized
		}
		results = append(results, each)
	}

	err = rows.Err()
	if err != nil {
		util.Log.Error().Msg(err.Error())
		return nil, errUnauthorized
	}

	if len(results) == 0 {
		return nil, errors.New("You can't access this data/s")
	}

	return results, nil
}

func IsUserLevelCannotAccess(claims Claims) bool {
	return strings.ToUpper(claims.KepolisianLevel) != "MABES" && !claims.IsSuperadmin
}
