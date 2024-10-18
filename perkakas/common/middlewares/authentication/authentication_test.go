package authentication

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/tigapilarmandiri/perkakas/common/test"
	"github.com/tigapilarmandiri/perkakas/configs"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
)

func TestIsAuthorizationOrDateEmpty(t *testing.T) {
	tests := []struct {
		name          string
		expected      bool
		authorization string
		date          string
	}{
		{"all empty", true, "", ""},
		{"authorization empty", true, "", "a"},
		{"date empty", true, "a", ""},
		{"not empty", false, "a", "a"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := isAuthorizationOrDateEmpty(tt.authorization, tt.date)
			if actual != tt.expected {
				t.Errorf("(%s, %s): expected %v, actual %v", tt.authorization, tt.date, tt.expected, actual)
			}
		})
	}
}

func TestGetHmacDate_JwtToken(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		err      error
		given    string
	}{
		{"empty", "", errTokenNotValid, ""},
		{"not valid", "", errTokenNotValid, "Bearerr asdf"},
		{"valid", "asdf", nil, "Bearer asdf"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual, err2 := getHmacDate_JwtToken(tt.given)
			if actual != tt.expected || err2 != tt.err {
				t.Errorf("(%s): expected %s, %s, actual %s, %s", tt.given, tt.expected, tt.err, actual, err2)
			}
		})
	}
}

// in this test, token is not validate
func TestValidateDate(t *testing.T) {
	key := "secret"
	configs.Config.JWT.DateKey = key
	configs.Config.JWT.SecretKey = key

	now := time.Now()
	dateExpired := strconv.Itoa(int(now.Add(-time.Hour * 24 * 8).UnixMilli()))
	dateAfter := strconv.Itoa(int(now.Add(time.Hour * 4).UnixMilli()))
	dateBefore := strconv.Itoa(int(now.Add(-time.Hour * 4).UnixMilli()))
	dateAfterExpired := strconv.Itoa(int(now.Add(time.Hour * 24 * 8).UnixMilli()))

	sig := hmac.New(sha256.New, []byte(key))
	sig.Write([]byte(dateExpired))

	hexExpired := hex.EncodeToString(sig.Sum(nil))

	date := strconv.Itoa(int(now.UnixMilli()))

	sig = hmac.New(sha256.New, []byte(key))
	sig.Write([]byte(date))

	hexValid := hex.EncodeToString(sig.Sum(nil))

	sig = hmac.New(sha256.New, []byte(key))
	sig.Write([]byte(dateAfter))

	hexAfter := hex.EncodeToString(sig.Sum(nil))

	sig = hmac.New(sha256.New, []byte(key))
	sig.Write([]byte(dateBefore))

	hexBefore := hex.EncodeToString(sig.Sum(nil))

	sig = hmac.New(sha256.New, []byte(key))
	sig.Write([]byte(dateAfterExpired))

	hexAfterExpired := hex.EncodeToString(sig.Sum(nil))

	sig = hmac.New(sha256.New, []byte(key))
	sig.Write([]byte("a" + date))

	hexDateNotValidEpoch := hex.EncodeToString(sig.Sum(nil))
	// in here, token is not validate
	token := "asdf"

	tests := []struct {
		name          string
		expected      string
		err           error
		date          string
		hmac_jwtToken string
	}{
		{"not split with _", "", errDateNotValid, date, hexValid + token},
		{"date and hex not same", "", errHmacNotValid, date, hexExpired + "_" + token},
		{"date before expired", "", errDateExpired, dateExpired, hexExpired + "_" + token},
		{"date after expired", "", errDateExpired, dateAfterExpired, hexAfterExpired + "_" + token},
		{"date is not epoch", "", errDateIsNotEpoch, "a" + date, hexDateNotValidEpoch + "_" + token},
		{"date after", token, nil, dateAfter, hexAfter + "_" + token},
		{"date before", token, nil, dateBefore, hexBefore + "_" + token},
		{"valid", token, nil, date, hexValid + "_" + token},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual, err2 := validateDate(tt.date, tt.hmac_jwtToken)
			if actual != tt.expected || err2 != tt.err {
				t.Errorf("(%s, %s): expected (%s, %s), actual (%s, %s)", tt.date, tt.hmac_jwtToken, tt.expected, tt.err, actual, err2)
			}
		})
	}
}

