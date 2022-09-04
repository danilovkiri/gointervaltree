package gointervaltree

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sort"
	"testing"
)

// Tests

func TestNewIntervalTreeFailedBoundaries(t *testing.T) {
	_, err := NewIntervalTree(30, 25)
	assert.EqualError(t, err, "interval tree start must be numerically less than its end")
}

func TestIntervalTree_AddIntervalBad(t *testing.T) {
	tree, _ := NewIntervalTree(20, 25)
	_ = tree.AddInterval(23, 22, nil)
	assert.Equal(t, 0, tree.Len())
}

func TestIntervalTree_QueryEmptyTree(t *testing.T) {
	tree, _ := NewIntervalTree(10, 50)
	assert.Equal(t, []resultInterval[int](nil), tree.Query(1))
}

func TestIntervalTree_LenEmptyTree(t *testing.T) {
	tree, _ := NewIntervalTree(10, 50)
	assert.Equal(t, 0, tree.Len())
}

func TestIntervalTree_IterEmptyTree(t *testing.T) {
	tree, _ := NewIntervalTree(10, 50)
	assert.Equal(t, []resultInterval[int](nil), tree.Iter())
}

func TestNewIntervalTree(t *testing.T) {
	constants := struct {
		treeMin     int
		treeMax     int
		intervals   [][]int
		queryPoints []int
	}{
		treeMin:     0,
		treeMax:     100,
		intervals:   [][]int{{10, 20}, {20, 30}, {21, 31}, {30, 40}, {45, 55}, {45, 56}, {46, 57}, {55, 56}, {58, 59}, {50, 51}},
		queryPoints: []int{-1, 0, 1, 10, 11, 19, 20, 21, 24, 25, 26, 30, 40, 41, 48, 49, 50, 51, 52, 60, 74, 75, 76, 90, 100, 1000},
	}
	doTest(t, constants.treeMin, constants.treeMax, constants.intervals, constants.queryPoints)

}

func doTest(t *testing.T, min, max int, intervals [][]int, queryPoints []int) {
	tree, _ := NewIntervalTree(min, max)
	for _, interval := range intervals {
		_ = tree.AddInterval(interval[0], interval[1], nil)
	}
	tree.Sort()
	for _, q := range queryPoints {
		r := tree.Query(q)
		sort.Slice(r, func(i, j int) bool {
			if r[i].start != r[j].start {
				return r[i].start > r[j].start
			}
			return r[i].end > r[j].end
		})
		var trueR []resultInterval[int]
		for _, interval := range intervals {
			if (interval[0] <= q) && (q < interval[1]) {
				trueR = append(trueR, resultInterval[int]{interval[0], interval[1], nil})
			}
		}
		sort.Slice(trueR, func(i, j int) bool {
			if trueR[i].start != trueR[j].start {
				return trueR[i].start > trueR[j].start
			}
			return trueR[i].end > trueR[j].end
		})
		assert.Equal(t, trueR, r)
	}
	observedLength := tree.Len()
	expectedLength := 0
	for _, interval := range intervals {
		if interval[0] < interval[1] {
			expectedLength += 1
		}
	}
	assert.Equal(t, expectedLength, observedLength)
	assert.Equal(t, expectedLength, len(tree.Iter()))
}

// Benchmarks

func BenchmarkIntervalTree_Query(b *testing.B) {
	var (
		treeMin     = 0
		treeMax     = 100
		intervals   = [][]int{{10, 20}, {20, 30}, {21, 31}, {30, 40}, {45, 55}, {45, 56}, {46, 57}, {55, 56}, {58, 59}, {50, 51}}
		queryPoints = []int{-1, 0, 1, 10, 11, 19, 20, 21, 24, 25, 26, 30, 40, 41, 48, 49, 50, 51, 52, 60, 74, 75, 76, 90, 100, 1000}
	)
	tree, _ := NewIntervalTree(treeMin, treeMax)
	for _, interval := range intervals {
		_ = tree.AddInterval(interval[0], interval[1], nil)
	}
	tree.Sort()
	b.ResetTimer()
	b.Run("benchmark-tree-query", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer() // run with -benchtime=1000x, otherwise becnhmarking stalls
			q := queryPoints[rand.Intn(len(queryPoints))]
			b.StartTimer()
			_ = tree.Query(q)
		}
	})
}

// Examples

func ExampleNewIntervalTree() {
	t, _ := NewIntervalTree(0, 100)
	_ = t.AddInterval(1, 10, []string{"a", "b"})
	_ = t.AddInterval(20, 30, []bool{true, false})
	_ = t.AddInterval(32, 35, []int{1, 2, 3})
	_ = t.AddInterval(32, 38, nil)
	t.Sort()
	fmt.Println(t.Len())
	fmt.Println(t.Iter())
	fmt.Println(t.Query(33))

	// Output:
	// 4
	// [{1 10 [a b]} {32 35 [1 2 3]} {32 38 <nil>} {20 30 [true false]}]
	// [{32 35 [1 2 3]} {32 38 <nil>}]
}
