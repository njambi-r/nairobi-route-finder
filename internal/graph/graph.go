package graph

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Node struct {
	Name string `json:"name,omitempty"`
}

type Line struct {
	Name  string `json:"name"`
	Nodes []Node `json:"nodes"`
}

type TransitData struct {
	Lines []Line `json:"lines"`
}

type Graph struct {
	Adj map[string][]string
}

func NewGraph() *Graph {
	return &Graph{
		Adj: make(map[string][]string),
	}
}

// Normalizes station names (removes \n, trims, title-cases)
func normalize(name string) string {
	name = strings.ReplaceAll(name, "\n", " ")
	name = strings.TrimSpace(name)
	return strings.Title(strings.ToLower(name))
}

// LoadGraphFromFile builds a graph from schematic JSON, skipping unnamed (schematic-only) nodes.
func LoadGraphFromFile(path string) (*Graph, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var data TransitData
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	graph := NewGraph()

	for _, line := range data.Lines {
		var prev string
		for _, node := range line.Nodes {
			if node.Name == "" {
				continue
			}
			curr := normalize(node.Name)
			if prev != "" {
				graph.Adj[prev] = append(graph.Adj[prev], curr)
				graph.Adj[curr] = append(graph.Adj[curr], prev)
			}
			prev = curr
		}
	}

	return graph, nil
}

// BFS returns all shortest routes from start to goal
func (g *Graph) FindAllRoutes(start, goal string) ([][]string, error) {
	start = normalize(start)
	goal = normalize(goal)

	if _, ok := g.Adj[start]; !ok {
		return nil, errors.New("start station not found")
	}
	if _, ok := g.Adj[goal]; !ok {
		return nil, errors.New("goal station not found")
	}

	var results [][]string
	visited := map[string]int{start: 0}
	queue := [][]string{{start}}
	minLength := -1

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]

		last := path[len(path)-1]
		if minLength != -1 && len(path) > minLength {
			break // all further paths are longer than the shortest found
		}
		if last == goal {
			if minLength == -1 {
				minLength = len(path)
			}
			results = append(results, path)
			continue
		}

		for _, neighbor := range g.Adj[last] {
			if prevLen, seen := visited[neighbor]; !seen || len(path) <= prevLen {
				visited[neighbor] = len(path)
				newPath := append([]string{}, path...)
				newPath = append(newPath, neighbor)
				queue = append(queue, newPath)
			}
		}
	}

	if len(results) == 0 {
		return nil, errors.New("no route found")
	}
	return results, nil
}
