package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	chess2 "github.com/girvel/chess2/src"
	"github.com/girvel/chess2/src/iosystem"
)

func main() {
	iosystem.Init()
	defer iosystem.Deinit()
	board := chess2.EmptyBoard()
	ai := chess2.CreateAi(*board)

	for {
		iosystem.Draw(*board)

		playerMove, shouldClose := iosystem.ReadInput(*board)
		if playerMove != nil {
			board.Move(*playerMove)
			ai.PushMove(*playerMove)
		}

		if shouldClose {
			break
		}

		if board.Turn == chess2.SideBlack {
			if m := ai.PopResponse(); m != nil {
				rl.TraceLog(rl.LogInfo, "AI: %s", m)
				board.Move(*m)
			}
		}
	}
}
