package chess2

func Abs[T float64|int](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

func Sign[T float64|int](x T) int {
	switch {
	case x == 0: return 0
	case x < 0: return -1
	default: return 1
	}
}
