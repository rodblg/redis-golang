package server

import (
	"bytes"
	"sync"
	"testing"

	"bufio"
)

func TestHandleRespInput(t *testing.T) {

	reader := bufio.NewReader(bytes.NewReader([]byte("*1\r\n$4\r\nPING\r\n")))

	s := Server{
		store: &Memory{
			mu:    sync.RWMutex{},
			KVMap: make(map[string]string),
		},
	}

	got, err := s.handleResp(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "+PONG" {
		t.Errorf("Message = %q, Want = %q", got, "PONG")
	}

}

func TestHandleRespEcho(t *testing.T) {
	reader := bufio.NewReader(bytes.NewReader([]byte("*3\r\n$4\r\nECHO\r\n$5\r\nHello\r\n$5\r\nWorld\r\n")))

	s := Server{
		store: &Memory{
			mu:    sync.RWMutex{},
			KVMap: make(map[string]string),
		},
	}

	got, err := s.handleResp(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "+Hello World" {
		t.Errorf("Message = %q, Want = %q", got, "Hello World")
	}

}

func TestHandleRespSet(t *testing.T) {
	reader := bufio.NewReader(bytes.NewReader([]byte("*3\r\n$3\r\nSET\r\n$4\r\nName\r\n$4\r\nJohn\r\n")))

	s := Server{
		store: &Memory{
			mu:    sync.RWMutex{},
			KVMap: make(map[string]string),
		},
	}

	got, err := s.handleResp(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "+OK" {
		t.Errorf("Response = %q, Want = %q", got, "OK")
	}
}
