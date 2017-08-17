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
	checkParserErrors(t, p)
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

func checkParserErrors(t *testing.T, p *Parser){
	errors := p.Errors()
	if len(errors) == 0{
		return
	}
	t.Errorf("The Parser have %s errors", len(errors))

	for _, msg := range errors {
		t.Errorf("parser error: %s", msg)
	}
	t.FailNow()
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

func TestReturnStatements(t *testing.T){
	input := `
	return 5;
	return 10;
	return 993322;`

	l := lexer.New(input)
	p := New(l)

	program := p.ParserProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program statement does not contain 3 statements, got %d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got %T", stmt)
			continue
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.Tokenliteral not 'return', got %q", returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T)  {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParserProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1{
		t.Fatalf("program has not enough statements got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("expression not *ast.Identifier. got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s, got=%s", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", ident.TokenLiteral())
	}
}


func TestIntegerLiteralExpression(t *testing.T){
	input := "5;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParserProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] is not ast.ExpressionStatement, got %T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("stmt is not ast.IntegerLiteral, got %T", stmt)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("Literal.TokenLiteral not %s. got=%s", "5", literal.TokenLiteral())
	}

	if literal.Value != 5 {
		t.Errorf("Literal.Value not 5, got=%d", literal.Value)
	}
}