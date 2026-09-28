package server

import (
	"bytes"
	"strconv"
	"sync"
	"testing"
	"time"

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

func TestKeySetExpiryTimestamp(t *testing.T) {
	//we want to store keys
	// we receive a reader with options PX EX
	setReader := bufio.NewReader(bytes.NewReader([]byte("*5\r\n$3\r\nSET\r\n$4\r\nName\r\n$4\r\nJohn\r\n$2\r\nEX\r\n$1\r\n5\r\n")))
	s := Server{
		store: &Memory{
			mu:       sync.RWMutex{},
			KVMap:    make(map[string]string),
			Expirity: make(map[string]time.Time),
		},
	}

	got, err := s.handleResp(setReader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "+OK" {
		t.Errorf("Response = %q, Want = %q", got, "OK")
	}

	getReader := bufio.NewReader(bytes.NewReader([]byte("*2\r\n$3\r\nGET\r\n$4\r\nName\r\n")))

	// esperamos n tiempo
	time.Sleep(5 * time.Second)
	// hacemos un get y visualizamos que se borra
	got, err = s.handleResp(getReader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "-1" {
		t.Errorf("Response = %q, Want = %q", got, "-1")
	}
}

func TestTTLKeyExpirity(t *testing.T) {

	s := Server{
		store: &Memory{
			mu:       sync.RWMutex{},
			KVMap:    make(map[string]string),
			Expirity: make(map[string]time.Time),
		},
	}

	setReader := bufio.NewReader(bytes.NewReader([]byte("*5\r\n$3\r\nSET\r\n$4\r\nName\r\n$4\r\nJohn\r\n$2\r\nEX\r\n$1\r\n5\r\n")))
	//the key is added
	got, err := s.handleResp(setReader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "+OK" {
		t.Errorf("Response = %q, Want = %q", got, "OK")
	}
	time.Sleep(3 * time.Second)

	ttlReader := bufio.NewReader(bytes.NewReader([]byte("*2\r\n$3\r\nTTL\r\n$4\r\nName\r\n")))

	got, err = s.handleResp(ttlReader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ttl, _ := strconv.Atoi(got)
	if ttl < 1 || ttl > 2 {
		t.Errorf("TTL = %d, want between 1 and 2", ttl)
	}

}
