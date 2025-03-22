package main

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

const (
	smtpServer = "smtp.gmail.com"
	port       = "587"
)

var (
	from       string
	to         string
	password   string
	msg        string
	imgPath    string
	encodedImg string
)

func main() {
	flag.StringVar(&from, "from", "dummy_mail@aboba.com", "Sender email")
	flag.StringVar(&to, "to", from, "Recepient email")
	flag.StringVar(&password, "pass", "strong_password", "Sender mail password")
	flag.StringVar(&msg, "msg", "This is the email body.", "Email body")
	flag.StringVar(&imgPath, "image", "", "Path to image file")
	flag.Parse()

	if imgPath != "" {
		imgData, err := os.ReadFile(imgPath)
		if err != nil {
			fmt.Println("Failed to read image file:", err)
			return
		}
		encodedImg = base64.StdEncoding.EncodeToString(imgData)
	}

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
	sendMsg(w, tlsConn, "DATA\r\n")

	boundary := "my-boundary-12345"
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: Test email with image\r\n", from, to)
	message += "MIME-Version: 1.0\r\n"
	message += fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n", boundary)
	message += "\r\n"

	message += fmt.Sprintf("--%s\r\n", boundary)
	message += "Content-Type: text/plain; charset=\"UTF-8\"\r\n"
	message += "\r\n"
	message += msg + "\r\n"
	message += "\r\n"

	if encodedImg != "" {
		message += fmt.Sprintf("--%s\r\n", boundary)
		message += "Content-Type: image/png\r\n"
		message += "Content-Transfer-Encoding: base64\r\n"
		message += "Content-Disposition: attachment; filename=\"image\"\r\n"
		message += "\r\n"
		message += encodedImg + "\r\n"
		message += "\r\n"
	}
	message += fmt.Sprintf("--%s--\r\n", boundary)
	message += ".\r\n"

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
