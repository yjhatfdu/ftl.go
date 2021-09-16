%{
package expr
%}

%union {
  offset int
  node *AstNode
  text string
}

%token NULL
%token INT
%token STR
%token RAWSTR
%token BOOL
%token FLOAT
%token CONST
%token OR
%token AND
%token NOT
%token LIKE
%token NEQ
%token GT
%token LT
%token GTE
%token LTE
%token EQ
%token ADD
%token MINUS
%token MUL
%token DIV
%token CONTAINS
%token ID
%token IND
%token COMMA
%token ANY
%token FUNC
%token LP
%token RP
%token DOLLAR
%token VAR
%token CAST
%token PIPE
%token WHEN
%token CASE
%token THEN
%token IN
%token END
%token ELSE
%token IS
%token DOT
%token LBR
%token RBR
%token QM
%token SUB
%token ATTR
%token ARR
%token CALL

%left OR
%left AND
%left NOT
%left GT LT GTE LTE EQ NEQ
%left LIKE CONTAINS
%left ADD MINUS
%left MUL DIV
%left PIPE
%right IS
%right IN
%right LP LBR
%right DOLLAR
%right CAST
%right DOT
%%

input:    e      { yylex.(*Lexer).parseResult=$1.node;};

e:  INT          { $$.node =newAst(CONST,$1.text,INT,$1.offset); }
       | STR     { $$.node =newAst(CONST,$1.text,STR,$1.offset); }
       | RAWSTR  { $$.node =newAst(CONST,$1.text,STR,$1.offset); }
       | FLOAT   { $$.node =newAst(CONST,$1.text,FLOAT,$1.offset); }
       | BOOL    { $$.node =newAst(CONST,$1.text,BOOL,$1.offset); }
       | ID      { $$.node =newAst(VAR,$1.text,0,$1.offset); }
       | e DOT ID{ $$.node =newAst(ATTR,$3.text,0,$2.offset,$1.node);}
       | e LBR e RBR { $$.node =newAst(SUB,"",0,$2.offset,$1.node,$3.node);}
       | LBR arr_cnt RBR {
       		$$.node =newAst(ARR,"",0,$1.offset,$2.node.Children...)
       }
       | e AND e          { $$.node =newAst(FUNC,"and",0,$2.offset,$1.node,$3.node); }
       | e OR e           { $$.node =newAst(FUNC,"or",0,$2.offset,$1.node,$3.node); }
       | e ADD e          { $$.node =newAst(FUNC,"add",0,$2.offset,$1.node,$3.node); }
       | e MINUS e        { $$.node =newAst(FUNC,"minus",0,$2.offset,$1.node,$3.node); }
       | e DIV e          { $$.node =newAst(FUNC,"div",0,$2.offset,$1.node,$3.node); }
       | e MUL e          { $$.node =newAst(FUNC,"mul",0,$2.offset,$1.node,$3.node); }
       | NOT e            { $$.node =newAst(FUNC,"not",0,$1.offset,$2.node); }
       | e GT e           { $$.node =newAst(FUNC,"gt",0,$2.offset,$1.node,$3.node); }
       | e GTE e          { $$.node =newAst(FUNC,"gte",0,$2.offset,$1.node,$3.node); }
       | e LT e           { $$.node =newAst(FUNC,"lt",0,$2.offset,$1.node,$3.node); }
       | e LTE e          { $$.node =newAst(FUNC,"lte",0,$2.offset,$1.node,$3.node); }
       | e EQ e           { $$.node =newAst(FUNC,"eq",0,$2.offset,$1.node,$3.node); }
       | e NEQ e          { $$.node =newAst(FUNC,"neq",0,$2.offset,$1.node,$3.node); }
       | e LP arr_cnt RP  { $$.node =newAst(CALL,"",0,$1.offset,append([]*AstNode{$1.node},$3.node.Children...)...); }
       | e LP RP          { $$.node =newAst(CALL,"",0,$1.offset,$1.node) }

arr_cnt:
	e {
	$$.node=newAst(CONST,"",0,$1.offset,$1.node);
	}
	|e COMMA arr_cnt {
	$$.node=$3.node
	$$.node.Children=append([]*AstNode{$1.node},$3.node.Children...)
	}

