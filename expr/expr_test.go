package expr

import (
	"fmt"
	"ftl.go/ctx"
	"testing"
)

func TestParse(t *testing.T) {
	l := NewLexer("a.b+c")
	yyParse(l)
	fmt.Println(l.parseResult)
}

func TestRun(t *testing.T) {
	code := "a.b"
	p, err := Compile(code)
	if err != nil {
		panic(err)
	}
	var root = struct {
		A map[string]interface{}
	}{A: map[string]interface{}{"b": 1}}
	out, err := p.Run(&ctx.Global{
		Variables: map[string]interface{}{},
		Root:      root,
	})
	if err != nil {
		panic(err)
	}
	t.Log(out)
}
