package server

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
)

func TestParseBulkString(t *testing.T) {

	reader := bufio.NewReader(bytes.NewReader([]byte("$3\r\nfoo\r\n")))

	got, err := ParseBulkString(reader)
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

	reader := bufio.NewReader(bytes.NewReader([]byte("$3\r n")))

	_, err := ParseBulkString(reader)
	if err == nil {
		t.Fatalf("should produce error")
	}
}

func TestParseArray(t *testing.T) {

	reader := bufio.NewReader(bytes.NewReader([]byte("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n")))

	//Find number of elements+\r\n
	got, err := ParseArray(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fmt.Println(got)

	// if !reflect.DeepEqual(got.Elements, []any{"hello", "world"}) {
	// 	t.Errorf("Elements = %q, want %q", got.Elements, []string{"hello", "world"})
	// }

}
