package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Use port 0 to get any available port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close() // Close so run() can bind to it

	errChan := make(chan error, 1)

	// Run server in background
	go func() {
		errChan <- run(ctx, fmt.Sprintf("%d", port), "data/store.db")
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Verify server is responding
	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/api/product", port))
	if err != nil {
		t.Fatalf("Server not responding: %v", err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Cancel context to trigger shutdown
	cancel()

	// Wait for server to shutdown
	select {
	case err := <-errChan:
		if err != nil {
			t.Errorf("run() returned error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Server did not shutdown in time")
	}
}
