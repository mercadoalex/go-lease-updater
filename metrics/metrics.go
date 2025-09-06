package metrics

import (
	"fmt"
	"time"
)

// Metrics holds the summary data for the run
type Metrics struct {
	StartTime   time.Time
	EndTime     time.Time
	DeadNodes   int
	ChaosEvents int
}

// PrintSummaryTable outputs the insights summary table
func PrintSummaryTable(m Metrics) {
	duration := m.EndTime.Sub(m.StartTime)
	fmt.Println("\nSummary Table")
	fmt.Println("Metric        | Value")
	fmt.Println("--------------|---------------------")
	fmt.Printf("Duration      | ~%v\n", formatDuration(duration))
	fmt.Printf("Dead Nodes    | %d\n", m.DeadNodes)
	fmt.Printf("Chaos Events  | %d\n", m.ChaosEvents)
}

// formatDuration returns a human-readable duration string
func formatDuration(d time.Duration) string {
	mins := int(d.Minutes())
	secs := int(d.Seconds()) % 60
	return fmt.Sprintf("%d min %d seconds", mins, secs)
}
