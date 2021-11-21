package main

import (
	"bufio"
	"fmt"
	"os"

	"whitford.io/danlisp/interpreter"
	"whitford.io/danlisp/lexer"
	"whitford.io/danlisp/parser"
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
	fmt.Print(">>> ")
	for scanner.Scan() {
		line := scanner.Text()
		tokens, _ := lexer.GetTokens(line)
		exprs, _ := parser.GetExpressions(tokens)
		res, _ := interpreter.Interpret(exprs)
		fmt.Printf("%v\n", res)

		fmt.Print(">>> ")
	}
}
