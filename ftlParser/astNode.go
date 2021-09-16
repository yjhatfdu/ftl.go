package ftlParser

import (
	"fmt"
	"ftl.go/ctx"
	"io"
)

type astNode struct {
	T         int
	V         []string
	C         []*astNode
	processor func(w io.Writer, global *ctx.Global) error
}

func (n *astNode) compile() error {
	for _, c := range n.C {
		if err := c.compile(); err != nil {
			return err
		}
	}
	df := directiveFactoryMap[n.T]
	if df == nil {
		return fmt.Errorf("directive %v not defined", n.T)
	}
	var err error
	n.processor, err = df(n.V, n.C)
	return err
}
