package cons

type ConsCell struct {
	Car interface{}
	Cdr *ConsCell
}

func Cons(car interface{}, cdr *ConsCell) ConsCell {
	return ConsCell{Car: car, Cdr: cdr}
}
