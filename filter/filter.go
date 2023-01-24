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
	Operand  *FilterKV `json:"Operand,omitempty"`
	Operator Operator  `json:"Operator,omitempty"`
	Children []*Node   `json:"Children,omitempty"`
}

type Filter struct {
	Root *Node
}

func New(rootOp FilterOp) *Filter {
	f := &Filter{}

	rootOp(f, nil)

	return f
}

func (f Filter) Encode() []byte {
	e, _ := json.Marshal(f)
	return e
}

type FilterKV struct {
	Key   string `json:"Key,omitempty"`
	Value []byte `json:"Value,omitempty"`
}

func (kv FilterKV) String() string {
	return fmt.Sprintf("{ %s : %s }", kv.Key, string(kv.Value))
}

type FilterOp func(*Filter, *Node)

// Operators
func And(ops ...FilterOp) FilterOp {
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

			n.Children = append(n.Children, new)
		}

		for _, op := range ops {
			op(f, new)
		}
	}
}

func Or(ops ...FilterOp) FilterOp {
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

			n.Children = append(n.Children, new)
		}

		for _, op := range ops {
			op(f, new)
		}
	}
}

// Operands
func To(to string) FilterOp {
	return func(f *Filter, n *Node) {
		var new *Node
		if n == nil {
			new = &Node{
				Operand: &FilterKV{"to", common.HexToAddress(to).Bytes()},
			}

			f.Root = new
		} else {
			n.Children = append(n.Children, &Node{
				Operand: &FilterKV{"to", common.HexToAddress(to).Bytes()},
			})
		}
	}
}

func From(from string) FilterOp {
	return func(f *Filter, n *Node) {
		var new *Node
		if n == nil {
			new = &Node{
				Operand: &FilterKV{"from", common.HexToAddress(from).Bytes()},
			}

			f.Root = new
		} else {
			n.Children = append(n.Children, &Node{
				Operand: &FilterKV{"from", common.HexToAddress(from).Bytes()},
			})
		}
	}
}

func MethodID(id string) FilterOp {
	return func(f *Filter, n *Node) {
		var new *Node
		if n == nil {
			new = &Node{
				Operand: &FilterKV{"method", common.FromHex(id)},
			}

			f.Root = new
		} else {
			n.Children = append(n.Children, &Node{
				Operand: &FilterKV{"method", common.FromHex(id)},
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
			n.Children = append(n.Children, &Node{
				Operand: &FilterKV{"value", v.Bytes()},
			})
		}
	}
}
