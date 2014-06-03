package histogram

import (
	"fmt"
	"io"
	"math"
	"strings"
	"text/tabwriter"
)

// FormatFunc formats a float into the proper string form. Used to
// print meaningful axe labels.
type FormatFunc func(v float64) string

var blocks = []string{
	"▏", "▎", "▍", "▌", "▋", "▊", "▉", "█",
}

var barchar = func(v float64) string {
	digit := (v - math.Floor(v)) * 10.0
	diff := int(10.0 * digit / 8.0)
	return blocks[diff]
}

// FPrint prints a unicode histogram on the io.Writer, using
// scale s. This code:
//
// 	hist := Hist(9, data)
// 	err := FPrint(os.Stdout, hist, Linear(5))
//
// ... yields the graph:
//
//	0.1-0.2  5%   ▉       [1/20]
//	0.2-0.3  25%  ██▉     [5/20]
//	0.3-0.4  0%   ▏       [0/20]
//	0.4-0.5  5%   ▉       [1/20]
//	0.5-0.6  50%  █████▏  [10/20]
//	0.6-0.7  0%   ▏       [0/20]
//	0.7-0.8  0%   ▏       [0/20]
//	0.8-0.9  5%   ▉       [1/20]
//	0.9-1    10%  █▏      [2/20]
func FPrint(w io.Writer, h Histogram, s ScaleFunc) error {
	return fprintf(w, h, s, func(v float64) string {
		return fmt.Sprintf("%.4g", v)
	})
}

// FPrintf is the same as FPrint, but applies f to the axis labels.
func FPrintf(w io.Writer, h Histogram, s ScaleFunc, f FormatFunc) error {
	return fprintf(w, h, s, f)
}

func fprintf(w io.Writer, h Histogram, s ScaleFunc, f FormatFunc) error {
	tabw := tabwriter.NewWriter(w, 2, 2, 2, byte(' '), 0)
	for i, bkt := range h.Buckets {
		sz := h.Scale(s, i)
		fmt.Fprintf(tabw, "%s-%s\t%.3g%%\t%s\t[%d/%d]\n",
			f(bkt.Min), f(bkt.Max),
			float64(bkt.Count)*100.0/float64(h.Count),
			strings.Repeat("█", int(sz))+barchar(sz),
			bkt.Count,
			h.Count,
		)
	}

	return tabw.Flush()
}
