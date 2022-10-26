package filter

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Operator = int

const (
	AND Operator = 1
	OR  Operator = 2
)

type Node struct {
	Operand  *FilterKV
	Operator Operator
	Nodes    []*Node
}

type Filter struct {
	Root *Node
}

func NewFilter(ops ...FilterOp) *Filter {
	f := &Filter{}

	for _, op := range ops {
		op(f, nil)
	}

	return f
}

func (f Filter) Encode() []byte {
	e, _ := json.Marshal(f)
	return e
}

type FilterKV struct {
	Key   string `json:"key,omitempty"`
	Value []byte `json:"value,omitempty"`
}

func (kv FilterKV) String() string {
	return fmt.Sprintf("{ %s : %s }", kv.Key, string(kv.Value))
}

type FilterOp func(*Filter, *Node)

// Operators
func And(opts ...FilterOp) FilterOp {
	return func(f *Filter, n *Node) {
		// If we're at the first node (empty)
		var new *Node
		if n == nil {
			new = &Node{
				Operator: AND,
			}

			f.Root = new
		} else {
			new = &Node{
				Operator: AND,
			}

			n.Nodes = append(n.Nodes, new)
		}

		for _, op := range opts {
			op(f, new)
		}
	}
}

func Or(opts ...FilterOp) FilterOp {
	return func(f *Filter, n *Node) {
		// If we're at the first node (empty)
		var new *Node
		if n == nil {
			new = &Node{
				Operator: OR,
			}

			f.Root = new
		} else {
			new = &Node{
				Operator: OR,
			}

			n.Nodes = append(n.Nodes, new)
		}

		for _, op := range opts {
			op(f, new)
		}
	}
}

// Operands
func To(to common.Address) FilterOp {
	return func(f *Filter, n *Node) {
		var new *Node
		if n == nil {
			new = &Node{
				Operand: &FilterKV{"to", to.Bytes()},
			}

			f.Root = new
		} else {
			n.Nodes = append(n.Nodes, &Node{
				Operand: &FilterKV{"to", to.Bytes()},
			})
		}
	}
}

func From(from common.Address) FilterOp {
	return func(f *Filter, n *Node) {
		var new *Node
		if n == nil {
			new = &Node{
				Operand: &FilterKV{"from", from.Bytes()},
			}

			f.Root = new
		} else {
			n.Nodes = append(n.Nodes, &Node{
				Operand: &FilterKV{"from", from.Bytes()},
			})
		}
	}
}

func MethodID(id []byte) FilterOp {
	return func(f *Filter, n *Node) {
		var new *Node
		if n == nil {
			new = &Node{
				Operand: &FilterKV{"method", id},
			}

			f.Root = new
		} else {
			n.Nodes = append(n.Nodes, &Node{
				Operand: &FilterKV{"method", id},
			})
		}
	}
}

func Value(v *big.Int) FilterOp {
	return func(f *Filter, n *Node) {
		var new *Node
		if n == nil {
			new = &Node{
				Operand: &FilterKV{"value", v.Bytes()},
			}

			f.Root = new
		} else {
			n.Nodes = append(n.Nodes, &Node{
				Operand: &FilterKV{"value", v.Bytes()},
			})
		}
	}
}
