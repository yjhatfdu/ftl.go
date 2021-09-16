package ftlParser

import (
	"bytes"
	"fmt"
	"ftl.go/ctx"
	"io"
)

type Renderer struct {
	ast *astNode
}

func (r *Renderer) Render(root interface{}, wr io.Writer) error {
	g := &ctx.Global{
		Root: root,
	}
	return r.ast.processor(wr, g)
}

func (r *Renderer) RenderString(root interface{}) (string, error) {
	bw := new(bytes.Buffer)
	err := r.Render(root, bw)
	return bw.String(), err
}

func Compile(code string) (r *Renderer, err error) {
	defer func() {
		e := recover()
		if e != nil {
			err = fmt.Errorf("%v", err)
		}
	}()
	yyErrorVerbose = true
	l := newLex(code)
	yyParse(l)
	ast := l.parseResult
	err = ast.compile()
	if err != nil {
		return
	}
	return &Renderer{ast: ast}, nil
}
