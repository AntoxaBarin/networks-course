package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
)

const (
	username = "dlpuser"
	password = "rNrKYTX9g7z3RgJRmxWuGHbeu"
	host     = "ftp.dlptest.com"
)

func main() {
	c, err := ftp.Dial(host+":21", ftp.DialWithTimeout(1*time.Second))
	if err != nil {
		log.Fatal(err)
	}

	err = c.Login(username, password)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println(`Available commands:
- pwd       - show current directory
- cd <dir>  - change directory
- mkdir <dir> - create directory
- ls        - list files
- exit      - quit program
- put <path> - send file to the server
- get <path> - get file from the server`)

	for {
		fmt.Print("ftp> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		command := parts[0]
		args := parts[1:]

		switch command {
		case "exit", "quit", "bye":
			fmt.Println("Goodbye!")
			return

		case "pwd":
			if curDir, err := c.CurrentDir(); err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println(curDir)
			}

		case "cd":
			if len(args) < 1 {
				fmt.Println("Usage: cd <directory>")
				continue
			}
			if err := c.ChangeDir(args[0]); err != nil {
				fmt.Println("Error:", err)
			}

		case "mkdir":
			if len(args) < 1 {
				fmt.Println("Usage: mkdir <directory>")
				continue
			}
			if err := c.MakeDir(args[0]); err != nil {
				fmt.Println("Error:", err)
			}

		case "ls", "dir":
			curDir, err := c.CurrentDir()
			if err != nil {
				log.Fatal("Failed to get current directorty")
			}
			if entries, err := c.List(curDir); err != nil {
				fmt.Println("Error:", err)
			} else {
				for _, e := range entries {
					fmt.Println(e.Type.String(), e.Name)
				}
			}

		case "put":
			if len(args) < 1 {
				fmt.Println("Usage: put <local-file> [remote-file]")
				continue
			}

			localFile := args[0]
			remoteFile := args[0]
			if len(args) > 1 {
				remoteFile = args[1]
			}

			file, err := os.Open(localFile)
			if err != nil {
				fmt.Printf("Error opening file %s: %v\n", localFile, err)
				continue
			}
			defer file.Close()

			err = c.Stor(remoteFile, file)
			if err != nil {
				fmt.Printf("Error uploading file %s: %v\n", remoteFile, err)
			} else {
				fmt.Printf("File %s successfully uploaded as %s\n", localFile, remoteFile)
			}

		case "get":
			if len(args) < 1 {
				fmt.Println("Usage: get <remote-file> [local-file]")
				continue
			}

			remoteFile := args[0]
			localFile := filepath.Base(remoteFile)
			if len(args) > 1 {
				localFile = args[1]
			}

			if _, err := os.Stat(localFile); err == nil {
				fmt.Printf("Error: File %s already exists\n", localFile)
				continue
			}

			r, err := c.Retr(remoteFile)
			if err != nil {
				fmt.Printf("Error downloading file %s: %v\n", remoteFile, err)
				continue
			}
			defer r.Close()

			file, err := os.Create(localFile)
			if err != nil {
				fmt.Printf("Error creating file %s: %v\n", localFile, err)
				continue
			}
			defer file.Close()

			_, err = io.Copy(file, r)
			if err != nil {
				fmt.Printf("Error saving file %s: %v\n", localFile, err)
				os.Remove(localFile)
			} else {
				fmt.Printf("File %s successfully downloaded as %s\n", remoteFile, localFile)
			}

		default:
			fmt.Println("Unknown command. Available: pwd, cd, mkdir, ls, exit")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}

	if err := c.Quit(); err != nil {
		log.Fatal(err)
	}
}
