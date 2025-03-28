package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func echo(conn net.Conn, message string) error {
	dataLength := len(message)
	writeData := fmt.Sprintf("$%d\r\n%s\r\n", dataLength, message)
	_, err := conn.Write([]byte(writeData))
	return err
}

func ping(conn net.Conn) error {
	_, err := conn.Write([]byte("+PONG\r\n"))
	return err
}

func defaultCase(conn net.Conn, message string) error {
	_, err := conn.Write([]byte(fmt.Sprintf("-ERR unknown command '%s'\r\n", message)))
	return err
}

func set(conn net.Conn, store map[string]StoreValue, inputs []string) error {
	key := inputs[1]
	value := inputs[2]
	var expiryTime time.Time
	if len(inputs) > 4 && inputs[3] == "px" {
		duration, err := strconv.Atoi(inputs[4])
		if err == nil {
			expiryTime = time.Now().Add(time.Duration(duration) * time.Millisecond)
		}
	} else {
		expiryTime = time.Time{}
	}

	// var valueData StoreValue
	valueData := StoreValue{
		value:     value,
		expiresAt: expiryTime,
	}
	store[key] = valueData
	_, err := conn.Write([]byte("+OK\r\n"))
	return err
}

func get(conn net.Conn, store map[string]StoreValue, inputs []string) error {
	key := inputs[1]
	value := store[key]
	valueData := value.value
	if !value.expiresAt.IsZero() && value.expiresAt.Before(time.Now()) {
		valueData = ""
	}

	dataLength := len(valueData)
	writeData := fmt.Sprintf("$%d\r\n%s\r\n", dataLength, valueData)
	_, err := conn.Write([]byte(writeData))
	return err
}

func config(conn net.Conn, inputs []string) error {
	// function := inputs[1]
	param := inputs[2]

	if param == "dir" {
		value := os.Getenv(RdpDirPath)
		writeData := fmt.Sprintf("*2\r\n$3\r\ndir\r\n$16\r\n%s\r\n", value)
		_, err := conn.Write([]byte(writeData))
		fmt.Println(err, writeData)
		return err
	}
	if param == "dbfilename" {
		value := os.Getenv(RdpFileName)
		writeData := fmt.Sprintf("*2\r\n$3\r\ndir\r\n$16\r\n%s\r\n", value)
		_, err := conn.Write([]byte(writeData))
		return err
	}
	return nil
}
