package main

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"flag"
	"fmt"
	"net"
	"time"
)

const (
	smtpServer = "smtp.gmail.com"
	port       = "587"
)

var (
	from     string
	to       string
	password string
	msg      string
)

func main() {
	flag.StringVar(&from, "from", "dummy_mail@aboba.com", "Sender email")
	flag.StringVar(&to, "to", from, "Recepient email")
	flag.StringVar(&password, "pass", "strong_password", "Sender mail password")
	flag.StringVar(&msg, "msg", "This is the email body.", "Email body")
	flag.Parse()

	conn, err := net.Dial("tcp", smtpServer+":"+port)
	if err != nil {
		fmt.Println("Failed to connect to smtp server: ", err)
		return
	}
	defer conn.Close()

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	response, _ := r.ReadString('\n')
	fmt.Print(response)

	sendMsg(w, conn, "HELO localhost\r\n")
	sendMsg(w, conn, "STARTTLS\r\n")

	tlsConfig := &tls.Config{
		ServerName: smtpServer,
	}
	tlsConn := tls.Client(conn, tlsConfig)
	err = tlsConn.Handshake()
	if err != nil {
		fmt.Println("Failed to set TLS:", err)
		return
	}
	defer tlsConn.Close()

	w = bufio.NewWriter(tlsConn)

	sendMsg(w, tlsConn, "EHLO localhost\r\n")
	sendMsg(w, tlsConn, "AUTH PLAIN "+encryptPassword()+"\r\n")

	sendMsg(w, tlsConn, fmt.Sprintf("MAIL FROM:<%s>\r\n", from))
	sendMsg(w, tlsConn, fmt.Sprintf("RCPT TO:<%s>\r\n", to))
	sendMsg(w, tlsConn, "DATA\n")
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: Test email\r\n\r\n%s\r\n.\r\n", from, to, msg)
	sendMsg(w, tlsConn, message)
	sendMsg(w, tlsConn, "QUIT\r\n")
}

func sendMsg(writer *bufio.Writer, conn net.Conn, command string) {
	writer.WriteString(command)
	writer.Flush()
	fmt.Print("-> ", command)

	r := bufio.NewReader(conn)
	response, err := r.ReadString('\n')
	if err != nil {
		fmt.Println("Failed to read response from server:", err)
		return
	}
	fmt.Print("<- ", response)
	time.Sleep(200 * time.Millisecond)
}

func encryptPassword() string {
	authPlain := "\x00" + from + "\x00" + password
	return base64.StdEncoding.EncodeToString([]byte(authPlain))
}
