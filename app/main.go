package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	ln, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer ln.Close()

	fmt.Println("Listening on 0.0.0.0:6379")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	var inputs []string
	for { // Add a loop to handle multiple commands
		readData, err := reader.ReadString('\n') // Read until newline
		if err != nil {
			fmt.Println("Error reading:", err)
			return // Connection closed or error occurred
		}
		message := string(readData)
		message = strings.TrimSpace(message)

		if message == "" { // Handle empty input (e.g., just a newline)
			continue
		}
		input := parseInput(message)
		if len(input) > 0 {
			inputs = append(inputs, parseInput(message))
		}
		fmt.Println(inputs)

		if len(inputs) == 0 {
			fmt.Println("Empty input")
			continue
		}
		command := strings.ToLower(inputs[0])
		var errWrite error
		switch command {
		case "echo":
			if len(inputs) > 1 {
				errWrite = echo(conn, inputs[1])
			} else {
				continue
			}
		case "ping":
			errWrite = ping(conn)
		default:
			errWrite = defaultCase(conn, command)

		}
		if errWrite != nil {
			fmt.Println("Error writing:", errWrite)
			return
		}
	}
}
