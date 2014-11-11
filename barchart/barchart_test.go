package barchart

import (
	"github.com/dustin/go-humanize"
	"math/rand"
	"os"
	"time"
)

func ExampleBarChartXYs() {
	data := [][2]int{
		{0, 1},
		{1, 3},
		{2, 4},
		{3, 6},
		{4, 8},
		// nil,
		// nil,
		{7, 15},
		{8, 10},
		{9, 7},
		{10, 5},
		{11, 3},
		{12, 2},
		{13, 1},
		{14, 0},
		{15, 20},
	}

	plot := BarChartXYs(data)

	if err := Fprint(os.Stdout, plot, Linear(19)); err != nil {
		panic(err)
	}
	// Output:
	// 0   █ 1
	// 1   ██▉ 3
	// 2   ███▉ 4
	// 3   █████▊ 6
	// 4   ███████▋ 8
	// 5   nil
	// 6   nil
	// 7   ██████████████▎ 15
	// 8   █████████▋ 10
	// 9   ██████▋ 7
	// 10  ████▊ 5
	// 11  ██▉ 3
	// 12  ██ 2
	// 13  █ 1
	// 14  ▏ 0
	// 15  ███████████████████▏ 20
}

func ExampleBarChartXYs_time() {
	base := time.Date(1984, 01, 01, 00, 00, 00, 00, time.UTC)
	baseplus := func(i int) int {
		return int(base.Add(time.Duration(i) * time.Millisecond).UnixNano())
	}
	mb := 1000000

	r := rand.New(rand.NewSource(42))

	data := [][2]int{
		{baseplus(0), r.Intn(20) * mb},
		{baseplus(100), r.Intn(20) * mb},
		{baseplus(200), r.Intn(20) * mb},
		{baseplus(300), r.Intn(20) * mb},
		{baseplus(400), r.Intn(20) * mb},
		// nil,
		// nil,
		{baseplus(700), r.Intn(20) * mb},
		{baseplus(800), r.Intn(20) * mb},
		{baseplus(900), r.Intn(20) * mb},
		{baseplus(1000), r.Intn(20) * mb},
		{baseplus(1100), r.Intn(20) * mb},
		{baseplus(1200), r.Intn(20) * mb},
		{baseplus(1300), r.Intn(20) * mb},
		{baseplus(1400), r.Intn(20) * mb},
		{baseplus(1500), r.Intn(20) * mb},
	}

	xfmt := func(x float64) string {
		return time.Unix(base.Unix(), int64(x)).Format("15:04:05.000")
	}
	yfmt := func(y float64) string {
		return humanize.Bytes(uint64(y))
	}

	plot := BarChartXYs(data)

	if err := Fprintf(os.Stdout, plot, len(data), Linear(13), xfmt, yfmt); err != nil {
		panic(err)
	}
	// 19:00:00.000  ██████▏ 12MB
	// 19:00:00.115  ███▌ 8.0MB
	// 19:00:00.230  ████▊ 10MB
	// 19:00:00.346  ▏ 3.0MB
	// 19:00:00.461  nil
	// 19:00:00.576  nil
	// 19:00:00.692  █████████████▏ 22MB
	// 19:00:00.807  ████████▉ 16MB
	// 19:00:00.923  ███▌ 8.0MB
	// 19:00:01.038  ▏ 3.0MB
	// 19:00:01.153  ████▏ 9.0MB
	// 19:00:01.269  ██▊ 7.0MB
	// 19:00:01.384  ████████▎ 15MB
	// 19:00:01.500  ███▌ 8.0MB
}
