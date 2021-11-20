package expr

type Expr interface{}

type Atom struct {
	Value interface{}
}

type Seq struct {
	Exprs []Expr
}

type Symbol struct {
	Name string
}