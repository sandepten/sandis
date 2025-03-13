package main

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"testing"
	"time"
)

func TestEchoCommand(t *testing.T) {
	// Start the server in a goroutine
	go func() {
		main()
	}()

	// Give the server time to start
	time.Sleep(100 * time.Millisecond)

	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		t.Fatalf("Failed to connect to Redis server: %v", err)
	}
	defer conn.Close()

	// List of test strings
	testStrings := [10]string{
		"hello",
		"world",
		"mangos",
		"apples",
		"oranges",
		"watermelons",
		"grapes",
		"pears",
		"horses",
		"elephants",
	}

	// Seed the random generator
	rand.Seed(time.Now().UnixNano())

	// Select a random string
	randomString := testStrings[rand.Intn(10)]

	// Send ECHO command
	command := fmt.Sprintf("ECHO %s\r\n", randomString)
	_, err = conn.Write([]byte(command))
	if err != nil {
		t.Fatalf("Failed to send command: %v", err)
	}

	// Read response
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	// Parse the response
	response := string(buffer[:n])
	// The response format should be $<length>\r\n<content>\r\n
	expectedPrefix := fmt.Sprintf("$%d\r\n", len(randomString))
	expectedSuffix := "\r\n"

	if !strings.HasPrefix(response, expectedPrefix) ||
		!strings.HasSuffix(response, expectedSuffix) ||
		!strings.Contains(response, randomString) {
		t.Fatalf("Expected response containing %s, got %s", randomString, response)
	}

	// Extract the actual content from the response
	content := strings.TrimSuffix(strings.TrimPrefix(response, expectedPrefix), expectedSuffix)

	if content != randomString {
		t.Fatalf("Expected %s, got %s", randomString, content)
	}
}
