package cons

import (
	"fmt"
)

type ConsCell struct {
	Car interface{}
	Cdr interface{}
}

func Cons(car interface{}, cdr interface{}) ConsCell {
	return ConsCell{Car: car, Cdr: cdr}
}

func Register(env map[string]interface{}) {
	env["cons"] = func(argv []interface{}) (interface{}, error) {
		switch cdr := argv[1].(type) {
		case interface{}:
			return Cons(argv[0], cdr), nil
		case nil:
			return Cons(argv[0], nil), nil
		}
		return nil, fmt.Errorf("could not cons, the value %v was not a ConsCell but a %T", argv[1], argv[1])
	}

	env["car"] = func(argv []interface{}) (interface{}, error) {
		switch cons := argv[0].(type) {
		case ConsCell:
			return cons.Car, nil
		case nil:
			return nil, nil
		}
		return nil, fmt.Errorf("can only car a cons cell but not %v, which is %t", argv[0], argv[0])
	}

	env["cdr"] = func(argv []interface{}) (interface{}, error) {
		switch cons := argv[0].(type) {
		case ConsCell:
			if cons.Cdr == nil {
				return nil, nil
			} else {
				return cons.Cdr, nil
			}
		case nil:
			return nil, nil
		}
		return nil, fmt.Errorf("can only cdr a cons cell but not %v, which is %t", argv[0], argv[0])
	}
}
