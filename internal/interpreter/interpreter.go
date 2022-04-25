package interpreter

import (
	"fmt"
	"strings"

	"github.com/danwhitford/danlisp/internal/callable"
	"github.com/danwhitford/danlisp/internal/expr"
	"github.com/danwhitford/danlisp/internal/stdlib/danreflect"
	"github.com/danwhitford/danlisp/internal/stdlib/datastructures/cons"
	"github.com/danwhitford/danlisp/internal/stdlib/datastructures/list"
	stringswrapper "github.com/danwhitford/danlisp/internal/stdlib/wrappers"
)

type Interpreter struct {
	environment map[string]interface{}
}

func NewInterpreter() Interpreter {
	return Interpreter{NewEnvironment()}
}

func (interpreter *Interpreter) Interpret(exprs []expr.Expr) (interface{}, error) {
	var retval interface{}
	var err error
	for _, ex := range exprs {
		retval, err = interpreter.eval(ex)
		if err != nil {
			return nil, err
		}
	}
	return retval, nil
}

func (interpreter *Interpreter) eval(ex expr.Expr) (interface{}, error) {

	switch v := ex.(type) {
	case expr.Atom:
		return evalAtom(v), nil
	case expr.Seq:
		return interpreter.evalSeq(v)
	case expr.Symbol:
		return interpreter.evalSymbol(v)
	case expr.Set:
		return interpreter.evalSet(v)
	case expr.If:
		return interpreter.evalIf(v)
	case expr.While:
		return interpreter.evalWhile(v)
	case expr.Defn:
		return interpreter.evalDefun(v)
	}

	return nil, fmt.Errorf("don't know how to eval this thing %v of type %T", ex, ex)
}

func (interpreter *Interpreter) evalDefun(ex expr.Defn) (interface{}, error) {
	arglist := []string{}
	for _, a := range ex.Arglist {
		arglist = append(arglist, a.Name)
	}
	callable := callable.Callable{Arity: len(arglist), Args: arglist, Body: ex.Body}
	interpreter.environment[ex.Name.Name] = callable
	return nil, nil
}

func (interpreter *Interpreter) evalWhile(ex expr.While) (interface{}, error) {
	var retval interface{}

	for {
		c, err := interpreter.eval(ex.Cond)
		if err != nil {
			return nil, err
		}
		if !isTruthy(c) {
			break
		}
		for _, line := range ex.Body {
			val, er := interpreter.eval(line)
			if er != nil {
				return nil, er
			}
			retval = val
		}
	}
	return retval, nil
}

func evalAtom(ex expr.Atom) interface{} {
	return ex.Value
}

func (interpreter *Interpreter) evalSymbol(ex expr.Symbol) (interface{}, error) {
	val, ok := interpreter.environment[ex.Name]
	if !ok {
		return nil, fmt.Errorf("runtime error. Could not find symbol '%v'", ex.Name)
	}
	return val, nil
}

func (interpreter *Interpreter) evalSeq(ex expr.Seq) (interface{}, error) {
	symbol, err := interpreter.eval(ex.Exprs[0])
	if err != nil {
		return symbol, err
	}
	args := []interface{}{}
	for _, argex := range ex.Exprs[1:] {
		arg, err := interpreter.eval(argex)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
	//TODO convert builtins to Callables
	switch s := symbol.(type) {
	case callable.Callable:
		return interpreter.call(s, args)
	case func([]interface{}) interface{}:
		return s(args), nil
	case func([]interface{}) (interface{}, error):
		return s(args)
	}
	return nil, fmt.Errorf("did not know how to evaluate seq %v", ex)
}

func (interpreter *Interpreter) evalSet(ex expr.Set) (interface{}, error) {
	val, err := interpreter.eval(ex.Value)
	interpreter.environment[ex.Var.Name] = val
	return nil, err
}

func (interpreter *Interpreter) evalIf(iff expr.If) (interface{}, error) {
	cond, err := interpreter.eval(iff.Cond)
	if err != nil {
		return nil, err
	}
	var expr expr.Expr
	if isTruthy(cond) {
		expr = iff.TrueBranch
	} else {
		expr = iff.FalseBranch
	}
	return interpreter.eval(expr)
}

func NewEnvironment() map[string]interface{} {
	env := make(map[string]interface{})

	// Built in vars
	env["t"] = true

	// Basic operators
	env["+"] = func(argv []interface{}) interface{} { return argv[0].(float64) + argv[1].(float64) }
	env["-"] = func(argv []interface{}) interface{} { return argv[0].(float64) - argv[1].(float64) }
	env["*"] = func(argv []interface{}) interface{} { return argv[0].(float64) * argv[1].(float64) }
	env["/"] = func(argv []interface{}) interface{} { return argv[0].(float64) / argv[1].(float64) }
	env["mod"] = func(argv []interface{}) interface{} { return float64(int(argv[0].(float64)) % int(argv[1].(float64))) }

	// Bitwise ops
	env["&"] = func(argv []interface{}) interface{} { return float64(int(argv[0].(float64)) & int(argv[1].(float64))) }
	env["|"] = func(argv []interface{}) interface{} { return float64(int(argv[0].(float64)) | int(argv[1].(float64))) }
	env["^"] = func(argv []interface{}) interface{} { return float64(int(argv[0].(float64)) ^ int(argv[1].(float64))) }
	env["&^"] = func(argv []interface{}) interface{} { return float64(int(argv[0].(float64)) &^ int(argv[1].(float64))) }
	env[">>"] = func(argv []interface{}) interface{} { return float64(int(argv[0].(float64)) >> int(argv[1].(float64))) }
	env["<<"] = func(argv []interface{}) interface{} { return float64(int(argv[0].(float64)) << int(argv[1].(float64))) }

	// Boleans
	env["="] = func(argv []interface{}) interface{} { return argv[0] == argv[1] }
	env["and"] = func(argv []interface{}) interface{} { return isTruthy(argv[0]) && isTruthy(argv[1]) }
	env["or"] = func(argv []interface{}) interface{} { return isTruthy(argv[0]) || isTruthy(argv[1]) }

	// Comparison
	env["gt"] = func(argv []interface{}) interface{} { return argv[0].(float64) > argv[1].(float64) }
	env["lt"] = func(argv []interface{}) interface{} { return argv[0].(float64) < argv[1].(float64) }

	// Utility
	env["prn"] = func(argv []interface{}) interface{} {
		strs := []string{}
		for _, v := range argv {
			strs = append(strs, fmt.Sprintf("%v", v))
		}
		p, _ := fmt.Println(strings.Join(strs, " "))
		return p
	}

	cons.Import(env)
	stringswrapper.Import(env)
	danreflect.Import(env)
	list.Import(env)

	return env
}

func isTruthy(v interface{}) bool {
	if b, ok := v.(bool); ok {
		return b
	}
	return v != nil
}

// TODO check arity and stuff
func (context *Interpreter) call(callable callable.Callable, argv []interface{}) (interface{}, error) {
	for i, a := range argv {
		context.environment[callable.Args[i]] = a
	}
	var retval interface{}
	var err error
	for _, e := range callable.Body {
		retval, err = context.eval(e)
		if err != nil {
			return nil, err
		}
	}
	return retval, nil
}
