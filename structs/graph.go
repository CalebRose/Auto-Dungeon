package structs

// Graph - Map Structure for missions
type Graph struct {
	AdjList map[string]*Room
}

// AddVertex - Adding Rooms to the map
func (g *Graph) AddVertex(v1 *Room) {
	name := v1.Key
	_, found := g.AdjList[name]
	if found == true {
		return
	}
	g.AdjList[name] = v1
}

// func (g *Graph) AddEdge (v1, v2 *Room){
// 	n1 := v1.Name
// 	n2 := v2.Name
// 	_, found1 := g.AdjList[n1]
// 	_, found2 := g.AdjList[n2]
// 	if found1 == false || found2 == false {
// 		fmt.Println("One of the two rooms does not exist within the Graph")
// 		return
// 	}
// 		g.AdjList[n1].Edges = append(g.AdjList[n1].Edges, n2)
// 		g.AdjList[n2].Edges = append(g.AdjList[n2].Edges, n1)
// }
