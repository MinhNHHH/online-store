package store

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
)

func Test_app_GetCategories(t *testing.T) {
	var tests = []struct {
		name               string
		queryParams        string
		expectedStatusCode int
	}{
		{
			name:               "get all categories",
			queryParams:        "",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "get categories with name filter",
			queryParams:        "?category_name=Electronics",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "get categories with pagination",
			queryParams:        "?page=1&page_size=5",
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, e := range tests {
		req, _ := http.NewRequest("GET", "/api/v1/categories"+e.queryParams, nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(app.GetCategories)
		handler.ServeHTTP(rr, req)

		if e.expectedStatusCode != rr.Code {
			t.Errorf("%s: returned wrong status code; expected %d but got %d", e.name, e.expectedStatusCode, rr.Code)
		}
	}
}

func Test_app_CreateCategory(t *testing.T) {
	var tests = []struct {
		name               string
		requestBody        string
		expectedStatusCode int
	}{
		{
			name:               "valid category",
			requestBody:        `{"name": "Electronics", "description": "Electronic devices and accessories"}`,
			expectedStatusCode: http.StatusCreated,
		},
	}

	for _, e := range tests {
		var reader io.Reader
		reader = bytes.NewBufferString(e.requestBody)
		req, _ := http.NewRequest("POST", "/api/v1/categories", reader)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(app.CreateCategory)
		handler.ServeHTTP(rr, req)

		if e.expectedStatusCode != rr.Code {
			t.Errorf("%s: returned wrong status code; expected %d but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		if e.expectedStatusCode == http.StatusCreated {
			var response struct {
				ID int `json:"id"`
			}
			err := json.NewDecoder(rr.Body).Decode(&response)
			if err != nil {
				t.Errorf("%s: failed to decode response: %v", e.name, err)
			}
			if response.ID <= 0 {
				t.Errorf("%s: expected positive ID, got %d", e.name, response.ID)
			}
		}
	}
}

func Test_app_UpdateCategory(t *testing.T) {
	var tests = []struct {
		name               string
		categoryID         string
		requestBody        string
		expectedStatusCode int
	}{
		{
			name:               "valid update",
			categoryID:         "1",
			requestBody:        `{"id": 1, "name": "Updated Electronics", "description": "Updated description"}`,
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, e := range tests {
		var reader io.Reader
		reader = bytes.NewBufferString(e.requestBody)
		req, _ := http.NewRequest("PUT", "/api/v1/categories/"+e.categoryID, reader)
		rr := httptest.NewRecorder()

		// Create a new chi router context
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", e.categoryID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		handler := http.HandlerFunc(app.UpdateCategory)
		handler.ServeHTTP(rr, req)

		if e.expectedStatusCode != rr.Code {
			t.Errorf("%s: returned wrong status code; expected %d but got %d", e.name, e.expectedStatusCode, rr.Code)
		}
	}
}

func Test_app_DeleteCategory(t *testing.T) {
	var tests = []struct {
		name               string
		categoryID         string
		expectedStatusCode int
	}{
		{
			name:               "valid delete",
			categoryID:         "1",
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, e := range tests {
		req, _ := http.NewRequest("DELETE", "/api/v1/categories/"+e.categoryID, nil)
		rr := httptest.NewRecorder()

		// Create a new chi router context
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", e.categoryID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		handler := http.HandlerFunc(app.DeleteCategory)
		handler.ServeHTTP(rr, req)

		if e.expectedStatusCode != rr.Code {
			t.Errorf("%s: returned wrong status code; expected %d but got %d", e.name, e.expectedStatusCode, rr.Code)
		}
	}
}
