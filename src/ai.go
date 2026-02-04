package chess2

import (
	"context"
	"iter"
	"time"
)

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

func GetAllMoves(b Board) iter.Seq[Move] {
	return func(yield func(Move) bool) {
		for x := range BoardSize {
			for y := range BoardSize {
				piece := *b.At(x, y)
				if piece.Side() != b.Turn {
					continue
				}

				for _, m := range b.GetMoves(x, y) {
					if !yield(m) {
						return
					}
				}
			}
		}
	}
}

func BestMove(b Board, depth int) Move {
	var result Move
	var bestScore float64
	isMaximizing := b.Turn == SideWhite
	if isMaximizing {
		bestScore = -1000000
	} else {
		bestScore = 1000000
	}

	for m := range GetAllMoves(b) {
		nextBoard := *b.Apply(m)
		if depth > 0 {
			nextBoard = *nextBoard.Apply(BestMove(nextBoard, depth - 1))
		}
		score := Evaluate(nextBoard)

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

	return result
}

func SearchBestMove(b Board, out chan Move) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 2)
	defer cancel()

	var result Move
	depth := 0
	loop: for {
		depth += 1
		resultChannel := make(chan Move)
		go func() {
			resultChannel <- BestMove(b, depth)
		}()

		select {
		case <-ctx.Done():
			depth--
			break loop
		case m := <-resultChannel:
			result = m
		}
	}
	out <- result
	println("Depth", depth)
}
