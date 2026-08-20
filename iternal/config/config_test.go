package config

import (
	"os"
	"strings"
	"testing"
)

func TestLoad_Success(t *testing.T) {
	t.Setenv("REDIS_ADDR", "localhost:6379")
	t.Setenv("REDIS_PASSWORD", "secret")
	t.Setenv("REDIS_DB", "2")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg == nil {
		t.Fatal("expected config not to be nil")
	}

	if cfg.RedisAddr != "localhost:6379" {
		t.Errorf(
			"expected RedisAddr %q, got %q",
			"localhost:6379",
			cfg.RedisAddr,
		)
	}

	if cfg.RedisPassword != "secret" {
		t.Errorf(
			"expected RedisPassword %q, got %q",
			"secret",
			cfg.RedisPassword,
		)
	}

	if cfg.RedisDb != 2 {
		t.Errorf(
			"expected RedisDb %d, got %d",
			2,
			cfg.RedisDb,
		)
	}
}

func TestLoad_EmptyPasswordAllowed(t *testing.T) {
	t.Setenv("REDIS_ADDR", "localhost:6379")
	t.Setenv("REDIS_PASSWORD", "")
	t.Setenv("REDIS_DB", "0")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.RedisPassword != "" {
		t.Errorf(
			"expected empty password, got %q",
			cfg.RedisPassword,
		)
	}

	if cfg.RedisDb != 0 {
		t.Errorf("expected RedisDb 0, got %d", cfg.RedisDb)
	}
}

func TestLoad_MissingRedisAddr(t *testing.T) {
	t.Setenv("REDIS_ADDR", "")
	t.Setenv("REDIS_PASSWORD", "")
	t.Setenv("REDIS_DB", "0")

	cfg, err := Load()

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if cfg != nil {
		t.Fatal("expected config to be nil")
	}

	if !strings.Contains(err.Error(), "REDIS_ADDR") {
		t.Errorf(
			"expected error to contain REDIS_ADDR, got %q",
			err.Error(),
		)
	}
}

func TestLoad_InvalidRedisDB(t *testing.T) {
	t.Setenv("REDIS_ADDR", "localhost:6379")
	t.Setenv("REDIS_PASSWORD", "")
	t.Setenv("REDIS_DB", "invalid")

	cfg, err := Load()

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if cfg != nil {
		t.Fatal("expected config to be nil")
	}

	if !strings.Contains(err.Error(), "REDIS_DB") {
		t.Errorf(
			"expected error to contain REDIS_DB, got %q",
			err.Error(),
		)
	}
}

func TestLoad_EmptyRedisDB(t *testing.T) {
	t.Setenv("REDIS_ADDR", "localhost:6379")
	t.Setenv("REDIS_PASSWORD", "")
	t.Setenv("REDIS_DB", "")

	cfg, err := Load()

	if err == nil {
		t.Fatal("expected error for empty REDIS_DB")
	}

	if cfg != nil {
		t.Fatal("expected config to be nil")
	}
}

func TestLoad_NegativeRedisDB(t *testing.T) {
	t.Setenv("REDIS_ADDR", "localhost:6379")
	t.Setenv("REDIS_PASSWORD", "")
	t.Setenv("REDIS_DB", "-1")

	cfg, err := Load()

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if cfg != nil {
		t.Fatal("expected config to be nil")
	}

	if !strings.Contains(err.Error(), "non-negative") {
		t.Errorf(
			"expected non-negative error, got %q",
			err.Error(),
		)
	}
}

func TestLoad_ZeroRedisDB(t *testing.T) {
	t.Setenv("REDIS_ADDR", "localhost:6379")
	t.Setenv("REDIS_PASSWORD", "")
	t.Setenv("REDIS_DB", "0")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.RedisDb != 0 {
		t.Errorf("expected RedisDb 0, got %d", cfg.RedisDb)
	}
}

func TestLoad_LargeRedisDB(t *testing.T) {
	t.Setenv("REDIS_ADDR", "localhost:6379")
	t.Setenv("REDIS_PASSWORD", "")
	t.Setenv("REDIS_DB", "100")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.RedisDb != 100 {
		t.Errorf("expected RedisDb 100, got %d", cfg.RedisDb)
	}
}

func TestLoad_UsesEnvironmentVariables(t *testing.T) {
	const (
		addr     = "redis.example.com:6380"
		password = "my-password"
	)

	t.Setenv("REDIS_ADDR", addr)
	t.Setenv("REDIS_PASSWORD", password)
	t.Setenv("REDIS_DB", "5")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.RedisAddr != addr {
		t.Errorf("expected address %q, got %q", addr, cfg.RedisAddr)
	}

	if cfg.RedisPassword != password {
		t.Errorf(
			"expected password %q, got %q",
			password,
			cfg.RedisPassword,
		)
	}

	if cfg.RedisDb != 5 {
		t.Errorf("expected DB 5, got %d", cfg.RedisDb)
	}
}

// Проверяем, что отсутствие .env файла само по себе
// не является ошибкой, если необходимые переменные окружения заданы.
func TestLoad_NoDotEnvFileStillWorks(t *testing.T) {
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}

	tempDir := t.TempDir()

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	t.Cleanup(func() {
		_ = os.Chdir(originalDir)
	})

	t.Setenv("REDIS_ADDR", "localhost:6379")
	t.Setenv("REDIS_PASSWORD", "")
	t.Setenv("REDIS_DB", "0")

	cfg, err := Load()
	if err != nil {
		t.Fatalf(
			"expected Load to work without .env file, got error: %v",
			err,
		)
	}

	if cfg == nil {
		t.Fatal("expected config not to be nil")
	}
}
