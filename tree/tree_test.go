package tree

import (
	"encoding/json"
	"io"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoot(t *testing.T) {

	node1 := &node[uint, int]{primary: 1}

	var tests = map[string]struct {
		tree *Tree[uint, int]
		exp  Node[uint, int]
	}{
		"nil root": {
			tree: Empty[uint, int](),
			exp:  nil,
		},
		"non-nil root": {
			tree: &Tree[uint, int]{
				root: node1,
			},
			exp: node1,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := tt.tree.Root()
			assert.Equal(t, tt.exp, got)
		})
	}
}

type addInput struct {
	nodeID   uint
	parentID uint
}

func TestAdd(t *testing.T) {

	var tests = map[string]struct {
		prep      func() *Tree[uint, int]
		add       addInput
		expAdded  bool
		expExists bool
	}{
		"primary exists": {
			prep: func() *Tree[uint, int] {
				n := &node[uint, int]{primary: 1}
				return &Tree[uint, int]{root: n, primary: &index[uint, int]{1: n}}
			},
			add:       addInput{1, 0},
			expAdded:  false,
			expExists: true,
		},
		"root is nil": {
			prep: func() *Tree[uint, int] {
				return &Tree[uint, int]{primary: &index[uint, int]{}}
			},
			add:       addInput{1, 0},
			expAdded:  true,
			expExists: false,
		},
		"re-root": {
			prep: func() *Tree[uint, int] {
				n := &node[uint, int]{primary: 1, parentID: 2}
				return &Tree[uint, int]{root: n, primary: &index[uint, int]{1: n}}
			},
			add:       addInput{2, 3},
			expAdded:  true,
			expExists: false,
		},
		"re-root with cycle": {
			prep: func() *Tree[uint, int] {
				n := &node[uint, int]{primary: 1, parentID: 2}
				return &Tree[uint, int]{root: n, primary: &index[uint, int]{1: n}}
			},
			add:       addInput{2, 1},
			expAdded:  false,
			expExists: false,
		},
		"parent does not exist": {
			prep: func() *Tree[uint, int] {
				n := &node[uint, int]{primary: 1}
				return &Tree[uint, int]{root: n, primary: &index[uint, int]{1: n}}
			},
			add:       addInput{2, 3},
			expAdded:  false,
			expExists: false,
		},
		"added": {
			prep: func() *Tree[uint, int] {
				n := &node[uint, int]{primary: 1}
				return &Tree[uint, int]{root: n, primary: &index[uint, int]{1: n}}
			},
			add:       addInput{2, 1},
			expAdded:  true,
			expExists: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tree := tt.prep()
			gotAdded, gotExists := tree.Add(tt.add.nodeID, tt.add.parentID, 0)

			assert.Equal(t, tt.expAdded, gotAdded)
			assert.Equal(t, tt.expExists, gotExists)
		})
	}
}

func TestAddResults(t *testing.T) {

	var tests = map[string]struct {
		adds   []addInput
		expBFC []uint
		expDFC []uint
	}{
		"three level parent-child": {
			adds: []addInput{
				{1, 0},
				{2, 1},
				{3, 2},
			},
			expBFC: []uint{1, 2, 3},
			expDFC: []uint{1, 2, 3},
		},
		"three level multi-children": {
			adds: []addInput{
				{1, 0},
				{2, 1},
				{3, 2},
				{4, 1},
			},
			expBFC: []uint{1, 2, 4, 3},
			expDFC: []uint{1, 2, 3, 4},
		},
		"re-root with a new subtree": {
			adds: []addInput{
				{1, 2},
				{3, 1},
				{2, 0},
				{4, 2},
			},
			expBFC: []uint{2, 1, 4, 3},
			expDFC: []uint{2, 1, 3, 4},
		},
		"failed inserts": {
			adds: []addInput{
				{1, 0},
				{2, 1},
				{3, 2},
				{2, 1},
				{4, 5},
			},
			expBFC: []uint{1, 2, 3},
			expDFC: []uint{1, 2, 3},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tree := Empty[uint, int]()
			for _, input := range tt.adds {
				tree.Add(input.nodeID, input.parentID, 0)
			}

			assert.Equal(t, tt.expBFC, bfc([]Node[uint, int]{tree.root}, []uint{}))
			assert.Equal(t, tt.expDFC, dfc(tree.root, []uint{}))

			for _, key := range tt.expBFC {
				k := tree.primary.find(key)
				if assert.NotNil(t, k, "Expected value for %d not to be nil", key) {
					assert.Equal(t, key, k.GetID())
				}
			}
		})

	}
}

