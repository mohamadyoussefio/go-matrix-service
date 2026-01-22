package matrix

import (
	"math/rand"
)

type Matrix struct {
	Size int
	Data []float64
}

func NewRandom(size int, seed int64) *Matrix {
	src := rand.NewSource(seed)
	rnd := rand.New(src)
	data := make([]float64, size*size)

	for i := range data {
		data[i] = rnd.Float64()
	}

	return &Matrix{
		Size: size,
		Data: data,
	}
}

func (m *Matrix) Checksum() float64 {
	sum := 0.0
	for _, v := range m.Data {
		sum += v
	}
	return sum
}
