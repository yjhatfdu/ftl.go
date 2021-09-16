package expr

type AstNode struct {
	Children  []*AstNode
	NodeType  int
	Value     string
	Offset    int
	Length    int
	ValueType int
}

func newAst(t int, v string, vt int, pos int, children ...*AstNode) *AstNode {
	return &AstNode{
		Children:  children,
		NodeType:  t,
		Value:     v,
		ValueType: vt,
		Offset:    pos,
	}
}

func (an *AstNode) SetOffset(offset int) *AstNode {
	an.Offset = offset
	return an
}
