package main

import (
	"encoding/json"
	"fmt"
	"go-matrix-service/internal/matrix"
	"go-matrix-service/internal/protocol"
	"log"
	"net"
	"time"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server listening on port 8080...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	dec := json.NewDecoder(conn)
	enc := json.NewEncoder(conn)

	var req protocol.Request
	if err := dec.Decode(&req); err != nil {
		return
	}

	fmt.Printf("[%s] Job started: %dx%d (%d workers)\n", conn.RemoteAddr(), req.MatrixSize, req.MatrixSize, req.Workers)

	matA := matrix.NewRandom(req.MatrixSize, req.Seed)
	matB := matrix.NewRandom(req.MatrixSize, req.Seed+1)

	start := time.Now()
	matA.MultiplySequential(matB)
	seqTime := time.Since(start).Seconds()

	progressChan := make(chan int)
	go func() {
		total := 0
		for n := range progressChan {
			total += n
			enc.Encode(protocol.Response{
				Type:          "progress",
				RowsProcessed: total,
				TotalRows:     req.MatrixSize,
			})
		}
	}()

	start = time.Now()
	resConc := matA.MultiplyConcurrent(matB, req.Workers, req.ChunkSize, progressChan)
	concTime := time.Since(start).Seconds()

	enc.Encode(protocol.Response{
		Type:     "result",
		RowsProcessed: req.MatrixSize, 
		TotalRows: req.MatrixSize,
		SeqTime:  seqTime,
		ConcTime: concTime,
		Speedup:  seqTime / concTime,
		Checksum: resConc.Checksum(),
	})

	fmt.Printf("[%s] Job finished\n", conn.RemoteAddr())
}	
