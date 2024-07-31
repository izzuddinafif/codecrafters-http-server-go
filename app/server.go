package main

import (
	"fmt"
	"net"
	"os"
	"log"
	"strings"
)

func main() {
	fmt.Println("This is my own HTTP Server, It's now up and running! :D")

	listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		log.Println("Failed to bind to port 4221: ", err)
		os.Exit(1)
	}
	defer listener.Close() // Ensure we close the server when the program exists

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err)
			continue
		}
		// Handle a client connection
		handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	// Ensure we terminate the connection after we're done
	defer conn.Close()
	var (
		err error
		res string
		fullReq []string
		n int
	)
	buf := make([]byte, 1024)

	// Read request
	n, err = conn.Read(buf)
	rec := string(buf[:n])
	if err != nil {
		log.Println("Error reading data: ", err)
		return
	}
	log.Println("HTTP Request: ", rec)


	fullReq = strings.Split(rec, "\r\n")
	req := strings.Split(fullReq[0], " ")
	//method := req[0]
	path := req[1]
	//host := strings.CutPrefix(fullReq[1], "Host: ")
	//userAgent := strings.CutPrefix(fullreq[2], "User-Agent: ")
	fmt.Println(path)

	switch path {
	case "/":
		// Write ok response
		res = "HTTP/1.1 200 OK\r\n\r\n"
		_, err = conn.Write([]byte(res))
		if err != nil {
			log.Println("Error writing response: ", err)
			return
		}
		log.Println("Response sent: \"", res, "\"")
	default:
		res = "HTTP/1.1 404 Not Found\r\n\r\n"
		_, err = conn.Write([]byte(res))
		if err != nil {
			log.Println("Error writing response: ", err)
			return
		}
		log.Println("Response sent: \"", res, "\"")
	}
	
}