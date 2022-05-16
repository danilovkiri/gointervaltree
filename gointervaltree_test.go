package gointervaltree

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

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
	assert.Equal(t, []interface{}(nil), tree.Query(1))
}

func TestIntervalTree_LenEmptyTree(t *testing.T) {
	tree, _ := NewIntervalTree(10, 50)
	assert.Equal(t, 0, tree.Len())
}

func TestIntervalTree_IterEmptyTree(t *testing.T) {
	tree, _ := NewIntervalTree(10, 50)
	assert.Equal(t, []interface{}(nil), tree.Iter())
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
			if r[i].([]interface{})[0].(int) != r[j].([]interface{})[0].(int) {
				return r[i].([]interface{})[0].(int) > r[j].([]interface{})[0].(int)
			}
			return r[i].([]interface{})[1].(int) > r[j].([]interface{})[1].(int)
		})
		var trueR []interface{}
		for _, interval := range intervals {
			if (interval[0] <= q) && (q < interval[1]) {
				trueR = append(trueR, []interface{}{interval[0], interval[1], nil})
			}
		}
		sort.Slice(trueR, func(i, j int) bool {
			if trueR[i].([]interface{})[0].(int) != trueR[j].([]interface{})[0].(int) {
				return trueR[i].([]interface{})[0].(int) > trueR[j].([]interface{})[0].(int)
			}
			return trueR[i].([]interface{})[1].(int) > trueR[j].([]interface{})[1].(int)
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
