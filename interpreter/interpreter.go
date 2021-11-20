package interpreter

import (
	"fmt"

	"whitford.io/danlisp/expr"
)

var environment map[string]interface{}

func Interpret(exprs []expr.Expr) (interface{}, error) {
	environment = newEnvironment()

	var retval interface{}
	var err error
	for _, ex := range exprs {
		retval, err = eval(ex)
		if err != nil {
			return nil, err
		}
	}
	return retval, nil
}

func eval(ex expr.Expr) (interface{}, error) {
	atom, ok := ex.(expr.Atom)
	if ok {
		return evalAtom(atom), nil
	}

	seq, ok := ex.(expr.Seq)
	if ok {
		return evalSeq(seq)
	}

	symbol, ok := ex.(expr.Symbol)
	if ok {
		return symbol, nil
	}

	panic("Don't know how to eval this thing")
}

func evalAtom(ex expr.Atom) interface{} {
	return ex.Value
}

func evalSeq(ex expr.Seq) (interface{}, error) {
	symbol, _ := eval(ex.Exprs[0])
	fname := symbol.(expr.Symbol).Name
	args := []interface{}{}
	for _, argex := range ex.Exprs[1:] {
		arg, _ := eval(argex)
		args = append(args, arg)
	}
	applyer, ok := environment[fname].(func(argv []interface{}) interface{})
	if !ok {
		return nil, fmt.Errorf("Runtime error. Could not find function '%v'.", fname)
	}
	return applyer(args), nil
}

func newEnvironment() map[string]interface{} {
	env := make(map[string]interface{})
	env["+"] = func(argv []interface{}) interface{} { return argv[0].(float64) + argv[1].(float64) }
	env["-"] = func(argv []interface{}) interface{} { return argv[0].(float64) - argv[1].(float64) }
	env["*"] = func(argv []interface{}) interface{} { return argv[0].(float64) * argv[1].(float64) }
	env["/"] = func(argv []interface{}) interface{} { return argv[0].(float64) / argv[1].(float64) }
	env["mod"] = func(argv []interface{}) interface{} { return int(argv[0].(float64)) % int(argv[1].(float64)) }

	return env
}
