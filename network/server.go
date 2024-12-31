package network

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Please provide a server ip")
	}
	ip := os.Args[1]

	fmt.Println("Starting the server")

	listen, err := net.Listen("tcp", ip)
	if err != nil {
		panic(err)
	}
	g := &Graph{}

	for {

		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Something happend")
			continue
		}
		fmt.Println("Getting a connection")

		go connection(conn, g)
	}
}

func connection(conn net.Conn, g *Graph) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error accourd in reading a massage")
		return
	}

	parts := strings.Split(message, ",")

	if len(parts) != 2 {
		panic("Massage has problem!")
	}

	username, ip := parts[0], parts[1]

	//Create a node
	node := &Node{}
	node.username = username
	node.ip = ip
	node.indx = make(chan int)

	g.addNode(node)

	matched := <-node.indx

	fmt.Println(matched)

}
