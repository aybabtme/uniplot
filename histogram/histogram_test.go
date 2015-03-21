package histogram

import (
	"os"
	"time"
)

func ExampleHist() {
	data := []float64{
		0.1,
		0.2, 0.21, 0.22, 0.22,
		0.3,
		0.4,
		0.5, 0.51, 0.52, 0.53, 0.54, 0.55, 0.56, 0.57, 0.58,
		0.6,
		// 0.7 is empty
		0.8,
		0.9,
		1.0,
	}

	hist := Hist(9, data)

	if err := Fprint(os.Stdout, hist, Linear(5)); err != nil {
		panic(err)
	}
	// 0.1-0.2  5%   ▋       1
	// 0.2-0.3  25%  ██▊     5
	// 0.3-0.4  0%   ▏
	// 0.4-0.5  5%   ▋       1
	// 0.5-0.6  45%  █████▏  9
	// 0.6-0.7  5%   ▋       1
	// 0.7-0.8  0%   ▏
	// 0.8-0.9  5%   ▋       1
	// 0.9-1    10%  █▏      2
}

func ExampleHist_duration() {
	data := []float64{
		float64(time.Millisecond * 100),
		float64(time.Millisecond * 200),
		float64(time.Millisecond * 210),
		float64(time.Millisecond * 220),
		float64(time.Millisecond * 221),
		float64(time.Millisecond * 222),
		float64(time.Millisecond * 223),
		float64(time.Millisecond * 300),
		float64(time.Millisecond * 400),
		float64(time.Millisecond * 500),
		float64(time.Millisecond * 510),
		float64(time.Millisecond * 520),
		float64(time.Millisecond * 530),
		float64(time.Millisecond * 540),
		float64(time.Millisecond * 550),
		float64(time.Millisecond * 560),
		float64(time.Millisecond * 570),
		float64(time.Millisecond * 580),
		float64(time.Millisecond * 600),
		// no 0.7s
		float64(time.Millisecond * 800),
		float64(time.Millisecond * 900),
		float64(time.Millisecond * 1000),
	}

	hist := Hist(9, data)

	err := Fprintf(os.Stdout, hist, Linear(5), func(v float64) string {
		return time.Duration(v).String()
	})
	if err != nil {
		panic(err)
	}
	// 100ms-200ms  4.55%  ▋       1
	// 200ms-300ms  27.3%  ███▍    6
	// 300ms-400ms  4.55%  ▋       1
	// 400ms-500ms  4.55%  ▋       1
	// 500ms-600ms  40.9%  █████▏  9
	// 600ms-700ms  4.55%  ▋       1
	// 700ms-800ms  0%     ▏
	// 800ms-900ms  4.55%  ▋       1
	// 900ms-1s     9.09%  █▏      2
}