func TestFind(t *testing.T) {

	var tests = map[string]struct {
		prep      func() *Tree[uint, string]
		argID     uint
		expNodeID uint
		expOK     bool
	}{
		"primary does not exist": {
			prep: func() *Tree[uint, string] {
				n := &node[uint, string]{primary: 1}
				t := &Tree[uint, string]{root: n, primary: &index[uint, string]{1: n}}
				t.Add(2, 1, "")
				return t
			},
			argID:     3,
			expNodeID: 0,
			expOK:     false,
		},
		"primary exists - branch end": {
			prep: func() *Tree[uint, string] {
				n := &node[uint, string]{primary: 1}
				t := &Tree[uint, string]{root: n, primary: &index[uint, string]{1: n}}
				t.Add(2, 1, "")
				t.Add(3, 2, "")
				t.Add(4, 1, "")
				return t
			},
			argID:     3,
			expNodeID: 3,
			expOK:     true,
		},
		"primary exists - mid tree": {
			prep: func() *Tree[uint, string] {
				n := &node[uint, string]{primary: 1}
				t := &Tree[uint, string]{root: n, primary: &index[uint, string]{1: n}}
				t.Add(2, 1, "")
				t.Add(3, 2, "")
				t.Add(4, 1, "")
				return t
			},
			argID:     2,
			expNodeID: 2,
			expOK:     true,
		},
		"primary exists - root": {
			prep: func() *Tree[uint, string] {
				n := &node[uint, string]{primary: 1}
				t := &Tree[uint, string]{root: n, primary: &index[uint, string]{1: n}}
				t.Add(2, 1, "")
				t.Add(3, 2, "")
				t.Add(4, 1, "")
				return t
			},
			argID:     1,
			expNodeID: 1,
			expOK:     true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tree := tt.prep()
			gotNode, gotOK := tree.Find(tt.argID)

			var gotNodeID uint = 0
			if gotNode != nil {
				gotNodeID = gotNode.GetID()
			}

			assert.Equal(t, tt.expOK, gotOK)
			assert.Equal(t, tt.expNodeID, gotNodeID)
		})
	}
}

func TestFindParents(t *testing.T) {

	var tests = map[string]struct {
		prep       func() *Tree[uint, string]
		argID      uint
		expNodeIDs []uint
		expOK      bool
	}{
		"primary does not exist": {
			prep: func() *Tree[uint, string] {
				n := &node[uint, string]{primary: 1}
				t := &Tree[uint, string]{root: n, primary: &index[uint, string]{1: n}}
				t.Add(2, 1, "")
				return t
			},
			argID:      3,
			expNodeIDs: []uint{},
			expOK:      false,
		},
		"primary exists - branch end": {
			prep: func() *Tree[uint, string] {
				n := &node[uint, string]{primary: 1}
				t := &Tree[uint, string]{root: n, primary: &index[uint, string]{1: n}}
				t.Add(2, 1, "")
				t.Add(3, 2, "")
				t.Add(4, 1, "")
				return t
			},
			argID:      3,
			expNodeIDs: []uint{2, 1},
			expOK:      true,
		},
		"primary exists - mid tree": {
			prep: func() *Tree[uint, string] {
				n := &node[uint, string]{primary: 1}
				t := &Tree[uint, string]{root: n, primary: &index[uint, string]{1: n}}
				t.Add(2, 1, "")
				t.Add(3, 2, "")
				t.Add(4, 1, "")
				return t
			},
			argID:      2,
			expNodeIDs: []uint{1},
			expOK:      true,
		},
		"primary exists - root": {
			prep: func() *Tree[uint, string] {
				n := &node[uint, string]{primary: 1}
				t := &Tree[uint, string]{root: n, primary: &index[uint, string]{1: n}}
				t.Add(2, 1, "")
				t.Add(3, 2, "")
				t.Add(4, 1, "")
				return t
			},
			argID:      1,
			expNodeIDs: []uint{},
			expOK:      true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tree := tt.prep()
			gotNodes, gotOK := tree.FindParents(tt.argID)

			var gotNodeIDs = make([]uint, len(gotNodes))
			if gotNodes != nil {
				for i, n := range gotNodes {
					gotNodeIDs[i] = n.GetID()
				}
			}

			assert.Equal(t, tt.expOK, gotOK)
			assert.Equal(t, tt.expNodeIDs, gotNodeIDs)
		})
	}
}

