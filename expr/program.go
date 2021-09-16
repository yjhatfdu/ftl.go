package expr

import (
	"errors"
	"fmt"
	"ftl.go/ctx"
	"ftl.go/functions"
	"reflect"
	"strconv"
	"sync"
)

type op struct {
	t    int
	v    interface{}
	argc int
}

type Program struct {
	code []op
}

type buildCtx struct {
	code []op
}

func (c *buildCtx) append(code op) {
	c.code = append(c.code, code)
}

func (n *AstNode) compile(ctx *buildCtx) error {
	for _, c := range n.Children {
		if err := c.compile(ctx); err != nil {
			return err
		}
	}
	c := op{
		t:    n.NodeType,
		v:    n.Value,
		argc: len(n.Children),
	}
	var err error
	if n.NodeType == CONST {
		switch n.ValueType {
		case INT:
			c.v, err = strconv.ParseInt(n.Value, 10, 64)
		case FLOAT:
			c.v, err = strconv.ParseFloat(n.Value, 64)
		case BOOL:
			c.v, err = strconv.ParseBool(n.Value)
		case STR:
			c.v = n.Value
		}
	}
	ctx.append(c)
	return err
}

func compile(node *AstNode) (*Program, error) {
	c := &buildCtx{make([]op, 0)}
	err := node.compile(c)
	if err != nil {
		return nil, err
	}
	return &Program{code: c.code}, nil
}

func Compile(code string) (p *Program, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	l := NewLexer(code)
	yyErrorVerbose = true
	yyParse(l)
	ast := l.parseResult
	return compile(ast)
}

var stackPool = sync.Pool{New: func() interface{} {
	return &Stack{
		s: make([]interface{}, 4),
	}
}}

func NewStack() *Stack {
	return stackPool.Get().(*Stack)
}

type Stack struct {
	s []interface{}
}

func (s *Stack) Close() {
	s.s = s.s[0:0]
	stackPool.Put(s)
}
func (s *Stack) Push(item interface{}) {
	s.s = append(s.s, item)
}

func (s *Stack) Pop() interface{} {
	l := len(s.s) - 1
	item := s.s[l]
	s.s = s.s[0:l]
	return item
}
func (s *Stack) PopN(n int) []interface{} {
	// 不用处理n>len(s.s)，不应该发生
	l := len(s.s) - n
	items := s.s[l:len(s.s)]
	s.s = s.s[0:l]
	return items
}

var errorType = reflect.TypeOf((*error)(nil)).Elem()

func (p *Program) Run(g *ctx.Global) (interface{}, error) {
	stack := NewStack()
	defer stack.Close()
	for _, op := range p.code {
		switch op.t {
		case CONST:
			stack.Push(op.v)
		case VAR:
			if v, err := g.Get(op.v.(string)); err != nil {
				return nil, err
			} else {
				stack.Push(v)
			}
		case ATTR:
			val := stack.Pop()
			if v, err := ctx.GetField(val, op.v.(string)); err != nil {
				return nil, err
			} else {
				stack.Push(v)
			}
		case FUNC:
			if f, ok := functions.FuncMap[op.v.(string)]; !ok {
				return nil, fmt.Errorf("function %s not defined", op.v.(string))
			} else {
				args := stack.PopN(op.argc)
				if ret, err := f(args...); err != nil {
					return nil, err
				} else {
					stack.Push(ret)
				}
			}
		case CALL:
			args := stack.PopN(op.argc - 1)
			fn := reflect.ValueOf(stack.Pop())
			if !fn.IsValid() {
				stack.Push(nil)
				break
			}
			argsValue := make([]reflect.Value, len(args))
			for i := range args {
				argsValue[i] = reflect.ValueOf(args[i])
			}
			var out []reflect.Value
			if fn.Type().IsVariadic() {
				out = fn.CallSlice(argsValue)
			} else {
				out = fn.Call(argsValue)
			}
			if len(out) == 0 {
				return nil, errors.New("method should have at least one argument")
			}
			if len(out) == 2 && out[1].Type().Implements(errorType) {
				err := out[1].Interface().(error)
				return nil, err
			}
			stack.Push(out[0].Interface())
		}
	}
	return stack.Pop(), nil
}
