package main

import (
	"cmp"
	"fmt"
	"math/rand"
)

type SkipList[K cmp.Ordered, V any] struct {
	head *node[K, V]
}

type node[K cmp.Ordered, V any] struct {
	key   K
	val   V
	nexts []*node[K, V]
}

func (s *SkipList[K, V]) search(key K) *node[K, V] {
	move := s.head

	for level := len(s.head.nexts) - 1; level >= 0; level-- {
		for move.nexts[level] != nil && move.nexts[level].key < key {
			move = move.nexts[level]
		}

		if move.nexts[level] != nil && move.nexts[level].key == key {
			return move.nexts[level]
		}
	}
	return nil
}

func (s *SkipList[K, V]) Get(key K) (*V, bool) {
	if _node := s.search(key); _node != nil {
		return &_node.val, true
	}
	return nil, false
}

func (s *SkipList[K, V]) roll() int {
	level := 0
	for rand.Intn(2)-1 > 0 {
		level++
	}
	return level
}

func (s *SkipList[K, V]) Put(key K, val V) {
	if _node := s.search(key); _node != nil {
		_node.val = val
		return
	}

	level := s.roll()

	for len(s.head.nexts)-1 < level {
		s.head.nexts = append(s.head.nexts, nil)
	}

	newNode := node[K, V]{
		key:   key,
		val:   val,
		nexts: make([]*node[K, V], level+1),
	}
	move := s.head

	for l := level; l >= 0; l-- {
		for move.nexts[level] != nil && move.nexts[level].key < key {
			move = move.nexts[level]
		}
		newNode.nexts[level] = move.nexts[level]
		move.nexts[level] = &newNode
		//move -> newNode -> move.next
	}
}

func (s *SkipList[K, V]) Del(key K) {
	if _node := s.search(key); _node == nil {
		return
	}

	move := s.head
	for level := len(s.head.nexts) - 1; level >= 0; level-- {
		for move.nexts[level] != nil && move.nexts[level].key < key {
			move = move.nexts[level]
		}

		if move.nexts[level] == nil || move.nexts[level].key > key {
			continue
		}

		// move.key == key
		move.nexts[level] = move.nexts[level].nexts[level]
	}

	// update SkipList hight
	dif := 0
	for level := len(s.head.nexts) - 1; level > 0 && s.head.nexts[level] == nil; level-- {
		dif++
	}
	s.head.nexts = s.head.nexts[:len(s.head.nexts)-dif]
}

func (s *SkipList[K, V]) Range(start K, end K) []struct {
	Key K
	Val V
} {
	ceilNode := s.ceiling(start)
	if ceilNode == nil {
		return make([]struct {
			Key K
			Val V
		}, 0)
	}
	var res []struct {
		Key K
		Val V
	}
	for move := ceilNode; move != nil && move.key <= end; move = move.nexts[0] {
		res = append(res, struct {
			Key K
			Val V
		}{move.key, move.val})
	}
	return res
}

func (s *SkipList[K, V]) ceiling(target K) *node[K, V] {
	move := s.head

	for level := len(s.head.nexts) - 1; level >= 0; level-- {
		for move.nexts[level] != nil && move.nexts[level].key < target {
			move = move.nexts[level]
		}

		if move.nexts[level] != nil && move.nexts[level].key == target {
			return move.nexts[level]
		}
	}
	return move.nexts[0]
}

func (s *SkipList[K, V]) Floor(target K) (*struct {
	Key K
	Val V
}, bool) {
	if floorNode := s.floor(target); floorNode != nil {
		return &struct {
			Key K
			Val V
		}{floorNode.key, floorNode.val}, true
	}
	return nil, false
}

func (s *SkipList[K, V]) floor(target K) *node[K, V] {
	move := s.head

	for level := len(s.head.nexts) - 1; level >= 0; level-- {
		for move.nexts[level] != nil && move.nexts[level].key < target {
			move = move.nexts[level]
		}
		if move.nexts[level] != nil && move.nexts[level].key == target {
			return move.nexts[level]
		}
	}

	return move
}

func main() {
	// Initialize SkipList
	list := SkipList[int, string]{head: &node[int, string]{}}

	// Insert elements
	list.Put(1, "One")
	list.Put(2, "Two")
	list.Put(3, "Three")

	// Search for an element
	if val, found := list.Get(2); found {
		fmt.Println("Value found:", *val) // Output: Value found: Two
	}

	// Delete an element
	list.Del(2)

	// Range query
	results := list.Range(1, 3)
	for _, res := range results {
		fmt.Println("Key:", res.Key, "Value:", res.Val)
	}
}
