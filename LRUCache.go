package main

import "fmt"

type ListNode struct {
	key   int
	value int
	next  *ListNode
}

// @lc code=start
type LRUCache struct {
	m     map[int]*ListNode
	begin *ListNode
	end   *ListNode
	cap   int
}

func Constructor(capacity int) LRUCache {
	res := LRUCache{
		cap:   capacity,
		begin: nil,
		end:   nil,
		m:     make(map[int]*ListNode)}
	return res
}

func (cache *LRUCache) Get(key int) int {
	_, ok := cache.m[key]
	if ok {
		prevNode := cache.m[key]
		if prevNode == nil {
			return cache.begin.value
		} else {
			node := prevNode.next
			prevNode.next = node.next
			node.next = cache.begin
			cache.begin = node
			return node.value
		}
	}
	return -1
}

func (cache *LRUCache) Put(key int, value int) {
	_, ok := cache.m[key]
	if ok {
		prev := cache.m[key]
		if prev == nil {
			cache.begin.value = value
		} else {
			prev.next.value = value
		}
		cache.Get(key)
	} else {
		if len(cache.m) == cache.cap {
			prev := cache.m[cache.end.key]
			cache.m[key] = prev

			delete(cache.m, cache.end.key)

			cache.end.key = key
			cache.end.value = value
		} else {
			if cache.end == nil {
				cache.begin = &ListNode{key: key, value: value}
				cache.end = cache.begin
				cache.m[key] = nil
			} else {
				newEnd := &ListNode{key: key, value: value}
				cache.end.next = newEnd
				cache.m[key] = cache.end
				cache.end = newEnd
			}
		}
	}
}

func test_LRUCache() {
	lruCache := Constructor(1)
	lruCache.Put(1, 1)
	fmt.Println(lruCache.Get(1))
	lruCache.Put(1, 2)
	fmt.Println(lruCache.Get(1))
	lruCache.Put(2, 2)
	fmt.Println(lruCache.Get(1))
	fmt.Println(lruCache.Get(2))
}
