package chess2

var cost = []float64{
	0,
	1, -1,
	2.5, -2.5,
	3, -3,
	5, -5,
	9, -9,
	1000, -1000,
}

func Evaluate(b Board) float64 {
	switch b.Winner {
	case SideWhite: return 1000
	case SideBlack: return -1000
	}

	var result float64
	for _, piece := range b.inner {
		result += cost[piece]
	}

	return result
}

func BestMove(b Board, depth int) Move {
	var result Move
	var bestScore float64
	isMaximizing := b.Turn == SideWhite
	if isMaximizing {
		bestScore = -1000
	} else {
		bestScore = 1000
	}

	for x := range BoardSize {
		for y := range BoardSize {
			piece := *b.At(x, y)
			if piece.Side() != b.Turn {
				continue
			}

			for _, m := range b.GetMoves(x, y) {
				score := Evaluate(*b.Apply(m))
				var condition bool
				if isMaximizing {
					condition = score > bestScore
				} else {
					condition = score < bestScore
				}
				if condition {
					bestScore = score
					result = m
				}
			}
		}
	}

	return result
}
