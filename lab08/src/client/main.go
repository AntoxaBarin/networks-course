package main

import (
	"flag"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"
)

const PACKET_SIZE = 128

func main() {
	serverPort := flag.String("server-port", "", "server port")
	port := flag.String("port", "", "port to run client on")
	dataPath := flag.String("data", "", "path to the data")
	timeout := flag.Duration("timeout", time.Second, "timeout for ACK")
	flag.Parse()

	data, err := os.ReadFile(*dataPath)
	if err != nil {
		log.Fatalf("Failed to read data from file %s", *dataPath)
	}

	conn, err := net.Dial("udp", ":"+*serverPort)
	if err != nil {
		log.Fatalf("Failed to establish connection with server, error: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go receive(*port)
	go send(conn, data, *timeout)

	wg.Wait()
}

func send(conn net.Conn, data []byte, timeout time.Duration) {
	time.Sleep(3 * time.Second)

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
			log.Println("Writing into connection...")
			_, err := conn.Write(packet[:payloadSize+1])
			if err != nil {
				continue
			}
		} else {
			log.Println("Packet lost")
			generateNewPacket = false
			continue
		}

		ack := make([]byte, 1)
		for {
			err := conn.SetReadDeadline(time.Now().Add(timeout))
			if err != nil {
				log.Println("Failed to set timeout on read from connection")
				generateNewPacket = false
				continue
			}
			_, err = conn.Read(ack)
			if err != nil || ack[0] != packet[0] {
				log.Println(err)
				_, err = conn.Write(packet)
				if err != nil {
					log.Println("Failed to write to connection")
					continue
				}
			} else {
				log.Printf("Received ACK: %b", ack[0])
				break
			}
		}

		if restBytes == 0 {
			break
		}
	}
}

func receive(port string) {
	conn, err := net.ListenPacket("udp", ":"+port)
	if err != nil {
		log.Fatal("ABOA", err)
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
