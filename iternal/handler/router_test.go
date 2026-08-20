package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRouter_HealthRoute(t *testing.T) {
	counterMock := &mockCounter{}

	handler := NewCounterHandler(counterMock)
	router := NewRouter(*handler)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}

	expectedBody := "{\"status\":\"ok\"}\n"

	if rec.Body.String() != expectedBody {
		t.Errorf(
			"expected body %q, got %q",
			expectedBody,
			rec.Body.String(),
		)
	}
}

func TestNewRouter_GetCounterRoute(t *testing.T) {
	counterMock := &mockCounter{
		getCountFunc: func(ctx context.Context) (int64, error) {
			return 25, nil
		},
	}

	handler := NewCounterHandler(counterMock)
	router := NewRouter(*handler)

	req := httptest.NewRequest(http.MethodGet, "/counter", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}

	if rec.Body.String() != "25\n" {
		t.Errorf(
			"expected body %q, got %q",
			"25\n",
			rec.Body.String(),
		)
	}
}

func TestNewRouter_IncrementCounterRoute(t *testing.T) {
	counterMock := &mockCounter{
		incrCountFunc: func(ctx context.Context) (int64, error) {
			return 26, nil
		},
	}

	handler := NewCounterHandler(counterMock)
	router := NewRouter(*handler)

	req := httptest.NewRequest(
		http.MethodPost,
		"/counter/increment",
		nil,
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}

	if rec.Body.String() != "26\n" {
		t.Errorf(
			"expected body %q, got %q",
			"26\n",
			rec.Body.String(),
		)
	}
}

func TestNewRouter_UnknownRoute_ReturnsNotFound(t *testing.T) {
	counterMock := &mockCounter{}

	handler := NewCounterHandler(counterMock)
	router := NewRouter(*handler)

	req := httptest.NewRequest(
		http.MethodGet,
		"/unknown",
		nil,
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusNotFound,
			rec.Code,
		)
	}
}

func TestNewRouter_HealthWrongMethod_ReturnsMethodNotAllowed(t *testing.T) {
	counterMock := &mockCounter{}

	handler := NewCounterHandler(counterMock)
	router := NewRouter(*handler)

	req := httptest.NewRequest(
		http.MethodPost,
		"/health",
		nil,
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusMethodNotAllowed,
			rec.Code,
		)
	}
}

func TestNewRouter_GetCounterWrongMethod_ReturnsMethodNotAllowed(t *testing.T) {
	counterMock := &mockCounter{}

	handler := NewCounterHandler(counterMock)
	router := NewRouter(*handler)

	req := httptest.NewRequest(
		http.MethodPost,
		"/counter",
		nil,
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusMethodNotAllowed,
			rec.Code,
		)
	}
}

func TestNewRouter_IncrementWrongMethod_ReturnsMethodNotAllowed(t *testing.T) {
	counterMock := &mockCounter{}

	handler := NewCounterHandler(counterMock)
	router := NewRouter(*handler)

	req := httptest.NewRequest(
		http.MethodGet,
		"/counter/increment",
		nil,
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusMethodNotAllowed,
			rec.Code,
		)
	}
}

func TestNewRouter_CounterRouteCallsGetCount(t *testing.T) {
	called := false

	counterMock := &mockCounter{
		getCountFunc: func(ctx context.Context) (int64, error) {
			called = true
			return 1, nil
		},
	}

	handler := NewCounterHandler(counterMock)
	router := NewRouter(*handler)

	req := httptest.NewRequest(
		http.MethodGet,
		"/counter",
		nil,
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if !called {
		t.Fatal("expected GetCount to be called")
	}
}

func TestNewRouter_IncrementRouteCallsIncrCount(t *testing.T) {
	called := false

	counterMock := &mockCounter{
		incrCountFunc: func(ctx context.Context) (int64, error) {
			called = true
			return 1, nil
		},
	}

	handler := NewCounterHandler(counterMock)
	router := NewRouter(*handler)

	req := httptest.NewRequest(
		http.MethodPost,
		"/counter/increment",
		nil,
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if !called {
		t.Fatal("expected IncrCount to be called")
	}
}

func TestNewRouter_RoutesReturnJSONContentType(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
	}{
		{
			name:   "health",
			method: http.MethodGet,
			path:   "/health",
		},
		{
			name:   "get counter",
			method: http.MethodGet,
			path:   "/counter",
		},
		{
			name:   "increment counter",
			method: http.MethodPost,
			path:   "/counter/increment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counterMock := &mockCounter{
				getCountFunc: func(ctx context.Context) (int64, error) {
					return 1, nil
				},
				incrCountFunc: func(ctx context.Context) (int64, error) {
					return 2, nil
				},
			}

			handler := NewCounterHandler(counterMock)
			router := NewRouter(*handler)

			req := httptest.NewRequest(
				tt.method,
				tt.path,
				nil,
			)

			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			contentType := rec.Header().Get("Content-Type")

			if contentType != "application/json" {
				t.Errorf(
					"expected Content-Type %q, got %q",
					"application/json",
					contentType,
				)
			}
		})
	}
}
