package matrix

import (
	"sync"
)

func (a *Matrix) MultiplySequential(b *Matrix) *Matrix {
	result := &Matrix{
		Size: a.Size,
		Data: make([]float64, a.Size*a.Size),
	}

	for i := 0; i < a.Size; i++ {
		for k := 0; k < a.Size; k++ {
			temp := a.Data[i*a.Size+k]
			for j := 0; j < a.Size; j++ {
				result.Data[i*a.Size+j] += temp * b.Data[k*a.Size+j]
			}
		}
	}
	return result
}

type Job struct {
	Start int
	End   int
}

func (a *Matrix) MultiplyConcurrent(b *Matrix, workers, chunkSize int, progress chan int) *Matrix {
	result := &Matrix{
		Size: a.Size,
		Data: make([]float64, a.Size*a.Size),
	}

	jobs := make(chan Job, workers)
	var wg sync.WaitGroup

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				for i := job.Start; i < job.End; i++ {
					for k := 0; k < a.Size; k++ {
						temp := a.Data[i*a.Size+k]
						for j := 0; j < a.Size; j++ {
							result.Data[i*a.Size+j] += temp * b.Data[k*a.Size+j]
						}
					}
				}
				if progress != nil {
					progress <- (job.End - job.Start)
				}
			}
		}()
	}

	go func() {
		for i := 0; i < a.Size; i += chunkSize {
			end := i + chunkSize
			if end > a.Size {
				end = a.Size
			}
			jobs <- Job{Start: i, End: end}
		}
		close(jobs)
	}()

	wg.Wait()
	if progress != nil {
		close(progress)
	}

	return result
}
