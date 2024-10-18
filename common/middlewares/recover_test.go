package middlewares

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tigapilarmandiri/perkakas/common/constant"
	"github.com/tigapilarmandiri/perkakas/common/http_response"
	"github.com/tigapilarmandiri/perkakas/common/test"
)

func panicingHandler(http.ResponseWriter, *http.Request) { panic("foo") }

func TestRecover(t *testing.T) {
	responseError := http_response.HttpResponse{
		Status:  http.StatusForbidden,
		Message: constant.MSG_FORBIDDEN_ACCESS,
		Data:    nil,
	}

	b, err := json.Marshal(responseError)
	if err != nil {
		t.Fatal(err)
	}

	r := chi.NewRouter()

	r.Use(Recover)

	r.Get("/", panicingHandler)

	ts := httptest.NewServer(r)
	defer ts.Close()

	if status, resp := test.TestRequest(t, ts, "GET", "/", nil, nil); status != http.StatusForbidden || resp != string(b) {
		t.Fatalf(resp)
	}
}

func TestRecoverWithLog(t *testing.T) {
	responseError := http_response.HttpResponse{
		Status:  http.StatusForbidden,
		Message: constant.MSG_FORBIDDEN_ACCESS,
		Data:    nil,
	}

	b, err := json.Marshal(responseError)
	if err != nil {
		t.Fatal(err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(Recover)

	r.Get("/", panicingHandler)

	ts := httptest.NewServer(r)
	defer ts.Close()

	if status, resp := test.TestRequest(t, ts, "GET", "/", nil, nil); status != http.StatusForbidden || resp != string(b) {
		t.Fatalf(resp)
	}
}
