package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexFind(t *testing.T) {

	node1 := &node[uint, int]{primary: 1}
	node2 := &node[uint, int]{primary: 2}

	tests := map[string]struct {
		index   index[uint, int]
		argID   uint
		expNode Node[uint, int]
	}{
		"nil index": {
			index:   nil,
			argID:   1,
			expNode: nil,
		},
		"not in index": {
			index: index[uint, int]{
				1: node1,
				2: node2,
			},
			argID:   3,
			expNode: nil,
		},
		"success": {
			index: index[uint, int]{
				1: node1,
				2: node2,
			},
			argID:   2,
			expNode: node2,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {

			gotNode := tt.index.find(tt.argID)
			assert.Equal(t, tt.expNode, gotNode)
		})
	}
}
