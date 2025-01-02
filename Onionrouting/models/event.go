package models

func MatchMaking(g *Graph) {
	g.Match.Lock()
	defer g.Match.Unlock()
	if g.NotReserved() <= 1 {
		return
	}
	for _, node := range g.Nodes {
		if node.IsReserved {
			continue
		}
		for _, e := range node.Edges {
			if !e.IsReserved {
				node.IsReserved = true
				e.IsReserved = true
				e.Indx <- node
				node.Indx <- e
				node.Turn = true
				e.Turn = false
			}
		}
	}
}
