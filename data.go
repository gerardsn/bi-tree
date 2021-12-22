package bi_tree

import (
	"encoding/hex"
	"errors"
)

type Data interface {
	// Insert a transaction reference in place
	Insert(ref TxRef) error
	// Subtract data in place. Returns an error when the internal data structures do not match
	Subtract(data Data) error
	// Clone returns a copy of the data
	Clone() Data
}

type TxRef [32]byte

func NewTxRef() Data {
	return new(TxRef)
}

func (r *TxRef) Insert(ref TxRef) error {
	r.xor(ref)
	return nil
}

func (r *TxRef) xor(ref TxRef) {
	for i := range r {
		r[i] ^= ref[i]
	}
}

func (r *TxRef) Clone() Data {
	c := new(TxRef)
	copy(c[:], r[:])
	return c
}

func (r *TxRef) Subtract(data Data) error {
	ref, ok := data.(*TxRef)
	if !ok {
		return errors.New("internal data structures do not match")
	}
	r.xor(*ref)
	return nil
}

func (r TxRef) String() string {
	return hex.EncodeToString(r[:])
}
