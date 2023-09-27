package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBFS(t *testing.T) {

	tests := map[string]struct {
		tree      func() *Tree[uint, int]
		traversal TraversalType
		expSearch []uint
	}{
		"success": {
			tree: func() *Tree[uint, int] {
				node6 := &node[uint, int]{primary: 6}
				node5 := &node[uint, int]{primary: 5}
				node4 := &node[uint, int]{primary: 4}
				node3 := &node[uint, int]{primary: 3, children: []Node[uint, int]{node4, node5}}
				node2 := &node[uint, int]{primary: 2, children: []Node[uint, int]{node6}}
				node1 := &node[uint, int]{primary: 1, children: []Node[uint, int]{node2, node3}}
				return &Tree[uint, int]{root: node1}
			},
			traversal: TraverseBreadthFirst,
			expSearch: []uint{1, 2, 3, 6, 4, 5},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			i := 0
			for g := range tt.tree().Traverse(tt.traversal) {
				assert.Equal(t, tt.expSearch[i], g.GetID())
				i = i + 1
			}

		})
	}
}
