package hierarch

import (
	"encoding/json"
	"sync"
	"sync/atomic"
)

type version struct {
	v *uint64
}

func (v *version) add(delta uint64) {
	atomic.AddUint64(v.v, delta)
}

func (v *version) set(nv uint64) {
	atomic.StoreUint64(v.v, nv)
}

type Tree struct {
	sync.RWMutex
	root    *Node
	version *version
}

func NewTree() *Tree {
	var v uint64
	return &Tree{
		root:    NewNode("root"),
		version: &version{&v},
	}
}

func (tree *Tree) Version() uint64 {
	return atomic.LoadUint64(tree.version.v)
}

func (tree *Tree) SetNodeStatus(values map[string]interface{}, recursive bool, bubbleUp bool, path ...string) error {

	n, err := tree.root.ChildByPath(path...)
	if err != nil {
		return err
	}

	n.SetStatus(values, recursive)

	if bubbleUp {
		parentPath := []string{}
		for i := 0; i < len(path)-1; i++ {
			parentPath = append(parentPath, path[i])
			n, err := tree.root.ChildByPath(parentPath...)
			if err != nil {
				return err
			}

			n.SetStatus(values, false)
		}
	}

	tree.version.add(1)
	return nil
}

func (tree *Tree) NodeStatus(path ...string) (map[string]interface{}, error) {
	tree.RLock()
	defer tree.RUnlock()

	n, err := tree.root.ChildByPath(path...)
	if err != nil {
		return nil, err
	}

	n.RLock()
	defer n.RUnlock()
	tree.version.add(1)
	return n.status, nil
}

func (tree *Tree) NewNode(id string, path ...string) {
	nn := NewNode(id)
	tree.root.SetChildByPath(nn, path...)
	tree.version.add(1)
}

func (tree *Tree) Node(path ...string) (*Node, error) {
	n, err := tree.root.ChildByPath(path...)
	if err != nil {
		return nil, err
	}
	tree.version.add(1)
	return n, nil
}

func (tree *Tree) DeleteNode(path ...string) error {
	tree.root.DeletePath(path...)
	tree.version.add(1)
	// TODO: if delete did not do anything, don't change version
	return nil
}

func (tree *Tree) String() string {
	out, _ := json.MarshalIndent(tree.root, "", "   ")
	return string(out)
}
