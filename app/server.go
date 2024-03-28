package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

var directory string

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	flag.StringVar(&directory, "directory", "", "a directory")
	flag.Parse()

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

func handleClient(clientConn net.Conn) {
	buffer := make([]byte, 1024)

	for {
		n, err := clientConn.Read(buffer)
		if err != nil {
			clientConn.Close()
			return
		}

		request := strings.Split(string(buffer[:n]), "\r\n")
		requestLine := request[0]
		path := strings.Split(requestLine, " ")[1]

		switch {
		case path == "/":
			clientConn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		case strings.HasPrefix(path, "/echo/"):
			content := formatPlainTextContent(strings.TrimPrefix(path, "/echo/"))
			clientConn.Write([]byte(content))
		case path == "/user-agent":
			userAgent := strings.Split(request[2], " ")[1]
			content := formatPlainTextContent(userAgent)
			clientConn.Write([]byte(content))
		case strings.HasPrefix(path, "/files/"):
			handleFileRequest(clientConn, path, request, directory)
		default:
			clientConn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		}
	}
}

func handleFileRequest(conn net.Conn, path string, request []string, directory string) {
	method := strings.Split(request[0], " ")[0]
	fileName := strings.TrimPrefix(path, "/files/")
	filePath := filepath.Join(directory, fileName)

	switch {
	case method == "GET":
		handleFileGetRequest(conn, filePath)
	case method == "POST":
		handleFilePostRequest(conn, filePath, request)
	}
}

func handleFileGetRequest(conn net.Conn, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		conn.Close()
		return
	}
	defer file.Close()

	fileContent, _ := os.ReadFile(filePath)

	conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
	conn.Write([]byte("Content-Type: application/octet-stream\r\n"))
	conn.Write([]byte(fmt.Sprintf("Content-Length: %d\r\n", len(fileContent))))
	conn.Write([]byte("\r\n"))
	conn.Write(fileContent)
}

func handleFilePostRequest(conn net.Conn, filePath string, request []string) {
	err := os.WriteFile(filePath, []byte(request[len(request)-1]), 0644)
	if err != nil {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		conn.Close()
		return
	}

	conn.Write([]byte("HTTP/1.1 201 Created\r\n\r\n"))
}
