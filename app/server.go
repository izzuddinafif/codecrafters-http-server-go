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

	for { // Infinite loop to make sure server keeps listening
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err)
			continue
		}
		go handleClient(conn) // Handle a client connection
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close() // Ensure we terminate the connection after we're done
	var (
		err error
		res, msg, path, userAgent string
		fullReq, req []string
		n int
	)
	buf := make([]byte, 1024)

	n, err = conn.Read(buf) // Read request
	rec := string(buf[:n])
	if err != nil {
		log.Println("Error reading data: ", err)
		return
	}
	log.Println("HTTP Request: ", rec)

	fullReq = strings.Split(rec, "\r\n")
	
	req = strings.Split(fullReq[0], " ")
	// method := req[0]
	path = req[1]
	for i, v := range fullReq{
		fmt.Println(i, v)
		if strings.HasPrefix(v, "User-Agent: "){
			userAgent, _ = strings.CutPrefix(v, "User-Agent: ")
		}
	}
	switch {
	case path == "/":
		res = "HTTP/1.1 200 OK\r\n\r\n"
		_, err = conn.Write([]byte(res)) // Send response
		if err != nil {
			log.Println("Error writing response: ", err)
			return
		}
		log.Println("Response sent: \"", res, "\"")
	case strings.HasPrefix(path, "/echo/"):
		msg, _ = strings.CutPrefix(path, "/echo/") // Extract the echo string
		res = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(msg), msg)
		_, err = conn.Write([]byte(res))
		if err != nil {
			log.Println("Error writing response: ", err)
			return
		}
		log.Println("Response sent: \"", res, "\"")
	case path == "/user-agent":
		msg = userAgent
		res = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(msg), msg)
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