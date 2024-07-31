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
		handleClient(conn) // Handle a client connection
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close() // Ensure we terminate the connection after we're done
	var (
		err error
		res string
		msg string
		userAgent string
		fullReq []string
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
	
	req := strings.Split(fullReq[0], " ")
	fmt.Println(req)
	// method := req[0]
	path := req[1]
	fmt.Println(path)
	for i, v := range fullReq{
		fmt.Println(i, v)
		if strings.HasPrefix(v, "User-Agent: "){
			userAgent, _ = strings.CutPrefix(v, "User-Agent: ")
			fmt.Println(userAgent)
		}
	}
	
	// userAgent, _ := strings.CutPrefix(fullReq[3], "User-Agent: ")
	// fmt.Println(userAgent)

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
		fmt.Println(msg)
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