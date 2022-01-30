package gointervaltree

import (
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestNewIntervalTree(t *testing.T) {
	assert := assert.New(t)
	intervals := [][]int{[]int{10, 20}, []int{20, 30}, []int{21, 31}, []int{30, 40}, []int{45, 55}, []int{45, 56},
		[]int{46, 57}, []int{55, 56}, []int{58, 59}, []int{50, 51}}
	queryPoints := []int{-1, 0, 1, 10, 11, 19, 20, 21, 24, 25, 26, 30, 40, 41, 48, 49, 50, 51, 52, 60,
		74, 75, 76, 90, 100, 1000}
	doTest(0, 100, intervals, queryPoints, assert)

}

func doTest(min int, max int, intervals [][]int, queryPoints []int, assert *assert.Assertions) {
	tree := NewIntervalTree(min, max)
	for _, interval := range intervals {
		tree.AddInterval(interval[0], interval[1], nil)
	}
	tree.Sort()
	for _, q := range queryPoints {
		r := tree.Query(q)
		sort.Slice(r, func(i, j int) bool {
			if r[i].([]interface{})[0].(int) != r[j].([]interface{})[0].(int) {
				return r[i].([]interface{})[0].(int) > r[j].([]interface{})[0].(int)
			} else {
				return r[i].([]interface{})[1].(int) > r[j].([]interface{})[1].(int)
			}
		})
		var true_r []interface{}
		for _, interval := range intervals {
			if (interval[0] <= q) && (q < interval[1]) {
				true_r = append(true_r, []interface{}{interval[0], interval[1], nil})
			}
		}
		sort.Slice(true_r, func(i, j int) bool {
			if true_r[i].([]interface{})[0].(int) != true_r[j].([]interface{})[0].(int) {
				return true_r[i].([]interface{})[0].(int) > true_r[j].([]interface{})[0].(int)
			} else {
				return true_r[i].([]interface{})[1].(int) > true_r[j].([]interface{})[1].(int)
			}
		})
		assert.Equal(true_r, r)
	}
	observedLength := tree.Len()
	expectedLength := 0
	for _, interval := range intervals {
		if interval[0] < interval[1] {
			expectedLength += 1
		}
	}
	assert.Equal(expectedLength, observedLength)
	assert.Equal(expectedLength, len(tree.Iter()))
}
