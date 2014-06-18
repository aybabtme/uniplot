# spark

Makes sparklines... animated sparklines!

![example-sparklines](https://cloud.githubusercontent.com/assets/1189716/3317255/09998004-f70b-11e3-8aab-597bd848c467.gif)


# Usage

Create a `SparkStream` object:

```go
sprk := spark.Spark(time.Millisecond * 60) // prints every 60ms
sprk.Start()

timeout := time.NewTicker(time.Second * 10)
for {
    select {
    case <-time.After(time.Millisecond * 20):
        val := float64(rand.Intn(10000000000))
        sprk.Add(val)

    case <-timeout.C:
        sprk.Stop() // stops printing
        return
    }
}
```

That's it!  But you can also configure a few things:

```go
sprk.Out = os.Stderr
```

Or change the units it prints

```go
sprk.Units = spark.Bytes // print as {KB/MB/GB/TB}/s, i.e `100MB/s`
sprk.Units = "unicorns"  // will print as `unicorn/s`
```

# Docs?

[Godocs](http://godoc.org/github.com/aybabtme/uniplot/spark)!

# License

MIT license.
