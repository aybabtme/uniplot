# barchart

Makes barcharts.

# Usage

Pass `[][2]int` values to `BarChartXYs`:

```go
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
```

You can use the `Fprint` utility to create this Unicode graph (the graph looks better with fonts that
draw unicode blocks properly):

```go
if err := Fprint(os.Stdout, plot, Linear(19)); err != nil {
    panic(err)
}
```

Yields :

```
0   █ 1
1   ██▉ 3
2   ███▉ 4
3   █████▊ 6
4   ███████▋ 8
5   nil
6   nil
7   ██████████████▎ 15
8   █████████▋ 10
9   ██████▋ 7
10  ████▊ 5
11  ██▉ 3
12  ██ 2
13  █ 1
14  ▏ 0
15  ███████████████████▏ 20
```

If your console font is Monaco, one of the blocks look weird. Use Menlo. =)

# Docs?

[Godocs](http://godoc.org/github.com/aybabtme/uniplot/barchart)!

# License

MIT license.
