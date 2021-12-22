package bi_tree

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestNew(t *testing.T) {
	// TODO: how do you test this?
	leafSize := LC(4)
	dataInitFn := func() Data { return new(TxRef) }

	emptyTree := New(leafSize, dataInitFn).(*tree)

	// tree
	assert.Equal(t, leafSize, emptyTree.leafSize)
	assert.Equal(t, leafSize, emptyTree.maxLC)
	assert.NotNil(t, emptyTree.dataInitFn())

	// root
	assert.NotNil(t, emptyTree.root)
	assert.Nil(t, emptyTree.root.left)
	assert.Nil(t, emptyTree.root.right)
	assert.Equal(t, dataInitFn(), emptyTree.root.data) // TODO: is this testing what it should be testing?
}

func TestTree_Insert(t *testing.T) {
	leafSize := LC(4)

	t.Run("insert single Tx", func(t *testing.T) {
		ref := generateTxRef()
		tr := newTree(leafSize)

		_ = tr.Insert(0, ref)

		assert.Equal(t, ref, *tr.root.data.(*TxRef))
	})

	t.Run("insert single Tx out of Tree range", func(t *testing.T) {
		ref := generateTxRef()
		tr := newTree(leafSize)

		_ = tr.Insert(leafSize+1, ref)

		assert.Equal(t, ref, *tr.root.data.(*TxRef))
		assert.Equal(t, ref, *tr.root.right.data.(*TxRef))
		assert.Equal(t, *new(TxRef), *tr.root.left.data.(*TxRef))
	})

	t.Run("insert multiple Tx", func(t *testing.T) {
		tr, c0, c1, c2, p0, p1, r := filledTree(leafSize)

		assert.NotEqual(t, c0, tr.root.data) // sanity check
		assert.Equal(t, c0, tr.root.left.left.data)
		assert.Equal(t, c1, tr.root.left.right.data)
		assert.Equal(t, c2, tr.root.right.left.data)
		assert.Nil(t, tr.root.right.right)

		assert.Equal(t, p0, tr.root.left.data)
		assert.Equal(t, p1, tr.root.right.data)

		assert.Equal(t, r, tr.root.data)
	})
}

func TestTree_GetRoot(t *testing.T) {
	leafSize := LC(4)
	t.Run("root data is zero", func(t *testing.T) {
		tr := newTree(leafSize)

		assert.Equal(t, *new(TxRef), *tr.root.data.(*TxRef))
	})

	t.Run("root data is zero", func(t *testing.T) {
		tr := newTree(leafSize)
		ref := generateTxRef()

		_ = tr.Insert(0, ref)

		assert.Equal(t, ref, *tr.root.data.(*TxRef))
	})

	t.Run("root after re-rooting", func(t *testing.T) {
		tr := newTree(leafSize)
		ref := generateTxRef()

		_ = tr.Insert(leafSize, ref)

		assert.Equal(t, ref, *tr.root.data.(*TxRef))
	})

	t.Run("root of many Tx", func(t *testing.T) {
		tr := newTree(4)

		allRefs := new(TxRef)
		N := leafSize * 3
		for i := LC(0); i < N; i++ {
			ref := generateTxRef()
			allRefs.xor(ref)
			_ = tr.Insert(N-i, ref)
		}

		assert.Equal(t, *allRefs, *tr.root.data.(*TxRef))
	})
}

func TestTree_GetZeroTo(t *testing.T) {
	leafSize := LC(4)
	tr, c0, _, _, p0, _, r0 := filledTree(leafSize)

	c0t, _ := tr.GetZeroTo(0 * leafSize)
	p0t, _ := tr.GetZeroTo(1 * leafSize)
	r0t, _ := tr.GetZeroTo(2 * leafSize)
	root, _ := tr.GetZeroTo(2 * tr.maxLC)

	assert.Equal(t, c0, c0t)
	assert.Equal(t, p0, p0t)
	assert.Equal(t, r0, r0t)
	assert.Equal(t, r0, root)
}

