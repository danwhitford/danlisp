package list

import "github.com/shaftoe44/danlisp/internal/stdlib/datastructures/cons"

func Import(env map[string]interface{}) {
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
}
