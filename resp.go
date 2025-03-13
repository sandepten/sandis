package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

// parseRESP reads a full command using the RESP protocol.
func parseRESP(reader *bufio.Reader) ([]string, error) {
	// Read the first line. It should start with '*' if using RESP arrays.
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return nil, fmt.Errorf("empty input")
	}

	// If it doesn't start with '*', fallback to splitting by space.
	if line[0] != '*' {
		return strings.Split(line, " "), nil
	}

	// Parse the number of arguments in the array.
	count, err := strconv.Atoi(line[1:])
	if err != nil {
		return nil, fmt.Errorf("invalid array count: %v", err)
	}

	tokens := make([]string, 0, count)
	for i := 0; i < count; i++ {
		// Next, read the bulk string header (should start with '$').
		header, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		header = strings.TrimSpace(header)
		if len(header) == 0 || header[0] != '$' {
			return nil, fmt.Errorf("expected bulk string header, got: %s", header)
		}
		// Optionally, parse the length of the bulk string.
		_, err = strconv.Atoi(header[1:])
		if err != nil {
			return nil, fmt.Errorf("invalid bulk string length: %v", err)
		}
		// Now read the actual bulk string data.
		data, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		data = strings.TrimRight(data, "\r\n")
		tokens = append(tokens, data)
	}
	return tokens, nil
}
