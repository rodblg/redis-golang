package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strconv"
	"strings"
)

func main() {

	// arguments := os.Args
	// if len(arguments) == 1 {
	// 	fmt.Println("Please provide port number")
	// 	return
	// }
	//default redis port
	PORT := ":" + "6379"
	//listen on server socket, network endpoint IP+Port
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	slog.Info("server is listening on port :6379")

	//We implement multiple clients concurrently
	for {

		conn, err := listener.Accept()
		if err != nil {
			slog.Error("listener#error", "error", err)
			continue
		}
		slog.Info("accepted#connection", "socket", conn.LocalAddr())

		go handleConnection(conn)

	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		//Read and process data from the client
		//The receiver buffers the received segments
		msg, err := reader.ReadString('\n')
		if err != nil {
			slog.Error("error#buffer", "error", err)
			return
		}

		//Write data back to the client
		if strings.TrimSpace(string(msg)) == "STOP" {
			fmt.Println("Exiting TCP server!")
			return
		}

		fmt.Print("-> ", string(msg))
		// t := time.Now()
		// myTime := t.Format(time.RFC3339) + "\n"
		conn.Write([]byte("received\n"))
	}
}

// Write a function that reads $\r\n\r\n and returns the string. Read the length first, then read exactly that many bytes.
// Test with a hardcoded byte slice before wiring it to the socket

type BulkString struct {
	Length  int
	Message string
}

// a bulk string represent a single binary string. RESP encodes as $<length>\r\n<data>\r\n
func ParseBulkString(input []byte) (*BulkString, error) {

	reader := bufio.NewReader(bytes.NewReader(input))

	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	line = strings.TrimSuffix(line, "\r\n")
	line = strings.TrimPrefix(line, "$")
	fmt.Println(line)

	length, err := strconv.Atoi(line)
	if err != nil {
		return nil, err
	}

	msgBuff := make([]byte, length)
	_, err = io.ReadFull(reader, msgBuff)
	if err != nil {
		return nil, err
	}

	return &BulkString{Length: length, Message: string(msgBuff)}, nil
	//Return BulkString
}

// Clients send commands to redis as RESP arrays
// *<number-of-elements>\r\n<element-1>...<element-n>
// EXAMPLE *2\r\n$4\r\nLLEN\r\n$6\r\nmylist\r\n
func parseArray() {}
