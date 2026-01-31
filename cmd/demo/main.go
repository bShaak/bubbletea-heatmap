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

	// customTheme := []string{"#e0f7fa", "#b2ebf2", "#80deea", "#4dd0e1", "#26c6da"}

	// Create the model
	// Use a date that covers the test data (ends March 2023)
	refDate := time.Date(2023, 3, 4, 0, 0, 0, 0, time.UTC)
	// m := heatmap.New(calData, refDate, 0, heatmap.WithTheme("light"))
	m := heatmap.New(calData, refDate, 0, heatmap.WithTheme("dark"))
	// m := heatmap.New(calData, refDate, 0, heatmap.WithTheme(customTheme))

	// Run the bubble tea program
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}

