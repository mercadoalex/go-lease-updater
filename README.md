# Go Lease Updater

This project is a Go application that simulates Kubernetes lease updates and chaos injection. It updates the `renewTime` fields in a YAML file every 10 seconds, injects random chaos, and occasionally simulates dead nodes by expiring leases. The application outputs a summary table with metrics at the end of each run.

## Project Structure

```
go-lease-updater
├── main.go
├── updater
│   └── updater.go
├── metrics
│   └── metrics.go
├── types
│   └── lease.go
├── leases.yaml
├── go.mod
└── README.md
```

## Requirements

- Go 1.16 or later
- A YAML file (`leases.yaml`) containing lease data

## Installation

1. Clone the repository:
   ```
   git clone <repository-url>
   cd go-lease-updater
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

## Configuration

- Place your `leases.yaml` file in the project root directory.
- Update the file path in `main.go` if you use a different location.

## Running the Application

To run the application, execute the following command from the project root:

```
go run main.go
```

You will be prompted for:
- The number of minutes to run the simulation
- Chaos minimum and maximum seconds

## Logic Overview

- Every 10 seconds, the application updates all lease `renewTime` fields.
- One random lease receives a chaos injection (randomly extended renew time).
- With a low probability (10%), any lease may be randomly expired to simulate a dead node.
- Dead nodes and chaos events are tracked independently.
- At the end, a summary table is printed showing duration, chaos events, and dead nodes detected.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.