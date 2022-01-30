<h1 align="center">GoIntervalTree</h1>

<p align="center">
  <a href="/LICENSE"><img src="https://img.shields.io/badge/license-MIT-black"/></a>
  <a href="https://app.travis-ci.com/github/danilovkiri/gointervaltree"><img src="https://app.travis-ci.com/danilovkiri/gointervaltree.svg?branch=main"/></a>
  <a href="https://app.codecov.io/gh/danilovkiri/gointervaltree"><img src="https://codecov.io/gh/danilovkiri/gointervaltree/branch/main/graph/badge.svg" /></a>
  <a href="https://goreportcard.com/report/github.com/danilovkiri/gointervaltree"><img src="https://goreportcard.com/badge/github.com/danilovkiri/gointervaltree"/></a>
  <a href="https://pkg.go.dev/github.com/danilovkiri/gointervaltree"><img src="https://godoc.org/github.com/danilovkiri/gointervaltree?status.svg"/></a>
  <a href="https://app.fossa.com/projects/git%2Bgithub.com%2Fdanilovkiri%2Fgointervaltree?ref=badge_shield" alt="FOSSA Status"><img src="https://app.fossa.com/api/projects/git%2Bgithub.com%2Fdanilovkiri%2Fgointervaltree.svg?type=shield"/></a>

</p>

<p align="center">
  An IntervalTree package for Go
</p>

<br>

> Inspired by Centered Interval Tree
> [implementation](https://github.com/konstantint/pyliftover/blob/master/pyliftover/intervaltree.py) in Python 

This package provides functionality for indexing a set of integer intervals, e.g. [start, end) based on
[Wikipedia reference](http://en.wikipedia.org/wiki/Interval_tree). No interval removal is implemented.

## License

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fdanilovkiri%2Fgointervaltree.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fdanilovkiri%2Fgointervaltree?ref=badge_large)

## Installation
```shell
go get github.com/danilovkiri/gointervaltree
```

## Usage

```go
package main
import (
	"fmt"
	"github.com/danilovkiri/gointervaltree"
)

func main() {
	t := gointervaltree.NewIntervalTree(0, 100)
	t.AddInterval(1, 10, []string{"a", "b"})
	t.AddInterval(20, 30, []bool{true, false})
	t.AddInterval(32, 35, []int{1, 2, 3})
	t.AddInterval(32, 38, nil)
	t.Sort()
	fmt.Println(t.Len())
	// 4
	fmt.Println(t.Iter())
	// [[1 10 [a b]] [32 35 [1 2 3]] [32 38 <nil>] [20 30 [true false]]]
	fmt.Println(t.Query(33))
	// [[32 35 [1 2 3]] [32 38 <nil>]]
}
```
