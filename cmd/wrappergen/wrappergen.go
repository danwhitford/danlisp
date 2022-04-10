package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

var TEMPLATE_DIR = "tmpl"

type WrapperFnArgs struct {
	Name string
	Typ  string
}

type WrapperFns struct {
	Name  string
	ArgTypes []string
}

type Wrapper struct {
	Pkg string
	Fns []WrapperFns
}

func main() {
	stringsWrapper := Wrapper{
		Pkg: "strings",
		Fns: []WrapperFns{{
			Name:  "Contains",
			ArgTypes: []string{"string", "string"},
		}, {
			Name:  "Join",
			ArgTypes: []string{"[]string", "string"},
		}},
	}

	fs := template.FuncMap{"arglist": func(argtypes []string) string {
		argvs := make([]string, 0)
		for i, el := range argtypes {
			argvs = append(argvs, fmt.Sprintf("argv[%d].(%s)", i, el))
		}
		return strings.Join(argvs, ",")
	},
		"lower": strings.ToLower}
	t, err := template.New("wrapper.tmpl").Funcs(fs).ParseFiles("tmpl/wrapper.tmpl")
	if err != nil {
		log.Panic(err)
	}

	f, err := os.Create(stringsWrapper.Pkg + "wrapper.go")
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()
	b := bufio.NewWriter(f)
	err = t.Execute(b, stringsWrapper)
	if err != nil {
		log.Panic(err)
	}
	b.Flush()
}
