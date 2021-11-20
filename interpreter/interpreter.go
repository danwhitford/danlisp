package interpreter

import (
	"whitford.io/danlisp/expr"
)

var environment map[string]interface{}

func Interpret(exprs []expr.Expr) interface{} {
	environment = newEnvironment()

	var retval interface{}
	for _, ex := range exprs {
		retval = eval(ex)
	}
	return retval
}

func eval(ex expr.Expr) interface{} {
	atom, ok := ex.(expr.Atom)
	if ok {
		return evalAtom(atom)
	}

	seq, ok := ex.(expr.Seq)
	if ok {
		return evalSeq(seq)
	}

	symbol, ok := ex.(expr.Symbol)
	if ok {
		return symbol
	}

	panic("Don't know how to eval this thing")
}

func evalAtom(ex expr.Atom) interface{} {
	return ex.Value
}

func evalSeq(ex expr.Seq) interface{} {
	fname := eval(ex.Exprs[0]).(expr.Symbol).Name
	args := []interface{}{}
	for _, argex := range ex.Exprs[1:] {
		args = append(args, eval(argex))
	}
	applyer := environment[fname].(func(argv []interface{}) interface{})
	return applyer(args)
}

func newEnvironment() map[string]interface{} {
	env := make(map[string]interface{})
	env["+"] = func(argv []interface{}) interface{} { return argv[0].(float64) + argv[1].(float64) }

	return env
}
