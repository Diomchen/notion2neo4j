package main

import (
	"fmt"
	"notion2neo4j/process"
	"path/filepath"
)

func main() {
	// Configuration
	parentFolder := "C:\\Users\\thgy\\Downloads\\aee088f0-0835-432b-b3d1-01cff038238d_Export-a5f5a71a-f4af-4607-befc-47391251e0fe"
	csvPath := filepath.Join(parentFolder, "网吧 1b4ee31ce44080279e25e167f564e0bb.csv")
	dataFolder := filepath.Join(parentFolder, "网吧 1b4ee31ce44080279e25e167f564e0bb")

	// 1. Parse the CSV file
	entries, err := process.ProcessCSV(csvPath)
	if err != nil {
		fmt.Printf("Error parsing CSV: %v\n", err)
		return
	}

	// 2. Load markdown files
	articles, err := process.ProcessArticle(dataFolder, entries)
	if err != nil {
		fmt.Printf("Error loading articles: %v\n", err)
		return
	}

	// 3. Create Neo4j graph
	err = process.CreateNeo4jGraph(entries, articles)
	if err != nil {
		fmt.Printf("Error creating Neo4j graph: %v\n", err)
		return
	}
}
