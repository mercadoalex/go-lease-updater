package updater

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"

	"go-lease-updater/types"

	"gopkg.in/yaml.v2"
)

var cycleCount int
var deadNodesInWindow int
var reportedDeadNodes = make(map[string]bool)

// UpdateRenewTimeWithChaos updates leases and injects chaos using user-defined chaosMin and chaosMax.
// It also randomly expires leases with a low probability, simulating dead nodes.
// Dead node reporting is limited to a minimum of 1 and a maximum of 4 per 2-minute window.
func UpdateRenewTimeWithChaos(filePath string, chaosMin, chaosMax int) (deadNodes int, chaosEvents int) {
	// Read the YAML file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("failed to read file: %v\n", err)
		return
	}

	// Unmarshal YAML into LeaseList structure
	var leaseList types.LeaseList
	if err := yaml.Unmarshal(data, &leaseList); err != nil {
		fmt.Printf("failed to unmarshal YAML: %v\n", err)
		return
	}

	rand.Seed(time.Now().UnixNano())
	chaosIdx := rand.Intn(len(leaseList.Items)) // Select a random lease for chaos
	now := time.Now().UTC()

	cycleCount++
	if cycleCount%12 == 1 { // Reset every 2 minutes (12 cycles)
		deadNodesInWindow = 0
		reportedDeadNodes = make(map[string]bool)
	}

	maxDeadPerWindow := 4
	minDeadPerWindow := 1

	// Collect indices of expired leases
	expiredIndices := []int{}

	// Limit to at most 1 random expiration per cycle
	maxRandomExpirations := 1
	randomExpirations := 0

	injectChaos := rand.Float64() < 0.8 // 80% chance to inject chaos this cycle

	for i := range leaseList.Items {
		if injectChaos && i == chaosIdx {
			chaosSeconds := rand.Intn(chaosMax-chaosMin+1) + chaosMin
			leaseList.Items[i].Spec.RenewTime = now.Add(time.Duration(chaosSeconds) * time.Second).UTC().Format(time.RFC3339Nano)
			fmt.Println("Chaos!!!")
			chaosEvents++
		} else {
			lease := &leaseList.Items[i]
			lease.Spec.LeaseDurationSeconds = 40

			// Always use now for renewTime updates to avoid accumulation of expired leases
			if randomExpirations < maxRandomExpirations && rand.Float64() < 0.03 { // 3% chance
				lease.Spec.RenewTime = now.Add(-41 * time.Second).UTC().Format(time.RFC3339Nano)
				fmt.Println("Randomly expired node!")
				randomExpirations++
			} else {
				lease.Spec.RenewTime = now.Add(10 * time.Second).UTC().Format(time.RFC3339Nano)
			}

			fmt.Printf("Node: %s, Updated renewTime: %s\n", lease.Metadata.Name, lease.Spec.RenewTime)

			expiry, err := time.Parse(time.RFC3339Nano, lease.Spec.RenewTime)
			if err != nil {
				continue
			}
			if now.After(expiry.Add(40 * time.Second)) {
				expiredIndices = append(expiredIndices, i)
			}
		}
	}

	// Limit dead nodes per window
	if len(expiredIndices) > maxDeadPerWindow {
		expiredIndices = expiredIndices[:maxDeadPerWindow]
	}
	if len(expiredIndices) < minDeadPerWindow && len(expiredIndices) > 0 {
		expiredIndices = expiredIndices[:minDeadPerWindow]
	}

	for _, i := range expiredIndices {
		if deadNodesInWindow < maxDeadPerWindow {
			nodeName := leaseList.Items[i].Metadata.Name
			if !reportedDeadNodes[nodeName] {
				deadNodes++
				deadNodesInWindow++
				reportedDeadNodes[nodeName] = true
				fmt.Printf("##### DEAD NODE DETECTED: %s #####\n", nodeName)
			}
		}
	}

	// Marshal the updated LeaseList back to YAML
	updatedData, err := yaml.Marshal(leaseList)
	if err != nil {
		fmt.Printf("failed to marshal updated leases: %v\n", err)
		return
	}

	// Write the updated YAML back to the file
	if err := ioutil.WriteFile(filePath, updatedData, 0644); err != nil {
		fmt.Printf("failed to write updated file: %v\n", err)
	}

	return
}
