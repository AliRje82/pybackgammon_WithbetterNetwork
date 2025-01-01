package network

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
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

	go func() {
		for {
			time.Sleep(45 * time.Second)
			fmt.Println("Running the matchmaking algorithm")
			matchMaking(g)
			fmt.Println("Finishing up the algorithm")
		}

	}()

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
	node.indx = make(chan *Node)
	node.message = make(chan string)
	node.isReserved = false
	g.addNode(node)
	//Matchmaking...
	for {
		matched := <-node.indx
		g.match.RLock()
		fmt.Println("Found a match indx", matched)
		message := matched.username + "," + matched.ip
		conn.Write([]byte(message))

		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Something happend!")
			node.isReserved = false
			matched.message <- "ERR"
			g.match.RUnlock()
			continue
		}
		matched.message <- message
		other_message := <-node.message
		if other_message == "Accept" && message == "Accept" {
			//play the game
			g.match.RUnlock()
			break
		} else if other_message == "ERR" {
			node.isReserved = false
			g.match.RUnlock()
		} else {
			node.RemoveEdge(matched)
			node.isReserved = false
			g.match.RUnlock()
		}
	}
	close(node.indx)
	close(node.message)

	// Rolling dice

}
