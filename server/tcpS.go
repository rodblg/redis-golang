package server

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"strings"

	"sync"
)

type Memory struct {
	mu    sync.RWMutex
	KVMap map[string]string
}

func (m *Memory) SetKey(array *RespArray) {
	//we get the value
	message := echoMessage(array)
	keyValue := strings.SplitN(message, " ", 2)

	m.mu.Lock()

	m.KVMap[keyValue[0]] = keyValue[1]

	m.mu.Unlock()
}

func (m *Memory) GetKey(array *RespArray) (string, bool) {
	//we get the value
	message := echoMessage(array)
	key := strings.SplitN(message, " ", 2)

	m.mu.RLock()

	v, ok := m.KVMap[key[0]]

	m.mu.RUnlock()

	return v, ok
}

type Server struct {
	store *Memory
}

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

	server := Server{store: &Memory{
		mu:    sync.RWMutex{},
		KVMap: make(map[string]string),
	},
	}

	//We implement multiple clients concurrently
	for {

		conn, err := listener.Accept()
		if err != nil {
			slog.Error("listener#error", "error", err)
			continue
		}
		slog.Info("accepted#connection", "socket", conn.LocalAddr())

		go server.handleConnection(conn)

	}

}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, err := s.handleResp(reader)
		if err != nil {
			conn.Write([]byte("-ERR " + err.Error() + "\r\n"))
			continue
		}
		conn.Write([]byte(message + "\r\n"))
	}
}

func (s *Server) handleResp(reader *bufio.Reader) (string, error) {

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
	case "SET":
		s.store.SetKey(respArray)
		return "+OK", nil
	case "GET":
		value, ok := s.store.GetKey(respArray)
		if !ok {
			return fmt.Sprintf("+%s", "key not found"), nil
		}
		return fmt.Sprintf("+%s", value), nil
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