func TestNoJwt(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	r := chi.NewRouter()
	key := "secretKey"
	configs.Config.JWT.SecretKey = key
	configs.Config.JWT.DateKey = key

	r.Use(Authentication(false))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	// sending empty date
	respEmptyToken := `{"status":403,"message":"forbidden access","data":null,"debug":{"error":true,"error_message":"token or date is empty"}}`

	if status, resp := test.TestRequest(t, ts, "GET", "/", nil, nil); status != http.StatusForbidden || resp != respEmptyToken {
		t.Fatalf(resp)
	}

	// sending authorized request and expired date
	message := strconv.Itoa(int(time.Now().Add(-time.Hour * 24 * 8).UnixMilli()))

	sig := hmac.New(sha256.New, []byte(key))
	sig.Write([]byte(message))

	ttoken := hex.EncodeToString(sig.Sum(nil)) + "_"
	h := make(http.Header)
	h.Set("Authorization", "Bearer "+ttoken)
	h.Set("Dates", message)
	respReq := `{"status":403,"message":"forbidden access","data":null,"debug":{"error":true,"error_message":"date is expired"}}`
	if status, resp := test.TestRequest(t, ts, "GET", "/", h, nil); status != http.StatusForbidden || resp != respReq {
		t.Fatalf(resp)
	}

	// sending authorized request and  date
	message = strconv.Itoa(int(time.Now().Add(time.Hour * 24).UnixMilli()))

	sig = hmac.New(sha256.New, []byte(key))
	sig.Write([]byte(message))

	ttoken = hex.EncodeToString(sig.Sum(nil)) + "_"
	h = make(http.Header)
	h.Set("Authorization", "Bearer "+ttoken)
	h.Set("Dates", message)
	if status, resp := test.TestRequest(t, ts, "GET", "/", h, nil); status != http.StatusOK || resp != "welcome" {
		t.Fatalf("resp: " + resp)
	}
}

