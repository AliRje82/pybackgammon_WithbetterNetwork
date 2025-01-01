package network

func matchMaking(g *Graph) {
	g.match.Lock()
	defer g.match.Unlock()
	if g.NotReserved() <= 1 {
		return
	}
	for _, node := range g.nodes {
		if node.isReserved {
			continue
		}
		for _, e := range node.edges {
			if !e.isReserved {
				node.isReserved = true
				e.isReserved = true
				e.indx <- node
				node.indx <- e
			}
		}
	}
}
