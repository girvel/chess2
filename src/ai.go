package chess2

import (
	"context"
	"iter"
	"sync"
	"time"
)

type Ai struct {
	board *Board
	nextMove Move
	responseChannel chan map[Move]Move
	responseCtx context.Context
	responseCancel context.CancelFunc
	lastMoveTime time.Time
}

func CreateAi(board Board) *Ai {
	ai := Ai{
		responseChannel: make(chan map[Move]Move),
		board: &board,
		lastMoveTime: time.Now(),
	}
	ai.launchSearch()
	return &ai
}

func (ai *Ai) PushMove(m Move) {
	time.AfterFunc(min(time.Since(ai.lastMoveTime), time.Second * 10), func() {
		ai.responseCancel()
	})
	ai.nextMove = m
}

func (ai *Ai) PopResponse() *Move {
	select {
	case responses := <-ai.responseChannel:
		result := responses[ai.nextMove]
		ai.board.Move(ai.nextMove)
		ai.board.Move(result)
		ai.launchSearch()
		ai.lastMoveTime = time.Now()
		return &result
	default:
		return nil
	}
}

var cost = []float64{
	0,
	1, -1,
	2.5, -2.5,
	3, -3,
	5, -5,
	9, -9,
	1000, -1000,
}

func evaluate(b *Board) float64 {
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

func getAllMoves(b *Board) iter.Seq[Move] {
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

func alphaBeta(b *Board, depth int, alpha, beta float64) float64 {
	if depth == 0 || b.Winner != SideNone {
		return evaluate(b)
	}

	isMaximizing := b.Turn == SideWhite
	if isMaximizing {
		maxEval := -1000000.
		for m := range getAllMoves(b) {
			nextBoard := b.Apply(m)
			eval := alphaBeta(nextBoard, depth - 1, alpha, beta)
			maxEval = max(maxEval, eval)
			alpha = max(alpha, eval)
			if beta <= alpha {
				break
			}
		}
		return maxEval
	} else {
		minEval := 1000000.
		for m := range getAllMoves(b) {
			nextBoard := b.Apply(m)
			eval := alphaBeta(nextBoard, depth - 1, alpha, beta)
			minEval = min(minEval, eval)
			beta = min(beta, eval)
			if beta <= alpha {
				break
			}
		}
		return minEval
	}
}

func searchBestResponse(b *Board, out chan map[Move]Move, ctx context.Context) {
	var currentResult, lastResult map[Move]Move
	depth := 0
	search: for {
		depth += 1

		type movePair struct {
			move, response Move
		}

		results := make(chan movePair, 10)
		var wg sync.WaitGroup
		for m := range getAllMoves(b) {
			wg.Go(func() {
				nextBoard := b.Apply(m)
				bestScore := 1000000.
				var bestResponse Move
				for response := range getAllMoves(nextBoard) {
					score := alphaBeta(nextBoard.Apply(response), depth, -1000000., 1000000)
					if score < bestScore {
						bestResponse = response
						bestScore = score
					}
				}
				select {
				case <-ctx.Done():
					return
				case results <- movePair{ move: m, response: bestResponse }:
				}
			})
		}

		go func() {
			wg.Wait()
			close(results)
		}()

		currentResult = make(map[Move]Move)
		build: for {
			select {
			case <-ctx.Done():
				depth--
				break search

			case pair, ok := <-results:
				if !ok {
					break build
				}
				currentResult[pair.move] = pair.response
			}
		}
		lastResult = currentResult
	}
	out <- lastResult
	println("Depth", depth)
}

func (ai *Ai) launchSearch() {
	ai.responseCtx, ai.responseCancel = context.WithCancel(context.Background())
	go searchBestResponse(ai.board, ai.responseChannel, ai.responseCtx)
}
