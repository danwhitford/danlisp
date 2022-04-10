package cons

import "fmt"

type ConsCell struct {
	Car interface{}
	Cdr *ConsCell
}

func Cons(car interface{}, cdr *ConsCell) ConsCell {
	return ConsCell{Car: car, Cdr: cdr}
}

func Import(env map[string]interface{}) {
	env["cons"] = func(argv []interface{}) (interface{}, error) {
		switch cdr := argv[1].(type) {
		case ConsCell:
			return Cons(argv[0], &cdr), nil
		case nil:
			return Cons(argv[0], nil), nil
		}
		return nil, fmt.Errorf("could not cons, the value %v was not a ConsCell but a %T", argv[1], argv[1])
	}

	env["car"] = func(argv []interface{}) (interface{}, error) {
		switch cons := argv[0].(type) {
		case ConsCell:
			return cons.Car, nil
		}
		return nil, fmt.Errorf("can only car a cons cell but not %v, which is %t", argv[0], argv[0])
	}

	env["cdr"] = func(argv []interface{}) (interface{}, error) {
		switch cons := argv[0].(type) {
		case ConsCell:
			return *cons.Cdr, nil
		}
		return nil, fmt.Errorf("can only cdr a cons cell but not %v, which is %t", argv[0], argv[0])
	}

	env["list"] = func(argv []interface{}) (interface{}, error) {
		var val interface{}
		var outer ConsCell
		for l := len(argv) - 1; l >= 0; l-- {
			switch cdr := val.(type) {
			case nil:
				outer = Cons(argv[l], nil)
			case ConsCell:
				outer = Cons(argv[l], &cdr)

			}
			val = outer
		}
		return outer, nil
	}
}
