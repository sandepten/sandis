package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

const (
	// RDP
	RdpDirPath  = "RDB_DIR_PATH"
	RdpFileName = "RDB_FILE_NAME"
)

type StoreValue struct {
	value     string
	expiresAt time.Time
}

func main() {
	ln, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer ln.Close()

	fmt.Println("Listening on 0.0.0.0:6379")

	// redis store
	store := make(map[string]StoreValue)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue
		}
		go handleConnection(conn, store)
	}
}

func handleConnection(conn net.Conn, store map[string]StoreValue) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		// Parse the complete RESP command.
		inputs, err := parseRESP(reader)
		if err != nil {
			fmt.Println("Error parsing RESP:", err)
			return
		}
		if len(inputs) == 0 {
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
		case "set":
			if len(inputs) > 2 {
				errWrite = set(conn, store, inputs)
			} else {
				continue
			}
		case "get":
			if len(inputs) > 1 {
				errWrite = get(conn, store, inputs)
			} else {
				continue
			}
		case "config":
			if len(inputs) > 1 {
				errWrite = config(conn, inputs)
			} else {
				continue
			}
		default:
			errWrite = defaultCase(conn, command)
		}
		if errWrite != nil {
			fmt.Println("Error writing:", errWrite)
			return
		}
	}
}
