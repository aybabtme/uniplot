package barchart

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

var barstring = func(v float64) string {
	decimalf := (v - math.Floor(v)) * 10.0
	decimali := math.Floor(decimalf)
	charIdx := int(decimali / 10.0 * 8.0)
	return strings.Repeat("█", int(v)) + blocks[charIdx]
}

// Fprint plots p as a Unicode XY plot, using scale s. The results
// is something like:
//
//    0   █ 1
//    1   ██▉ 3
//    2   ███▉ 4
//    3   █████▊ 6
//    4   ███████▋ 8
//    5   nil
//    6   nil
//    7   ██████████████▎ 15
//    8   █████████▋ 10
//    9   ██████▋ 7
//    10  ████▊ 5
//    11  ██▉ 3
//    12  ██ 2
//    13  █ 1
//    14  ▏ 0
//    15  ███████████████████▏ 20
func Fprint(w io.Writer, p BarChart, s ScaleFunc) error {
	fmtFunc := func(v float64) string { return fmt.Sprintf("%v", v) }
	width := p.MaxX - p.MinX + 1
	return fprintf(w, p, width, s, fmtFunc, fmtFunc)
}

// Fprintf plots p as a Unicode XY plot of width, scaling Y values with
// s and rendering axis labels with x and y format funcs.  As in the
// example, you can render pretty things like:
//    12:06:00.098  █████▋ 12MB
//    12:06:00.213  ████▋ 8.0MB
//    12:06:00.329  ██████▋ 10MB
//    12:06:00.444  ▏ 3.0MB
//    12:06:00.559  nil
//    12:06:00.675  nil
//    12:06:00.790  ██████████████▉ 22MB
//    12:06:00.906  ████████████▏ 16MB
//    12:06:01.021  ████▋ 8.0MB
//    12:06:01.136  ▏ 3.0MB
//    12:06:01.252  █████▋ 9.0MB
//    12:06:01.367  ███▊ 7.0MB
//    12:06:01.483  ███████████▏ 15MB
//    12:06:01.598  ████▋ 8.0MB
func Fprintf(w io.Writer, p BarChart, width int, s ScaleFunc, x, y FormatFunc) error {
	return fprintf(w, p, width, s, x, y)
}

func fprintf(w io.Writer, p BarChart, width int, s ScaleFunc, xfmt, yfmt FormatFunc) error {
	tabw := tabwriter.NewWriter(w, 2, 2, 2, byte(' '), 0)
	for _, xy := range p.ScaleXYs(width, s) {
		xstr := xfmt(xy.X)
		var ystr string
		var bar string
		if xy.Y == nil {
			ystr = ""
			bar = "nil"
		} else {
			bar = barstring(*xy.ScaledY)
			ystr = " " + yfmt(*xy.Y)
		}

		fmt.Fprintf(tabw, "%s\t%s\n",
			xstr,
			bar+ystr,
		)
	}

	return tabw.Flush()
}
