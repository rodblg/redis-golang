package server

import (
	"bytes"
	"testing"

	"bufio"
)

func TestHandleRespInput(t *testing.T) {

	reader := bufio.NewReader(bytes.NewReader([]byte("*1\r\n$4\r\nPING\r\n")))

	got, err := handleResp(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "PONG" {
		t.Errorf("Message = %q, Want = %q", got, "PONG")
	}

}

func TestHandleRespEcho(t *testing.T) {
	reader := bufio.NewReader(bytes.NewReader([]byte("*3\r\n$4\r\nECHO\r\n$5\r\nHello\r\n$5\r\nWorld\r\n")))

	got, err := handleResp(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "Hello World" {
		t.Errorf("Message = %q, Want = %q", got, "Hello World")
	}

}
