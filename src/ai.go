package chess2

import (
	"context"
	"slices"
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
	if cost == nil {
		var empty [BoardSize * BoardSize]float64
		cost = append(cost, empty)

		for _, e := range rawCost {
			cost = append(cost, e)
			var mirror [BoardSize * BoardSize]float64
			for i, v := range e {
				x := i % BoardSize
				y := BoardSize - i / BoardSize - 1
				mirror[x + y * BoardSize] = -v
			}
			cost = append(cost, mirror)
		}
	}

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

var rawCost = [][BoardSize * BoardSize]float64{
	{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 
	 2.0, 2.0, 2.0, 2.0, 2.0, 2.0, 2.0, 2.0, 
	 1.5, 1.5, 1.5, 1.5, 1.5, 1.5, 1.5, 1.5, 
	 1.1, 1.0, 1.0, 1.5, 1.5, 1.0, 1.0, 1.1, 
	 0.7, 1.0, 1.0, 1.5, 1.5, 1.0, 1.0, 0.7, 
	 0.9, 1.0, 1.0, 1.3, 1.3, 1.0, 1.0, 0.9, 
	 0.8, 1.0, 1.0, 1.2, 1.2, 1.0, 1.0, 0.8, 
	 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,},

	{2.3, 2.3, 2.3, 2.3, 2.3, 2.3, 2.3, 2.3, 
	 2.3, 2.5, 2.5, 2.5, 2.5, 2.5, 2.5, 2.3, 
	 2.3, 2.5, 2.5, 2.5, 2.5, 2.5, 2.5, 2.3, 
	 2.3, 2.5, 2.5, 2.8, 2.8, 2.5, 2.5, 2.3, 
	 2.3, 2.5, 2.5, 2.8, 2.8, 2.5, 2.5, 2.3, 
	 2.3, 2.5, 2.5, 2.5, 2.5, 2.5, 2.5, 2.3, 
	 2.3, 2.5, 2.5, 2.5, 2.5, 2.5, 2.5, 2.3, 
	 2.3, 2.3, 2.3, 2.3, 2.3, 2.3, 2.3, 2.3,},

	{3.0, 3.0, 3.0, 3.0, 3.0, 3.0, 3.0, 3.0, 
	 3.0, 3.2, 3.0, 3.0, 3.0, 3.0, 3.2, 3.0, 
	 3.0, 3.0, 3.2, 3.0, 3.0, 3.2, 3.0, 3.0, 
	 3.0, 3.0, 3.0, 3.2, 3.2, 3.0, 3.0, 3.0, 
	 3.0, 3.0, 3.0, 3.2, 3.2, 3.0, 3.0, 3.0, 
	 3.0, 3.0, 3.2, 3.0, 3.0, 3.2, 3.0, 3.0, 
	 3.0, 3.2, 3.0, 3.0, 3.0, 3.0, 3.2, 3.0, 
	 3.0, 2.8, 2.8, 2.8, 2.8, 2.8, 2.8, 3.0,},

	{5.0, 5.0, 5.0, 5.0, 5.0, 5.0, 5.0, 5.0, 
	 5.5, 5.5, 5.5, 5.5, 5.5, 5.5, 5.5, 5.5, 
	 5.0, 5.0, 5.0, 5.0, 5.0, 5.0, 5.0, 5.0, 
	 5.0, 5.0, 5.0, 5.0, 5.0, 5.0, 5.0, 5.0, 
	 5.0, 5.0, 5.0, 5.0, 5.0, 5.0, 5.0, 5.0, 
	 5.0, 5.0, 5.0, 5.0, 5.0, 5.0, 5.0, 5.0, 
	 5.0, 5.0, 5.0, 5.0, 5.0, 5.0, 5.0, 5.0, 
	 4.8, 5.0, 5.0, 5.0, 5.0, 5.0, 5.0, 4.8,},

	{9.0, 9.0, 9.0, 9.0, 9.0, 9.0, 9.0, 9.0, 
	 9.0, 9.0, 9.0, 9.0, 9.0, 9.0, 9.0, 9.0, 
	 9.0, 9.0, 9.0, 9.0, 9.0, 9.0, 9.0, 9.0, 
	 9.0, 9.0, 9.0, 9.0, 9.0, 9.0, 9.0, 9.0, 
	 9.0, 9.0, 9.0, 9.0, 9.0, 9.0, 9.0, 9.0, 
	 9.0, 9.0, 9.0, 9.0, 9.0, 9.0, 9.0, 9.0, 
	 9.0, 9.0, 9.0, 9.0, 9.0, 9.0, 9.0, 9.0, 
	 9.0, 9.0, 9.0, 9.0, 9.0, 9.0, 9.0, 9.0,},

	{1000.0, 1000.0, 1000.0, 1000.0, 1000.0, 1000.0, 1000.0, 1000.0, 
	 1000.0, 1000.0, 1000.0, 1000.0, 1000.0, 1000.0, 1000.0, 1000.0, 
	 1000.0, 1000.0, 1000.0, 1000.0, 1000.0, 1000.0, 1000.0, 1000.0, 
	 1000.0, 1000.0, 1000.0, 1000.0, 1000.0, 1000.0, 1000.0, 1000.0, 
	 1000.0, 1000.0, 1000.0, 1000.0, 1000.0, 1000.0, 1000.0, 1000.0, 
	 1000.0, 1000.0, 1000.0, 1000.0, 1000.0, 1000.0, 1000.0, 1000.0, 
	 1000.0, 1000.0, 1000.0, 1000.0, 1000.0, 1000.0, 1000.0, 1000.0, 
	 1001.0, 1000.5, 1000.3, 1000.0, 1000.0, 1000.0, 1000.5, 1001.0,},
}

var cost [][BoardSize * BoardSize]float64

func evaluate(b *Board) float64 {
	switch b.Winner {
	case SideWhite: return 1000
	case SideBlack: return -1000
	}

	var result float64
	for i, piece := range b.inner {
		result += cost[piece][i]
	}

	return result
}

func getAllMoves(b *Board) []Move {
	result := make([]Move, 0, 16)

	for x := range BoardSize {
		for y := range BoardSize {
			piece := *b.At(x, y)
			if piece.Side() != b.Turn {
				continue
			}

			result = append(result, b.GetMoves(x, y)...)
		}
	}

	score := func(m Move) float64 {
		capture := *b.At(m.X2, m.Y2)
		if capture == PieceNone {
			return 0
		}

		attacker := *b.At(m.X1, m.Y1)
		return 1000000 + 100 * Abs(cost[capture][m.X2 + BoardSize * m.Y2]) - Abs(cost[attacker][m.X1 + BoardSize * m.Y1])
	}

	slices.SortFunc(result, func(a, b Move) int { return Sign(score(b) - score(a)) })

	return result
}

func alphaBeta(b *Board, depth int, alpha, beta float64, onlyCaptures bool) float64 {
	if depth <= 0 || b.Winner != SideNone {
		return evaluate(b)
	}

	isMaximizing := b.Turn == SideWhite
	if isMaximizing {
		maxEval := -1000000.
		for _, m := range getAllMoves(b) {
			if onlyCaptures && !m.IsCapture(b) { continue }
			nextBoard := b.Apply(m)
			var eval float64
			if depth == 1 && m.IsCapture(b) {
				eval = alphaBeta(nextBoard, 1, alpha, beta, true)
			} else {
				eval = alphaBeta(nextBoard, depth - 1, alpha, beta, onlyCaptures)
			}
			maxEval = max(maxEval, eval)
			alpha = max(alpha, eval)
			if beta <= alpha {
				break
			}
		}
		return maxEval
	} else {
		minEval := 1000000.
		for _, m := range getAllMoves(b) {
			if onlyCaptures && !m.IsCapture(b) { continue }
			nextBoard := b.Apply(m)
			var eval float64
			if depth == 1 && m.IsCapture(b) {
				eval = alphaBeta(nextBoard, 1, alpha, beta, true)
			} else {
				eval = alphaBeta(nextBoard, depth - 1, alpha, beta, onlyCaptures)
			}
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
		for _, m := range getAllMoves(b) {
			wg.Go(func() {
				nextBoard := b.Apply(m)
				bestScore := 1000000.
				var bestResponse Move
				for _, response := range getAllMoves(nextBoard) {
					score := alphaBeta(nextBoard.Apply(response), depth, -1000000., 1000000, false)
					if score < bestScore {
						bestResponse = response
						bestScore = score
					}
				}
				select {
				case <-ctx.Done(): if depth > 1 { return }
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
				if depth > 1 {
					depth--
					break search
				}

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
