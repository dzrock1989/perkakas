package http_response

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/dzrock1989/perkakas/common/constant"
	"github.com/dzrock1989/perkakas/common/test"
	"github.com/dzrock1989/perkakas/common/util"
	"github.com/dzrock1989/perkakas/configs"
)

type testValidation struct {
	Id    int    `json:"id" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

func newServer() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SendSuccess(w, "success", nil, nil)
	}))

	router.HandleFunc("/forbidden", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		unsupportedValues := []any{
			math.NaN(),
		}

		SendSuccess(w, unsupportedValues, nil, nil)
	}))

	router.HandleFunc("/not-found", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SendNotFoundResponse(w, errors.New("just message"))
	}))

	router.HandleFunc("/with-meta", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		meta := &Meta{
			Page:      1,
			TotalPage: 1,
			TotalData: 2,
		}
		SendSuccess(w, []string{"Hello", "World"}, meta, nil)
	}))

	router.HandleFunc("/validation", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tv testValidation
		// tv.Id =
		tv.Email = "not valid"
		err := util.ValidateStruct(r.Context(), tv)
		SendForbiddenResponse(w, err)
	}))

	return router
}

func TestHttpResponse(t *testing.T) {
	svr := httptest.NewServer(newServer())
	defer svr.Close()

	i, s := test.TestRequest(t, svr, "GET", "/", nil, nil)

	var resp HttpResponse
	if err := json.Unmarshal([]byte(s), &resp); err != nil {
		require.Equal(t, nil, err, "Should not error")
	}

	require.Equal(t, http.StatusOK, resp.Status, "Should return 200")
	require.Equal(t, http.StatusOK, i, "Should return 200")
	require.Equal(t, constant.MSG_SUCCESS, resp.Message, "Should return success")
	require.Equal(t, "success", resp.Data, "Should return success")
	require.Equal(t, (*Meta)(nil), resp.Meta, "Should return nil")

	// Test on failed to marshal data
	i, s = test.TestRequest(t, svr, "GET", "/forbidden", nil, nil)

	resForbidden := `{"status":403,"message":"forbidden access","data":null}`

	require.Equal(t, http.StatusForbidden, i, "Should return 403")
	require.Equal(t, resForbidden, s, "Should return forbidden access")

	// Test get meta
	i, s = test.TestRequest(t, svr, "GET", "/with-meta", nil, nil)

	if err := json.Unmarshal([]byte(s), &resp); err != nil {
		require.Equal(t, nil, err, "Should not error")
	}

	require.Equal(t, http.StatusOK, resp.Status, "Should return 200")
	require.Equal(t, http.StatusOK, i, "Should return 200")
	require.Equal(t, "success", resp.Message, "Should return success")
	require.NotEqual(t, nil, resp.Meta, "Should not return nil")
	require.Equal(t, 1, resp.Meta.Page, "Should return page 1")
	require.Equal(t, 1, resp.Meta.TotalPage, "Should return total page 1")
	require.Equal(t, 2, resp.Meta.TotalData, "Should return total data 2")

	// Test on not found
	i, s = test.TestRequest(t, svr, "GET", "/not-found", nil, nil)

	resNotFound := `{"status":404,"message":"not found","data":null,"debug":{"error":true,"error_message":"just message"}}`

	require.Equal(t, http.StatusNotFound, i, "Should return 404")
	require.Equal(t, resNotFound, s, "Should return not found")

	// Test on validation
	i, s = test.TestRequest(t, svr, "GET", "/validation", nil, nil)

	resValidation := `{"status":403,"message":"forbidden access","data":null,"debug":{"error":true,"error_message":["Id wajib diisi","Email harus berupa alamat email yang valid"]}}`

	require.Equal(t, http.StatusForbidden, i, "Should return 403")
	require.Equal(t, resValidation, s, "Should return `"+resValidation+"`, got `"+s+"`")
}

func TestHttpResponseProduction(t *testing.T) {
	svr := httptest.NewServer(newServer())
	configs.Config.Env = "production"
	defer svr.Close()

	i, s := test.TestRequest(t, svr, "GET", "/", nil, nil)

	var resp HttpResponse
	if err := json.Unmarshal([]byte(s), &resp); err != nil {
		require.Equal(t, nil, err, "Should not error")
	}

	require.Equal(t, http.StatusOK, resp.Status, "Should return 200")
	require.Equal(t, http.StatusOK, i, "Should return 200")
	require.Equal(t, "success", resp.Message, "Should return success")
	require.Equal(t, (*Meta)(nil), resp.Meta, "Should return nil")

	// Test on failed to marshal data
	i, s = test.TestRequest(t, svr, "GET", "/forbidden", nil, nil)

	resForbidden := `{"status":403,"message":"forbidden access","data":null}`

	require.Equal(t, http.StatusForbidden, i, "Should return 403")
	require.Equal(t, http.StatusForbidden, i, "Should return 403")
	require.Equal(t, resForbidden, s, "Should return forbidden access")

	// Test get meta
	i, s = test.TestRequest(t, svr, "GET", "/with-meta", nil, nil)

	if err := json.Unmarshal([]byte(s), &resp); err != nil {
		require.Equal(t, nil, err, "Should not error")
	}

	require.Equal(t, http.StatusOK, resp.Status, "Should return 200")
	require.Equal(t, http.StatusOK, i, "Should return 200")
	require.Equal(t, "success", resp.Message, "Should return success")
	require.NotEqual(t, nil, resp.Meta, "Should not return nil")
	require.Equal(t, 1, resp.Meta.Page, "Should return page 1")
	require.Equal(t, 1, resp.Meta.TotalPage, "Should return total page 1")
	require.Equal(t, 2, resp.Meta.TotalData, "Should return total data 2")

	// Test on not found
	i, s = test.TestRequest(t, svr, "GET", "/not-found", nil, nil)

	resNotFound := `{"status":404,"message":"not found","data":null}`

	require.Equal(t, http.StatusNotFound, i, "Should return 404")
	require.Equal(t, http.StatusNotFound, i, "Should return 404")
	require.Equal(t, resNotFound, s, "Should return not found")

	// Test on validation
	i, s = test.TestRequest(t, svr, "GET", "/validation", nil, nil)

	require.Equal(t, http.StatusForbidden, i, "Should return 403")
	require.Equal(t, resForbidden, s, "Should return forbidden")
}
