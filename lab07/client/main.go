package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	port := flag.String("server-port", "8080", "server port")
	flag.Parse()
	serverAddr := "localhost:" + *port

	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for i := 1; i <= 10; i++ {
		currentTime := time.Now()
		message := fmt.Sprintf("Ping %d %s", i, currentTime.Format("15:04:05.000"))

		_, err = conn.Write([]byte(message))
		if err != nil {
			log.Printf("Error sending packet %d: %v\n", i, err)
			continue
		}

		err = conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		if err != nil {
			log.Printf("Error setting timeout for packet %d: %v\n", i, err)
			continue
		}

		buffer := make([]byte, 1024)
		start := time.Now()
		n, err := conn.Read(buffer)

		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Printf("Request %d timed out\n", i)
			} else {
				log.Printf("Error reading response for packet %d: %v\n", i, err)
			}
			continue
		}

		response := string(buffer[:n])
		fmt.Printf("%s, RTT=%.5f seconds\n", response, time.Since(start).Seconds())
	}
}
