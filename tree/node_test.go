package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetID(t *testing.T) {

	tests := map[string]struct {
		n     Node[uint, int]
		expID uint
	}{
		"trivial": {
			n:     &node[uint, int]{primary: 1},
			expID: 1,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotID := tt.n.GetID()

			assert.Equal(t, tt.expID, gotID)
		})
	}
}

func TestGetParentID(t *testing.T) {

	tests := map[string]struct {
		n   Node[uint, int]
		exp uint
	}{
		"trivial": {
			n:   &node[uint, int]{parentID: 1},
			exp: 1,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := tt.n.GetParentID()

			assert.Equal(t, tt.exp, got)
		})
	}
}

func TestGetChildren(t *testing.T) {

	node1 := &node[uint, int]{primary: 1}
	node2 := &node[uint, int]{primary: 2}

	tests := map[string]struct {
		n        Node[uint, int]
		expChild []Node[uint, int]
	}{
		"nil child array": {
			n:        &node[uint, int]{},
			expChild: nil,
		},
		"empty child array": {
			n:        &node[uint, int]{children: []Node[uint, int]{}},
			expChild: []Node[uint, int]{},
		},
		"success": {
			n:        &node[uint, int]{children: []Node[uint, int]{node1, node2}},
			expChild: []Node[uint, int]{node1, node2},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotChild := tt.n.GetChildren()

			assert.Equal(t, tt.expChild, gotChild)
		})
	}
}

func TestGetParent(t *testing.T) {

	node1 := &node[uint, int]{primary: 1}

	tests := map[string]struct {
		n       Node[uint, int]
		expNode Node[uint, int]
	}{
		"nil parent": {
			n:       &node[uint, int]{},
			expNode: nil,
		},
		"success": {
			n:       &node[uint, int]{parent: node1},
			expNode: node1,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotNode := tt.n.GetParent()

			assert.Equal(t, tt.expNode, gotNode)
		})
	}
}

func TestNodeAddChildren(t *testing.T) {

	node1 := &node[uint, int]{primary: 1}
	node2 := &node[uint, int]{primary: 2}
	node3 := &node[uint, int]{primary: 3}

	tests := map[string]struct {
		n        Node[uint, int]
		argNodes []Node[uint, int]
		expChild []Node[uint, int]
	}{
		"nil child array": {
			n:        &node[uint, int]{},
			argNodes: []Node[uint, int]{node1},
			expChild: []Node[uint, int]{node1},
		},
		"empty child array": {
			n:        &node[uint, int]{children: []Node[uint, int]{}},
			argNodes: []Node[uint, int]{node1},
			expChild: []Node[uint, int]{node1},
		},
		"add nil": {
			n:        &node[uint, int]{children: []Node[uint, int]{node1}},
			argNodes: nil,
			expChild: []Node[uint, int]{node1},
		},
		"add empty array": {
			n:        &node[uint, int]{children: []Node[uint, int]{node1}},
			argNodes: []Node[uint, int]{},
			expChild: []Node[uint, int]{node1},
		},
		"non-empty child array": {
			n:        &node[uint, int]{children: []Node[uint, int]{node1}},
			argNodes: []Node[uint, int]{node2, node3},
			expChild: []Node[uint, int]{node1, node2, node3},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.n.AddChildren(tt.argNodes...)

			assert.Equal(t, tt.expChild, tt.n.GetChildren())
		})
	}
}

func TestReplaceChildren(t *testing.T) {

	node1 := &node[uint, int]{primary: 1}
	node2 := &node[uint, int]{primary: 2}
	node3 := &node[uint, int]{primary: 3}

	tests := map[string]struct {
		n        Node[uint, int]
		argNodes []Node[uint, int]
		expChild []Node[uint, int]
	}{
		"nil child array": {
			n:        &node[uint, int]{},
			argNodes: []Node[uint, int]{node1},
			expChild: []Node[uint, int]{node1},
		},
		"empty child array": {
			n:        &node[uint, int]{children: []Node[uint, int]{}},
			argNodes: []Node[uint, int]{node1},
			expChild: []Node[uint, int]{node1},
		},
		"use nil": {
			n:        &node[uint, int]{children: []Node[uint, int]{node1}},
			argNodes: nil,
			expChild: []Node[uint, int]{},
		},
		"use empty array": {
			n:        &node[uint, int]{children: []Node[uint, int]{node1}},
			argNodes: []Node[uint, int]{},
			expChild: []Node[uint, int]{},
		},
		"non-empty replacement array": {
			n:        &node[uint, int]{children: []Node[uint, int]{node1}},
			argNodes: []Node[uint, int]{node2, node3},
			expChild: []Node[uint, int]{node2, node3},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.n.ReplaceChildren(tt.argNodes...)

			assert.Equal(t, tt.expChild, tt.n.GetChildren())
		})
	}

}

func TestSetParent(t *testing.T) {

	node1 := &node[uint, int]{primary: 1}
	node2 := &node[uint, int]{primary: 2}

	tests := map[string]struct {
		n           Node[uint, int]
		argParent   Node[uint, int]
		expParent   Node[uint, int]
		expParentID uint
	}{
		"set nil parent": {
			n:           &node[uint, int]{primary: 1},
			argParent:   nil,
			expParent:   nil,
			expParentID: 0,
		},
		"set circular ref parent": {
			n:           node1,
			argParent:   node1,
			expParent:   nil,
			expParentID: 0,
		},
		"success": {
			n:           &node[uint, int]{primary: 1},
			argParent:   node2,
			expParent:   node2,
			expParentID: 2,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			n := tt.n
			n.setParent(tt.argParent)

			assert.Equal(t, tt.expParent, n.GetParent())
			assert.Equal(t, tt.expParentID, n.GetParentID())
		})
	}

}
