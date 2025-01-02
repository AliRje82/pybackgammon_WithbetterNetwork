package myNet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

func MakePkt(message []byte) ([]byte, error) {
	messageLength := uint32(len(message))

	var lengthHeader bytes.Buffer

	err := binary.Write(&lengthHeader, binary.BigEndian, messageLength)
	if err != nil {
		return nil, fmt.Errorf("failed to write length header: %v", err)
	}

	packet := append(lengthHeader.Bytes(), message...)

	return packet, nil
}
func ReciveMessage(conn net.Conn) ([]byte, error) {
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
