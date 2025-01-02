package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
)

func reciveMessage(conn net.Conn) ([]byte, error) {
	buf := make([]byte, 1024)

	_, err := conn.Read(buf[:4])
	if err != nil {
		return nil, err
	}
	var messageLength uint32
	err = binary.Read(bytes.NewReader(buf[:4]), binary.BigEndian, &messageLength)
	if err != nil {
		return nil, fmt.Errorf("failed to read message length: %v", err)
	}

	message := make([]byte, messageLength)
	_, err = conn.Read(message)
	if err != nil {
		return nil, fmt.Errorf("failed to read message data: %v", err)
	}

	fmt.Printf("Received message: %s\n", message)

	return message, nil
}
func makePkt(message []byte) ([]byte, error) {
	messageLength := uint32(len(message))

	var lengthHeader bytes.Buffer

	err := binary.Write(&lengthHeader, binary.BigEndian, messageLength)
	if err != nil {
		return nil, fmt.Errorf("failed to write length header: %v", err)
	}

	packet := append(lengthHeader.Bytes(), message...)

	return packet, nil
}

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
	for {
		message, err := reciveMessage(conn)
		if err != nil {
			fmt.Println("Error accourd in reading a massage")
			return
		}
		if string(message) == "Roll" {
			dice1 := rand.Intn(6) + 1
			dice2 := rand.Intn(6) + 1
			responce, err := makePkt([]byte(string(dice1) + "," + string(dice2)))
			fmt.Println(dice1, dice2)
			if err != nil {
				fmt.Println("Something happend in making a pkt")
			}
			conn.Write(responce)
		}

	}

}
