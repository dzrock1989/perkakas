package authorization

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/tigapilarmandiri/perkakas/common/test"
	"github.com/tigapilarmandiri/perkakas/common/util"
	"github.com/tigapilarmandiri/perkakas/configs"
)

func setContext() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := Role{
				Uuid: "uuid-1",
				Name: "testing",
			}
			ctx := context.WithValue(r.Context(), util.ContextKey(util.ContextClaims), Claims{
				Roles: []Role{role},
			})
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func TestAuthorization(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	key := "secretKey"
	configs.Config.JWT.SecretKey = key
	configs.Config.JWT.DateKey = key

	authFunc := func(ctx context.Context, key string) ([]byte, error) {
		permissions := Permission{
			"uuid-1": map[string]string{
				"_users_*": "RU",
			},
		}

		return json.Marshal(permissions)
	}

	r := chi.NewRouter()
	r.Use(setContext())
	r.With(Authorization(authFunc)).Get("/users/{uuid}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.With(Authorization(authFunc)).Patch("/users/{uuid}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.With(Authorization(authFunc)).Delete("/users/{uuid}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.With(Authorization(authFunc)).Post("/users/{uuid}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	respUnauthorize := `{"status":403,"message":"forbidden access","data":null,"debug":{"error":true,"error_message":"you're not authorized"}}`

	if status, resp := test.TestRequest(t, ts, "POST", "/users/asdf", nil, nil); status != http.StatusForbidden || resp != respUnauthorize {
		t.Errorf(resp)
	}

	if status, resp := test.TestRequest(t, ts, "GET", "/users/asdf", nil, nil); status != http.StatusOK || resp != "welcome" {
		t.Errorf(resp)
	}

	if status, resp := test.TestRequest(t, ts, "PATCH", "/users/asdf", nil, nil); status != http.StatusOK || resp != "welcome" {
		t.Errorf(resp)
	}

	if status, resp := test.TestRequest(t, ts, "DELETE", "/users/asdf", nil, nil); status != http.StatusForbidden || resp != respUnauthorize {
		t.Errorf("%v %v \n", resp, status)
	}
}

func TestQueryValidation(t *testing.T) {
	tests := []struct {
		name     string
		expected error
		given1   []Wilayah
		given2   string
	}{
		{"empty", errUnauthorized, []Wilayah{{Name: "", Uuid: ""}}, ""},
		{"empty", errUserWilayahEmpty, nil, ""},
		{"uuid not valid", errUnauthorized, []Wilayah{{Name: "", Uuid: "asdhf"}}, "asdfaga"},
		{"uuid jwt not valid", errUnauthorized, []Wilayah{{Name: "", Uuid: "asdhf"}}, "63d9fefc-f566-4ff4-90f6-af42f57a6ab1"},
		{"uuid jwt 2 not valid", errUnauthorized, []Wilayah{{Uuid: "63d9fefc-f566-4ff4-90f6-af42f57a6ab1"}, {Name: "", Uuid: "asdhf"}}, "63d9fefc-f566-4ff4-90f6-af42f57a6ab1"},
		{"not permission", errUnauthorized, []Wilayah{{Uuid: "63d9fefc-f566-4ff4-90f6-af42f57a6ab1"}}, "0d51076f-1fdf-40eb-9bc0-1f939b75f376"},
		{"multiple not permission", errUnauthorized, []Wilayah{{Uuid: "63d9fefc-f566-4ff4-90f6-af42f57a6ab1"}, {Uuid: "63d9fefc-f566-4ff4-90f6-af42f57a6ab1"}}, "0d51076f-1fdf-40eb-9bc0-1f939b75f376"},
		{"valid", nil, []Wilayah{{Uuid: "63d9fefc-f566-4ff4-90f6-af42f57a6ab1"}}, "0c0477e6-341d-4062-afad-8cd0f37f8b27"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := QueryAuthorization(context.Background(), tt.given1, tt.given2)
			if actual != tt.expected {
				t.Errorf("(%s, %s): expected %v, actual %v", tt.given1, tt.given2, tt.expected, actual)
			}
		})
	}
}
