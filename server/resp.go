package server

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

type RespBulkString struct {
	Length  int
	Message string
}

// Write a function that reads $\r\n\r\n and returns the string. Read the length first, then read exactly that many bytes.
// Test with a hardcoded byte slice before wiring it to the socket
// a bulk string represent a single binary string. RESP encodes as $<length>\r\n<data>\r\n
func ParseBulkString(reader *bufio.Reader) (*RespBulkString, error) {

	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	line = strings.TrimSuffix(line, "\r\n")
	line = strings.TrimPrefix(line, "$")

	length, err := strconv.Atoi(line)
	if err != nil {
		return nil, err
	}
	//also skip the next \r\n in the msg
	msgBuff := make([]byte, length+2)
	_, err = io.ReadFull(reader, msgBuff)
	if err != nil {
		return nil, err
	}

	strMsgBuff := strings.TrimSuffix(string(msgBuff), "\r\n")

	return &RespBulkString{Length: length, Message: strMsgBuff}, nil
}

type RespArray struct {
	Length   int
	Elements []RespBulkString
}

// Clients send commands to redis as RESP arrays
// *<number-of-elements>\r\n<element-1>...<element-n>
// EXAMPLE *2\r\n$4\r\nLLEN\r\n$6\r\nmylist\r\n
func ParseArray(reader *bufio.Reader) (*RespArray, error) {

	//Find number of elements+\r\n
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	line = strings.TrimPrefix(line, "*")
	line = strings.TrimSuffix(line, "\r\n")

	nElements, err := strconv.Atoi(line)
	if err != nil {
		return nil, err
	}
	//If its not zero we continue
	if nElements == 0 {
		return &RespArray{}, nil
	}

	var elements []RespBulkString
	for i := 0; i < nElements; i++ {

		bulkString, err := ParseBulkString(reader)
		if err != nil {
			return &RespArray{}, err
		}
		elements = append(elements, *bulkString)
	}

	return &RespArray{
		Length:   nElements,
		Elements: elements,
	}, nil

}
