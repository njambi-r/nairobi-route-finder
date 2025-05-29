// internal/graph/graph.go
package graph

import (
	"encoding/json"
	"os"
	"strings"
)

type Station struct {
	Label    string   `json:"label"`
	Position Position `json:"position"`
}

type Position struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Node struct {
	Name   string    `json:"name"`
	Coords []float64 `json:"coords"`
}

type Line struct {
	Name  string `json:"name"`
	Color string `json:"color"`
	Nodes []Node `json:"nodes"`
}

type RawGraph struct {
	Stations map[string]Station `json:"stations"`
	Lines    []Line             `json:"lines"`
}

type Graph struct {
	Adjacency map[string][]string
}

func LoadGraphFromFile(filename string) (*Graph, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var raw RawGraph
	err = json.Unmarshal(file, &raw)
	if err != nil {
		return nil, err
	}

	adj := make(map[string][]string)

	for _, line := range raw.Lines {
		prev := ""
		for _, node := range line.Nodes {
			name := strings.TrimSpace(node.Name)
			if name == "" {
				continue // skip unnamed (schematic-only) nodes
			}
			if prev != "" && name != prev {
				adj[prev] = append(adj[prev], name)
				adj[name] = append(adj[name], prev)
			}
			prev = name
		}
	}

	return &Graph{Adjacency: adj}, nil
}

func (g *Graph) GetKeysForStation(name string) []string {
	keys := []string{}
	for k := range g.Adjacency {
		if strings.EqualFold(k, name) {
			keys = append(keys, k)
		}
	}
	return keys
}

func (g *Graph) FindRoutesBetweenStations(from, to string, maxDepth int, maxRoutes int) [][]string {
	results := [][]string{}
	visited := map[string]bool{}
	path := []string{}

	for _, startKey := range g.GetKeysForStation(from) {
		g.dfs(startKey, to, maxDepth, visited, path, &results, maxRoutes)
		if len(results) >= maxRoutes {
			break
		}
	}
	return results
}

func (g *Graph) dfs(current, target string, depth int, visited map[string]bool, path []string, results *[][]string, maxRoutes int) {
	if depth < 0 || len(*results) >= maxRoutes {
		return
	}
	visited[current] = true
	path = append(path, current)

	if strings.EqualFold(current, target) {
		copyPath := make([]string, len(path))
		copy(copyPath, path)
		*results = append(*results, copyPath)
	} else {
		for _, neighbor := range g.Adjacency[current] {
			if !visited[neighbor] {
				g.dfs(neighbor, target, depth-1, visited, path, results, maxRoutes)
			}
		}
	}
	visited[current] = false
}

func (g *Graph) FindShortestRoutesBFS(from, to string, maxRoutes int) [][]string {
	type Path struct {
		Current string
		Trace   []string
	}

	var results [][]string
	visited := map[string]bool{}
	queue := []Path{}

	// Initialize queue with all matching start keys
	for _, startKey := range g.GetKeysForStation(from) {
		queue = append(queue, Path{Current: startKey, Trace: []string{startKey}})
	}

	for len(queue) > 0 && len(results) < maxRoutes {
		curr := queue[0]
		queue = queue[1:]

		// Skip already visited nodes at this path depth
		if visitedKey := strings.Join(curr.Trace, "->"); visited[visitedKey] {
			continue
		}
		visited[strings.Join(curr.Trace, "->")] = true

		// Check if current node is a destination
		if matchesTarget(curr.Current, to) {
			results = append(results, curr.Trace)
			continue
		}

		for _, neighbor := range g.Adjacency[curr.Current] {
			// Prevent cycles
			if contains(curr.Trace, neighbor) {
				continue
			}
			newPath := append([]string{}, curr.Trace...)
			newPath = append(newPath, neighbor)
			queue = append(queue, Path{Current: neighbor, Trace: newPath})
		}
	}

	return results
}

// Helper: check if current station name matches the target (case-insensitive)
func matchesTarget(current, target string) bool {
	return strings.EqualFold(current, target)
}

// Helper: check if a station exists in the current path
func contains(path []string, station string) bool {
	for _, s := range path {
		if s == station {
			return true
		}
	}
	return false
}
