package histogram

import (
	"math"
)

type Histogram struct {
	Min     int
	Max     int
	Count   int
	Buckets []Bucket
}

type Bucket struct {
	Count int
	Min   float64
	Max   float64
}

func Hist(bins int, input []float64) Histogram {
	if len(input) == 0 || bins == 0 {
		return Histogram{}
	}
	buckets := make([]Bucket, bins)

	min, max := input[0], input[0]
	for _, val := range input {
		min = math.Min(min, val)
		max = math.Max(max, val)
	}

	scale := (max - min) / float64(bins)

	for i := range buckets {
		bmin, bmax := float64(i+1)*scale, float64(i+2)*scale
		buckets[i] = Bucket{Min: bmin, Max: bmax}
	}

	minC, maxC := 0, 0
	// next_value:
	for _, val := range input {
		bi := imin(int(val/scale), len(buckets)) - 1
		buckets[bi].Count++
		minC = imin(minC, buckets[bi].Count)
		maxC = imax(maxC, buckets[bi].Count)
	}

	return Histogram{
		Min:     minC,
		Max:     maxC,
		Count:   len(input),
		Buckets: buckets,
	}
}

func (h Histogram) Scale(s Scale, idx int) float64 {
	bkt := h.Buckets[idx]
	scale := s(h.Min, h.Max, bkt.Count)
	return scale
}

type Scale func(min, max, value int) float64

func Linear(width int) Scale {
	return func(min, max, value int) float64 {
		return float64(value-min) / float64(max-min) * float64(width)
	}
}

func imin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func imax(a, b int) int {
	if a > b {
		return a
	}
	return b
}