func TestMerge(t *testing.T) {

	var tests = map[string]struct {
		prepRoot  func() *Tree[uint, string]
		prepOther func() *Tree[uint, string]
		expOK     bool
		expBFC    []uint
		expDFC    []uint
	}{
		"other parent not in tree": {
			prepRoot: func() *Tree[uint, string] {
				n := &node[uint, string]{primary: 1}
				t := &Tree[uint, string]{root: n, primary: &index[uint, string]{1: n}}
				t.Add(2, 1, "")
				return t
			},
			prepOther: func() *Tree[uint, string] {
				n := &node[uint, string]{primary: 3}
				t := &Tree[uint, string]{root: n, primary: &index[uint, string]{3: n}}
				t.Add(4, 3, "")
				return t
			},
			expOK:  false,
			expBFC: []uint{1, 2},
			expDFC: []uint{1, 2},
		},
		"dulicate keys": {
			prepRoot: func() *Tree[uint, string] {
				n := &node[uint, string]{primary: 1}
				t := &Tree[uint, string]{root: n, primary: &index[uint, string]{1: n}}
				t.Add(2, 1, "")
				return t
			},
			prepOther: func() *Tree[uint, string] {
				t := Empty[uint, string]()
				t.Add(3, 1, "")
				t.Add(2, 3, "")
				return t
			},
			expOK:  false,
			expBFC: []uint{1, 2},
			expDFC: []uint{1, 2},
		},
		"merged - branch end": {
			prepRoot: func() *Tree[uint, string] {
				n := &node[uint, string]{primary: 1}
				t := &Tree[uint, string]{root: n, primary: &index[uint, string]{1: n}}
				t.Add(2, 1, "")
				t.Add(3, 2, "")
				t.Add(4, 2, "")
				t.Add(5, 1, "")
				return t
			},
			prepOther: func() *Tree[uint, string] {
				t := Empty[uint, string]()
				t.Add(6, 5, "")
				t.Add(7, 6, "")
				return t
			},
			expOK:  true,
			expBFC: []uint{1, 2, 5, 3, 4, 6, 7},
			expDFC: []uint{1, 2, 3, 4, 5, 6, 7},
		},
		"merged - mid tree": {
			prepRoot: func() *Tree[uint, string] {
				n := &node[uint, string]{primary: 1}
				t := &Tree[uint, string]{root: n, primary: &index[uint, string]{1: n}}
				t.Add(2, 1, "")
				t.Add(3, 2, "")
				t.Add(4, 2, "")
				t.Add(5, 1, "")
				return t
			},
			prepOther: func() *Tree[uint, string] {
				t := Empty[uint, string]()
				t.Add(6, 1, "")
				t.Add(7, 6, "")
				return t
			},
			expOK:  true,
			expBFC: []uint{1, 2, 5, 6, 3, 4, 7},
			expDFC: []uint{1, 2, 3, 4, 5, 6, 7},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tree := tt.prepRoot()
			other := tt.prepOther()
			gotOK := tree.Merge(other)

			assert.Equal(t, tt.expOK, gotOK)

			assert.Equal(t, tt.expBFC, bfc([]Node[uint, string]{tree.root}, []uint{}))
			assert.Equal(t, tt.expDFC, dfc(tree.root, []uint{}))

			for _, key := range tt.expBFC {
				k := tree.primary.find(key)
				if assert.NotNil(t, k, "Expceted value for %d not to be nil", key) {
					assert.Equal(t, key, k.GetID())
				}
			}
		})
	}
}

