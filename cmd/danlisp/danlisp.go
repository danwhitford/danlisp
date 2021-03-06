package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/danwhitford/danlisp/internal/interpreter"
	"github.com/danwhitford/danlisp/internal/lexer"
	"github.com/danwhitford/danlisp/internal/parser"
	"os"
	"strings"
)

var header string = `
______            _     _           
|  _  \          | |   (_)          
| | | |__ _ _ __ | |    _ ___ _ __  
| | | / _  | '_ \| |   | / __| '_ \ 
| |/ / (_| | | | | |___| \__ \ |_) |
|___/ \__,_|_| |_\_____/_|___/ .__/ 
                             | |    
                             |_|    `

func repl() {
	scanner := bufio.NewScanner(os.Stdin)
	intr := interpreter.NewInterpreter()
	var lxr lexer.Lexer
	var psr parser.Parser
	var buf strings.Builder

	fmt.Println(header)
	for {
		if buf.Len() > 0 {
			fmt.Print("... ")
		} else {
			fmt.Print(">>> ")
		}

		scanned := scanner.Scan()
		if !scanned {
			break
		}
		line := scanner.Text()

		if len(line) < 1 {
			if buf.Len() < 1 {
				continue
			}
			stmt := buf.String()
			lxr = lexer.NewLexer(stmt)
			tokens, err := lxr.GetTokens()
			if err != nil {
				fmt.Println(err.Error())
				buf.Reset()
				continue
			}
			psr = parser.NewParser(tokens)
			exprs, err := psr.GetExpressions()
			if err != nil {
				fmt.Println(err.Error())
				buf.Reset()
				continue
			}
			res, err := intr.Interpret(exprs)
			if err != nil {
				fmt.Println(err.Error())
				buf.Reset()
				continue
			}

			if res != nil {
				fmt.Printf("%v\n", res)
			}
			buf.Reset()
		} else {
			buf.WriteString(line)
		}
	}
}

func fromFile(filename string) {
	dat, err := os.ReadFile(filename)
	if err != nil {
		errorQuit(err)
	} else {
		lxr := lexer.NewLexer(string(dat))
		tokens, err := lxr.GetTokens()
		if err != nil {
			errorQuit(err)
		}

		prsr := parser.NewParser(tokens)
		ast, err := prsr.GetExpressions()
		if err != nil {
			errorQuit(err)
		}

		intr := interpreter.NewInterpreter()
		_, err = intr.Interpret(ast)
		if err != nil {
			errorQuit(err)
		}
	}
}

func errorQuit(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "%s [filename]\n", os.Args[0])
		fmt.Fprint(flag.CommandLine.Output(), "\t [filename] to run from source or blank to start REPL\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if filename := flag.Arg(0); filename != "" {
		fromFile(filename)
	} else {
		repl()
	}
}