func TestAuthentication(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	r := chi.NewRouter()
	key := "secretKey"
	configs.Config.JWT.SecretKey = key
	configs.Config.JWT.DateKey = key

	r.Use(Authentication(true))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	// sending empty token and date
	respEmptyToken := `{"status":403,"message":"forbidden access","data":null,"debug":{"error":true,"error_message":"token or date is empty"}}`

	if status, resp := test.TestRequest(t, ts, "GET", "/", nil, nil); status != http.StatusForbidden || resp != respEmptyToken {
		t.Fatalf(resp)
	}

	// sending wrong key and empty date
	h := http.Header{}
	jwtToken := jwt.New(jwt.SigningMethodHS256)
	token, err := jwtToken.SignedString([]byte("wrong"))
	if err != nil {
		t.Fatal(err)
	}
	h.Set("Authorization", token)
	respWrongKey := `{"status":403,"message":"forbidden access","data":null,"debug":{"error":true,"error_message":"token or date is empty"}}`
	if status, resp := test.TestRequest(t, ts, "GET", "/", h, nil); status != http.StatusForbidden || resp != respWrongKey {
		t.Fatalf(resp)
	}

	// sending wrong jwt token and empty date
	h.Set("Authorization", "asdf")
	respWrongJWT := `{"status":403,"message":"forbidden access","data":null,"debug":{"error":true,"error_message":"token or date is empty"}}`
	if status, resp := test.TestRequest(t, ts, "GET", "/", h, nil); status != http.StatusForbidden || resp != respWrongJWT {
		t.Fatalf(resp)
	}

	// sending expired token and empty date
	jwtToken = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"exp": time.Now().Add(-30 * time.Second).Unix()})
	respTokenExpired := `{"status":403,"message":"forbidden access","data":null,"debug":{"error":true,"error_message":"token or date is empty"}}`
	token, err = jwtToken.SignedString([]byte(key))
	if err != nil {
		t.Fatal(err)
	}
	h = make(http.Header)
	h.Set("Authorization", token)
	if status, resp := test.TestRequest(t, ts, "GET", "/", h, nil); status != http.StatusForbidden || resp != respTokenExpired {
		t.Fatalf(resp)
	}

	// sending authorized requests and empty date
	jwtToken = jwt.New(jwt.SigningMethodHS256)
	token, err = jwtToken.SignedString([]byte(key))
	if err != nil {
		t.Fatal(err)
	}
	h = make(http.Header)
	h.Set("Authorization", token)
	if status, resp := test.TestRequest(t, ts, "GET", "/", h, nil); status != http.StatusForbidden || resp != respTokenExpired {
		t.Fatalf(resp)
	}

	// sending unauthorized requests and date
	jwtToken = jwt.New(jwt.SigningMethodHS256)
	token, err = jwtToken.SignedString([]byte(key))
	if err != nil {
		t.Fatal(err)
	}
	h = make(http.Header)
	h.Set("Authorization", token)
	h.Set("Dates", strconv.Itoa(int(time.Now().UnixMilli())))
	respReq := `{"status":403,"message":"forbidden access","data":null,"debug":{"error":true,"error_message":"token not valid"}}`
	if status, resp := test.TestRequest(t, ts, "GET", "/", h, nil); status != http.StatusForbidden || resp != respReq {
		t.Fatalf(resp)
	}

	// sending authorized request and expired date
	message := strconv.Itoa(int(time.Now().Add(-time.Hour * 24 * 8).UnixMilli()))

	sig := hmac.New(sha256.New, []byte(key))
	sig.Write([]byte(message))

	ttoken := hex.EncodeToString(sig.Sum(nil)) + "_" + token
	h = make(http.Header)
	h.Set("Authorization", "Bearer "+ttoken)
	h.Set("Dates", message)
	respReq = `{"status":403,"message":"forbidden access","data":null,"debug":{"error":true,"error_message":"date is expired"}}`
	if status, resp := test.TestRequest(t, ts, "GET", "/", h, nil); status != http.StatusForbidden || resp != respReq {
		t.Fatalf(resp)
	}

	// sending invalid token
	message = strconv.Itoa(int(time.Now().Add(-time.Hour * 24 * 8).UnixMilli()))

	sig = hmac.New(sha256.New, []byte(key))
	sig.Write([]byte(message))

	ttoken = hex.EncodeToString(sig.Sum(nil)) + token
	h = make(http.Header)
	h.Set("Authorization", "Bearer "+ttoken)
	h.Set("Dates", message)
	respReq = `{"status":403,"message":"forbidden access","data":null,"debug":{"error":true,"error_message":"date not valid"}}`
	if status, resp := test.TestRequest(t, ts, "GET", "/", h, nil); status != http.StatusForbidden || resp != respReq {
		t.Fatalf(resp)
	}

	// sending invalid hmac
	message = strconv.Itoa(int(time.Now().Add(-time.Hour*24*8).UnixMilli())) + "a"

	sig = hmac.New(sha256.New, []byte(key))
	sig.Write([]byte(message))

	ttoken = hex.EncodeToString(sig.Sum(nil)) + "_" + token
	h = make(http.Header)
	h.Set("Authorization", "Bearer "+ttoken)
	h.Set("Dates", message)
	respReq = `{"status":403,"message":"forbidden access","data":null,"debug":{"error":true,"error_message":"date is not epoch"}}`
	if status, resp := test.TestRequest(t, ts, "GET", "/", h, nil); status != http.StatusForbidden || resp != respReq {
		t.Fatalf(resp)
	}

	// sending invalid date
	message = strconv.Itoa(int(time.Now().Add(-time.Hour * 24 * 8).UnixMilli()))

	sig = hmac.New(sha256.New, []byte(key))
	sig.Write([]byte(message))

	ttoken = hex.EncodeToString(sig.Sum(nil)) + "_" + token
	h = make(http.Header)
	h.Set("Authorization", "Bearer "+ttoken)
	h.Set("Dates", message+"a")
	respReq = `{"status":403,"message":"forbidden access","data":null,"debug":{"error":true,"error_message":"hmac not valid"}}`
	if status, resp := test.TestRequest(t, ts, "GET", "/", h, nil); status != http.StatusForbidden || resp != respReq {
		t.Fatalf(resp)
	}

	// sending expired token
	message2 := strconv.Itoa(int(time.Now().UnixMilli()))

	sig = hmac.New(sha256.New, []byte(key))
	sig.Write([]byte(message))

	jwtToken = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"exp": time.Now().Add(-30 * time.Second).Unix()})
	token2, err := jwtToken.SignedString([]byte(key))
	if err != nil {
		t.Fatal(err)
	}

	ttoken = hex.EncodeToString(sig.Sum(nil)) + "_" + token2
	h = make(http.Header)
	h.Set("Authorization", "Bearer "+ttoken)
	h.Set("Dates", message2)
	respReq = `{"status":403,"message":"forbidden access","data":null,"debug":{"error":true,"error_message":"hmac not valid"}}`
	if status, resp := test.TestRequest(t, ts, "GET", "/", h, nil); status != http.StatusForbidden || resp != respReq {
		t.Fatalf(resp)
	}

	// sending expired token
	message = strconv.Itoa(int(time.Now().UnixMilli()))

	sig = hmac.New(sha256.New, []byte(key))
	sig.Write([]byte(message))

	ttoken = hex.EncodeToString(sig.Sum(nil)) + "_" + "asdf"
	h = make(http.Header)
	h.Set("Authorization", "Bearer "+ttoken)
	h.Set("Dates", message)
	respReq = `{"status":403,"message":"forbidden access","data":null,"debug":{"error":true,"error_message":"token contains an invalid number of segments"}}`
	if status, resp := test.TestRequest(t, ts, "GET", "/", h, nil); status != http.StatusForbidden || resp != respReq {
		t.Fatalf(resp)
	}

	// sending authorized request and date
	message = strconv.Itoa(int(time.Now().UnixMilli()))

	sig = hmac.New(sha256.New, []byte(key))
	sig.Write([]byte(message))

	ttoken = hex.EncodeToString(sig.Sum(nil)) + "_" + token
	h = make(http.Header)
	h.Set("Authorization", "Bearer "+ttoken)
	h.Set("Dates", message)
	if status, resp := test.TestRequest(t, ts, "GET", "/", h, nil); status != 200 || resp != "welcome" {
		t.Fatalf(resp)
	}
}

