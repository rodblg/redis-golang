package main

import "testing"

func TestParseBulkString(t *testing.T) {
	got, err := ParseBulkString([]byte("$3\r\nfoo\r\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Message != "foo" {
		t.Errorf("Message = %q, want %q", got.Message, "foo")
	}

}
