package utils

type Tuple[A any, B any] struct {
	A A
	B B
}

func AbsInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
