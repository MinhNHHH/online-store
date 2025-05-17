package store

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MinhNHHH/online-store/pkg/databases/schema"
)

func Test_app_EnableCORS(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	app := OnlineStore{}
	var tests = []struct {
		name           string
		method         string
		expectedHeader bool
	}{
		{"get", http.MethodGet, true},
		{"post", http.MethodPost, true},
		{"put", http.MethodPut, true},
		{"delete", http.MethodDelete, true},
	}

	for _, e := range tests {
		handlerToTest := app.enableCORS(nextHandler)

		req := httptest.NewRequest(e.method, "http://testing", nil)
		rr := httptest.NewRecorder()

		handlerToTest.ServeHTTP(rr, req)
		if e.expectedHeader && rr.Header().Get("Access-Control-Allow-Origin") == "" {
			t.Errorf("%s: expected header, but did not find it", e.name)
		}

		if !e.expectedHeader && rr.Header().Get("Access-Control-Allow-Credentials") != "" {
			t.Errorf("%s: expected no header, but got one", e.name)
		}
	}
}

func Test_app_authRequired(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	app := OnlineStore{}

	testUser := schema.User{
		ID:       1,
		Name:     "Admin",
		Email:    "admin@example.com",
		Password: "password",
		IsAdmin:  true,
	}

	tokens, _ := app.generateTokenPair(&testUser)
	var tests = []struct {
		name             string
		token            string
		expectAuthorized bool
		setHeader        bool
	}{
		{name: "valid token", token: fmt.Sprintf("Bearer %s", tokens.Token), expectAuthorized: true, setHeader: true},
		{name: "no token", token: "", expectAuthorized: false, setHeader: false},
		{name: "invalid token", token: fmt.Sprintf("Bearer %s1", tokens.Token), expectAuthorized: false, setHeader: true},
	}
	for _, e := range tests {
		req, _ := http.NewRequest("GET", "/", nil)
		if e.setHeader {
			req.Header.Set("Authorization", e.token)
		}

		rr := httptest.NewRecorder()
		handlerToTest := app.authRequired(nextHandler)
		handlerToTest.ServeHTTP(rr, req)

		if e.expectAuthorized && rr.Code == http.StatusUnauthorized {
			t.Errorf("%s: got ode 401, and should not have", e.name)
		}

		if !e.expectAuthorized && rr.Code != http.StatusUnauthorized {
			t.Errorf("%s: did not get code 401, and should have", e.name)
		}
	}
}
