package spark_test

import (
	"github.com/aybabtme/uniplot/spark"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
)

func ExampleSpark_random() {
	sprk := spark.Spark(time.Millisecond * 30)
	sprk.Out = os.Stderr
	sprk.Units = spark.Bytes
	sprk.Start()

	timeout := time.NewTicker(time.Second * 5)

loop:
	for {
		select {
		case <-time.After(time.Millisecond * 20):
			sprk.Add(float64(rand.Intn(10000000000)))
		case <-timeout.C:
			log.Printf("DONE!")
			break loop
		}
	}
	sprk.Stop()

	// Output:
	//
}

func ExampleSpark_sine() {
	sprk := spark.Spark(time.Millisecond * 60)
	sprk.Out = os.Stderr
	sprk.Units = spark.Bytes
	sprk.Start()

	timeout := time.NewTicker(time.Second * 15)

	x := 0.0
loop:
	for {
		select {
		case <-time.After(time.Millisecond * 60):
			y := math.Sin(x)
			x += 0.1
			sprk.Add(y)
		case <-timeout.C:
			log.Printf("DONE!")
			break loop
		}
	}
	sprk.Stop()

	// Output:
	//
}

func ExampleSpark_saw() {
	sprk := spark.Spark(time.Millisecond * 60)
	sprk.Out = os.Stderr
	sprk.Units = spark.Bytes
	sprk.Start()

	timeout := time.NewTicker(time.Second * 15)

	x := 0.0
loop:
	for {
		select {
		case <-time.After(time.Millisecond * 60):
			x += 0.1
			sprk.Add(x)
			if x >= 1.0 {
				x = 0.0
			}
		case <-timeout.C:
			log.Printf("DONE!")
			break loop
		}
	}
	sprk.Stop()

	// Output:
	//
}
