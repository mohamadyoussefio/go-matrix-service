package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go-matrix-service/internal/protocol"
	"log"
	"net"
	"strings"
	"time"
)

func main() {
	size := flag.Int("n", 1000, "Matrix size")
	workers := flag.Int("w", 8, "Workers")
	chunk := flag.Int("c", 100, "Chunk size")
	host := flag.String("host", "localhost:8080", "Server address")
	flag.Parse()

	conn, err := net.Dial("tcp", *host)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	req := protocol.Request{
		MatrixSize: *size,
		Workers:    *workers,
		ChunkSize:  *chunk,
		Seed:       42,
	}

	if err := json.NewEncoder(conn).Encode(req); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Connected to %s. Job: %dx%d, %d workers\n", *host, *size, *size, *workers)
	fmt.Println("Waiting for server processing...")

	dec := json.NewDecoder(conn)
	start := time.Now()

	for {
		var resp protocol.Response
		if err := dec.Decode(&resp); err != nil {
			break
		}

		if resp.Type == "progress" {
			printBar(resp.RowsProcessed, resp.TotalRows, start)
		} else if resp.Type == "result" {
			printResult(resp)
			return
		}
	}
}

func printBar(curr, total int, start time.Time) {
	width := 30
	percent := float64(curr) / float64(total)
	filled := int(percent * float64(width))
	
	bar := strings.Repeat("=", filled) + strings.Repeat("-", width-filled)
	rate := float64(curr) / time.Since(start).Seconds()
	
	fmt.Printf("\r[%s] %.0f%% | %d/%d | %.0f rows/s", bar, percent*100, curr, total, rate)
}

func printResult(r protocol.Response) {
	fmt.Println("\n\n-----------------------------------")
	fmt.Println("        BENCHMARK RESULTS")
	fmt.Println("-----------------------------------")
	fmt.Printf("Sequential Time : %8.4f s\n", r.SeqTime)
	fmt.Printf("Concurrent Time : %8.4f s\n", r.ConcTime)
	fmt.Printf("Speedup         : %8.2f x\n", r.Speedup)
	fmt.Printf("Checksum        : %.4e\n", r.Checksum)
	fmt.Println("-----------------------------------")
}
