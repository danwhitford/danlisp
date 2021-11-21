package main

import (
	"bufio"
	"fmt"
	"os"

	"whitford.io/danlisp/internal/interpreter"
	"whitford.io/danlisp/internal/lexer"
	"whitford.io/danlisp/internal/parser"
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

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println(header)
	for {
		fmt.Print(">>> ")
		scanned := scanner.Scan()
		if !scanned {
			break
		}
		line := scanner.Text()
		
		if len(line) < 1 {
			continue
		}

		tokens, err := lexer.GetTokens(line)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		exprs, err := parser.GetExpressions(tokens)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		res, err := interpreter.Interpret(exprs)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		
		fmt.Printf("%v\n", res)
	}
}
