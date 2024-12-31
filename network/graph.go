package network

import "sync"

type Node struct {
	edges    []*Node
	ip       string
	username string
	indx     chan int
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
