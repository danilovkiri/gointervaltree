// Package gointervaltree provides functionality for indexing a set of integer intervals, e.g. [start, end)
// based on http://en.wikipedia.org/wiki/Interval_tree. Copyright 2022, Kirill Danilov. Licensed under MIT license.
package gointervaltree

import (
	"errors"
	"golang.org/x/exp/constraints"
	"sort"
)

// resultInterval is a node of an intervalTree without technical fields
type resultInterval[T constraints.Signed] struct {
	start T
	end   T
	data  any
}

// interval is a node of an intervalTree.
type interval[T constraints.Signed] struct {
	start   T
	end     T
	data    any
	blocked bool
}

// intervalTree struct defines data structure for indexing a set of integer intervals, e.g. [start, end).
type intervalTree[T constraints.Signed] struct {
	min              T
	max              T
	center           T
	singleInterval   *interval[T]
	leftSubtree      *intervalTree[T]
	rightSubtree     *intervalTree[T]
	midSortedByStart []*interval[T]
	midSortedByEnd   []*interval[T]
}

// NewIntervalTree creates and returns an IntervalTree object.
func NewIntervalTree[T constraints.Signed](min, max T) (*intervalTree[T], error) {
	tree := new(intervalTree[T])
	tree.min = min
	tree.max = max
	if !(tree.min < tree.max) {
		return nil, errors.New("interval tree start must be numerically less than its end")
	}
	tree.center = (min + max) / 2
	tree.singleInterval = nil
	tree.leftSubtree = nil
	tree.rightSubtree = nil
	tree.midSortedByStart = []*interval[T]{}
	tree.midSortedByEnd = []*interval[T]{}
	return tree, nil
}

// AddInterval method adds intervals to the tree without sorting them along the way.
func (tree *intervalTree[T]) AddInterval(start, end T, data any) error {
	if (end - start) <= 0 {
		return errors.New("interval start must be numerically less than its end")
	}
	if tree.singleInterval == nil {
		tree.singleInterval = &interval[T]{start, end, data, false}
	} else if !tree.singleInterval.blocked { // singleInterval is not blocked
		tree.addIntervalMain(tree.singleInterval.start, tree.singleInterval.end, tree.singleInterval.data)
		tree.singleInterval.blocked = true
		tree.addIntervalMain(start, end, data)
	} else { // singleInterval is blocked
		tree.addIntervalMain(start, end, data)
	}
	return nil
}

// addIntervalMain method is a technical method used inside AddInterval.
func (tree *intervalTree[T]) addIntervalMain(start, end T, data any) {
	if end <= tree.center {
		if tree.leftSubtree == nil {
			tree.leftSubtree, _ = NewIntervalTree(tree.min, tree.center)
		}
		_ = tree.leftSubtree.AddInterval(start, end, data)
	} else if start > tree.center {
		if tree.rightSubtree == nil {
			tree.rightSubtree, _ = NewIntervalTree(tree.center, tree.max)
		}
		_ = tree.rightSubtree.AddInterval(start, end, data)
	} else {
		tree.midSortedByStart = append(tree.midSortedByStart, &interval[T]{start, end, data, false})
		tree.midSortedByEnd = append(tree.midSortedByEnd, &interval[T]{start, end, data, false})
	}
}

// Sort method is used to sort intervals within the tree and must be invoked after adding intervals.
func (tree *intervalTree[T]) Sort() {
	if tree.singleInterval == nil || !tree.singleInterval.blocked {
		return
	}
	sort.Slice(tree.midSortedByStart, func(i, j int) bool {
		return tree.midSortedByStart[i].start < tree.midSortedByStart[j].start

	})
	sort.Slice(tree.midSortedByEnd, func(i, j int) bool {
		return tree.midSortedByEnd[i].end > tree.midSortedByEnd[j].end
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
func (tree *intervalTree[T]) Query(x T) []resultInterval[T] {
	var result []resultInterval[T]
	if tree.singleInterval == nil {
		return result
	} else if !tree.singleInterval.blocked {
		if tree.singleInterval.start <= x && x < tree.singleInterval.end {
			result = append(result, resultInterval[T]{start: (*tree.singleInterval).start, end: (*tree.singleInterval).end, data: (*tree.singleInterval).data})
		}
		return result
	} else if x < tree.center {
		if tree.leftSubtree != nil {
			result = append(result, tree.leftSubtree.Query(x)...)
		}
		for _, element := range tree.midSortedByStart {
			if element.start <= x {
				result = append(result, resultInterval[T]{start: (*element).start, end: (*element).end, data: (*element).data})
			} else {
				break
			}
		}
		return result
	} else {
		for _, element := range tree.midSortedByEnd {
			if element.end > x {
				result = append(result, resultInterval[T]{start: (*element).start, end: (*element).end, data: (*element).data})
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

// Len represents the number of intervals maintained in the tree, zero- or negative-size intervals are not registered.
func (tree *intervalTree[T]) Len() int {
	if tree.singleInterval == nil {
		return 0
	} else if !tree.singleInterval.blocked {
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
func (tree *intervalTree[T]) Iter() []resultInterval[T] {
	var result []resultInterval[T]
	if tree.singleInterval == nil {
		return result
	} else if !tree.singleInterval.blocked {
		result = append(result, resultInterval[T]{start: (*tree.singleInterval).start, end: (*tree.singleInterval).end, data: (*tree.singleInterval).data})
		return result
	} else {
		if tree.leftSubtree != nil {
			result = append(result, tree.leftSubtree.Iter()...)
		}
		if tree.rightSubtree != nil {
			result = append(result, tree.rightSubtree.Iter()...)
		}
		// cannot use `result = append(result, tree.midSortedByStart...)` due to explicit dereferencing
		for _, i := range tree.midSortedByStart {
			result = append(result, resultInterval[T]{start: (*i).start, end: (*i).end, data: (*i).data})
		}
		return result
	}
}
