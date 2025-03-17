package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestEndToEndCommands(t *testing.T) {
	// Define the path to the binary - this assumes you've built it
	binaryPath := "/tmp/redis-go"

	// Start the Redis server as a subprocess
	cmd := exec.Command(binaryPath)

	// Redirect output for debugging
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the process
	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start Redis server: %v", err)
	}

	// Ensure we kill the process when done
	defer func() {
		if err := cmd.Process.Kill(); err != nil {
			t.Logf("Failed to kill process: %v", err)
		}
		// Wait for the process to exit to prevent zombies
		_ = cmd.Wait()
	}()

	// Give the server time to start
	time.Sleep(1 * time.Second)

	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Test PING command
	t.Run("PING", func(t *testing.T) {
		// Send PING command using RESP protocol
		_, err := conn.Write([]byte("*1\r\n$4\r\nPING\r\n"))
		if err != nil {
			t.Fatalf("Failed to send PING command: %v", err)
		}

		// Read response
		reader := bufio.NewReader(conn)
		response, err := reader.ReadString('\n')
		if err != nil {
			t.Fatalf("Failed to read response: %v", err)
		}

		// Verify response
		expected := "+PONG\r\n"
		if response != expected {
			t.Errorf("Expected response %q, got %q", expected, response)
		}
	})

	// Test PING with plain text protocol
	t.Run("PING Plain Text", func(t *testing.T) {
		// Send PING command as plain text
		_, err := conn.Write([]byte("PING\r\n"))
		if err != nil {
			t.Fatalf("Failed to send PING command: %v", err)
		}

		// Read response
		reader := bufio.NewReader(conn)
		response, err := reader.ReadString('\n')
		if err != nil {
			t.Fatalf("Failed to read response: %v", err)
		}

		// Verify response
		expected := "+PONG\r\n"
		if response != expected {
			t.Errorf("Expected response %q, got %q", expected, response)
		}
	})

	// Test ECHO command
	t.Run("ECHO", func(t *testing.T) {
		testMessage := "Hello, Redis!"

		// Send ECHO command using RESP protocol
		cmd := fmt.Sprintf("*2\r\n$4\r\nECHO\r\n$%d\r\n%s\r\n", len(testMessage), testMessage)
		_, err := conn.Write([]byte(cmd))
		if err != nil {
			t.Fatalf("Failed to send ECHO command: %v", err)
		}

		// Read response
		reader := bufio.NewReader(conn)
		responseLine, err := reader.ReadString('\n')
		if err != nil {
			t.Fatalf("Failed to read response header: %v", err)
		}

		// Parse the bulk string length
		if !strings.HasPrefix(responseLine, "$") {
			t.Fatalf("Expected bulk string response, got: %s", responseLine)
		}

		// Read the message content
		response, err := reader.ReadString('\n')
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		// Verify response content (trim CRLF)
		responseContent := strings.TrimRight(response, "\r\n")
		if responseContent != testMessage {
			t.Errorf("Expected echo response %q, got %q", testMessage, responseContent)
		}
	})
}
