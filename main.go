package main

import (
	"context"

	rl "github.com/gen2brain/raylib-go/raylib"
	chess2 "github.com/girvel/chess2/src"
)

const scale int = 4
const cellSize int = 16
const totalCellSize int = scale * cellSize
const windowSize = chess2.BoardSize * totalCellSize

var colorWhiteSquare rl.Color = rl.GetColor(0xedededff)
var colorBlackSquare rl.Color = rl.GetColor(0x3a373dff)
var colorWhitePiece rl.Color = rl.GetColor(0xedededff)
var colorBlackPiece rl.Color = rl.GetColor(0x544747ff)
var colorSelected rl.Color = rl.GetColor(0xcfa867ff)

func LoadSprite(filepath string) rl.Texture2D {
	image := rl.LoadImage(filepath)
	defer rl.UnloadImage(image)
	rl.ImageResizeNN(image, image.Width * int32(scale), image.Height * int32(scale))
	return rl.LoadTextureFromImage(image)
}

func LoadSpriteColored(filepath string) rl.Texture2D {
	image := rl.LoadImage(filepath)
	defer rl.UnloadImage(image)
	rl.ImageColorReplace(image, colorWhitePiece, colorBlackPiece)
	rl.ImageResizeNN(image, image.Width * int32(scale), image.Height * int32(scale))
	return rl.LoadTextureFromImage(image)
}

func main() {
	rl.InitWindow(int32(windowSize), int32(windowSize), "girvel's chess app")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	board := chess2.EmptyBoard()

	sprites := []rl.Texture2D{
		LoadSprite("sprites/none.png"),
		LoadSprite("sprites/pawn.png"),
		LoadSpriteColored("sprites/pawn.png"),
		LoadSprite("sprites/knight.png"),
		LoadSpriteColored("sprites/knight.png"),
		LoadSprite("sprites/bishop.png"),
		LoadSpriteColored("sprites/bishop.png"),
		LoadSprite("sprites/rook.png"),
		LoadSpriteColored("sprites/rook.png"),
		LoadSprite("sprites/queen.png"),
		LoadSpriteColored("sprites/queen.png"),
		LoadSprite("sprites/king.png"),
		LoadSpriteColored("sprites/king.png"),
	}

	moveSprite := LoadSprite("sprites/move.png")
	winSprite := LoadSprite("sprites/win.png")
	lossSprite := LoadSprite("sprites/loss.png")

	var selectedX, selectedY int
	var potentialMoves []chess2.Move
	isSelected := false
	responseChannel := make(chan map[chess2.Move]chess2.Move)
	responseCtx, responseCancel := context.WithCancel(context.Background())
	go chess2.SearchBestResponse(*board, responseChannel, responseCtx)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		for x := range chess2.BoardSize {
			for y := range chess2.BoardSize {
				var squareColor rl.Color
				switch {
				case isSelected && x == selectedX && y == selectedY:
					squareColor = colorSelected
				case (x + y) % 2 == 0:
					squareColor = colorWhiteSquare
				default:
					squareColor = colorBlackSquare
				}

				render_x := int32(x * totalCellSize)
				render_y := int32(y * totalCellSize)
				rl.DrawRectangle(
					render_x, render_y,
					int32(totalCellSize), int32(totalCellSize),
					squareColor,
				)
				
				piece := *board.At(x, y)
				if piece != chess2.PieceNone {
					rl.DrawTexture(sprites[piece], render_x, render_y, rl.White)
				}
			}
		}

		if isSelected {
			for _, m := range potentialMoves {
				rl.DrawTexture(
					moveSprite,
					int32(m.X2 * totalCellSize), int32(m.Y2 * totalCellSize),
					rl.White,
				)
			}
		}

		if board.Winner != chess2.SideNone {
			var texture rl.Texture2D
			switch board.Winner {
			case chess2.SideWhite: texture = winSprite;
			case chess2.SideBlack: texture = lossSprite;
			}

			rl.DrawTexture(
				texture,
				(int32(windowSize) - texture.Width) / 2, (int32(windowSize) - texture.Height) / 2,
				rl.White,
			)
		}
		rl.EndDrawing()

		if board.Winner != chess2.SideNone {
			continue
		}

		if board.Turn == chess2.SideBlack {
			continue
		}

		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			x := int(rl.GetMouseX()) / totalCellSize
			y := int(rl.GetMouseY()) / totalCellSize

			if isSelected {
				move := chess2.NewMove(selectedX, selectedY, x, y)
				if board.IsMoveLegal(move) {
					board.Move(move)
					responseCancel()
					responses := <-responseChannel
					board.Move(responses[move])
					responseCtx, responseCancel = context.WithCancel(context.Background())
					go chess2.SearchBestResponse(*board, responseChannel, responseCtx)
				}
				isSelected = false
			} else {
				if board.At(x, y).Is(board.Turn) {
					selectedX = x
					selectedY = y
					isSelected = true
					potentialMoves = board.GetMoves(x, y)
				}
			}
		}
	}

	for _, sprite := range sprites {
		rl.UnloadTexture(sprite)
	}

	rl.UnloadTexture(moveSprite)
	rl.UnloadTexture(winSprite)
	rl.UnloadTexture(lossSprite)
}
