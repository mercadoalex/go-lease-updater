package updater

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"

	"gopkg.in/yaml.v2"
)

// LeaseList represents the top-level structure of the leases.yaml file.
type LeaseList struct {
	Items []Lease `yaml:"items"`
}

// Lease represents a single lease entry in the YAML file.
type Lease struct {
	Metadata struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"metadata"`
	Spec struct {
		RenewTime            string `yaml:"renewTime"`
		LeaseDurationSeconds int    `yaml:"leaseDurationSeconds"`
	} `yaml:"spec"`
}

// UpdateRenewTimeWithChaos updates leases and injects chaos using user-defined chaosMin and chaosMax.
// It also randomly expires leases with a low probability, simulating dead nodes.
func UpdateRenewTimeWithChaos(filePath string, chaosMin, chaosMax int) (deadNodes int, chaosEvents int) {
	// Read the YAML file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("failed to read file: %v\n", err)
		return 0, 0
	}

	// Unmarshal YAML into LeaseList structure
	var leaseList LeaseList
	if err := yaml.Unmarshal(data, &leaseList); err != nil {
		fmt.Printf("failed to unmarshal YAML: %v\n", err)
		return 0, 0
	}

	// Seed the random number generator for chaos injection
	rand.Seed(time.Now().UnixNano())
	chaosIdx := rand.Intn(len(leaseList.Items)) // Select a random lease for chaos

	now := time.Now().UTC() // Current time for dead node check
	maxDeadPerRun := 1      // Limit dead node reports per run
	deadCount := 0

	// Iterate over all leases
	for i := range leaseList.Items {
		rt := leaseList.Items[i].Spec.RenewTime
		leaseList.Items[i].Spec.LeaseDurationSeconds = 40 // Always set to 40 (Kubernetes default)

		// Parse the current renewTime
		t, err := time.Parse(time.RFC3339Nano, rt)
		if err != nil {
			continue
		}

		var newRenewTime string
		if i == chaosIdx {
			// Chaos: add random seconds between chaosMin and chaosMax to this lease's renewTime
			chaosSeconds := rand.Intn(chaosMax-chaosMin+1) + chaosMin
			newRenewTime = t.Add(time.Duration(chaosSeconds) * time.Second).UTC().Format(time.RFC3339Nano)
			leaseList.Items[i].Spec.RenewTime = newRenewTime
			fmt.Println("Chaos!!!")
			chaosEvents++
		} else {
			// Normal: add 10 seconds to renewTime
			newRenewTime = t.Add(10 * time.Second).UTC().Format(time.RFC3339Nano)
			leaseList.Items[i].Spec.RenewTime = newRenewTime
		}

		// Randomly expire a lease with a low probability (e.g., 10%)
		if rand.Float64() < 0.1 {
			leaseList.Items[i].Spec.RenewTime = time.Now().Add(-41 * time.Second).UTC().Format(time.RFC3339Nano)
			fmt.Println("Randomly expired node!")
		}

		// Print node name and updated renewTime
		fmt.Printf("Node: %s, Updated renewTime: %s\n", leaseList.Items[i].Metadata.Name, leaseList.Items[i].Spec.RenewTime)

		// Dead node detection: only increment if expired (renewTime + 40s < now)
		expiry, err := time.Parse(time.RFC3339Nano, leaseList.Items[i].Spec.RenewTime)
		if err != nil {
			continue
		}
		if now.After(expiry.Add(40 * time.Second)) {
			if deadCount < maxDeadPerRun {
				fmt.Println("---------------------------------------------------")
				fmt.Printf("##### DEAD NODE DETECTED: %s (expired at %s) #####\n",
					leaseList.Items[i].Metadata.Name,
					expiry.Add(40*time.Second).Format(time.RFC3339Nano))
				fmt.Println("---------------------------------------------------")
				deadCount++
				deadNodes++
				// Heal the lease so it won't be counted as dead in the next cycle
				leaseList.Items[i].Spec.RenewTime = time.Now().Add(10 * time.Second).UTC().Format(time.RFC3339Nano)
			}
		}
	}

	// Marshal the updated LeaseList back to YAML
	updatedData, err := yaml.Marshal(leaseList)
	if err != nil {
		fmt.Printf("failed to marshal updated leases: %v\n", err)
		return deadNodes, chaosEvents
	}

	// Write the updated YAML back to the file
	if err := ioutil.WriteFile(filePath, updatedData, 0644); err != nil {
		fmt.Printf("failed to write updated file: %v\n", err)
	}

	return deadNodes, chaosEvents
}
