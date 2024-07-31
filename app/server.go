package main

import (
	"fmt"
	"net"
	"os"
	"log"
	"strings"
	"flag"
	"path/filepath"
	"compress/gzip"
	"bytes"
	"slices"
)

func main() {
	fmt.Println("This is my own HTTP Server, It's now up and running! :D")
	var dir string
	flag.StringVar(&dir, "directory", "/tmp/", "Directory to write/serve the file from")
	flag.Parse()
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
		go handleClient(conn, dir) // Handle a client connection
	}
}

func handleClient(conn net.Conn, dir string) {
	defer conn.Close() // Ensure we terminate the connection after we're done
	var (
		err error
		res, msg, path, userAgent string
		content []byte
		fullReq, req, encodings []string
		n int
		encFlag bool
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
	method := req[0]
	path = req[1]
	body := fullReq[len(fullReq)-1]
	for i, v := range fullReq{
		// fmt.Println(i, v)
		if strings.HasPrefix(v, "User-Agent: "){
			userAgent, _ = strings.CutPrefix(v, "User-Agent: ")
		}
		if strings.HasPrefix(v, "Accept-Encoding: "){
			encFlag = true
			temp, _ := strings.CutPrefix(fullReq[i], "Accept-Encoding: ")
			encs := strings.Split(temp, ",")
			for _, w := range encs{
				encodings = append(encodings, strings.Trim(w, " "))
			}
			fmt.Println(encodings)
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
		var b bytes.Buffer
		res = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(msg), msg)
		if encFlag {
			if slices.Contains(encodings, "gzip"){
				fmt.Println("hey im here")
				gz := gzip.NewWriter(&b)
				defer gz.Close()
				_, err = gz.Write([]byte(msg)) // Compress the data
				if err != nil {
					log.Println("Error compressing data: ", err)
					return
				}
				if err = gz.Close(); err != nil {
					log.Println("Error closing gzip: ", err)
					return
				}
				msg = b.String()
				fmt.Println()
				enc := "gzip"
				res = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Encoding: %s\r\n\r\n%X", enc, b.Bytes())
			} else {
				res = "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\n"
			}
		}
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
	case strings.HasPrefix(path, "/files/"):
		fileName, _ := strings.CutPrefix(path, "/files/") // Extract the filename
		filePath := filepath.Join(dir, fileName)
		switch method {
		case "GET":
			_, err := os.Stat(filePath)
			if err == nil {
				content, err = os.ReadFile(filePath)
				if err != nil {
					log.Println("Failed to open file: ", err)
					return
				}
				res = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(content), string(content))
			} 
			if os.IsNotExist(err) {
				res = "HTTP/1.1 404 Not Found\r\n\r\n"
			}
			_, err = conn.Write([]byte(res))
				if err != nil {
					log.Println("Error writing response: ", err)
					return
				}
				log.Println("Response sent: \"", res, "\"")
		case "POST":
			content = []byte(body)
			err = os.WriteFile(filePath, content, 0644)
			if err != nil {
				log.Println("Error writing to file: ", err)
				return
			}
			res = "HTTP/1.1 201 Created\r\n\r\n"
			_, err = conn.Write([]byte(res))
				if err != nil {
					log.Println("Error writing response: ", err)
					return
				}
				log.Println("Response sent: \"", res, "\"")
		}
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