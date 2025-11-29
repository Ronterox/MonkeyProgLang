package execution

import (
	"bufio"
	"fmt"
	"io"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/token"
	"os"
)

const PROMPT = ">> "

func StartREPL(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			for _, msg := range p.Errors() {
				fmt.Println(msg)
			}
			continue
		}

		result := evaluator.Eval(program, env)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
		fmt.Println(program.String())
		fmt.Println(result.Inspect())
	}
}

func RunCode(in io.Reader, out io.Writer, file string) {
	env := object.NewEnvironment()

	text, err := os.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	l := lexer.New(string(text))
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		for _, msg := range p.Errors() {
			fmt.Println(msg)
		}
		return
	}

	result := evaluator.Eval(program, env)

	for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		fmt.Printf("%+v\n", tok)
	}

	fmt.Println(result.Inspect())
}
