package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/njambi-r/nairobi-route-finder/internal/graph"
)

func main() {
	// Load the graph from file
	graphPath := "data/NCR+BRT_v1.json"
	g, err := graph.LoadGraphFromFile(graphPath)
	if err != nil {
		log.Fatalf("Failed to load graph: %v", err)
	}

	// Handle route queries
	http.HandleFunc("/routes", func(w http.ResponseWriter, r *http.Request) {
		// Allow cross-origin requests from local frontend
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		from := r.URL.Query().Get("from")
		to := r.URL.Query().Get("to")
		maxDepthStr := r.URL.Query().Get("maxdepth")

		if from == "" || to == "" {
			http.Error(w, "`from` and `to` query parameters are required", http.StatusBadRequest)
			return
		}

		// Default values
		maxDepth := 30
		if maxDepthStr != "" {
			if md, err := strconv.Atoi(maxDepthStr); err == nil {
				maxDepth = md
			}
		}

		maxRoutes := 10

		// ✅ Pass both maxDepth and maxRoutes
		routes := g.FindRoutesBetweenStations(from, to, maxDepth, maxRoutes)

		resp := map[string]interface{}{
			"from":   from,
			"to":     to,
			"count":  len(routes),
			"routes": routes,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