func TestSerialize(t *testing.T) {

	type Serializable struct {
		SomeData  string
		OtherData []int
	}

	type embeddedSerializable struct {
		Serializable
		ExtraString string
	}

	type CannotSerialize struct {
		Unserializable func() error
	}

	var tests = map[string]struct {
		prep      func() *Tree[uint, any]
		traversal TraversalType
		expCount  int
		expErr    error
	}{
		"empty": {
			prep:      Empty[uint, any],
			traversal: TraverseBreadthFirst,
		},
		"breadth-first": {
			prep: func() *Tree[uint, any] {

				t := Empty[uint, any]()
				t.Add(1, 0, Serializable{"valuable data", []int{1, 2, 3, 4, 5, 6, 7, 8}})
				t.Add(2, 1, map[string]string{"us": "good", "them": "bad"})
				es := embeddedSerializable{
					Serializable: Serializable{"first", []int{1}},
					ExtraString:  "second",
				}
				t.Add(3, 2, es)
				t.Add(4, 1, "Plain ol' data")
				t.Add(5, 4, 1234)
				return t
			},
			traversal: TraverseBreadthFirst,
			expCount:  5,
		},
		// "cannot serialize": {
		// 	prep: func() *Tree[uint, any] {

		// 		t := Empty[uint, any]()
		// 		t.Add(2, 1, CannotSerialize{})
		// 		return t
		// 	},
		// 	traversal: TraverseBreadthFirst,
		// 	expErr: &json.UnsupportedTypeError{
		// 		Type: reflect.TypeOf(func() error { return nil }),
		// 	},
		// },
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {

			rdr, senderr := tt.prep().Serialize(tt.traversal)

			var gotCount int = 0
			var block interface{}

			var wg sync.WaitGroup
			wg.Add(1)

			go func() {
				defer rdr.Close()
				decoder := json.NewDecoder(rdr)
				for {
					err := decoder.Decode(&block)
					if err == io.EOF {
						// end of file
						//t.Logf("deserialize - eof")
						wg.Done()
						return
					}
					if err == nil {
						//t.Logf("deserialize - success")
						// successful decoding
						gotCount = gotCount + 1
					}
					//t.Logf("deserialize - error: %s", err)
					// if decoder throws an error, we ignore it; decoding is tested
					// in a separate unit test for Deserialize
				}
			}()

			// wait for sender to finish and see if any errors occur
			gotErr := <-senderr

			// wait for reciever to finish
			wg.Wait()

			if gotErr != nil || tt.expErr != nil {
				assert.Equal(t, tt.expErr, gotErr)
			} else {
				assert.Equal(t, tt.expCount, gotCount)
			}
		})

	}
}

