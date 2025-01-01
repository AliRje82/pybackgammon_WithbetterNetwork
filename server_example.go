package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	fmt.Println("Server is up")

	listen, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		fmt.Println("Err")
	}
	conn, err := listen.Accept()
	if err != nil {
		fmt.Println("ERR")
	}
	defer conn.Close()
	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error accourd in reading a massage")
		return
	}
	fmt.Println(message)
}
