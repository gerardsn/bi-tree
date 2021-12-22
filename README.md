
```golang
// Create a new tree with page size 4
leafSize := LC(4)
tree := New(leafSize, NewTxRef)
// Insert a transaction reference at the specified clock value.
err := tree.Insert(clock, ref)
// GetRoot returns the accumulated data for the entire Tree
data := tree.GetRoot()
// GetZeroTo Data for LC-range [0, ceil(clock/leafSize)*leafSize)
data, err := tree.GetZeroTo(clock)
// DropLeaves shrinks the Tree by dropping all leaves unless root is a leaf. The parent of a leaf will become the new leaf
DropLeaves()
```

![alt text](https://github.com/gerardsn/bi-tree/blob/master/tree.png?raw=true)
