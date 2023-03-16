package main

import "fmt"

type ListNode struct {
	key   int
	value int
	next  *ListNode
}

type NodeLocator struct {
	prev     *ListNode
	refCount int
}

type NodeBucket struct {
	begin *ListNode
	end   *ListNode
}

type LfuCache struct {
	m1  map[int]*NodeLocator
	m2  map[int]*NodeBucket
	cap int
}

func Constructor(capacity int) *LfuCache {
	return &LfuCache{m1: make(map[int]*NodeLocator), m2: make(map[int]*NodeBucket), cap: capacity}
}

func (cache *LfuCache) Put(key int, value int) {
	nodeLocator, ok := cache.m1[key]
	if ok {
		prev := nodeLocator.prev
		nodeBucket := cache.m2[nodeLocator.refCount]
		if prev == nil {
			nodeBucket.begin.value = value
		} else {
			prev.next.value = value
		}

		fmt.Printf("Put || update node key=%d, value=%d\n", key, value)
	} else {
		nodeBucket, ok2 := cache.m2[1]
		newNode := &ListNode{key: key, value: value, next: nil}
		var prevNode *ListNode = nil
		if ok2 {
			prevNode = nodeBucket.end
			nodeBucket.end.next = newNode
			nodeBucket.end = nodeBucket.end.next

			if len(cache.m1) == cache.cap {
				// drop node
				beginKey := nodeBucket.begin.key
				fmt.Printf("Put || drop node key=%d\n", beginKey)
				delete(cache.m1, beginKey)
				nodeBucket.begin = nodeBucket.begin.next
			}
		} else {
			if len(cache.m1) == cache.cap {
				// drop node
				for k := range cache.m2 {
					lfuBucket, _ := cache.m2[k]
					beginKey := lfuBucket.begin.key
					fmt.Printf("Put || drop node key=%d\n", beginKey)
					delete(cache.m1, beginKey)
					if lfuBucket.begin == lfuBucket.end {
						delete(cache.m2, k)
					} else {
						lfuBucket.begin = lfuBucket.begin.next
					}
					break
				}
			}

			cache.m2[1] = &NodeBucket{begin: newNode, end: newNode}
		}

		fmt.Printf("Put || insert new node key=%d, value=%d\n", key, value)
		cache.m1[key] = &NodeLocator{prev: prevNode, refCount: 1}
	}
}

func (cache *LfuCache) Get(key int) int {
	nodeLocator, ok := cache.m1[key]
	if ok {
		nodeBucket := cache.m2[nodeLocator.refCount]
		newCount := nodeLocator.refCount + 1
		fmt.Printf("Get || key=%d's refCount raise to %d\n", key, newCount)
		prev := nodeLocator.prev

		var node *ListNode = nil

		// m2 clear old node
		if prev == nil {
			node = nodeBucket.begin
			if nodeBucket.end == nodeBucket.begin {
				delete(cache.m2, key)
			} else {
				nodeBucket.begin = nodeBucket.begin.next
			}
		} else {
			node = prev.next
			prev.next = node.next
		}
		res := node.value
		node.next = nil

		newPrevNode := &ListNode{next: nil}
		// m2 insert new node
		nodeBucket2, ok2 := cache.m2[newCount]
		if ok2 {
			newPrevNode.next = nodeBucket2.end
			nodeBucket2.end.next = node
			nodeBucket2.end = node
		} else {
			nodeBucket2 = &NodeBucket{begin: node, end: node}
			cache.m2[newCount] = nodeBucket2
			newPrevNode = nil
		}

		// update node locator in m1
		nodeLocator.refCount++
		nodeLocator.prev = newPrevNode
		return res
	}
	return -1
}

func test_LFUCache() {
	cache := Constructor(2)
	cache.Put(1, 1)
	cache.Get(1)
	cache.Put(1, 2)
	cache.Get(1)
	cache.Put(2, 2)
	cache.Put(1, 2)

	cache.Put(3, 1)
}
