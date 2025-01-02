package main

import (
	"bytes"
	"crypto/aes"
	"fmt"
	"log"
	"net"
	"os"
	"pybackgammon_WithbetterNetwork/Onionrouting/myNet"
)

func main() {

	if len(os.Args) != 3 {
		fmt.Println("You should run this with 2 Arg for example 'IP:PORT' 'IP:PORT'")
	}
	node := os.Args[1]
	nexNode := os.Args[2]

	listen, err := net.Listen("tcp", node)

	if err != nil {
		fmt.Println("Error")
		os.Exit(1)
	}
	defer listen.Close()
	fmt.Println("Router ip is", node)

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Error In accept")
			continue
		}
		fmt.Println("A connection is accepted")

		go routing(conn, nexNode)

	}

}

func pad(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

// Removes padding from the plaintext
func unpad(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

// Encrypts data using AES-ECB
func encryptAES_ECB(key, plaintext []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("Failed to create cipher: %v", err)
	}

	if len(plaintext)%aes.BlockSize != 0 {
		plaintext = pad(plaintext, aes.BlockSize)
	}

	ciphertext := make([]byte, len(plaintext))

	for start := 0; start < len(plaintext); start += aes.BlockSize {
		end := start + aes.BlockSize
		block.Encrypt(ciphertext[start:end], plaintext[start:end])
	}

	return ciphertext
}

// Decrypts data using AES-ECB
func decryptAES_ECB(key, ciphertext []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("Failed to create cipher: %v", err)
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		log.Fatalf("Ciphertext is not a multiple of the block size")
	}

	plaintext := make([]byte, len(ciphertext))

	for start := 0; start < len(ciphertext); start += aes.BlockSize {
		end := start + aes.BlockSize
		block.Decrypt(plaintext[start:end], ciphertext[start:end])
	}

	plaintext = unpad(plaintext)
	return plaintext
}

func routing(conn net.Conn, next string) {
	defer conn.Close()

	key, err := myNet.ReciveMessage(conn)
	if err != nil {
		fmt.Println("Connection problem in reading!")
		return
	}

	fmt.Println("Key recived!")

	call, err := net.Dial("tcp", next)
	if err != nil {
		panic(err)
	}
	defer call.Close()

	errorChan := make(chan error, 2)

	//Sending a massage
	go func() {

		for {

			nextMassage, nextErr := myNet.ReciveMessage(call)

			if nextErr != nil {
				errorChan <- fmt.Errorf("next connection error: %v", err)
				return
			}

			//Encrypt
			message, err := myNet.MakePkt(encryptAES_ECB(key, nextMassage))
			if err != nil {
				errorChan <- err
				return
			}
			conn.Write(message)
		}
	}()

	//Getting a massage
	go func() {
		for {
			beforeMassage, beforeErr := myNet.ReciveMessage(conn)
			if beforeErr != nil {
				errorChan <- fmt.Errorf("next connection error: %v", err)
				return
			}
			//Decrypt
			message, err := myNet.MakePkt(decryptAES_ECB(key, beforeMassage))
			if err != nil {
				errorChan <- err
				return
			}
			call.Write(message)
		}
	}()

	err = <-errorChan
	fmt.Println("Connection closed due to error:", err)
}
