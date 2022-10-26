package filter

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func traverse(n *Node) {
	if n.Operand != nil {
		fmt.Println("Operand:", n.Operand)
	} else {
		fmt.Println("Operator:", n.Operator)
	}

	fmt.Println("Children:", len(n.Nodes))
	for _, child := range n.Nodes {
		traverse(child)
	}
}

func TestFilter(t *testing.T) {
	a1 := common.HexToAddress("0xdc6C276D357e82C7D38D73061CEeD2e33990E5bC")
	// a2 := common.HexToAddress("0x34Be5b8C30eE4fDe069DC878989686aBE9884470")
	a3 := common.HexToAddress("0x34Be5b8C30eE4fDe069DC87D989686aBE98abcde")

	f := NewFilter(
		And(
			To(a1),
			From(a3),
		),
	)

	traverse(f.Root)

	v, err := json.MarshalIndent(f, "", "    ")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(v))
}
