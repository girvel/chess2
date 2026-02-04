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

func evaluate(b Board) float64 {
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

func getAllMoves(b Board) iter.Seq[Move] {
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

func bestMove(b Board, depth int) Move {
	var result Move
	var bestScore float64
	isMaximizing := b.Turn == SideWhite
	if isMaximizing {
		bestScore = -1000000
	} else {
		bestScore = 1000000
	}

	for m := range getAllMoves(b) {
		nextBoard := *b.Apply(m)
		if depth > 0 {
			nextBoard = *nextBoard.Apply(bestMove(nextBoard, depth - 1))
		}
		score := evaluate(nextBoard)

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

func searchBestResponse(b Board, out chan map[Move]Move, ctx context.Context) {
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
				results <- movePair{
					move: m,
					response: bestMove(*b.Apply(m), depth),
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
	go searchBestResponse(*ai.board, ai.responseChannel, ai.responseCtx)
}
