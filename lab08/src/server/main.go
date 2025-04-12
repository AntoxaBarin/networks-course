package main

import (
	"flag"
	"log"
	"math/rand"
	"net"
	"os"
)

const PACKET_SIZE = 128

func main() {
	port := flag.String("port", "", "port to run server on")
	flag.Parse()

	conn, err := net.ListenPacket("udp", ":"+*port)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	out, err := os.OpenFile("received.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	buf := make([]byte, PACKET_SIZE)
	for {
		log.Println("Listening...")
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Println("Failed to read from connection")
			continue
		}
		if n < 1 {
			log.Println("Received empty packet")
			continue
		}

		if rand.Float64() < 0.3 {
			log.Println("Packet lost")
			continue
		}

		_, err = conn.WriteTo([]byte{buf[0]}, addr)
		if err != nil {
			log.Println("Failed to send an ACK")
			continue
		}

		out.WriteString(string(buf[1:n]))
	}
}
