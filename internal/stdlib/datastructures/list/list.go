package list

import "github.com/danwhitford/danlisp/internal/stdlib/datastructures/cons"

func Register(env map[string]interface{}) {
	env["list"] = func(argv []interface{}) (interface{}, error) {
		if len(argv) < 1 {
			return nil, nil
		}

		var val interface{}
		var outer cons.ConsCell
		for l := len(argv) - 1; l >= 0; l-- {
			switch cdr := val.(type) {
			case nil:
				outer = cons.Cons(argv[l], nil)
			case interface{}:
				outer = cons.Cons(argv[l], cdr)

			}
			val = outer
		}
		return outer, nil
	}

	env["nth"] = func(argv []interface{}) (interface{}, error) {
		hd := argv[0].(cons.ConsCell)
		nth := argv[1].(float64)
		for nth > 0 {
			hd = hd.Cdr.(cons.ConsCell)
		}
		return hd.Car, nil
	}
}
