package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

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

func formatResponseContent(content string) string {
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
			content := formatResponseContent(randomString)
			conn.Write([]byte(content))
			return
		}
		// task 5
		if path == "/user-agent" {
			user_agent := strings.Split(bufferArr[2], " ")[1]
			content := formatResponseContent(user_agent)
			conn.Write([]byte(content))
			return
		}
		// task 3
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}