func TestAuthenticationProduction(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	r := chi.NewRouter()
	key := "secretKey"
	configs.Config.JWT.SecretKey = key
	configs.Config.JWT.DateKey = key
	configs.Config.Env = "production"

	r.Use(Authentication(true))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	// sending empty token and date
	respEmptyToken := `{"status":403,"message":"forbidden access","data":null}`

	if status, resp := test.TestRequest(t, ts, "GET", "/", nil, nil); status != http.StatusForbidden || resp != respEmptyToken {
		t.Fatalf(resp)
	}

	// sending wrong key and empty date
	h := http.Header{}
	jwtToken := jwt.New(jwt.SigningMethodHS256)
	token, err := jwtToken.SignedString([]byte("wrong"))
	if err != nil {
		t.Fatal(err)
	}
	h.Set("Authorization", token)
	respWrongKey := `{"status":403,"message":"forbidden access","data":null}`
	if status, resp := test.TestRequest(t, ts, "GET", "/", h, nil); status != http.StatusForbidden || resp != respWrongKey {
		t.Fatalf(resp)
	}

	// sending wrong jwt token and empty date
	h.Set("Authorization", "asdf")
	respWrongJWT := `{"status":403,"message":"forbidden access","data":null}`
	if status, resp := test.TestRequest(t, ts, "GET", "/", h, nil); status != http.StatusForbidden || resp != respWrongJWT {
		t.Fatalf(resp)
	}

	// sending expired token and empty date
	jwtToken = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"exp": time.Now().Add(-30 * time.Second).Unix()})
	respTokenExpired := `{"status":403,"message":"forbidden access","data":null}`
	token, err = jwtToken.SignedString([]byte(key))
	if err != nil {
		t.Fatal(err)
	}
	h = make(http.Header)
	h.Set("Authorization", token)
	if status, resp := test.TestRequest(t, ts, "GET", "/", h, nil); status != http.StatusForbidden || resp != respTokenExpired {
		t.Fatalf(resp)
	}

	// sending authorized requests and empty date
	jwtToken = jwt.New(jwt.SigningMethodHS256)
	token, err = jwtToken.SignedString([]byte(key))
	if err != nil {
		t.Fatal(err)
	}
	h = make(http.Header)
	h.Set("Authorization", token)
	if status, resp := test.TestRequest(t, ts, "GET", "/", h, nil); status != http.StatusForbidden || resp != respTokenExpired {
		t.Fatalf(resp)
	}

	// sending unauthorized requests and date
	jwtToken = jwt.New(jwt.SigningMethodHS256)
	token, err = jwtToken.SignedString([]byte(key))
	if err != nil {
		t.Fatal(err)
	}
	h = make(http.Header)
	h.Set("Authorization", token)
	h.Set("Dates", strconv.Itoa(int(time.Now().UnixMilli())))
	respReq := `{"status":403,"message":"forbidden access","data":null}`
	if status, resp := test.TestRequest(t, ts, "GET", "/", h, nil); status != http.StatusForbidden || resp != respReq {
		t.Fatalf(resp)
	}

	// sending authorized request and expired date
	message := strconv.Itoa(int(time.Now().Add(-time.Hour * 24 * 8).UnixMilli()))

	sig := hmac.New(sha256.New, []byte(key))
	sig.Write([]byte(message))

	ttoken := hex.EncodeToString(sig.Sum(nil)) + "_" + token
	h = make(http.Header)
	h.Set("Authorization", "Bearer "+ttoken)
	h.Set("Dates", message)
	respReq = `{"status":403,"message":"forbidden access","data":null}`
	if status, resp := test.TestRequest(t, ts, "GET", "/", h, nil); status != http.StatusForbidden || resp != respReq {
		t.Fatalf(resp)
	}

	// sending invalid token
	message = strconv.Itoa(int(time.Now().Add(-time.Hour * 24 * 8).UnixMilli()))

	sig = hmac.New(sha256.New, []byte(key))
	sig.Write([]byte(message))

	ttoken = hex.EncodeToString(sig.Sum(nil)) + token
	h = make(http.Header)
	h.Set("Authorization", "Bearer "+ttoken)
	h.Set("Dates", message)
	respReq = `{"status":403,"message":"forbidden access","data":null}`
	if status, resp := test.TestRequest(t, ts, "GET", "/", h, nil); status != http.StatusForbidden || resp != respReq {
		t.Fatalf(resp)
	}

	// sending invalid hmac
	message = strconv.Itoa(int(time.Now().Add(-time.Hour*24*8).UnixMilli())) + "a"

	sig = hmac.New(sha256.New, []byte(key))
	sig.Write([]byte(message))

	ttoken = hex.EncodeToString(sig.Sum(nil)) + "_" + token
	h = make(http.Header)
	h.Set("Authorization", "Bearer "+ttoken)
	h.Set("Dates", message)
	respReq = `{"status":403,"message":"forbidden access","data":null}`
	if status, resp := test.TestRequest(t, ts, "GET", "/", h, nil); status != http.StatusForbidden || resp != respReq {
		t.Fatalf(resp)
	}

	// sending invalid date
	message = strconv.Itoa(int(time.Now().Add(-time.Hour * 24 * 8).UnixMilli()))

	sig = hmac.New(sha256.New, []byte(key))
	sig.Write([]byte(message))

	ttoken = hex.EncodeToString(sig.Sum(nil)) + "_" + token
	h = make(http.Header)
	h.Set("Authorization", "Bearer "+ttoken)
	h.Set("Dates", message+"a")
	respReq = `{"status":403,"message":"forbidden access","data":null}`
	if status, resp := test.TestRequest(t, ts, "GET", "/", h, nil); status != http.StatusForbidden || resp != respReq {
		t.Fatalf(resp)
	}

	// sending expired token
	message2 := strconv.Itoa(int(time.Now().UnixMilli()))

	sig = hmac.New(sha256.New, []byte(key))
	sig.Write([]byte(message))

	jwtToken = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"exp": time.Now().Add(-30 * time.Second).Unix()})
	token2, err := jwtToken.SignedString([]byte(key))
	if err != nil {
		t.Fatal(err)
	}

	ttoken = hex.EncodeToString(sig.Sum(nil)) + "_" + token2
	h = make(http.Header)
	h.Set("Authorization", "Bearer "+ttoken)
	h.Set("Dates", message2)
	respReq = `{"status":403,"message":"forbidden access","data":null}`
	if status, resp := test.TestRequest(t, ts, "GET", "/", h, nil); status != http.StatusForbidden || resp != respReq {
		t.Fatalf(resp)
	}

	// sending expired token
	message = strconv.Itoa(int(time.Now().UnixMilli()))

	sig = hmac.New(sha256.New, []byte(key))
	sig.Write([]byte(message))

	ttoken = hex.EncodeToString(sig.Sum(nil)) + "_" + "asdf"
	h = make(http.Header)
	h.Set("Authorization", "Bearer "+ttoken)
	h.Set("Dates", message)
	respReq = `{"status":403,"message":"forbidden access","data":null}`
	if status, resp := test.TestRequest(t, ts, "GET", "/", h, nil); status != http.StatusForbidden || resp != respReq {
		t.Fatalf(resp)
	}

	// sending authorized request and date
	message = strconv.Itoa(int(time.Now().UnixMilli()))

	sig = hmac.New(sha256.New, []byte(key))
	sig.Write([]byte(message))

	ttoken = hex.EncodeToString(sig.Sum(nil)) + "_" + token
	h = make(http.Header)
	h.Set("Authorization", "Bearer "+ttoken)
	h.Set("Dates", message)
	if status, resp := test.TestRequest(t, ts, "GET", "/", h, nil); status != 200 || resp != "welcome" {
		t.Fatalf(resp)
	}
}
