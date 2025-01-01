package network

import "sync"

type Node struct {
	edges      []*Node
	ip         string
	username   string
	indx       chan *Node
	isReserved bool
	message    chan string
}

func (n *Node) RemoveEdge(node *Node) {
	for indx, nod := range n.edges {
		if nod == node {
			n.edges = append(n.edges[:indx], n.edges[indx+1:]...)
			break
		}
	}
}

type Graph struct {
	nodes []*Node
	m     *map[string]Node
	mutex sync.Mutex
	match sync.RWMutex
}

func (g *Graph) addNode(node *Node) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	node.edges = append(node.edges, (g.nodes)...)
	g.nodes = append(g.nodes, node)
}

func (g *Graph) NotReserved() int {
	counter := 0
	for _, n := range g.nodes {
		if !n.isReserved {
			counter++
		}
	}
	return counter
}
