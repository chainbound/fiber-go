package client

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestFilter(t *testing.T) {
	a1 := common.HexToAddress("0xdc6C276D357e82C7D38D73061CEeD2e33990E5bC")
	a2 := common.HexToAddress("0x34Be5b8C30eE4fDe069DC878989686aBE9884470")
	a3 := common.HexToAddress("0x34Be5b8C30eE4fDe069DC87D989686aBE9804470")

	f := NewFilter().And(
		From(a1),
		To(a2),
	).Or(To(a3))

	fmt.Println(f.inner)
}
