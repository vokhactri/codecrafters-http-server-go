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

	start_line := strings.Split(string(buffer[:n]), "\r\n")[0]
	path := strings.Split(start_line, " ")[1]
	if path == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		return
	}
	if strings.HasPrefix(path, "/echo/") {
		random_string := strings.TrimPrefix(path, "/echo/")
		conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
		conn.Write([]byte("Content-Type: text/plain\r\n"))
		conn.Write([]byte(fmt.Sprintf("Content-Length: %d\r\n", len(random_string))))
		conn.Write([]byte("\r\n"))
		conn.Write([]byte(fmt.Sprintf("%s\r\n\r\n", random_string)))
		return
	}
	conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
}
