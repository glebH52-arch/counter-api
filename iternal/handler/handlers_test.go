package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockCounter struct {
	getCountFunc  func(ctx context.Context) (int64, error)
	incrCountFunc func(ctx context.Context) (int64, error)
}

func (m *mockCounter) GetCount(ctx context.Context) (int64, error) {
	if m.getCountFunc != nil {
		return m.getCountFunc(ctx)
	}

	return 0, nil
}

func (m *mockCounter) IncrCount(ctx context.Context) (int64, error) {
	if m.incrCountFunc != nil {
		return m.incrCountFunc(ctx)
	}

	return 0, nil
}

func TestNewCounterHandler(t *testing.T) {
	counterMock := &mockCounter{}

	handler := NewCounterHandler(counterMock)

	if handler == nil {
		t.Fatal("expected handler not to be nil")
	}

	if handler.Counter != counterMock {
		t.Fatal("expected counter dependency to be assigned")
	}
}

func TestCounterHandler_GetCount_Success(t *testing.T) {
	counterMock := &mockCounter{
		getCountFunc: func(ctx context.Context) (int64, error) {
			return 10, nil
		},
	}

	handler := NewCounterHandler(counterMock)

	req := httptest.NewRequest(http.MethodGet, "/counter", nil)
	rec := httptest.NewRecorder()

	handler.GetCount(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusOK,
			res.StatusCode,
		)
	}

	if contentType := res.Header.Get("Content-Type"); contentType != "application/json" {
		t.Errorf(
			"expected Content-Type application/json, got %q",
			contentType,
		)
	}

	if got := rec.Body.String(); got != "10\n" {
		t.Errorf("expected body %q, got %q", "10\n", got)
	}
}

func TestCounterHandler_GetCount_Zero(t *testing.T) {
	counterMock := &mockCounter{
		getCountFunc: func(ctx context.Context) (int64, error) {
			return 0, nil
		},
	}

	handler := NewCounterHandler(counterMock)

	req := httptest.NewRequest(http.MethodGet, "/counter", nil)
	rec := httptest.NewRecorder()

	handler.GetCount(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}

	if got := rec.Body.String(); got != "0\n" {
		t.Errorf("expected body %q, got %q", "0\n", got)
	}
}

func TestCounterHandler_GetCount_NegativeValue(t *testing.T) {
	counterMock := &mockCounter{
		getCountFunc: func(ctx context.Context) (int64, error) {
			return -5, nil
		},
	}

	handler := NewCounterHandler(counterMock)

	req := httptest.NewRequest(http.MethodGet, "/counter", nil)
	rec := httptest.NewRecorder()

	handler.GetCount(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}

	if got := rec.Body.String(); got != "-5\n" {
		t.Errorf("expected body %q, got %q", "-5\n", got)
	}
}

func TestCounterHandler_GetCount_ServiceError(t *testing.T) {
	expectedErr := errors.New("redis unavailable")

	counterMock := &mockCounter{
		getCountFunc: func(ctx context.Context) (int64, error) {
			return 0, expectedErr
		},
	}

	handler := NewCounterHandler(counterMock)

	req := httptest.NewRequest(http.MethodGet, "/counter", nil)
	rec := httptest.NewRecorder()

	handler.GetCount(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusServiceUnavailable,
			rec.Code,
		)
	}

	expectedBody := http.StatusText(http.StatusServiceUnavailable) + "\n"

	if got := rec.Body.String(); got != expectedBody {
		t.Errorf(
			"expected body %q, got %q",
			expectedBody,
			got,
		)
	}
}

func TestCounterHandler_GetCount_PassesRequestContext(t *testing.T) {
	type contextKey string

	const key contextKey = "test-key"

	contextReceived := false

	counterMock := &mockCounter{
		getCountFunc: func(ctx context.Context) (int64, error) {
			if ctx.Value(key) == "test-value" {
				contextReceived = true
			}

			return 1, nil
		},
	}

	handler := NewCounterHandler(counterMock)

	req := httptest.NewRequest(http.MethodGet, "/counter", nil)
	req = req.WithContext(
		context.WithValue(req.Context(), key, "test-value"),
	)

	rec := httptest.NewRecorder()

	handler.GetCount(rec, req)

	if !contextReceived {
		t.Fatal("expected handler to pass request context to counter service")
	}
}

