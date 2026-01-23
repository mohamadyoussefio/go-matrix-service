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

const (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
	Gray   = "\033[90m"
)

func main() {
	size := flag.Int("n", 1000, "Matrix size (NxN)")
	workers := flag.Int("w", 8, "Number of workers")
	chunk := flag.Int("c", 100, "Chunk size")
	host := flag.String("host", "localhost:8080", "Server address")
	flag.Parse()

	printHeader()

	fmt.Printf("%s[+] Connecting to server at %s...%s\n", Gray, *host, Reset)
	conn, err := net.Dial("tcp", *host)
	if err != nil {
		log.Fatalf("%s[!] Connection Failed: %v%s", Red, err, Reset)
	}
	defer conn.Close()
	fmt.Printf("%s[+] Connected!%s\n", Green, Reset)

	req := protocol.Request{
		MatrixSize: *size,
		Workers:    *workers,
		ChunkSize:  *chunk,
		Seed:       42,
	}
	if err := json.NewEncoder(conn).Encode(req); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s[+] Job Config:%s %dx%d Matrix | %d Workers\n", Cyan, Reset, *size, *size, *workers)
	fmt.Printf("%s[.] Waiting for server (Data Gen + Sequential Benchmark)...%s\n", Yellow, Reset)

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

func printHeader() {
	fmt.Print("\033[H\033[2J") // Clear screen
	fmt.Println(Cyan + Bold + `
   __  __      _       _        
  |  \/  |    | |     (_)       
  | \  / | __ | |_ ___ _ __  __ 
  | |\/| |/ _'| __/ __| |\ \/ / 
  | |  | | (_| | || (__| | >  <  
  |_|  |_|\__,_|\__\___|_|/_/\_\ 
      HIGH PERF CONCURRENCY      
` + Reset)
}

func printBar(curr, total int, start time.Time) {
	const width = 40
	percent := float64(curr) / float64(total)
	filled := int(percent * float64(width))

	color := Yellow
	if curr == total {
		color = Green
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("-", width-filled)
	
	elapsed := time.Since(start).Seconds()
	rate := float64(curr) / elapsed

	fmt.Printf("\r%s[%s] %3.0f%% %s| %d/%d | %.0f rows/s   ", 
		color, bar, percent*100, Reset, curr, total, rate)
}

func printResult(r protocol.Response) {
	fmt.Println("\n") // Drop down from progress bar
	
	fmt.Println(Gray + "┌──────────────────────────────────────────┐" + Reset)
	fmt.Printf("%s│           BENCHMARK RESULTS              │%s\n", Bold, Reset)
	fmt.Println(Gray + "├──────────────────────────────────────────┤" + Reset)
	
	fmt.Printf("│ Sequential Time : %s%8.4f s%s             │\n", Red, r.SeqTime, Reset)
	fmt.Printf("│ Concurrent Time : %s%8.4f s%s             │\n", Green, r.ConcTime, Reset)
	fmt.Println(Gray + "├──────────────────────────────────────────┤" + Reset)
	
	fmt.Printf("│ Speedup Factor  : %s%8.2f x%s             │\n", Cyan+Bold, r.Speedup, Reset)
	fmt.Printf("│ Integrity Check : %s%9.2e%s              │\n", Yellow, r.Checksum, Reset)
	fmt.Println(Gray + "└──────────────────────────────────────────┘" + Reset)
	fmt.Println("")

	GenerateHTML(r)
}
