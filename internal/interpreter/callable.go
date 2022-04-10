package interpreter

import (
	"github.com/shaftoe44/danlisp/internal/expr"
)

type ICallable interface {
	call(*Interpreter, []interface{}) (interface{}, error)
}

type Callable struct {
	Arity int
	Args  []string
	Body  []expr.Expr
}

// TODO check arity and stuff
func (callable Callable) call(context *Interpreter, argv []interface{}) (interface{}, error) {
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

type BuiltIn struct {
	Body func(argv []interface{}) interface{}
}

func (builtin BuiltIn) call(context *Interpreter, argv []interface{}) (interface{}, error) {
	return builtin.Body(argv), nil
}
