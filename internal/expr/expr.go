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

type Def struct {
	Var   Symbol
	Value Expr
}

type If struct {
	Cond        Expr
	TrueBranch  Expr
	FalseBranch Expr
}

type When struct {
	Cond Expr
	Body []Expr
}
