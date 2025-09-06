package main

import (
	"fmt"
	"go-lease-updater/metrics"
	"go-lease-updater/updater"
	"os"
	"time"
)

const yamlFilePath = "leases.yaml" // match the file location

func main() {
	var minutes int
	var chaosMin, chaosMax int

	// Metrics tracking
	metricsData := metrics.Metrics{
		StartTime: time.Now(),
	}

	// Get minutes input (only integer 1-20, default 2)
	for {
		fmt.Print("How many minutes should the program run? (integer 1-20, default 2): ")
		n, err := fmt.Scanf("%d", &minutes)
		if n == 0 || err != nil {
			minutes = 2
			break
		}
		if minutes < 1 || minutes > 20 {
			fmt.Println("Value out of range. Please enter an integer between 1 and 20.")
			minutes = 2
			break
		}
		break
	}

	// Get chaos range input (minimum and maximum seconds)
	for {
		fmt.Print("Enter chaos minimum seconds (default 11): ")
		n, err := fmt.Scanf("%d", &chaosMin)
		if n == 0 || err != nil || chaosMin < 1 {
			chaosMin = 11
		}
		fmt.Print("Enter chaos maximum seconds (default 60): ")
		n, err = fmt.Scanf("%d", &chaosMax)
		if n == 0 || err != nil || chaosMax <= chaosMin {
			chaosMax = 60
		}
		if chaosMax > chaosMin {
			break
		}
		fmt.Println("Chaos maximum must be greater than chaos minimum.")
	}

	fmt.Printf("Running for %d minute(s)... Chaos range: %d-%d seconds\n", minutes, chaosMin, chaosMax)

	end := time.Now().Add(time.Duration(minutes) * time.Minute)
	for time.Now().Before(end) {
		// Update leases and collect metrics
		deads, chaos := updater.UpdateRenewTimeWithChaos(yamlFilePath, chaosMin, chaosMax)
		metricsData.DeadNodes += deads
		metricsData.ChaosEvents += chaos

		time.Sleep(10 * time.Second)
	}
	metricsData.EndTime = time.Now()

	fmt.Println("Program finished.")
	fmt.Println("OK this the end")
	metrics.PrintSummaryTable(metricsData)
	os.Exit(0)
}
