package chess2

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func Sign(x int) int {
	switch {
	case x == 0: return 0
	case x < 0: return -1
	default: return 1
	}
}
