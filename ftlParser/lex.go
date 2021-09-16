package ftlParser

import (
	"regexp"
	"strings"
)

const (
	none = iota
	directive
	tmpl
)

var directiveMap = map[string]int{
	"#if":     IF,
	"#else":   ELSE,
	"#elseif": ELSE_IF,
	"/#if":    END_IF,
	"#list":   LIST,
	"/#list":  END_LIST,
}

type lex struct {
	data        string
	pos         int
	state       int
	parseResult *astNode
}

func newLex(s string) *lex {
	return &lex{
		data: s,
		pos:  0,
	}
}

var regexp1 = regexp.MustCompile(`(\${.*?}|</?#.*?>)`)
var directiveReg = regexp.MustCompile(`(</?#[a-zA-Z]+?[ >])`)

func (l *lex) Lex(lval *yySymType) int {
	for {
		if l.pos >= len(l.data) {
			return 0
		}
		switch l.state {
		case none:
			idx := regexp1.FindStringIndex(l.data[l.pos:])
			if idx == nil {
				lval.text = l.data[l.pos:]
				lval.offset = l.pos
				l.pos = len(l.data)
				return TEXT
			}
			s := l.data[l.pos+idx[0] : l.pos+idx[1]]
			switch {
			case strings.HasPrefix(s, "${"):
				l.state = tmpl
			case strings.HasPrefix(s, "<"):
				l.state = directive
			}
			if l.pos+idx[0] > l.pos {
				lval.text = l.data[l.pos : l.pos+idx[0]]
				lval.offset = l.pos
				l.pos = l.pos + idx[0]
				return TEXT
			}
		case tmpl:
			var state int
			var escape bool
			var idx = l.pos
		loop:
			for i := 0; l.pos+i < len(l.data); i++ {
				char := l.data[l.pos+i]
				switch state {
				case 0:
					switch char {
					case '}':
						break loop
					case '"':
						state = 1
					case '\'':
						state = 2
					}
				case 1:
					if char == '\\' && !escape {
						escape = true
					} else {
						escape = false
					}
					if char == '"' && !escape {
						state = 0
					}
				case 2:
					if char == '\\' && !escape {
						escape = true
					} else {
						escape = false
					}
					if char == '\'' && !escape {
						state = 0
					}
				}
				idx++
			}
			lval.text = l.data[l.pos+2 : idx]
			lval.offset = l.pos
			l.pos = idx + 1
			l.state = none
			return TMPL
		case directive:
			idx := regexp1.FindStringIndex(l.data[l.pos:])
			lval.text = l.data[l.pos+idx[0] : l.pos+idx[1]]
			lval.offset = l.pos
			l.pos = l.pos + idx[1]
			for {
				if l.pos < len(l.data)-1 && (l.data[l.pos] == ' ' || l.data[l.pos] == '\r' || l.data[l.pos] == '\n') {
					l.pos++
				} else {
					break
				}
			}
			l.state = none
			dir := directiveReg.FindString(lval.text)
			dir = dir[1 : len(dir)-1]
			tok, ok := directiveMap[dir]
			if !ok {
				l.Error("unknown directive " + dir)
			}
			return tok
		}
	}
}

func (l *lex) Error(s string) {
	panic(s)
}
