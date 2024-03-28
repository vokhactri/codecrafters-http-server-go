package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	go handleConnection(conn)
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

func handleConnection(conn net.Conn) {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		os.Exit(1)
	}

	bufferArr := strings.Split(string(buffer[:n]), "\r\n")
	startLine := bufferArr[0]
	path := strings.Split(startLine, " ")[1]
	if path == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		return
	}
	if strings.HasPrefix(path, "/echo/") {
		randomString := strings.TrimPrefix(path, "/echo/")
		content := formatResponseContent(randomString)
		conn.Write([]byte(content))
		return
	}
	if path == "/user-agent" {
		user_agent := strings.Split(bufferArr[2], " ")[1]
		content := formatResponseContent(user_agent)
		conn.Write([]byte(content))
		return
	}
	conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
}
