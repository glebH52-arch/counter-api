package redis_client

import (
	"bufio"
	"fmt"
	"net"
	"testing"
)

func startFakeRedisServer(t *testing.T, response string) (string, func()) {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start fake redis server: %v", err)
	}

	done := make(chan struct{})

	go func() {
		defer close(done)

		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		reader := bufio.NewReader(conn)

		// Redis client may send CLIENT SETINFO before PING.
		for {
			_, err := reader.ReadString('\n')
			if err != nil {
				return
			}

			// Read RESP command lines until we encounter PING.
			// For our tests it is sufficient to inspect incoming data
			// and respond with +OK for setup commands.
			buffered := reader.Buffered()
			if buffered > 0 {
				data, _ := reader.Peek(buffered)

				if containsPing(string(data)) {
					fmt.Fprint(conn, response)
					return
				}
			}

			fmt.Fprint(conn, "+OK\r\n")
		}
	}()

	cleanup := func() {
		_ = listener.Close()
		<-done
	}

	return listener.Addr().String(), cleanup
}

func containsPing(s string) bool {
	for i := 0; i+4 <= len(s); i++ {
		if s[i:i+4] == "PING" || s[i:i+4] == "ping" {
			return true
		}
	}

	return false
}

func TestNewRedisClient_ConnectionError(t *testing.T) {
	client, err := NewRedisClient(
		"127.0.0.1:1",
		"",
		0,
	)

	if err == nil {
		t.Fatal("expected connection error, got nil")
	}

	if client != nil {
		t.Fatal("expected client to be nil when connection fails")
	}
}

func TestNewRedisClient_InvalidAddress(t *testing.T) {
	client, err := NewRedisClient(
		"invalid-address",
		"",
		0,
	)

	if err == nil {
		t.Fatal("expected error for invalid address")
	}

	if client != nil {
		t.Fatal("expected nil client")
	}
}
