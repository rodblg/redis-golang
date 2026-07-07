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
