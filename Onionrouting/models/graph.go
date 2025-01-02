package models

import "sync"

type Node struct {
	Edges       []*Node
	Ip          string
	Username    string
	Indx        chan *Node
	IsReserved  bool
	Message     chan string
	Turn        bool
	OtherPlayer *Node
	MatchEnd    bool
}

func (n *Node) RemoveEdge(node *Node) {
	for indx, nod := range n.Edges {
		if nod == node {
			n.Edges = append(n.Edges[:indx], n.Edges[indx+1:]...)
			break
		}
	}
}

type Graph struct {
	Nodes []*Node
	M     *map[string]Node
	Mutex sync.Mutex
	Match sync.RWMutex
}

func (g *Graph) AddNode(node *Node) {
	g.Mutex.Lock()
	defer g.Mutex.Unlock()
	node.Edges = append(node.Edges, (g.Nodes)...)
	g.Nodes = append(g.Nodes, node)
}

func (g *Graph) NotReserved() int {
	counter := 0
	for _, n := range g.Nodes {
		if !n.IsReserved {
			counter++
		}
	}
	return counter
}
