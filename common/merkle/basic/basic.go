// Copyright © 2018 Zhao Ming <mint.zhao.chiu@gmail.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package basic

import (
	"bytes"
	"fmt"

	"github.com/mintzhao/topachain/common/crypto/hasher"
)

// ErrEmptyValues indicated that merkletree init with empty values
type ErrEmptyValues struct{}

func (err *ErrEmptyValues) Error() string {
	return fmt.Sprint("empty values")
}

// ErrInvalidValueType indicated the received value isn't Value interface
type ErrInvalidValueType struct{}

func (err *ErrInvalidValueType) Error() string {
	return fmt.Sprint("invalid value type")
}

// Value represents the data that is stored in the tree.
// Any type which implements the interface can stored in the MerkleTree
type Value interface {
	Hash() ([]byte, error)
	Equals(other interface{}) bool
}

type node struct {
	parent *node
	left   *node
	right  *node
	value  Value
	hash   []byte
	isleaf bool
}

func (n *node) verify(h hasher.Hasher) ([]byte, error) {
	if n.isleaf {
		return n.value.Hash()
	}

	rhash, err := n.right.verify(h)
	if err != nil {
		return nil, err
	}

	lhash, err := n.left.verify(h)
	if err != nil {
		return nil, err
	}

	return h.Hash(append(lhash, rhash...))
}

func (n *node) nodeHash(h hasher.Hasher) ([]byte, error) {
	if n.isleaf {
		return n.value.Hash()
	}

	return h.Hash(append(n.left.hash, n.right.hash...))
}

type BasicMerkletree struct {
	root  *node
	leafs []*node
	h     hasher.Hasher
}

func New(vals []Value, h hasher.Hasher) (*BasicMerkletree, error) {
	tree := &BasicMerkletree{
		leafs: make([]*node, 0),
		h:     h,
	}

	valsInterface := make([]interface{}, len(vals))
	for k, v := range vals {
		valsInterface[k] = v
	}

	if err := tree.Reconstruct(valsInterface...); err != nil {
		return nil, err
	}

	return tree, nil
}

func initValues(vals []Value) ([]*node, error) {
	if len(vals) == 0 {
		return nil, &ErrEmptyValues{}
	}

	leafs := make([]*node, 0)
	for _, val := range vals {
		hash, err := val.Hash()
		if err != nil {
			return nil, err
		}

		leafs = append(leafs, &node{
			hash:   hash,
			value:  val,
			isleaf: true,
		})
	}
	if len(leafs)%2 == 1 {
		duplicate := &node{
			hash:   leafs[len(leafs)-1].hash,
			value:  leafs[len(leafs)-1].value,
			isleaf: true,
		}
		leafs = append(leafs, duplicate)
	}

	return leafs, nil
}

func constructTree(leafs []*node, h hasher.Hasher) (*node, error) {
	nodes := make([]*node, 0)
	for i := 0; i < len(leafs); i += 2 {
		left, right := i, i+1
		if i+1 == len(leafs) {
			right = i
		}
		lrhashBytes := append(leafs[left].hash, leafs[right].hash...)
		lrhash, err := h.Hash(lrhashBytes)
		if err != nil {
			return nil, err
		}
		n := &node{
			left:  leafs[left],
			right: leafs[right],
			hash:  lrhash,
		}

		leafs[left].parent = n
		leafs[right].parent = n
		nodes = append(nodes, n)

		if len(leafs) == 2 {
			return n, nil
		}
	}

	return constructTree(nodes, h)
}

// Root return merkletree's root hash
func (t *BasicMerkletree) Root() []byte {
	return t.root.hash
}

// Reconstruct construct merkletree again
func (t *BasicMerkletree) Reconstruct(valsInterface ...interface{}) error {
	vals := make([]Value, 0)

	if len(valsInterface) == 0 {
		vals = t.values()
	} else {
		for _, v := range valsInterface {
			val, ok := v.(Value)
			if !ok {
				return &ErrInvalidValueType{}
			}

			vals = append(vals, val)
		}
	}

	leafs, err := initValues(vals)
	if err != nil {
		return err
	}

	root, err := constructTree(leafs, t.h)
	if err != nil {
		return err
	}

	t.leafs = leafs
	t.root = root

	return nil
}

func (t *BasicMerkletree) values() []Value {
	vals := make([]Value, 0)
	for _, leaf := range t.leafs {
		vals = append(vals, leaf.value)
	}

	return vals
}

func (t *BasicMerkletree) VerifyRoot() bool {
	calcroot, err := t.root.verify(t.h)
	if err != nil {
		return false
	}

	return bytes.Equal(t.Root(), calcroot)
}

func (t *BasicMerkletree) VerifyValue(exceptedRoot []byte, valInterface interface{}) bool {
	if !bytes.Equal(exceptedRoot, t.Root()) {
		return false
	}

	val, ok := valInterface.(Value)
	if !ok {
		return false
	}

	for _, l := range t.leafs {
		if !l.value.Equals(val) {
			continue
		}

		parent := l.parent
		for parent != nil {
			rhash, err := parent.right.nodeHash(t.h)
			if err != nil {
				return false
			}

			lhash, err := parent.left.nodeHash(t.h)
			if err != nil {
				return false
			}

			phash, err := t.h.Hash(append(lhash, rhash...))
			if err != nil {
				return false
			}

			if !bytes.Equal(phash, parent.hash) {
				return false
			}

			parent = parent.parent
		}

		return true
	}

	return false
}