func TestTree_DropLeaves(t *testing.T) {
	leafSize := LC(4)

	t.Run("leaf root should not be dropped", func(t *testing.T) {
		tr := newTree(leafSize)

		tr.DropLeaves()

		assert.NotNil(t, tr.root)
		assert.Nil(t, tr.root.left)
		assert.Nil(t, tr.root.right)
		assert.Equal(t, leafSize, tr.leafSize)
	})

	t.Run("drop leaves 2->1", func(t *testing.T) {
		tr := newTree(leafSize)
		tr.reRoot()

		tr.DropLeaves()

		assert.NotNil(t, tr.root)
		assert.Nil(t, tr.root.left)
		assert.Nil(t, tr.root.right)
		assert.Equal(t, 2*leafSize, tr.leafSize)
	})

	t.Run("drop leaves 3->2", func(t *testing.T) {
		tr := newTree(leafSize)
		tr.reRoot()
		tr.reRoot()

		tr.DropLeaves()

		assert.NotNil(t, tr.root)
		assert.NotNil(t, tr.root.left)
		assert.Nil(t, tr.root.right)
		assert.Equal(t, 2*leafSize, tr.leafSize)
	})

	t.Run("drop leaves 3->1", func(t *testing.T) {
		tr := newTree(leafSize)
		tr.reRoot()
		tr.reRoot()

		tr.DropLeaves()
		tr.DropLeaves()

		assert.NotNil(t, tr.root)
		assert.Nil(t, tr.root.left)
		assert.Nil(t, tr.root.right)
		assert.Equal(t, 4*leafSize, tr.leafSize)
	})
}

func TestTree_reRoot(t *testing.T) {
	leafSize := LC(4)

	t.Run("single re-root", func(t *testing.T) {
		tr := newTree(leafSize)

		tr.reRoot()

		assert.NotNil(t, tr.root)
		assert.NotNil(t, tr.root.left)
		assert.Nil(t, tr.root.right)
		assert.Equal(t, tr.maxLC, 2*leafSize)
	})

	t.Run("double re-root", func(t *testing.T) {
		tr := newTree(leafSize)

		tr.reRoot()
		tr.reRoot()

		assert.NotNil(t, tr.root)
		assert.NotNil(t, tr.root.left)
		assert.NotNil(t, tr.root.left.left)
		assert.Nil(t, tr.root.right)
		assert.Nil(t, tr.root.left.right)
		assert.Equal(t, tr.maxLC, 4*leafSize)
	})
}

func newTree(leafSize LC) *tree {
	return New(leafSize, func() Data { return new(TxRef) }).(*tree)
}

func generateTxRef() TxRef {
	ref := new(TxRef)
	rand.Read(ref[:]) // math rand is not random -> test is deterministic
	return *ref
}

func filledTree(leafSize LC) (tr *tree, c0, c1, c2, p0, p1, r Data) {
	//         0 + 4 + 5 + 8                 r
	//			/         \                 / \
	//     0 + 4 + 5       8             p0     p1
	//      /    \        / \            / \   / \
	//    0    4 + 5     8  nil         c0 c1 c2 nil
	tr = newTree(leafSize)

	ref0 := generateTxRef()
	ref4 := generateTxRef()
	ref5 := generateTxRef()
	ref8 := generateTxRef()

	// node values
	c0 = ref0.Clone()
	c1 = ref4.Clone()
	_ = c1.Insert(ref5)
	c2 = ref8.Clone()
	p0 = c1.Clone()
	_ = p0.Insert(ref0)
	p1 = c2.Clone()
	r = p0.Clone()
	_ = r.Insert(ref8)

	_ = tr.Insert(0, ref0)
	_ = tr.Insert(4, ref4)
	_ = tr.Insert(8, ref8)
	_ = tr.Insert(5, ref5)
	return
}
