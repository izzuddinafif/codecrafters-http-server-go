package main

import (
	"fmt"
	"strings"
)

func main() {
	a := "GET / HTTP/1.1\r\nHost: localhost:4221\r\nUser-Agent: curl/7.64.1\r\nAccept: */*\r\n\r\n"
	b := strings.Split(a, "\r\n")
	//c := strings.Split(b[0], " ")



	for i := range b {
		fmt.Println(i, b[i])
	}










	
}