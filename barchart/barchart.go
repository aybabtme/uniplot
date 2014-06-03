package barchat

import (
	"math"
)

// XY point in a 2D plot.
type XY struct{ X, Y int }

// XYf point in a 2D plot.
type XYf struct {
	X          float64
	Y, ScaledY *float64
}

// BarChart of XY points.
type BarChart struct {
	MinX, MaxX, MinY, MaxY int
	xy                     []XY
}

// Add a XY point to the plot.
func (p *BarChart) Add(x, y int) { p.xy = append(p.xy, XY{x, y}) }

// XYs aggregates the XY values together in a dense form.
// Empty slots between two Xs are left nil to represent
// the absence of data.
func (p *BarChart) XYs() []*XY {
	xys := make([]*XY, p.MaxX-p.MinX)
	for _, xy := range p.xy {
		slot := xys[xy.X-p.MinX]
		if slot == nil {
			slot = &XY{xy.X, xy.Y}
		} else {
			slot.Y += xy.Y
		}
	}
	return xys
}

// ScaleXYs aggregates the XY values together in a dense form.
// Empty slots between two Xs are left nil to represent
// the absence of data. The values are scaled using s.
func (p *BarChart) ScaleXYs(xWidth int, s ScaleFunc) []XYf {

	diff := p.MaxX - p.MinX
	scaleX := float64(diff) / float64(xWidth-1)

	buckets := make([]XYf, xWidth)
	for i := range buckets {
		buckets[i] = XYf{
			X:       float64(i)*scaleX + float64(p.MinX),
			Y:       nil,
			ScaledY: nil,
		}
	}

	miny, maxy := float64(p.xy[0].Y), float64(p.xy[0].Y)

	for _, val := range p.xy {
		minx := float64(p.MinX)
		xdiff := float64(val.X) - minx
		bi := int(xdiff / scaleX)

		slot := buckets[bi]
		if slot.Y == nil {
			slot.Y = new(float64)
			slot.ScaledY = new(float64)
		}
		y := float64(val.Y)
		*slot.Y += y
		miny = math.Min(*slot.Y, miny)
		maxy = math.Max(*slot.Y, maxy)
		buckets[bi] = slot
	}

	for _, val := range buckets {
		if val.Y == nil {
			continue
		}
		*val.ScaledY = s(miny, maxy, *val.Y)
	}

	return buckets
}

// BarChartXYs builds a BarChart using pairwise X and Y []float64.
func BarChartXYs(xys [][2]int) BarChart {
	if len(xys) == 0 {
		return BarChart{}
	}
	minx, maxx, miny, maxy := xys[0][0], xys[0][0], xys[0][1], xys[0][1]
	plot := BarChart{
		xy: make([]XY, len(xys)),
	}
	var x, y int
	for i := range plot.xy {
		x = xys[i][0]
		y = xys[i][1]
		minx = imin(x, minx)
		maxx = imax(x, maxx)
		miny = imin(y, miny)
		maxy = imax(y, maxy)
		plot.xy[i] = XY{x, y}
	}
	plot.MinX = minx
	plot.MaxX = maxx
	plot.MinY = miny
	plot.MaxY = maxy
	return plot
}

// ScaleFunc is the type to implement to scale an histogram.
type ScaleFunc func(min, max, value float64) float64

// Linear builds a ScaleFunc that will linearly scale the values of
// an histogram so that they do not exceed width.
func Linear(width int) ScaleFunc {
	return func(min, max, value float64) float64 {
		diff := max - min
		offset := value - min
		ratio := offset / diff
		return ratio * float64(width)
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
