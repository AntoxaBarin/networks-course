package main

import (
	"flag"
	"log"
	"math/rand"
	"net"
	"strings"
)

func main() {
	port := flag.String("port", "8080", "port to run server on")
	flag.Parse()

	addr, err := net.ResolveUDPAddr("udp", ":"+*port)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("Error reading from UDP: %v\n", err)
			continue
		}

		message := string(buf[:n])
		log.Printf("Received from %v: %s", clientAddr, message)

		if rand.Intn(10) < 2 {
			log.Printf("Packet lost (not responding to %s)", message)
			continue
		}
		upMessage := strings.ToUpper(message)
		_, err = conn.WriteToUDP([]byte(upMessage), clientAddr)
		if err != nil {
			log.Printf("Error writing to UDP: %v\n", err)
			continue
		}

		log.Printf("Sent to %v: %s", clientAddr, upMessage)
	}
}
