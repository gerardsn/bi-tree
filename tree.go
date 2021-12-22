package bi_tree

import "fmt"

// LC is the Lamport Clock value
type LC uint32

type Tree interface {
	// Insert a transaction reference at the specified clock value.
	Insert(clock LC, ref TxRef) error
	// GetRoot returns the accumulated data for the entire Tree
	GetRoot() Data
	// GetZeroTo Data for LC-range [0, ceil(clock/leafSize)*leafSize)
	GetZeroTo(clock LC) (Data, error)
	// DropLeaves shrinks the Tree by dropping all leaves. The parent of a leaf will become the new leaf
	DropLeaves()
	// LeafSize returns the size of the current leaves
	LeafSize() LC
}

// tree creates a binary tree, where the leaves contain Data over a fixed range of Lamport Clock (LC) values.
// The Data of the parent is the sum of that of its children. The root contains the sum of all Data in the tree.
// Since the leaves are of fixed size, a new root is created when added something to an LC outside of the current root range.
// Whenever a new branch is created, a string of left nodes is created all the way to the leaf.
// TODO: There is some redundancy in storing the left & right leaf + their parent. Dropping left leafs saves ~25% of memory.
type tree struct {
	root     *node
	maxLC    LC
	leafSize LC
	// dataInitFn produces Data with a properly initialized internal data structure.
	dataInitFn func() Data
}

func New(leafSize LC, dataInitFn func() Data) Tree {
	return &tree{
		root:       newNode(leafSize, leafSize, dataInitFn()),
		maxLC:      leafSize,
		leafSize:   leafSize,
		dataInitFn: dataInitFn,
	}
}

func (t *tree) newBranch(start, stop LC) *node {
	split := (stop + start) / 2
	n := newNode(split, stop, t.dataInitFn())
	if stop-start > t.leafSize {
		n.left = t.newBranch(start, split)
	}
	return n
}

func (t *tree) Insert(clock LC, ref TxRef) error {
	// grow tree if needed
	for clock >= t.maxLC {
		t.reRoot()
	}

	// insert ref in all nodes from root to leave
	next := t.root
	for next != nil {
		err := next.data.Insert(ref)
		if err != nil {
			return fmt.Errorf("insert failed for node with split %d: %w", next.split, err)
		}
		next = t.getNextNode(next, clock)
	}

	return nil
}

// reRoot creates a new root with data from the current root and adds current root as its left branch.
func (t *tree) reRoot() {
	newRoot := newNode(t.maxLC, 2*t.maxLC, t.root.data.Clone())
	newRoot.left = t.root
	t.root = newRoot
	t.maxLC *= 2
}

// getNextNode retrieves the next node based on the clock value. If the node does not exist it is created.
func (t *tree) getNextNode(current *node, clock LC) *node {
	// return nil if current is a leaf
	if current.left == nil {
		return nil
	}

	if clock < current.split {
		return current.left
	} else {
		if current.right == nil {
			current.right = t.newBranch(current.split, current.max)
		}
		return current.right
	}
}

func (t *tree) GetZeroTo(clock LC) (Data, error) {
	data := t.root.data.Clone()
	next := t.root
	for next != nil {
		if clock < next.split {
			if next.right != nil {
				// TODO: should this error be checked?
				// Only fails when data structures do not match, which should not happen for data managed by the tree.
				_ = data.Subtract(next.right.data)
			}
			next = next.left
		} else {
			next = next.right
		}
	}
	return data, nil
}

func (t *tree) GetRoot() Data {
	return t.root.data.Clone()
}

func (t *tree) LeafSize() LC {
	return t.leafSize
}

func (t *tree) DropLeaves() {
	if t.root.left != nil {
		dropLeaves(t.root)
		t.leafSize *= 2
	}
}

func dropLeaves(current *node) {
	// Nothing to do if a current.right was nil || should not drop leaf if it is the root.
	if current == nil || current.left == nil {
		return
	}
	// if current.left is a leaf, make current node a leaf
	if current.left.left == nil {
		current.left = nil
		current.right = nil
		return
	}
	dropLeaves(current.left)
	dropLeaves(current.right)
}

// node
type node struct {
	// split point for left / right node
	split LC
	// max clock value for the node
	max LC
	// data held by the node. Should not be nil
	data Data

	// child nodes. if left == nil, current node is a leaf
	left  *node
	right *node
}

func newNode(split, max LC, data Data) *node {
	return &node{
		split: split,
		max:   max,
		data:  data,
	}
}
