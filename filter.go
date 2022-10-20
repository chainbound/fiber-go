package client

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

type Filter struct {
	inner string
}

type FilterField struct {
	Key   string
	Value string
}

type FilterOp func(*Filter)

func NewFilter() *Filter {
	return &Filter{
		inner: "",
	}
}

func (f *Filter) And(ops ...FilterOp) *Filter {
	f.inner += "("
	for i, op := range ops {
		op(f)
		if i < len(ops)-1 {
			f.inner += " AND "
		}
	}
	f.inner += ")"

	return f
}

func (f *Filter) Or(ops ...FilterOp) *Filter {
	f.inner += "("
	for i, op := range ops {
		op(f)
		if i < len(ops)-1 {
			f.inner += " OR "
		}
	}
	f.inner += ")"

	return f
}

func To(to common.Address) FilterOp {
	return func(f *Filter) {
		f.inner += fmt.Sprintf("to:%s", to)
	}
}

func From(from common.Address) FilterOp {
	return func(f *Filter) {
		f.inner += fmt.Sprintf("from:%s", from)
	}
}
