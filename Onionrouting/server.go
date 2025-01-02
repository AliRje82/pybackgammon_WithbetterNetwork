package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"pybackgammon_WithbetterNetwork/Onionrouting/models"
	"pybackgammon_WithbetterNetwork/Onionrouting/myNet"
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
	g := &models.Graph{}

	go func() {
		for {
			time.Sleep(45 * time.Second)
			fmt.Println("Running the matchmaking algorithm")
			models.MatchMaking(g)
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

func connection(conn net.Conn, g *models.Graph) {
	defer conn.Close()

	message, err := myNet.ReciveMessage(conn)
	if err != nil {
		fmt.Println("Error accourd in reading a massage")
		return
	}

	parts := strings.Split(string(message), ",")

	if len(parts) != 2 {
		panic("Massage has problem!")
	}

	username, ip := parts[0], parts[1]
	fmt.Println("username: ", username, "ip: ", ip)
	//Create a node
	node := &models.Node{}
	node.Username = username
	node.Ip = ip
	node.Indx = make(chan *models.Node)
	node.Message = make(chan string)
	node.IsReserved = false
	g.AddNode(node)
	//Matchmaking...
	for {
		matched := <-node.Indx
		g.Match.RLock()
		fmt.Println("Found a match indx", matched)
		message := matched.Username + "," + matched.Ip
		conn.Write([]byte(message))

		me, err := myNet.ReciveMessage(conn)
		message = string(me)
		if err != nil {
			fmt.Println("Something happend!")
			node.IsReserved = false
			matched.Message <- "ERR"
			g.Match.RUnlock()
			continue
		}
		matched.Message <- message
		other_message := <-node.Message
		if other_message == "Accept" && message == "Accept" {
			//play the game
			node.OtherPlayer = matched
			g.Match.RUnlock()
			break
		} else if other_message == "ERR" {
			node.IsReserved = false
			g.Match.RUnlock()
		} else {
			node.RemoveEdge(matched)
			node.IsReserved = false
			g.Match.RUnlock()
		}
	}
	// close(node.Indx)
	// close(node.Message)
	node.MatchEnd = false

	// Rolling dice
	fmt.Println("Rolling phase!")
	for {
		message, err := myNet.ReciveMessage(conn)
		if err != nil {
			fmt.Println("Error accourd in reading a massage")
			return
		}

		if string(message) == "Roll" && node.Turn {
			dice1 := rand.Intn(6) + 1
			dice2 := rand.Intn(6) + 1
			responce, err := myNet.MakePkt([]byte(fmt.Sprint(dice1) + "," + fmt.Sprint(dice2)))
			if err != nil {
				fmt.Println("Something happend in making a pkt")
			}
			node.Turn = false
			node.OtherPlayer.Turn = true
			fmt.Println(dice1, dice2)
			conn.Write(responce)
		} else if string(message) == "End" {
			node.MatchEnd = true
		}
		//TODO make a thread that check if the match has ended!
		if node.MatchEnd && node.OtherPlayer.MatchEnd {
			fmt.Println("Ending the match")
			break
		}
	}
}
