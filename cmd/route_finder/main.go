package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/njambi-r/nairobi-route-finder/internal/graph"
)

func main() {
	from := flag.String("from", "", "Source station")
	to := flag.String("to", "", "Destination station")
	jsonOut := flag.Bool("json", false, "Output results in JSON format")
	maxDepth := flag.Int("maxdepth", 30, "Maximum route search depth") //removed since only using BFS
	maxRoutes := flag.Int("maxroutes", 10, "Maximum number of routes to find")
	flag.Parse()

	if *from == "" || *to == "" {
		fmt.Println("Usage: route_finder --from='Station A' --to='Station B' [--json] [--maxdepth=N] [--maxroutes=N]")
		os.Exit(1)
	}

	// Determine the absolute path of the JSON graph file
	execDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}
	graphPath := filepath.Join(execDir, "data", "NCR+BRT_v1.json")

	// Load the graph
	g, err := graph.LoadGraphFromFile(graphPath)
	if err != nil {
		log.Fatalf("Failed to load graph from %s: %v", graphPath, err)
	}

	// Debug: print keys (station nodes) found for input stations
	fmt.Printf("Looking up start station '%s': keys found %v\n", *from, g.GetKeysForStation(*from))
	fmt.Printf("Looking up destination station '%s': keys found %v\n", *to, g.GetKeysForStation(*to))

	// Find all routes
	routes := g.FindRoutesBetweenStations(*from, *to, *maxDepth, *maxRoutes)

	if *jsonOut {
		output := map[string]interface{}{
			"from":   *from,
			"to":     *to,
			"routes": routes,
			"count":  len(routes),
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(output); err != nil {
			log.Fatalf("Failed to encode JSON output: %v", err)
		}
	} else {
		fmt.Printf("Found %d route(s):\n", len(routes))
		for i, r := range routes {
			fmt.Printf("Route %d: %s\n", i+1, formatRoute(r))
		}
	}
}

// formatRoute formats a slice of station names as a string
func formatRoute(stations []string) string {
	return fmt.Sprintf("%v", stations)
}
