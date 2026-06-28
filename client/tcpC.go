package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"log/slog"
)

func main() {

	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		slog.Error("error#connection", "error", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	connReader := bufio.NewReader(conn)

	for {

		fmt.Println(">> ")
		text, _ := reader.ReadString('\n')

		fmt.Fprintf(conn, text+"\n")

		message, _ := connReader.ReadString('\n')
		fmt.Print("->: " + message)
		if strings.TrimSpace(string(text)) == "STOP" {
			fmt.Println("TCP client exiting")
			return
		}
	}
}

//The client needs to send bulks of strings (a single binary string)
