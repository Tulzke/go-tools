package tree

import (
	"fmt"
)

// Node is the interface for a node within a tree.
//
// This package implementes this interface with an internal struct
// representation. This interface cannot be implemented by other packages
// due to unexported fields.
//
// The properties of a Node are that it has a uint primary key, a parent
// Node and an array of child Nodes. Instantiating a Node through the
// Tree.Add() function allows the node to be created with a parentID. This
// represents the primary key ID of a parent that may not be added to the
// tree yet. This value can be specified to indentify parents added to the
// tree after their children.
//
// The node also has arbitrary data attached to it; any struct may be inserted
// as the data. However, because the data is not specified, this trait contains
// no functionality to modify the data in place. The data should be implemented
// as a pointer to avoid performance ramifications.
//
// The associated data attached to a node must be serializable
// using the json encoding. If either node or arbitrary data is not, the
// Tree.Serialize() function will throw serialization errors when called.
// There is no serialization checking at the time that a Node is added to a tree
// or that its data is updated.
//
// When this node is serialized, only the id and parent id, along with
// associated data are serialized. The internal pointers to parent and children
// are recreated with the nodes are deserialized.
type Node[K comparable, T any] interface {
	// GetID returns the primary key of this node.
	GetID() K
	// GetParentID returns the primary key of this node's parent.
	GetParentID() K

	// GetChildren returns an array of pointers to all children of this node.
	GetChildren() []Node[K, T]
	// GetParent returns a pointer to the parent node of this node.
	GetParent() Node[K, T]

	// AddChildren adds a list of Nodes as children of this node.
	AddChildren(...Node[K, T])
	// ReplaceChildren replaces the current list of children with a new list of
	// Nodes.
	ReplaceChildren(...Node[K, T])

	setParent(n Node[K, T])

	// GetData retruns this node's internal data.
	GetData() T
	// SetData replaces this nodes data with the argument. The argument may be any
	// type, but must be serializable via json.
	//
	// This function does not attempt to test json encoding when the data is set;
	// any error with encoding will only occur when the data is serialized
	// to a repository.
	SetData(T)
}

type node[K comparable, T any] struct {
	primary  K
	parentID K
	parent   Node[K, T]
	data     T
	children []Node[K, T]
}

func (n *node[K, T]) GetID() K {
	return n.primary
}

func (n *node[K, T]) GetParentID() K {
	return n.parentID
}

func (n *node[K, T]) GetChildren() []Node[K, T] {
	return n.children
}

func (n *node[K, T]) GetParent() Node[K, T] {
	return n.parent
}

func (n *node[K, T]) AddChildren(children ...Node[K, T]) {
	if n.children == nil {
		n.children = []Node[K, T]{}
	}
	n.children = append(n.children, children[:]...)
}

func (n *node[K, T]) ReplaceChildren(children ...Node[K, T]) {
	n.children = []Node[K, T]{}
	n.AddChildren(children...)
}

func (n *node[K, T]) setParent(parent Node[K, T]) {
	if parent == nil || parent.GetID() == n.GetID() {
		return
	}
	n.parent = parent
	n.parentID = parent.GetID()

}

func (n *node[K, T]) GetData() T {
	return n.data
}

func (n *node[K, T]) SetData(newdata T) {
	n.data = newdata
}

func (n *node[K, T]) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		fmt.Fprintf(f, "{primary: %v parentID: %v data:%+v children:[", n.primary, n.parentID, n.data)
		for i, n := range n.children {
			if i != 0 {
				fmt.Fprint(f, " ")
			}
			fmt.Fprintf(f, "%v", n.GetID())
		}
		fmt.Fprint(f, "]}")
	}
}
