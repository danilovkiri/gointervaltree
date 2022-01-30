// Package gointervaltree provides functionality for indexing a set of integer intervals, e.g. [start, end)
// based on http://en.wikipedia.org/wiki/Interval_tree. Copyright 2022, Kirill Danilov. Licensed under MIT license.
package gointervaltree

import (
	"log"
	"reflect"
	"sort"
)

// IntervalTree struct defines data structure for indexing a set of integer intervals, e.g. [start, end).
type IntervalTree struct {
	min              int
	max              int
	center           int
	singleInterval   []interface{}
	leftSubtree      *IntervalTree
	rightSubtree     *IntervalTree
	midSortedByStart []interface{}
	midSortedByEnd   []interface{}
}

// NewIntervalTree method instantiates an instance of IntervalTree struct creating a node for keeping intervals.
func NewIntervalTree(min int, max int) (tree *IntervalTree) {
	tree = new(IntervalTree)
	tree.min = min
	tree.max = max
	if !(tree.min < tree.max) {
		log.Panic("AssertionError: interval tree start must be numerically less than its end")
	}
	tree.center = (min + max) / 2
	tree.singleInterval = nil
	tree.leftSubtree = nil
	tree.rightSubtree = nil
	tree.midSortedByStart = []interface{}{}
	tree.midSortedByEnd = []interface{}{}
	return tree
}

// AddInterval method adds intervals to the tree without sorting them along the way.
func (tree *IntervalTree) AddInterval(start int, end int, data interface{}) {
	if (end - start) <= 0 {
		return
	}
	if tree.singleInterval == nil {
		tree.singleInterval = []interface{}{start, end, data}
	} else if reflect.DeepEqual(tree.singleInterval, []interface{}{0}) {
		tree.addIntervalMain(start, end, data)
	} else {
		tree.addIntervalMain(tree.singleInterval[0].(int), tree.singleInterval[1].(int), tree.singleInterval[2])
		tree.singleInterval = []interface{}{0}
		tree.addIntervalMain(start, end, data)
	}
}

// addIntervalMain method is a technical method used inside AddInterval.
func (tree *IntervalTree) addIntervalMain(start int, end int, data interface{}) {
	if end <= tree.center {
		if tree.leftSubtree == nil {
			tree.leftSubtree = NewIntervalTree(tree.min, tree.center)
		}
		tree.leftSubtree.AddInterval(start, end, data)
	} else if start > tree.center {
		if tree.rightSubtree == nil {
			tree.rightSubtree = NewIntervalTree(tree.center, tree.max)
		}
		tree.rightSubtree.AddInterval(start, end, data)
	} else {
		tree.midSortedByStart = append(tree.midSortedByStart, []interface{}{start, end, data})
		tree.midSortedByEnd = append(tree.midSortedByEnd, []interface{}{start, end, data})
	}
}

// Sort method is used to sort intervals within the tree and must be invoked after adding intervals.
func (tree *IntervalTree) Sort() {
	if tree.singleInterval == nil || !reflect.DeepEqual(tree.singleInterval, []interface{}{0}) {
		return
	}
	sort.Slice(tree.midSortedByStart, func(i, j int) bool {
		return tree.midSortedByStart[i].([]interface{})[0].(int) < tree.midSortedByStart[j].([]interface{})[0].(int)
	})
	sort.Slice(tree.midSortedByEnd, func(i, j int) bool {
		return tree.midSortedByEnd[i].([]interface{})[1].(int) > tree.midSortedByEnd[j].([]interface{})[1].(int)
	})
	if tree.leftSubtree != nil {
		tree.leftSubtree.Sort()
	}
	if tree.rightSubtree != nil {
		tree.rightSubtree.Sort()
	}
}

// Query method returns all intervals in the tree which overlap given point,
// i.e. all (start, end, data) records, for which (start <= x < end).
func (tree *IntervalTree) Query(x int) []interface{} {
	var result []interface{}
	if tree.singleInterval == nil {
		return result
	} else if !reflect.DeepEqual(tree.singleInterval, []interface{}{0}) {
		if tree.singleInterval[0].(int) <= x && x < tree.singleInterval[1].(int) {
			result = append(result, tree.singleInterval)
		}
		return result
	} else if x < tree.center {
		if tree.leftSubtree != nil {
			result = append(result, tree.leftSubtree.Query(x)...)
		}
		for _, element := range tree.midSortedByStart {
			if element.([]interface{})[0].(int) <= x {
				result = append(result, element)
			} else {
				break
			}
		}
		return result
	} else {
		for _, element := range tree.midSortedByEnd {
			if element.([]interface{})[1].(int) > x {
				result = append(result, element)
			} else {
				break
			}
		}
		if tree.rightSubtree != nil {
			result = append(result, tree.rightSubtree.Query(x)...)
		}
		return result
	}
}

// Len method represents the number of intervals maintained in the tree, zero- or negative-size intervals
// are not registered.
func (tree *IntervalTree) Len() int {
	if tree.singleInterval == nil {
		return 0
	} else if !reflect.DeepEqual(tree.singleInterval, []interface{}{0}) {
		return 1
	} else {
		size := len(tree.midSortedByStart)
		if tree.leftSubtree != nil {
			size += tree.leftSubtree.Len()
		}
		if tree.rightSubtree != nil {
			size += tree.rightSubtree.Len()
		}
		return size
	}
}

// Iter method returns a slice of all intervals maintained in the tree.
func (tree *IntervalTree) Iter() []interface{} {
	var result []interface{}
	if tree.singleInterval == nil {
		return result
	} else if !reflect.DeepEqual(tree.singleInterval, []interface{}{0}) {
		result = append(result, tree.singleInterval)
		return result
	} else {
		if tree.leftSubtree != nil {
			result = append(result, tree.leftSubtree.Iter()...)
		}
		if tree.rightSubtree != nil {
			result = append(result, tree.rightSubtree.Iter()...)
		}
		for _, element := range tree.midSortedByStart {
			result = append(result, element)
		}
		return result
	}
}
