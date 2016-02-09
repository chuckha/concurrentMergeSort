package main

import (
	"fmt"
)

func main() {
	input := []int{9, 3, 4, 5, 2, 1, 3, 4, 6, 7, 8, 3, 3}
	fmt.Println(input)
	s := NewSorter()
	s.Sort(input)
	close(s.merger)
	fmt.Println(<-s.result)
	close(s.result)
}

type Sorter struct {
	result chan []int
	merger chan []int
}

func NewSorter() *Sorter {
	sorter := &Sorter{
		result: make(chan []int),
		merger: make(chan []int),
	}
	go sorter.merge()
	return sorter
}

// gets called 1+nlgn
func (s *Sorter) Sort(x []int) {
	//	fmt.Println(fmt.Sprintf("incoming x: %v", x))
	//	fmt.Println(fmt.Sprintf("using this slice: x[:%v/2] or x[:%v]", len(x), len(x)/2))
	if len(x) == 0 {
		return
	}
	if len(x) == 1 {
		//		fmt.Println(fmt.Sprintf("pushing %v onto merger", x))
		s.merger <- x
		return
	}
	s.Sort(x[:len(x)/2])
	s.Sort(x[len(x)/2:])
}

// runges mergeArrays 1+nlgn times
// which means this is (n+m) * (1+nlgn)
// n + m*n*lgn+m+n^2lgn=O(n^2lgn)
func (s *Sorter) merge() {
	var next []int
	for {
		//		fmt.Println("pulling from merger")
		left, ok := <-s.merger
		//		fmt.Println(fmt.Sprintf("pulled: %v, %v", left, ok))
		if !ok {
			s.result <- next
			return
		}
		next = mergeArrays(next, left)
		//		fmt.Println(fmt.Sprintf("next: %v", next))
	}
}

// O(n+m)
func mergeArrays(a, b []int) []int {
	merged := make([]int, len(a)+len(b))
	var ai, bi, mi int
	for ai < len(a) && bi < len(b) {
		if a[ai] <= b[bi] {
			merged[mi] = a[ai]
			ai++
			mi++
		} else {
			merged[mi] = b[bi]
			bi++
			mi++
		}
	}
	for ai < len(a) {
		merged[mi] = a[ai]
		ai++
		mi++
	}
	for bi < len(b) {
		merged[mi] = b[bi]
		bi++
		mi++
	}
	return merged
}
