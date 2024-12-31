package routers

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"net"
	"os"
	"strings"
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

func routing(conn net.Conn, next string) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Connection problem in reading!")
		return
	}

	message = strings.TrimSpace(message)
	block, err := aes.NewCipher([]byte(message))
	if err != nil {
		panic(err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	nonceSize := aesGCM.NonceSize()

	fmt.Println("Key recived!")

	call, err := net.Dial("tcp", next)
	if err != nil {
		panic(err)
	}
	defer call.Close()

	errorChan := make(chan error, 2)
	callReader := bufio.NewReader(call)

	//Sending a massage
	go func() {

		for {

			nextMassage, nextErr := callReader.ReadString('\n')

			if nextErr != nil {
				errorChan <- fmt.Errorf("next connection error: %v", err)
				return
			}

			//Encrypt
			nextMassage = strings.TrimSpace(nextMassage)
			//Make nonce
			nonce := make([]byte, nonceSize)
			if _, err := rand.Read(nonce); err != nil {
				fmt.Println("Failed to generate nonce:", err)
				continue
			}

			encrypted := aesGCM.Seal(nil, nonce, []byte(nextMassage), nil)
			fullMessage := append(nonce, encrypted...)
			_, err := conn.Write(append(fullMessage, '\n'))

			if err != nil {
				fmt.Println("Failed to send message to user:", err)
				errorChan <- err
				return
			}
		}
	}()

	//Getting a massage
	go func() {
		for {
			beforeMassage, beforeErr := reader.ReadString('\n')
			if beforeErr != nil {
				errorChan <- fmt.Errorf("next connection error: %v", err)
				return
			}
			//Decrypt
			beforeMassage = strings.TrimSpace(beforeMassage)
			nonce, beforeMassage := beforeMassage[:nonceSize], beforeMassage[nonceSize:]
			decrypted, err := aesGCM.Open(nil, []byte(nonce), []byte(beforeMassage), nil)
			if err != nil {
				panic(err)
			}
			_, err = call.Write(append(decrypted, '\n'))
			if err != nil {
				fmt.Println("Failed to send message to next service:", err)
				errorChan <- err
				return
			}
		}
	}()

	err = <-errorChan
	fmt.Println("Connection closed due to error:", err)
}
