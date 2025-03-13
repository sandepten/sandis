package main

import (
	"fmt"
	"net"
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
	_, err := conn.Write([]byte(message + " is an unrecoginized command\n"))
	return err
}
