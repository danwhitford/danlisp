package callable

import (
	"github.com/shaftoe44/danlisp/internal/expr"
)

type Callable struct {
	Arity int
	Args  []string
	Body  []expr.Expr
}
