package histogram

import (
	"math"
)

// Histogram holds a count of values partionned over buckets.
type Histogram struct {
	// Min is the size of the smallest bucket.
	Min int
	// Max is the size of the biggest bucket.
	Max int
	// Count is the total size of all buckets.
	Count int
	// Buckets over which values are partionned.
	Buckets []Bucket
}

// Bucket counts a partion of values.
type Bucket struct {
	// Count is the number of values represented in the bucket.
	Count int
	// Min is the low, inclusive bound of the bucket.
	Min float64
	// Max is the high, exclusive bound of the bucket. If
	// this bucket is the last bucket, the bound is inclusive
	// and contains the max value of the histogram.
	Max float64
}

// Hist creates an histogram partionning input over `bins` buckets.
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
	for _, val := range input {
		minx := float64(min)
		xdiff := val - minx
		bi := imin(int(xdiff/scale), len(buckets)-1)

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

// Scale gives the scaled count of the bucket at idx, using the
// provided scale func.
func (h Histogram) Scale(s ScaleFunc, idx int) float64 {
	bkt := h.Buckets[idx]
	scale := s(h.Min, h.Max, bkt.Count)
	return scale
}

// ScaleFunc is the type to implement to scale an histogram.
type ScaleFunc func(min, max, value int) float64

// Linear builds a ScaleFunc that will linearly scale the values of
// an histogram so that they do not exceed width.
func Linear(width int) ScaleFunc {
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
