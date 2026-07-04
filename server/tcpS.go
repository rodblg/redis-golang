package server

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"strings"
)

func main() {

	// arguments := os.Args
	// if len(arguments) == 1 {
	// 	fmt.Println("Please provide port number")
	// 	return
	// }
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

	//We implement multiple clients concurrently
	for {

		conn, err := listener.Accept()
		if err != nil {
			slog.Error("listener#error", "error", err)
			continue
		}
		slog.Info("accepted#connection", "socket", conn.LocalAddr())

		go handleConnection(conn)

	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		//Read and process data from the client
		//The receiver buffers the received segments
		msg, err := reader.ReadString('\n')
		if err != nil {
			slog.Error("error#buffer", "error", err)
			return
		}

		//Write data back to the client
		if strings.TrimSpace(string(msg)) == "STOP" {
			fmt.Println("Exiting TCP server!")
			return
		}

		fmt.Print("-> ", string(msg))
		// t := time.Now()
		// myTime := t.Format(time.RFC3339) + "\n"
		conn.Write([]byte("received\n"))
	}
}
