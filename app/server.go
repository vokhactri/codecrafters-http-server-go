package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

var dirPointer string

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	flag.StringVar(&dirPointer, "directory", "", "a directory")
	flag.Parse()

	// task 6
	// Accept and handle incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		// Handle connection concurrently
		go handleClient(conn)
	}
}

func formatPlainTextContent(content string) string {
	return strings.Join(
		[]string{
			"HTTP/1.1 200 OK",
			"Content-Type: text/plain",
			"Content-Length: " + fmt.Sprint(len(content)),
			"",
			content + "\r\n",
		},
		"\r\n")
}

func handleClient(conn net.Conn) {
	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from connection: ", err.Error())
			os.Exit(1)
		}

		bufferArr := strings.Split(string(buffer[:n]), "\r\n")
		startLine := bufferArr[0]
		path := strings.Split(startLine, " ")[1]
		// task 2
		if path == "/" {
			conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
			return
		}
		// task 4
		if strings.HasPrefix(path, "/echo/") {
			randomString := strings.TrimPrefix(path, "/echo/")
			content := formatPlainTextContent(randomString)
			conn.Write([]byte(content))
			return
		}
		// task 5
		if path == "/user-agent" {
			user_agent := strings.Split(bufferArr[2], " ")[1]
			content := formatPlainTextContent(user_agent)
			conn.Write([]byte(content))
			return
		}
		// task 7
		if strings.HasPrefix(path, "/files/") {
			fileName := strings.TrimPrefix(path, "/files")
			fileDir := dirPointer + fileName
			file, err := os.Open(fileDir)
			if err != nil {
				fmt.Println("Error opening file: ", err.Error())
				conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
				return
			}
			defer file.Close()

			conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
			conn.Write([]byte("Content-Type: application/octet-stream\r\n"))
			conn.Write([]byte("\r\n"))

			// Create a bufferedReader to efficiently read file contents
			reader := bufio.NewReader(file)

			// Create a buffer to store data read from the file
			buffer := make([]byte, 1024)

			// Read file contents and write them to the TCP connection
			for {
				bytesRead, err := reader.Read(buffer)
				if err != nil {
					return
				}
				if bytesRead == 0 {
					break // End of file
				}
				conn.Write(buffer[:bytesRead])
			}
		}
		// task 3
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}