func TestCounterHandler_IncrCount_Success(t *testing.T) {
	counterMock := &mockCounter{
		incrCountFunc: func(ctx context.Context) (int64, error) {
			return 11, nil
		},
	}

	handler := NewCounterHandler(counterMock)

	req := httptest.NewRequest(http.MethodPost, "/counter", nil)
	rec := httptest.NewRecorder()

	handler.IncrCount(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusOK,
			res.StatusCode,
		)
	}

	if contentType := res.Header.Get("Content-Type"); contentType != "application/json" {
		t.Errorf(
			"expected Content-Type application/json, got %q",
			contentType,
		)
	}

	if got := rec.Body.String(); got != "11\n" {
		t.Errorf("expected body %q, got %q", "11\n", got)
	}
}

func TestCounterHandler_IncrCount_Zero(t *testing.T) {
	counterMock := &mockCounter{
		incrCountFunc: func(ctx context.Context) (int64, error) {
			return 0, nil
		},
	}

	handler := NewCounterHandler(counterMock)

	req := httptest.NewRequest(http.MethodPost, "/counter", nil)
	rec := httptest.NewRecorder()

	handler.IncrCount(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}

	if got := rec.Body.String(); got != "0\n" {
		t.Errorf("expected body %q, got %q", "0\n", got)
	}
}

func TestCounterHandler_IncrCount_ServiceError(t *testing.T) {
	expectedErr := errors.New("redis unavailable")

	counterMock := &mockCounter{
		incrCountFunc: func(ctx context.Context) (int64, error) {
			return 0, expectedErr
		},
	}

	handler := NewCounterHandler(counterMock)

	req := httptest.NewRequest(http.MethodPost, "/counter", nil)
	rec := httptest.NewRecorder()

	handler.IncrCount(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusServiceUnavailable,
			rec.Code,
		)
	}

	expectedBody := http.StatusText(http.StatusServiceUnavailable) + "\n"

	if got := rec.Body.String(); got != expectedBody {
		t.Errorf(
			"expected body %q, got %q",
			expectedBody,
			got,
		)
	}
}

func TestCounterHandler_IncrCount_PassesRequestContext(t *testing.T) {
	type contextKey string

	const key contextKey = "test-key"

	contextReceived := false

	counterMock := &mockCounter{
		incrCountFunc: func(ctx context.Context) (int64, error) {
			if ctx.Value(key) == "test-value" {
				contextReceived = true
			}

			return 1, nil
		},
	}

	handler := NewCounterHandler(counterMock)

	req := httptest.NewRequest(http.MethodPost, "/counter", nil)
	req = req.WithContext(
		context.WithValue(req.Context(), key, "test-value"),
	)

	rec := httptest.NewRecorder()

	handler.IncrCount(rec, req)

	if !contextReceived {
		t.Fatal("expected handler to pass request context to counter service")
	}
}

func TestCounterHandler_Health_Success(t *testing.T) {
	handler := NewCounterHandler(&mockCounter{})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	handler.Health(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusOK,
			res.StatusCode,
		)
	}

	if contentType := res.Header.Get("Content-Type"); contentType != "application/json" {
		t.Errorf(
			"expected Content-Type application/json, got %q",
			contentType,
		)
	}

	expectedBody := "{\"status\":\"ok\"}\n"

	if got := rec.Body.String(); got != expectedBody {
		t.Errorf(
			"expected body %q, got %q",
			expectedBody,
			got,
		)
	}
}

func TestCounterHandler_Health_ResponseIsJSON(t *testing.T) {
	handler := NewCounterHandler(&mockCounter{})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	handler.Health(rec, req)

	body := strings.TrimSpace(rec.Body.String())

	if body != `{"status":"ok"}` {
		t.Errorf(
			"expected JSON body %q, got %q",
			`{"status":"ok"}`,
			body,
		)
	}
}
