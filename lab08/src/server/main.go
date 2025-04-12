package main

import (
	"flag"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

const PACKET_SIZE = 128

func main() {
	clientPort := flag.String("client-port", "", "client port")
	port := flag.String("port", "", "port to run server on")
	dataPath := flag.String("data", "", "path to the data")
	timeout := flag.Duration("timeout", time.Second, "timeout for ACK")
	version := flag.String("v", "1", "server version")
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

	if *version == "2" {
		data, err := os.ReadFile(*dataPath)
		if err != nil {
			log.Fatalf("Failed to read data from file %s", *dataPath)
		}
		go send(data, *clientPort, *timeout)
	}

	for {
		log.Println("[RECEIVE]: Listening...")
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Println("[RECEIVE]: Failed to read from connection")
			continue
		}
		if n < 1 {
			log.Println("[RECEIVE]: Received empty packet")
			continue
		}

		if rand.Float64() < 0.3 {
			log.Println("[RECEIVE]: Packet lost")
			continue
		}

		_, err = conn.WriteTo([]byte{buf[0]}, addr)
		if err != nil {
			log.Println("[RECEIVE]: Failed to send an ACK")
			continue
		}

		out.WriteString(string(buf[1:n]))
	}
}

func send(data []byte, port string, timeout time.Duration) {
	time.Sleep(3 * time.Second)

	conn, err := net.Dial("udp", ":"+port)
	if err != nil {
		log.Fatalf("[SEND]: Failed to establish connection with client, error: %v", err)
	}
	log.Println("[SEND]: Sending...")

	restBytes := len(data)
	packetSeqNumber := 0

	var generateNewPacket bool = true
	packet := make([]byte, PACKET_SIZE)
	var payloadSize int

	for {
		if generateNewPacket {
			packet[0] = byte(packetSeqNumber)
			packetSeqNumber = 1 - packetSeqNumber

			payloadSize = copy(packet[1:], data)
			data = data[payloadSize:]
			log.Println(payloadSize, restBytes)
			restBytes -= payloadSize
		} else {
			generateNewPacket = true
		}

		if rand.Float64() >= 0.3 {
			log.Println("[SEND]: Writing into connection...")
			_, err := conn.Write(packet[:payloadSize+1])
			if err != nil {
				continue
			}
		} else {
			log.Println("[SEND]: Packet lost")
			generateNewPacket = false
			continue
		}

		ack := make([]byte, 1)
		for {
			err := conn.SetReadDeadline(time.Now().Add(timeout))
			if err != nil {
				log.Println("[SEND]: Failed to set timeout on read from connection")
				generateNewPacket = false
				continue
			}
			_, err = conn.Read(ack)
			if err != nil || ack[0] != packet[0] {
				log.Println(err)
				_, err = conn.Write(packet)
				if err != nil {
					log.Println("[SEND]: Failed to write to connection")
					continue
				}
			} else {
				log.Printf("[SEND]: Received ACK: %b", ack[0])
				break
			}
		}

		if restBytes == 0 {
			break
		}
	}
}
