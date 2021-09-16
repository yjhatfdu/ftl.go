package ftlParser

import (
	"errors"
	"fmt"
	"ftl.go/ctx"
	"ftl.go/expr"
	"ftl.go/utils"
	"io"
	"reflect"
	"strings"
)

type processor func(w io.Writer, global *ctx.Global) error

type directiveFactory func(values []string, children []*astNode) (processor, error)

var directiveFactoryMap = map[int]directiveFactory{
	GROUP: groupDirective,
	TEXT:  textDirective,
	IF:    ifDirective,
	TMPL:  tmplDirective,
	LIST:  listDirective,
}

func textDirective(values []string, children []*astNode) (processor, error) {
	return func(w io.Writer, global *ctx.Global) error {
		_, _ = w.Write([]byte(values[0]))
		return nil
	}, nil
}

func groupDirective(values []string, children []*astNode) (processor, error) {
	return func(w io.Writer, global *ctx.Global) error {
		for _, c := range children {
			if err := c.processor(w, global); err != nil {
				return err
			}
		}
		return nil
	}, nil
}

func ifDirective(values []string, children []*astNode) (processor, error) {
	programs := make([]*expr.Program, len(values))
	c0 := values[0]
	c0 = strings.TrimPrefix(c0, "<#if ")
	c0 = strings.TrimSuffix(c0, ">")
	var err error
	programs[0], err = expr.Compile(c0)
	if err != nil {
		return nil, err
	}
	for i := 1; i < len(values); i++ {
		c := values[i]
		c = strings.TrimPrefix(c, "<#elseif ")
		c = strings.TrimSuffix(c, ">")
		programs[i], err = expr.Compile(c)
		if err != nil {
			return nil, err
		}
	}
	return func(w io.Writer, global *ctx.Global) error {
		for i, p := range programs {
			ret, err := p.Run(global)
			if err != nil {
				return err
			}
			if utils.Truly(ret) {
				return children[i].processor(w, global)
			}
		}
		if len(children)-len(values) == 1 {
			return children[len(children)-1].processor(w, global)
		}
		return nil
	}, nil
}

func tmplDirective(values []string, children []*astNode) (processor, error) {
	p, err := expr.Compile(values[0])
	if err != nil {
		return nil, err
	}
	return func(w io.Writer, global *ctx.Global) error {
		ret, err := p.Run(global)
		if err != nil {
			return err
		}
		if ret == nil {
			_, _ = fmt.Fprint(w, "NULL")
		} else {
			_, _ = fmt.Fprint(w, ret)
		}
		return nil
	}, err
}

func listDirective(values []string, children []*astNode) (processor, error) {
	codes := strings.SplitN(strings.TrimSuffix(strings.TrimPrefix(values[0], "<#list "), ">"), " as ", 2)
	if len(codes) != 2 {
		return nil, errors.New("list directive should be 'sequence' as 'variable'")
	}
	code := strings.TrimSpace(codes[0])
	p, err := expr.Compile(code)
	if err != nil {
		return nil, err
	}
	varName := strings.TrimSpace(codes[1])
	return func(w io.Writer, global *ctx.Global) error {
		list, err := p.Run(global)
		if err != nil {
			return err
		}
		li := reflect.ValueOf(list)
		if li.Kind() != reflect.Slice {
			return fmt.Errorf("sequence '%s' is not list", code)
		}
		l := li.Len()
		if l == 0 && len(children) == 2 {
			return children[1].processor(w, global)
		}
		g := global.Fork()
		for i := 0; i < l; i++ {
			v := li.Index(i).Interface()
			g.Variables[varName] = v
			if err := children[0].processor(w, g); err != nil {
				return err
			}
		}
		return nil
	}, nil
}
