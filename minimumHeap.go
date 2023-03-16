package main

import "fmt"

type MinimumHeap struct {
	slice []int
}

func (heap *MinimumHeap) leftChildIndex(index int) int {
	return index*2 + 1
}

func (heap *MinimumHeap) rightChildIndex(index int) int {
	return index*2 + 2
}

// 0 1 2 3 4 5 6 7 8 9 10
// - 0 0 1 1 2 2 3 3 4 4
func (heap *MinimumHeap) parentIndex(index int) int {
	if index == 0 {
		return -1
	}
	if index%2 == 0 {
		return (index - 1) / 2
	} else {
		return index / 2
	}
}

func (heap *MinimumHeap) insert(num int) {
	heap.slice = append(heap.slice, num)
	index := len(heap.slice) - 1
	pIndex := heap.parentIndex(index)
	slice := heap.slice
	for index != 0 && slice[index] < slice[pIndex] {
		slice[index], slice[pIndex] = slice[pIndex], slice[index]
		index = pIndex
		pIndex = heap.parentIndex(index)
	}

	i := 0
	li := heap.leftChildIndex(i)
	for li < len(slice) && slice[i] > slice[li] {
		slice[i], slice[li] = slice[li], slice[i]
		i = li
		li = heap.leftChildIndex(i)
	}

	i = 0
	ri := heap.rightChildIndex(i)
	for ri < len(slice) && slice[i] > slice[ri] {
		slice[i], slice[ri] = slice[ri], slice[i]
		i = ri
		ri = heap.rightChildIndex(i)
	}
}

// 主要思想：把最后一个元素跟头部互换，然后跟左右哪个小换哪个
func (heap *MinimumHeap) pop() int {
	slice := heap.slice
	if len(slice) == 0 {
		return -1
	}

	top := slice[0]
	slice[0] = slice[len(slice)-1]
	heap.slice = slice[:len(slice)-1]

	i := 0
	li := heap.leftChildIndex(i)
	ri := heap.rightChildIndex(i)
	for {
		smallerThanLeft := true
		smallerThanRight := true
		if li < len(slice) && slice[i] > slice[li] {
			smallerThanLeft = false
		}

		if ri < len(slice) && slice[i] > slice[ri] {
			smallerThanLeft = false
		}

		if smallerThanLeft && smallerThanRight {
			break
		}

		if li < len(slice) && ri < len(slice) {
			if slice[li] < slice[ri] {
				slice[i], slice[li] = slice[li], slice[i]
				i = li
			} else {
				slice[i], slice[ri] = slice[ri], slice[i]
				i = ri
			}
		}

		li = heap.leftChildIndex(i)
		ri = heap.rightChildIndex(i)
	}

	return top
}

func Construct(arr []int) *MinimumHeap {
	heap := &MinimumHeap{}
	for _, num := range arr {
		heap.insert(num)
	}
	return heap
}

func test_MinimumHeap() {
	heap := Construct([]int{3, 12, 3, 213, 21, 4, 2, 1, 43, 21, 3, 21, 321})
	for num := heap.pop(); num != -1; num = heap.pop() {
		fmt.Println(num)
	}
}
