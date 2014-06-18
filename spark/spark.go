package spark

import (
	"bytes"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/eapache/queue"
	"github.com/olekukonko/ts"
	"io"
	"log"
	"math"
	"os"
	"sync"
	"time"
)

const debug = false

const (
	// Bytes is a predefined unit type, will pretty print using humanized values.
	Bytes = "bytes"
)

// Spark creates a stream of sparklines that will print every lines bucketize
// with the given resolution.
//
// By default, it prints to os.Stdout and is unitless.
func Spark(resolution time.Duration) *SparkStream {
	return &SparkStream{
		Out:   os.Stdout,
		buf:   bytes.NewBuffer(nil),
		queue: queue.New(),
		res:   resolution,
		tick:  time.NewTicker(resolution),
	}
}

// SparkStream prints sparklines for values it receives in real time. It
// scales up and down the visible window of values and prints the average
// values per second it has observed in recent history.
type SparkStream struct {
	Units string
	Out   io.Writer

	l sync.Mutex

	maxwidth int
	curwidth int

	buf *bytes.Buffer

	max   float64
	min   float64
	avg   float64
	cur   float64
	queue *queue.Queue

	res  time.Duration
	tick *time.Ticker
}

// Start starts the printing of sparklines.
func (s *SparkStream) Start() {
	go func() {
		for _ = range s.tick.C {
			s.printLines()
		}
	}()
}

// Stop stops the printing of sparklines.
func (s *SparkStream) Stop() { s.tick.Stop() }

// Add puts the value in the current bucket of sparklines. The value
// will appear part of the next update of the spark stream.
func (s *SparkStream) Add(v float64) {
	s.l.Lock()
	s.cur += v
	s.l.Unlock()
}

func (s *SparkStream) printLines() {
	size, err := ts.GetSize()
	if err != nil {
		log.Printf("SparkStream: get size error, %v", err)
		return
	}

	// ensure queue is as large as largest terminal seen so far
	s.curwidth = size.Col()
	if s.curwidth > s.maxwidth {
		s.maxwidth = s.curwidth
	}
	if s.queue.Length() == s.maxwidth {
		s.queue.Remove()
	}

	s.queue.Add(s.cur)
	s.cur = 0.0

	var sum float64
	for i := 0; i < s.queue.Length(); i++ {
		val := s.queue.Get(i).(float64)
		if i == 0 {
			s.max = val
			s.min = val
		} else {
			s.max = math.Max(s.max, val)
			s.min = math.Min(s.min, val)
		}
		sum += val
	}
	s.avg = sum / (s.res.Seconds() * float64(s.queue.Length()))

	var avgStr string
	switch s.Units {
	case Bytes:
		avgStr = " " + humanize.Bytes(uint64(s.avg)) + "/s"
	case "":
		if s.avg > 1000.0 {
			avgStr = " " + humanize.Comma(int64(s.avg)) + "/s"
		} else {
			avgStr = " " + fmt.Sprintf("%.3f", s.avg) + "/s"
		}
	default:
		if s.avg > 1000.0 {
			avgStr = " " + humanize.Comma(int64(s.avg)) + string(s.Units) + "/s"
		} else {
			avgStr = " " + fmt.Sprintf("%.3f", s.avg) + string(s.Units) + "/s"
		}
	}

	// visibleIdx is the index in the queue where the last value
	// visible on the terminal is found. before that, the queue
	// still tracks the data but we should not print it
	visibleIdx := imax(0, s.queue.Length()-s.curwidth+len(avgStr))

	if debug {
		log.Printf("visibleIdx=%d\tlen(avgStr)=%d\ts.queue.Length()=%d",
			visibleIdx, len(avgStr), s.queue.Length())
	}

	s.buf.Reset()
	_, _ = s.buf.WriteRune('\r')
	var runec int
	for i := visibleIdx; i < s.queue.Length(); i++ {
		val := s.queue.Get(i).(float64)
		runec++
		_, _ = s.buf.WriteRune(blockIdx(val, s.min, s.max))
	}
	for runec < s.curwidth-len(avgStr) {
		runec++
		_, _ = s.buf.WriteRune(' ')
	}
	_, _ = s.buf.WriteString(avgStr)
	_, err = s.buf.WriteTo(s.Out)
	if err != nil {
		log.Printf("SparkStream: writing to output: %v", err)
	}
}

var blocks = [...]rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

func blockIdx(val, min, max float64) rune {
	if val >= max {
		return blocks[len(blocks)-1]
	}

	width := (max - min)
	if width <= math.SmallestNonzeroFloat64 {
		return blocks[0]
	}

	parts := (val - min) / width
	fIdx := parts * float64(len(blocks))
	i := imin(int(fIdx), len(blocks)-1)

	if debug {
		log.Printf("val=%g\tmin=%g\tmax=%g\twidth=%g\tparts=%g\tfIdx=%g\ti=%d", val, min, max, width, parts, fIdx, i)
	}

	return blocks[i]
}

func imax(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func imin(a, b int) int {
	if a < b {
		return a
	}
	return b
}
