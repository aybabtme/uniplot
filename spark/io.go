package spark

import (
	"io"
	"os"
	"time"
)

type reader func(b []byte) (int, error)

func (r reader) Read(b []byte) (int, error) { return r(b) }

type writer func(b []byte) (int, error)

func (w writer) Write(b []byte) (int, error) { return w(b) }

func Reader(r io.Reader) io.Reader {
	sprk := Spark(time.Millisecond * 30)
	sprk.Units = Bytes
	started := false
	return reader(func(b []byte) (int, error) {
		if !started {
			sprk.Start()
			started = true
		}
		n, err := r.Read(b)
		if err == nil {
			sprk.Add(float64(n))
		} else {
			sprk.Stop()
		}
		return n, err
	})
}

func Writer(w io.Writer) io.Writer {

	if f, ok := w.(*os.File); ok && f == os.Stdout {
		panic("cannot use Writer with os.Stdout")
	}

	sprk := Spark(time.Millisecond * 30)
	sprk.Units = Bytes
	started := false
	return writer(func(b []byte) (int, error) {
		if !started {
			sprk.Start()
			started = true
		}
		n, err := w.Write(b)
		if err == nil {
			sprk.Add(float64(n))
		} else {
			sprk.Stop()
		}
		return n, err
	})
}