func TestDeserializeMap(t *testing.T) {

	type elem map[string][]int

	var tests = map[string]struct {
		prep       func() *Tree[uint, elem]
		traversal  TraversalType
		expErr     error
		expBFC     []uint
		expDFC     []uint
		dataAssert func(*testing.T, *Tree[uint, elem])
	}{
		"empty": {
			prep:      Empty[uint, elem],
			traversal: TraverseBreadthFirst,
			expBFC:    []uint{},
			expDFC:    []uint{},
		},
		"serialize breadth first": {
			prep: func() *Tree[uint, elem] {
				t := Empty[uint, elem]()
				t.Add(1, 0, elem{"one": {1, 2}, "two": {2}})
				t.Add(2, 1, elem{"two": {2}})
				t.Add(3, 2, elem{"three": {3}})
				t.Add(4, 1, elem{"four": {4}})
				t.Add(5, 4, elem{"five": {5}})
				return t
			},
			traversal: TraverseBreadthFirst,
			expErr:    nil,
			expBFC:    []uint{1, 2, 4, 3, 5},
			expDFC:    []uint{1, 2, 3, 4, 5},
			dataAssert: func(t *testing.T, gotTree *Tree[uint, elem]) {
				iter := gotTree.Traverse(TraverseBreadthFirst)
				expData := []elem{
					{"one": {1, 2}, "two": {2}},
					{"two": {2}},
					{"four": {4}},
					{"three": {3}},
					{"five": {5}},
				}
				gotData := []elem{}
				for e := range iter {
					gotData = append(gotData, e.GetData())
				}
				assert.Equal(t, expData, gotData)
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {

			// this test assumes that Serialize will throw no errors
			rdr, _ := tt.prep().Serialize(TraverseBreadthFirst)

			//t.Logf("Started to serialize")

			gotTree, gotErr := Deserialize[uint, elem](rdr)

			//t.Logf("Finished deserializing")

			assert.Equal(t, tt.expErr, gotErr)
			//t.Logf("Arguments: %+v\n", tt)
			t.Logf("Results: {tree: %+v, error %+v}\n", gotTree, gotErr)

			// only check the tree value if both expected and got errors are nil
			if gotErr == nil && tt.expErr == nil {
				assert.Equal(t, tt.expBFC, bfc([]Node[uint, elem]{gotTree.root}, []uint{}))
				assert.Equal(t, tt.expDFC, dfc(gotTree.root, []uint{}))

				if tt.dataAssert != nil {
					tt.dataAssert(t, gotTree)
				}
			}

		})
	}

}

func TestDeserializeStruct(t *testing.T) {

	type Serializable struct {
		SomeData  string
		OtherData []int
	}

	type embeddedSerializable struct {
		Serializable
		ExtraString string
	}

	embedData1 := embeddedSerializable{
		Serializable: Serializable{"valuable data", []int{1, 2, 3, 4, 5, 6, 7, 8}},
		ExtraString:  "first data point",
	}
	embedData2 := embeddedSerializable{
		Serializable: Serializable{"first", []int{37}},
		ExtraString:  "second",
	}

	var tests = map[string]struct {
		prep       func() *Tree[uint, embeddedSerializable]
		expErr     error
		expBFC     []uint
		dataAssert func(*testing.T, *Tree[uint, embeddedSerializable])
	}{
		"struct": {
			prep: func() *Tree[uint, embeddedSerializable] {

				t := Empty[uint, embeddedSerializable]()
				t.Add(1, 0, embedData1)
				t.Add(2, 1, embedData2)
				return t
			},
			expBFC: []uint{1, 2},
			dataAssert: func(t *testing.T, gotTree *Tree[uint, embeddedSerializable]) {
				iter := gotTree.Traverse(TraverseBreadthFirst)
				expData := []embeddedSerializable{embedData1, embedData2}
				gotData := []embeddedSerializable{}
				for e := range iter {
					gotData = append(gotData, e.GetData())
				}
				assert.Equal(t, expData, gotData)
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {

			prepTree := tt.prep()

			// this test assumes that Serialize will throw no errors
			rdr, _ := prepTree.Serialize(TraverseBreadthFirst)

			gotTree, gotErr := Deserialize[uint, embeddedSerializable](rdr)

			assert.Equal(t, tt.expErr, gotErr)

			// only check the tree value if both expected and got errors are nil
			if gotErr == nil && tt.expErr == nil {
				assert.Equal(t, tt.expBFC, bfc([]Node[uint, embeddedSerializable]{gotTree.root}, []uint{}))

				if tt.dataAssert != nil {
					tt.dataAssert(t, gotTree)
				}
			}

		})
	}
}
