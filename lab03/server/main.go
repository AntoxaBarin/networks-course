package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	RESPONSE_404     = "HTTP/1.1 404 Not Found\r\n\r\n"
	RESPONSE_500     = "HTTP/1.1 500 Internal Server Error\r\n\r\n"
	FMT_RESPONSE_200 = "HTTP/1.1 200 OK\r\nContent-Length: %d\r\n\r\n%s"
	PATH_TO_STORAGE  = "../local_storage/"
)

func main() {
	var taskFlag = flag.String("task", "A", "Task")
	PORT := os.Args[1]
	if PORT[0] != ':' {
		PORT = ":" + PORT
	}
	concurrencyLevel, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Incorrect Concurrency level.")
		os.Exit(1)
	}
	routines := make(chan struct{}, concurrencyLevel)

	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Server is listening on localhost" + PORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			os.Exit(1)
		}
		fmt.Println("Connected with", conn.RemoteAddr().String())

		if *taskFlag == "A" {
			handleRequest(conn)
		} else if *taskFlag == "B" {
			go handleRequest(conn)
		} else if *taskFlag == "D" {

			// Channel size is concurrency level
			routines <- struct{}{}

			go func(conn net.Conn) {
				defer func() {
					<-routines
					conn.Close()
				}()
				handleRequest(conn)
			}(conn)
		}

	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	request, err := r.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	httpParts := strings.Fields(request)
	if len(httpParts) < 2 {
		fmt.Println("Incorrect HTTP request.")
		return
	}

	method, path := httpParts[0], httpParts[1]
	if method != "GET" {
		conn.Write([]byte(RESPONSE_404))
		return
	}

	filePath := PATH_TO_STORAGE + path
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			conn.Write([]byte(RESPONSE_404))
			return
		}
		conn.Write([]byte(RESPONSE_500))
		return
	}
	defer file.Close()

	fileStats, err := file.Stat()
	if err != nil {
		conn.Write([]byte(RESPONSE_500))
		return
	}
	fileContent := make([]byte, fileStats.Size())
	_, err = file.Read(fileContent)
	if err != nil {
		conn.Write([]byte(RESPONSE_500))
		return
	}

	conn.Write([]byte(fmt.Sprintf(FMT_RESPONSE_200, fileStats.Size(), fileContent)))
}
