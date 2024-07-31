package main

import (
	"fmt"
	"net"
	"os"
	"log"
)

func main() {
	fmt.Println("This is my own HTTP Server Program, It's now up and running! :D")

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

	response := "HTTP/1.1 200 OK\r\n\r\n"

	/*// Read data
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("Error reading data: " err)
		return
	}
	log.Println("Received data: ",, buf[:n])*/

	// Write data
	n, err := conn.Write([]byte(response))
	if err != nil {
		log.Println("Error writing response: ", err)
		return
	}
	log.Println(n, "bytes sent:\"", response, "\"")
}