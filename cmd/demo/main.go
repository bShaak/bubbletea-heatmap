package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	heatmap "github.com/slinlee/bubbletea-heatmap"
)

func main() {
	// Read the test data
	// Assuming we run from the project root
	data, err := os.ReadFile("tests/test.json")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Parse the JSON data
	var calData []heatmap.CalDataPoint
	if err := json.Unmarshal(data, &calData); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	// Create the model
	// Use a date that covers the test data (ends March 2023)
	refDate := time.Date(2023, 3, 4, 0, 0, 0, 0, time.UTC)
	m := heatmap.New(calData, refDate, 12)

	// Run the bubble tea program
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}

