package parser

import (
	"testing"
	"GoClang/ast"
	"GoClang/lexer"
)

func TestLetStament(t *testing.T){
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParserProgram()
	if program == nil {
		t.Fatalf("ParserProgram return nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if	!testLetStatement(t, stmt, tt.expectedIdentifier){
			return
		}
	}
}

func testLetStatement(t *testing.T, stmt ast.Statement, name string) bool{
	if  stmt.TokenLiteral() != "let" {
		t.Errorf("statement.Tokenliteral not 'let'. got=%s", stmt.TokenLiteral())
		return false
	}

	letstmt,ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("stmt not ast.LetStatement. got=%t", stmt)
		return false
	}

	if letstmt.Name.Value != name {
		t.Errorf("letstmt.Name.Value not '%s', got='%s'", name, letstmt.Name.Value)
		return false
	}

	if letstmt.Name.TokenLiteral() != name {
		t.Errorf("letstmt.Token.Literal not '%s', got='%s'", name, letstmt.Name.TokenLiteral())
		return false
	}

	return true
}
