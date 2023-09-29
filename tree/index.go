package tree

import "log"

type index[K comparable, T any] map[K]Node[K, T]

func (idx *index[K, T]) find(id K) Node[K, T] {
	if idx == nil { // do we need an error check here?
		log.Println("Attempting to find in an undefined index")
		return nil
	}
	m := *idx
	val, exists := m[id]
	if !exists {
		return nil
	}
	return val
}

func (idx *index[K, T]) insert(id K, node Node[K, T]) bool {
	if idx == nil { // do we need an error check here?
		log.Println("Attempting to insert in an undefined index")
		return false
	}
	m := *idx
	m[id] = node
	return true
}
