package tree

// Minimal implementaiton of a breadth and depth first search for testing
// if a tree matches what is expected. Yields primary keys of all nodes
// in a tree in order of selected traversal.

// returns the primary keys of a tree in order of a breadth first search
func bfc[K comparable, T any](q []Node[K, T], iter []K) []K {
	if q == nil || len(q) == 0 || q[0] == nil {
		return []K{}
	}
	iter = append(iter, q[0].GetID())
	q = append(q[1:], q[0].GetChildren()...)
	if len(q) == 0 {
		return iter
	}
	return bfc(q, iter)
}

// returns the primary keys of a tree in order of a breadth first search
func dfc[K comparable, T any](n Node[K, T], iter []K) []K {
	if n == nil {
		return []K{}
	}
	iter = append(iter, n.GetID())
	for _, c := range n.GetChildren() {
		iter = dfc(c, iter)
	}
	return iter
}
