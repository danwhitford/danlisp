package stringswrapper

import "strings"

func Import(env map[string]interface{}) {

	env["strings/Contains"] = func(argv []interface{}) (interface{}, error) {
		return strings.Contains(argv[0].(string), argv[1].(string)), nil
	}

	env["strings/Join"] = func(argv []interface{}) (interface{}, error) {
		return strings.Join(argv[0].([]string), argv[1].(string)), nil
	}

}
