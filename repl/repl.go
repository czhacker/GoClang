package repl

import (
	"GoClang/parser"
	"GoClang/lexer"
	"GoClang/evaluator"
	"fmt"
	"bufio"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer)  {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParserProgram()
		if len(p.Errors()) != 0 {
			return
		}

		evaluated := evaluator.Eval(program)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}
