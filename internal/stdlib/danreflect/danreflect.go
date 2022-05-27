package danreflect

import "fmt"

func getType(val interface{}) string {
	return fmt.Sprintf("%T", val)
}

func Register(env map[string]interface{}) {
	env["type"] = func(argv []interface{}) (interface{}, error) {
		return getType(argv[0]), nil
	}
}
