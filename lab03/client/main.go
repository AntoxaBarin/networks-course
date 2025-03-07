package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: ./client <server_host> <server_port> <filename>")
		return
	}
	serverHost := os.Args[1]
	serverPort := os.Args[2]
	filename := os.Args[3]

	serverAddress := serverHost + ":" + serverPort
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Println("Failed to connect to the server: ", err)
		return
	}
	defer conn.Close()

	request := fmt.Sprintf("GET /%s HTTP/1.1\r\nHost: %s\r\nConnection: close\r\n\r\n", filename, serverHost)
	_, err = conn.Write([]byte(request))
	if err != nil {
		fmt.Println("Failed to send HTTP request: ", err)
		return
	}

	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	for err == nil {
		fmt.Print(response)
		response, err = reader.ReadString('\n')
	}

	if err.Error() != "EOF" {
		fmt.Println("Failed to read response: ", err)
	}
}
