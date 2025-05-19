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
	flag.Parse()

	if *from == "" || *to == "" {
		fmt.Println("Usage: route_finder --from='Station A' --to='Station B' [--json]")
		os.Exit(1)
	}

	// Load graph from the local data path
	execDir, _ := os.Getwd()
	graphPath := filepath.Join(execDir, "data", "NCR+BRT_v1.json")

	g, err := graph.LoadGraphFromFile(graphPath)
	if err != nil {
		log.Fatalf("Failed to load graph: %v", err)
	}

	routes, err := g.FindAllRoutes(*from, *to)
	if err != nil {
		log.Fatalf("Route finding failed: %v", err)
	}

	if *jsonOut {
		output := map[string]interface{}{
			"from":   *from,
			"to":     *to,
			"routes": routes,
			"count":  len(routes),
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(output)
	} else {
		fmt.Printf("Found %d route(s):\n", len(routes))
		for i, r := range routes {
			fmt.Printf("Route %d: %s\n", i+1, formatRoute(r))
		}
	}
}

func formatRoute(stations []string) string {
	return fmt.Sprintf("%s", stations)
}
