package counter

import (
	"context"
	"errors"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
)

func TestNewCounterRedisService(t *testing.T) {
	db, _ := redismock.NewClientMock()

	service := NewCounterRedisService(db)

	if service == nil {
		t.Fatal("expected service not to be nil")
	}

	if service.RedisClient != db {
		t.Fatal("expected RedisClient to be assigned")
	}
}

func TestCounterRedisService_IncrCount_Success(t *testing.T) {
	db, mock := redismock.NewClientMock()
	service := NewCounterRedisService(db)

	ctx := context.Background()

	mock.ExpectIncr(counterKey).SetVal(1)

	got, err := service.IncrCount(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 1 {
		t.Errorf("expected count 1, got %d", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet redis expectations: %v", err)
	}
}

func TestCounterRedisService_IncrCount_MultipleValue(t *testing.T) {
	db, mock := redismock.NewClientMock()
	service := NewCounterRedisService(db)

	ctx := context.Background()

	mock.ExpectIncr(counterKey).SetVal(42)

	got, err := service.IncrCount(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 42 {
		t.Errorf("expected count 42, got %d", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet redis expectations: %v", err)
	}
}

func TestCounterRedisService_IncrCount_RedisNil(t *testing.T) {
	db, mock := redismock.NewClientMock()
	service := NewCounterRedisService(db)

	ctx := context.Background()

	mock.ExpectIncr(counterKey).SetErr(redis.Nil)

	got, err := service.IncrCount(ctx)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if got != 0 {
		t.Errorf("expected count 0, got %d", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet redis expectations: %v", err)
	}
}

func TestCounterRedisService_IncrCount_RedisError(t *testing.T) {
	db, mock := redismock.NewClientMock()
	service := NewCounterRedisService(db)

	ctx := context.Background()

	expectedErr := errors.New("redis unavailable")

	mock.ExpectIncr(counterKey).SetErr(expectedErr)

	got, err := service.IncrCount(ctx)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}

	if got != 0 {
		t.Errorf("expected count 0, got %d", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet redis expectations: %v", err)
	}
}

func TestCounterRedisService_IncrCount_ContextCanceled(t *testing.T) {
	db, _ := redismock.NewClientMock()
	service := NewCounterRedisService(db)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	got, err := service.IncrCount(ctx)

	if err == nil {
		t.Fatal("expected context error, got nil")
	}

	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}

	if got != 0 {
		t.Errorf("expected count 0, got %d", got)
	}
}

func TestCounterRedisService_GetCount_Success(t *testing.T) {
	db, mock := redismock.NewClientMock()
	service := NewCounterRedisService(db)

	ctx := context.Background()

	mock.ExpectGet(counterKey).SetVal("10")

	got, err := service.GetCount(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 10 {
		t.Errorf("expected count 10, got %d", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet redis expectations: %v", err)
	}
}

func TestCounterRedisService_GetCount_Zero(t *testing.T) {
	db, mock := redismock.NewClientMock()
	service := NewCounterRedisService(db)

	ctx := context.Background()

	mock.ExpectGet(counterKey).SetVal("0")

	got, err := service.GetCount(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 0 {
		t.Errorf("expected count 0, got %d", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet redis expectations: %v", err)
	}
}

func TestCounterRedisService_GetCount_NegativeValue(t *testing.T) {
	db, mock := redismock.NewClientMock()
	service := NewCounterRedisService(db)

	ctx := context.Background()

	mock.ExpectGet(counterKey).SetVal("-5")

	got, err := service.GetCount(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != -5 {
		t.Errorf("expected count -5, got %d", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet redis expectations: %v", err)
	}
}

func TestCounterRedisService_GetCount_RedisNil(t *testing.T) {
	db, mock := redismock.NewClientMock()
	service := NewCounterRedisService(db)

	ctx := context.Background()

	mock.ExpectGet(counterKey).RedisNil()

	got, err := service.GetCount(ctx)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if got != 0 {
		t.Errorf("expected count 0, got %d", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet redis expectations: %v", err)
	}
}

func TestCounterRedisService_GetCount_RedisError(t *testing.T) {
	db, mock := redismock.NewClientMock()
	service := NewCounterRedisService(db)

	ctx := context.Background()

	expectedErr := errors.New("redis connection error")

	mock.ExpectGet(counterKey).SetErr(expectedErr)

	got, err := service.GetCount(ctx)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}

	if got != 0 {
		t.Errorf("expected count 0, got %d", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet redis expectations: %v", err)
	}
}

func TestCounterRedisService_GetCount_ContextCanceled(t *testing.T) {
	db, _ := redismock.NewClientMock()
	service := NewCounterRedisService(db)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	got, err := service.GetCount(ctx)

	if err == nil {
		t.Fatal("expected context error, got nil")
	}

	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}

	if got != 0 {
		t.Errorf("expected count 0, got %d", got)
	}
}

func TestCounterRedisService_GetCount_InvalidInteger(t *testing.T) {
	db, mock := redismock.NewClientMock()
	service := NewCounterRedisService(db)

	ctx := context.Background()

	mock.ExpectGet(counterKey).SetVal("not-a-number")

	got, err := service.GetCount(ctx)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if got != 0 {
		t.Errorf("expected count 0, got %d", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet redis expectations: %v", err)
	}
}
