package filter

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
)

func traverse(n *Node) {
	if n.Operand != nil {
		fmt.Println("Operand:", n.Operand)
	} else {
		fmt.Println("Operator:", n.Operator)
	}

	fmt.Println("Children:", len(n.Children))
	for _, child := range n.Children {
		traverse(child)
	}
}

func TestAndFilter(t *testing.T) {
	f := New(
		And(
			To("0xdc6C276D357e82C7D38D73061CEeD2e33990E5bC"),
			From("0x34Be5b8C30eE4fDe069DC87D989686aBE98abcde"),
		),
	)

	traverse(f.Root)

	v, err := json.MarshalIndent(f, "", "    ")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(v))
}

func TestOrFilter(t *testing.T) {
	f := New(
		Or(
			To("0xdc6C276D357e82C7D38D73061CEeD2e33990E5bC"),
			To("0x34Be5b8C30eE4fDe069DC87D989686aBE98abcde"),
		),
	)

	traverse(f.Root)

	v, err := json.MarshalIndent(f, "", "    ")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(v))
}

func TestValueEqFilter(t *testing.T) {
	f := New(ValueEq(big.NewInt(10000)))

	traverse(f.Root)

	v, err := json.MarshalIndent(f, "", "    ")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(v))
}

func TestValueLteFilter(t *testing.T) {
	f := New(ValueLte(big.NewInt(10000)))

	traverse(f.Root)

	v, err := json.MarshalIndent(f, "", "    ")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(v))
}

func TestValueGteFilter(t *testing.T) {
	f := New(ValueGte(big.NewInt(10000)))

	traverse(f.Root)

	v, err := json.MarshalIndent(f, "", "    ")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(v))
}
