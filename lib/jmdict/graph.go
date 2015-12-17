package jmdict

// Graph struct
type Node struct {
	id   string
	data interface{}
	// weight or level in graph
	weight int
	mark   bool
}

// Constructor with default value
func makeSNode(data interface{}) Node {
	// Assign node id for tracking
	id := time.Now().Format("20060102150405")
	return Node{id, data, 3, false}
}
func makeKNode(data KanjiElement) Node {
	return Node{data.Keb, data, 1, false}
}
func makeRNode(data ReadingElement) Node {
	return Node{data.Reb, data, 2, false}
}

type Points []Node

func (points Points) Swap(i, j int) {
	points[i], points[j] = points[j], points[i]
}

func (points Points) Len() int {
	return len(points)
}

type XSortablePoints struct{ Points }

func (points XSortablePoints) Less(i, j int) bool {
	return points.Points[i].id < points.Points[j].id
}

// Return a directed graph & reverse one
func buildGraphs(entry JapEng) (map[string][]Node, map[string][]Node) {
	graph := map[string][]Node{}
	nodes := map[string]Node{}
	revGraph := map[string][]Node{}

	// Build K & R nodes
	for _, k := range entry.Kanji {
		n := makeKNode(k)
		nodes[n.id] = n
	}
	for _, r := range entry.Reading {
		n := makeRNode(r)
		nodes[n.id] = n
	}

	// Sort out distinct edges & common edges. Direction S -> R -> K
	// Make vertexes
	var rEdges []Node
	for _, s := range entry.Sense {
		node := makeSNode(s)
		if len(s.RestrictedToKanji) > 0 {
			// A special edge to K, mark it on node
			node.mark = true
			for _, k := range s.RestrictedToKanji {
				graph[k] = append(graph[k], node)
				revGraph[node.id] = append(revGraph[node.id], nodes[k])
			}
		}
		if len(s.RestrictedToReading) > 0 {
			// Edge to R, not belong to common edges
			for _, r := range s.RestrictedToReading {
				graph[r] = append(graph[r], node)
				revGraph[node.id] = append(revGraph[node.id], nodes[r])
			}
			continue
		}
		rEdges = append(rEdges, node)
	}

	var kEdges []Node
	for _, r := range entry.Reading {
		node := nodes[r.Reb]
		// Update common S-R vertexes
		graph[node.id] = append(graph[node.id][:0], append(rEdges, graph[node.id][0:]...)...)
		// Update Reverse graph
		for _, c := range rEdges {
			revGraph[c.id] = append(revGraph[c.id], node)
		}

		// K edge
		if len(r.RestrictedTo) > 0 {
			// Must be a K node - connect to this node only
			for _, k := range r.RestrictedTo {
				graph[k] = append(graph[k], node)
				revGraph[node.id] = append(revGraph[node.id], nodes[k])
			}
		} else {
			// Common edge to K nodes
			kEdges = append(kEdges, node)
		}

	}

	for _, k := range entry.Kanji {
		node := nodes[k.Keb]
		// Update common R-K vertexes
		graph[node.id] = append(graph[node.id][:0], append(kEdges, graph[node.id][0:]...)...)
		// Update Reverse graph
		for _, c := range kEdges {
			revGraph[c.id] = append(revGraph[c.id], node)
		}
	}

	return graph, revGraph
}

func getChildNodes(node Node, g map[string][]Node) []Node {
	return g[node.id]
}

// Impose contraint on collectable nodes
func collectValidNode(x Node, origin Node, nodeCollects map[int][]Node,
	graph map[string][]Node) {
	if x.id != origin.id {
		if markNode, ok := graph["mark_"+origin.id]; ok && x.weight == 1 && markNode[0].id != x.id {
			return
		}
		if otherNodes, ok := graph[x.id]; ok {
			for _, n := range otherNodes {
				if n.id != origin.id && n.weight == origin.weight {
					// Don't collect node with distinct edge to
					// another on same  level
					return
				}
			}
		}
		nodeCollects[x.weight] = append(nodeCollects[x.weight], x)
	}
}

// Return map of level number & list of nodes connected on that level
func DFS(node Node, g map[string][]Node) map[int][]Node {
	S := []Node{node}
	visited := map[string]bool{}
	// Track nodes type - <Type: Value1, Value 2..>
	nodeCollects := map[int][]Node{}
	var x Node
	for len(S) > 0 {
		// Pop a node
		x, S = S[len(S)-1], S[:len(S)-1]
		if _, ok := visited[x.id]; !ok {
			visited[x.id] = true
			// Collect valid node
			collectValidNode(x, node, nodeCollects, g)
			// Node not visited yet
			for _, edge := range getChildNodes(x, g) {
				S = append(S, edge)
			}
		}
	}
	return nodeCollects
}
