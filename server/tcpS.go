package server

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"net"
	"strconv"
	"strings"
	"time"

	"sync"
)

type Memory struct {
	mu       sync.RWMutex
	KVMap    map[string]string
	Expirity map[string]time.Time
}

func (m *Memory) SetKey(array *RespArray) error {
	//we get the value
	//message := echoMessage(array)

	if array.Length < 3 {
		return errors.New("wrong number of arguments for 'set' command")
	}
	var expEn time.Time
	for i := 3; i < array.Length; i++ {
		opt := strings.ToUpper(array.Elements[i].Message)
		if (opt == "EX" || opt == "PX") && (i+1 < array.Length) {
			duration, err := strconv.Atoi(array.Elements[i+1].Message)
			if err != nil {
				return err
			}
			var value time.Duration
			switch opt {
			case "EX":
				value = time.Second
			case "PX":
				value = time.Millisecond
			}

			expEn = time.Now().Add(time.Duration(duration) * value)

		}

	}

	key, value := array.Elements[1].Message, array.Elements[2].Message

	m.mu.Lock()
	defer m.mu.Unlock()

	m.KVMap[key] = value

	if !expEn.IsZero() {
		m.Expirity[key] = expEn
	}

	return nil
}

func (m *Memory) GetKey(array *RespArray) (string, bool) {
	//we get the value
	//message := echoMessage(array)

	//check if there is a expiration

	m.mu.Lock()
	defer m.mu.Unlock()
	key := array.Elements[1].Message
	expVar, exists := m.Expirity[key]

	if exists && time.Now().After(expVar) {
		delete(m.KVMap, key)
		delete(m.Expirity, key)
		return "-1", true
	}

	v, ok := m.KVMap[key]

	return fmt.Sprintf("+%s", v), ok
}

func (m *Memory) GetTTL(array *RespArray) string {

	if array.Length < 2 {
		return "-ERR invalid command"
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	key := array.Elements[1].Message
	_, ok := m.KVMap[key]
	if !ok {
		return "-2"
	}
	expIn, hasExp := m.Expirity[key]
	if ok && !hasExp {
		return "-1"
	}

	if hasExp && time.Now().After(expIn) {
		delete(m.KVMap, key)
		delete(m.Expirity, key)
		return "-2"
	}

	//convert time to expire in left seconds
	leftS := int(math.Round(time.Until(expIn).Seconds()))

	return fmt.Sprintf("%d", leftS)

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
		err := s.store.SetKey(respArray)
		if err != nil {
			return "-ERR " + err.Error(), nil
		}
		return "+OK", nil
	case "GET":
		value, ok := s.store.GetKey(respArray)
		if !ok {
			return fmt.Sprintf("+%s", "key not found"), nil
		}
		return fmt.Sprintf("%s", value), nil

	case "TTL":
		s := s.store.GetTTL(respArray)
		return s, nil
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
