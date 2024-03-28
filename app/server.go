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

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		os.Exit(1)
	}

	buffer_arr := strings.Split(string(buffer[:n]), "\r\n")
	start_line := buffer_arr[0]
	path := strings.Split(start_line, " ")[1]
	if path == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		return
	}
	if strings.HasPrefix(path, "/echo/") {
		random_string := strings.TrimPrefix(path, "/echo/")
		content := formatResponseContent(random_string)
		conn.Write([]byte(content))
		return
	}
	if path == "/user-agent" {
		user_agent := strings.Split(buffer_arr[2], " ")[1]
		content := formatResponseContent(user_agent)
		conn.Write([]byte(content))
		return
	}
	conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
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
