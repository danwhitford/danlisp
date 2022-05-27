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

type Set struct {
	Var   Symbol
	Value Expr
}

type If struct {
	Cond        Expr
	TrueBranch  Expr
	FalseBranch Expr
}

type While struct {
	Cond Expr
	Body []Expr
}

type Defn struct {
	Name    Symbol
	Arglist []Symbol
	Body    []Expr
}

type For struct {
	Initialiser Expr
	Cond        Expr
	Step        Expr
	Body        []Expr
}
