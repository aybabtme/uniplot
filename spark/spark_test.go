package spark_test

import (
	"github.com/aybabtme/uniplot/spark"
	"log"
	"math/rand"
	"os"
	"time"
)

func ExampleSpark() {
	sprk := spark.Spark(time.Millisecond * 60)
	sprk.Out = os.Stderr
	sprk.Units = spark.Bytes
	sprk.Start()

	timeout := time.NewTicker(time.Second * 10)

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
