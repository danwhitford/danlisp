package callable

import (
	"github.com/danwhitford/danlisp/internal/expr"
)

type Callable struct {
	Arity int
	Args  []string
	Body  []expr.Expr
}
