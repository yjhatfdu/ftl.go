%{
package ftlParser
%}

%union {
  offset int
  text string
  node *astNode
}

%token TEXT
%token IF
%token ELSE
%token ELSE_IF
%token END_IF
%token TMPL
%token GROUP
%token LIST
%token END_LIST
%right IF
%right ELSE_IF
%%

input: e_group {yylex.(*lex).parseResult=$1.node}

e:
TEXT {
	$$.node=&astNode{T:TEXT,V:[]string{$1.text}}
	}
   |TMPL {
   	$$.node=&astNode{T:TMPL,V:[]string{$1.text}}
   	}
   | IF e_group END_IF{
      $$.node=&astNode{T:IF,V:[]string{$1.text},C:[]*astNode{$2.node}}
   }
   | IF e_group ELSE e_group END_IF{
         $$.node=&astNode{T:IF,V:[]string{$1.text},C:[]*astNode{$2.node,$3.node}}
      }
   | IF e_group elseif_clauses END_IF {
   $$.node=&astNode{T:IF,V:append([]string{$1.text},$3.node.V...),C:append([]*astNode{$2.node},$3.node.C...)}
   }
   | IF e_group elseif_clauses ELSE e_group END_IF {
      $$.node=&astNode{T:IF,V:append([]string{$1.text},$3.node.V...),C:append([]*astNode{$2.node},$3.node.C...)}
      $$.node.C=append($$.node.C,($5.node))
   }
   | LIST e_group END_LIST {
      $$.node=&astNode{T:LIST,V:[]string{$1.text},C:[]*astNode{$2.node}}
   }
   | LIST e_group ELSE e_group END_LIST {
       $$.node=&astNode{T:LIST,V:[]string{$1.text},C:[]*astNode{$2.node,$4.node}}
   }


e_group: e {
	$$.node=&astNode{T:GROUP,C:[]*astNode{$1.node}}
	}
	| e_group e{
        $$.node.C=append($$.node.C,$2.node)
        }

elseif_clauses:
         ELSE_IF e_group{
        $$.node=&astNode{T:GROUP,V:[]string{$1.text},C:[]*astNode{$2.node}}
       }
       |  ELSE_IF e_group elseif_clauses{
        $$.node.C=append($3.node.C,$2.node)
        $$.node.V=append($3.node.V,$1.text)
       }
