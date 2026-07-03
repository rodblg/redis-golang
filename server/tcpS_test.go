package main

import (
	"reflect"
	"testing"
)

func TestParseBulkString(t *testing.T) {
	got, err := ParseBulkString([]byte("$3\r\nfoo\r\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Length != 3 {
		t.Errorf("Message = %q, want %q", got.Length, 3)
	}
	if got.Message != "foo" {
		t.Errorf("Message = %q, want %q", got.Message, "foo")
	}

}

func TestParseBulkString_NoNewLine(t *testing.T) {
	_, err := ParseBulkString([]byte("$3\r n"))
	if err == nil {
		t.Fatalf("should produce error")
	}
}

func TestParseArray(t *testing.T) {

	//Find number of elements+\r\n
	got, elements, err := ParseArray([]byte("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if elements != 2 {
		t.Errorf("Message = %q, want %q", got.Length, 3)
	}

	if !reflect.DeepEqual(got.Elements, []any{"hello", "world"}) {
		t.Errorf("Elements = %q, want %q", got.Elements, []string{"hello", "world"})
	}

}
