package server

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"net"
)

func Start() {

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
		message, err := handleResp(reader)
		if err != nil {
			conn.Write([]byte("-ERR " + err.Error() + "\r\n"))
			continue
		}
		conn.Write([]byte(message + "\r\n"))
	}
}

func handleResp(reader *bufio.Reader) (string, error) {

	respArray, err := ParseArray(reader)
	if err != nil {
		return "", err
	}

	if respArray.Length == 0 {
		return "", errors.New("invalid array")
	}

	command := respArray.Elements[0].Message

	switch command {
	case "PING":
		return "+PONG", nil
	case "ECHO":
		return fmt.Sprintf("+%s", echoMessage(respArray)), nil
	}
	return "-ERR unknown command '" + command + "'", nil

}

func echoMessage(array *RespArray) string {
	var resp string
	if array.Length > 1 {
		for i := 1; i < array.Length; i++ {
			if i == 1 {
				resp = array.Elements[i].Message
			} else {
				resp = fmt.Sprintf("%s %s", resp, array.Elements[i].Message)
			}
		}
	}

	return resp
}
