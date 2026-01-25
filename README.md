# Go Matrix Multiplication Service

A high-performance distributed system architecture designed to benchmark and visualize the efficiency of concurrent programming against sequential execution. This project implements a custom TCP protocol to offload computationally intensive matrix multiplication tasks to a dedicated server, utilizing Go's lightweight concurrency primitives (Goroutines) to achieve significant performance gains.

## Project Overview

The primary objective of this system is to demonstrate the mathematical and computational advantages of parallel processing. By splitting large matrix operations (O(N^3) complexity) across multiple worker threads, the system overcomes the limitations of single-threaded execution, providing a clear visualization of speedup factors through a generated web report.

### Key Features

* **Concurrent Computation Engine**: Utilizes worker pools and channels to distribute matrix rows for parallel processing.
* **Custom TCP Protocol**: Implements a robust JSON-based communication layer between Client and Server.
* **Real-Time Progress Tracking**: The server streams progress updates to the client during heavy computation.
* **Automated Performance Analysis**: Generates a self-contained HTML report with CSS visualizations and CSV export capabilities.
* **CLI Interface**: Fully configurable via command-line flags for matrix size, worker count, and chunk size.

## System Architecture

The project follows a strict Client-Server model:

1.  **Server (`cmd/server`)**:
    * Listens on a TCP port.
    * Accepts connection requests containing job parameters (Matrix Size, Worker Count).
    * Generates two random N x N matrices.
    * Executes both Sequential and Concurrent multiplication algorithms.
    * Streams results and performance metrics back to the client.

2.  **Client (`cmd/client`)**:
    * Connects to the server via TCP.
    * Sends configuration parameters.
    * Displays a real-time progress bar in the terminal.
    * Receives the final payload and generates a local `report.html` for analysis.

## Prerequisites

* **Go**: Version 1.18 or higher.
* **Make**: (Optional) For streamlined build commands.
* **Web Browser**: To view the generated performance reports.

## Getting Started

### 1. Clone the repo and navigate to the folder
```bash
git clone https://github.com/mohamadyoussefio/go-matrix-service.git
cd go-matrix-service
```

### 2. Build and Run the Server
The server must be running to accept incoming jobs.
```bash
# Using Makefile
make server
```

```bash
# Manual Execution
go run ./cmd/srver
```

### 3. Run the Client

Open a new terminal window to run the client. You can configure the workload using flags.

#### Default Run (1000x1000 Matrix, 8 Workers):
```bash
make client
```

#### Custom Configuration:
```bash
# Run with a 2000x2000 matrix and 16 concurrent workers
make client s=2000 w=16
```

## Performance Reporting

Upon completion of a job, the client automatically generates a report.html file in the project root. This report includes:

1. Executive Summary: Immediate visual comparison of Sequential vs. Concurrent execution time.

2. Speedup Factor: The calculated efficiency multiplier (e.g., 4.14x).

3. Data Export: A specialized button to download raw benchmark data as a .csv file for external analysis.

4. Integrity Check: A checksum hash to verify that the concurrent algorithm produced mathematically identical results to the sequential method.
