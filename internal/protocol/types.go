package protocol

type Request struct {
	MatrixSize int   `json:"matrix_size"`
	Workers    int   `json:"workers"`
	ChunkSize  int   `json:"chunk_size"`
	Seed       int64 `json:"seed"`
}

type Response struct {
	Type          string  `json:"type"`
	RowsProcessed int     `json:"rows_processed"`
	TotalRows     int     `json:"total_rows"`
	SeqTime       float64 `json:"seq_time"`
	ConcTime      float64 `json:"conc_time"`
	Speedup       float64 `json:"speedup"`
	Checksum      float64 `json:"checksum"`
}
