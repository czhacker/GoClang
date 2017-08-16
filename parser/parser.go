package parser

import (
	"GoClang/ast"
	"GoClang/token"
	"GoClang/lexer"
)

type Parser struct {
	l *lexer.Lexer

	curToken token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := Parser{l:l}
	p.nextToken()
	p.nextToken()
	return &p
}

func (p *Parser)nextToken(){
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParserProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parserStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parserStatement() ast.Statement{
	switch p.curToken.Type {
	case token.LET:
		return p.parserLetStatement()
	default:
		return nil

	}
}

func (p *Parser) parserLetStatement() ast.Statement {
	stmt := &ast.LetStatement{Token:p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token:p.curToken, Value:p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN){
		return nil
	}

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}


func (p *Parser)curTokenIs(t token.Tokentype) bool {
	return p.curToken.Type == t
}

func (p *Parser)peekTokenIs(t token.Tokentype) bool {
	return p.peekToken.Type == t
}

func (p *Parser)expectPeek(t token.Tokentype) bool {
	if p.peekTokenIs(t){
		p.nextToken()
		return true
	}else{
		return false
	}
}




