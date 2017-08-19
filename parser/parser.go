package parser

import (
	"GoClang/ast"
	"GoClang/token"
	"GoClang/lexer"
	"fmt"
	"strconv"
)

type Parser struct {
	l *lexer.Lexer

	errors []string
	curToken token.Token
	peekToken token.Token

	prefixParseFns map[token.Tokentype]prefixParseFn
	infixParseFns map[token.Tokentype]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := Parser{l:l,
	errors:[]string{}}
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.Tokentype]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parserIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBooleanExpression)
	p.registerPrefix(token.FALSE, p.parseBooleanExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)


	p.infixParseFns = make(map[token.Tokentype]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	return &p
}

func (p *Parser) registerPrefix(tokenType token.Tokentype, fn prefixParseFn){
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.Tokentype, fn infixParseFn){
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.Tokentype) {
	msg := fmt.Sprintf("expected next token to be %s; got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors,msg)
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
	case token.RETURN:
		return p.parserReturnStatement()
	default:
		return p.parserExpressionStatement()
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

func (p *Parser) parserReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{Token:p.curToken}

	p.nextToken()

	//TODO: Skipping the expression until we encounter a semicolon
	for p.curToken.Type != token.SEMICOLON {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parserExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{Token:p.curToken}
	stmt.Expression = p.parserExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON){
		p.nextToken()
	}
	return stmt
}

const (
	_ int = iota
	LOWEST
	EQUALS //==
	LESSGREATER //< or >
	SUM //+
	PRODUCT //-
	PREFIX //-X or !X
	CALL // myFunction(X)
)

func (p *Parser) noPrefixParseFnError(t token.Tokentype) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors,msg)
}

func (p *Parser) parserExpression (precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftexp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence(){
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftexp
		}
		p.nextToken()
		leftexp = infix(leftexp)
	}
	return leftexp
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
		p.peekError(t)
		return false
	}
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn func(ast.Expression) ast.Expression
)


func (p *Parser)parserIdentifier() ast.Expression {
	return &ast.Identifier{Token:p.curToken, Value:p.curToken.Literal}
}

func (p *Parser)parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token:p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser)parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{Token:p.curToken, Operator:p.curToken.Literal,}

	p.nextToken()

	expression.Right = p.parserExpression(PREFIX)

	return expression
}

func (p *Parser)parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{Token:p.curToken, Left:left, Operator:p.curToken.Literal}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parserExpression(precedence)

	return expression
}

var precedences = map[token.Tokentype]int{
	token.EQ: EQUALS,
	token.NOT_EQ: EQUALS,
	token.LT: LESSGREATER,
	token.GT: LESSGREATER,
	token.PLUS: SUM,
	token.MINUS: SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH : PRODUCT,
}

func (p *Parser)peekPrecedence() int{
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser)curPrecedence() int{
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}


func (p *Parser)parseBooleanExpression() ast.Expression {
	return &ast.Boolean{Token:p.curToken, Value:p.curTokenIs(token.TRUE)}
}

func (p *Parser)parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parserExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}






