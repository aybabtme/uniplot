package histogram

import (
	"fmt"
	"io"
	"math"
	"strings"
	"text/tabwriter"
)

var blocks = []string{
	"▏", "▎", "▍", "▌", "▋", "▊", "▉", "█",
}

var barchar = func(v float64) string {
	digit := (v - math.Floor(v)) * 10.0
	diff := int(10.0 * digit / 8.0)
	return blocks[diff]
}

func FPrint(w io.Writer, h Histogram, s Scale) error {
	tabw := tabwriter.NewWriter(w, 2, 2, 2, byte(' '), 0)
	for i, bkt := range h.Buckets {
		sz := h.Scale(s, i)
		fmt.Fprintf(tabw, "%.4g-%.4g\t%.3g%%\t%s\t[%d/%d]\n",
			bkt.Min, bkt.Max,
			float64(bkt.Count)*100.0/float64(h.Count),
			strings.Repeat("█", int(sz))+barchar(sz),
			bkt.Count,
			h.Count,
		)
	}

	return tabw.Flush()
}
