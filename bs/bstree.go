// Copyright (c) 2019, Yu Wu <yu.771991@gmail.com> All rights reserved.
//
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package bs

import (
	"github.com/yu31/gostructs/container"
	"github.com/yu31/gostructs/internal/tree"
)

// treeNode is used for Binary Search Tree.
//
// And it is also the implementation of interface container.Element and container.TreeNode
type treeNode struct {
	key   container.Key
	value container.Value
	left  *treeNode
	right *treeNode
}

// Key returns the key.
func (n *treeNode) Key() container.Key {
	return n.key
}

// Value returns the value.
func (n *treeNode) Value() container.Value {
	return n.value
}

// Left returns the left child of the TreeNode.
func (n *treeNode) Left() container.TreeNode {
	return n.left
}

// Right returns the right child of the TreeNode.
func (n *treeNode) Right() container.TreeNode {
	return n.right
}

// Tree implements the Binary Search Tree.
//
// And it is also the implementation of interface container.Container
type Tree struct {
	root *treeNode
	len  int
}

// New creates an Binary Search Tree.
func New() *Tree {
	return &Tree{
		root: nil,
		len:  0,
	}
}

// Root returns the root node of the tree.
func (tr *Tree) Root() container.TreeNode {
	return tr.root
}

// Len returns the number of elements.
func (tr *Tree) Len() int {
	return tr.len
}

// Insert inserts and returns an Element with given key and value if key doesn't exists.
// Or else, returns the existing Element for the key if present.
// The bool result is true if an Element was inserted, false if searched.
func (tr *Tree) Insert(k container.Key, v container.Value) (container.Element, bool) {
	node, _, ok := tr.insertOrSearch(k, v)
	return node, ok
}

// Delete removes and returns the Element of a given key.
// Returns nil if not found.
func (tr *Tree) Delete(k container.Key) container.Element {
	node, parent := tr.searchNode(k)
	if node == nil {
		return nil
	}
	node = tr.deleteNode(node, parent)
	return node
}

// Update updates an Element with the given key and value, And returns the old element.
// Returns nil if the key not be found.
func (tr *Tree) Update(k container.Key, v container.Value) container.Element {
	node, parent := tr.searchNode(k)
	if node != nil {
		tr.updateNode(node, parent, k, v)
	}
	return node
}

// Replace inserts or updates an Element by giving key and value.
// The bool result is true if an Element was inserted, false if an Element was updated.
//
// The operation are same as the Insert method if key not found,
// And are same as the Update method if key exists.
func (tr *Tree) Replace(k container.Key, v container.Value) (container.Element, bool) {
	node, parent, ok := tr.insertOrSearch(k, v)
	if !ok {
		tr.updateNode(node, parent, k, v)
	}
	return node, ok
}

// Search searches the Element of a given key.
// Returns nil if key not found.
func (tr *Tree) Search(k container.Key) container.Element {
	node, _ := tr.searchNode(k)
	return node
}

// Iter return an Iterator, it's a wrap for bs.Iterator.
func (tr *Tree) Iter(start container.Key, boundary container.Key) container.Iterator {
	return tree.NewIterator(tr.root, start, boundary)
}

// Range calls f sequentially each TreeNode present in the Tree.
// If f returns false, range stops the iteration.
func (tr *Tree) Range(start container.Key, boundary container.Key, f func(ele container.Element) bool) {
	tree.Range(tr.root, start, boundary, func(node container.TreeNode) bool {
		return f(node)
	})
}

// LastLT searches for the last node that less than the key.
func (tr *Tree) LastLT(k container.Key) container.Element {
	return tree.LastLT(tr.root, k)
}

// LastLE search for the last node that less than or equal to the key.
func (tr *Tree) LastLE(k container.Key) container.Element {
	return tree.LastLE(tr.root, k)
}

// FirstGT search for the first node that greater than to the key.
func (tr *Tree) FirstGT(k container.Key) container.Element {
	return tree.FirstGT(tr.root, k)
}

// FirstGE search for the first node that greater than or equal to the key.
func (tr *Tree) FirstGE(k container.Key) container.Element {
	return tree.FirstGE(tr.root, k)
}

// The insertOrSearch inserts and returns a new node with given key and value if key not exists.
// Or else, returns the exists node and its parent node for the key if present.
// The ok result is true if the node was inserted, false if searched.
func (tr *Tree) insertOrSearch(k container.Key, v container.Value) (node *treeNode, parent *treeNode, ok bool) {
	node = tr.root
	for node != nil {
		cmp := k.Compare(node.key)
		if cmp == 0 {
			// Found the exists key, returns it
			return
		}

		parent = node // The parent node of n.

		if cmp == -1 {
			if node.left == nil {
				node.left = tr.createNode(k, v)
				node = node.left
				break
			}
			node = node.left
		} else {
			if node.right == nil {
				node.right = tr.createNode(k, v)
				node = node.right
				break
			}
			node = node.right
		}
	}

	if node == nil {
		node = tr.createNode(k, v)
		tr.root = node
	}

	tr.len++
	ok = true
	return
}

// Helps to creates an tree node with given key and value.
func (tr *Tree) createNode(k container.Key, v container.Value) *treeNode {
	return &treeNode{
		key:   k,
		value: v,
		left:  nil,
		right: nil,
	}
}

// Helps to deletes the node, returns the node that actually deleted.
func (tr *Tree) deleteNode(node *treeNode, parent *treeNode) (d *treeNode) {
	d = node
	if d.left != nil && d.right != nil {
		xx := d
		x := d.left
		for x.right != nil {
			xx = x
			x = x.right
		}

		tr.swap(d, x)
		parent = xx
		d = x
	}

	var c *treeNode
	if d.left != nil {
		c = d.left
	} else {
		c = d.right
	}

	if parent == nil {
		tr.root = c
	} else if parent.left == d {
		parent.left = c
	} else {
		parent.right = c
	}

	// reset the unused field.
	d.left = nil
	d.right = nil

	tr.len--
	return d
}

// Help to creates a new tree node and instead of the node.
func (tr *Tree) updateNode(node *treeNode, parent *treeNode, k container.Key, v container.Value) {
	n0 := tr.createNode(k, v)
	n0.left = node.left
	n0.right = node.right

	if parent == nil {
		tr.root = n0
	} else if parent.left == node {
		parent.left = n0
	} else {
		parent.right = n0
	}

	// reset the unused field.
	node.left = nil
	node.right = nil
}

// Searches the node and its parent node of a given key.
func (tr *Tree) searchNode(k container.Key) (node *treeNode, parent *treeNode) {
	node = tr.root
	for node != nil {
		cmp := k.Compare(node.key)
		if cmp == 0 {
			// Found the node of key.
			return
		}

		parent = node // The parent node of n.

		if cmp == -1 {
			node = node.left
		} else {
			node = node.right
		}
	}
	return
}

func (tr *Tree) swap(n1, n2 *treeNode) {
	n1.key, n2.key = n2.key, n1.key
	n1.value, n2.value = n2.value, n1.value
}
