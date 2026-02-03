package main

import rl "github.com/gen2brain/raylib-go/raylib"

const scale int = 4
const cellSize int = 16
const totalCellSize int = scale * cellSize

const w int = 8
const h int = 8

type Piece int
const (
	PieceNone Piece = iota
	PiecePawn
)

type Board struct {
	inner [w * h]Piece
}

func (b *Board) At(x, y int) *Piece {
	return &b.inner[x + y * w]
}

func LoadSprite(filepath string) rl.Texture2D {
	image := rl.LoadImage(filepath)
	defer rl.UnloadImage(image)
	rl.ImageResizeNN(image, image.Width * int32(scale), image.Height * int32(scale))
	return rl.LoadTextureFromImage(image)
}

func main() {
	rl.InitWindow(int32(w * totalCellSize), int32(h * totalCellSize), "girvel's chess app")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	sprites := []rl.Texture2D{
		LoadSprite("sprites/none.png"),
		LoadSprite("sprites/pawn.png"),
	}

	var board Board
	*board.At(0, 6) = PiecePawn

	white := rl.GetColor(0xb5bdbaff)
	black := rl.GetColor(0x3a373dff)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		for x := range w {
			for y := range h {
				var squareColor rl.Color
				if (x + y) % 2 == 0 {
					squareColor = white
				} else {
					squareColor = black
				}
				render_x := int32(x * totalCellSize)
				render_y := int32(y * totalCellSize)
				rl.DrawRectangle(
					render_x, render_y,
					int32(totalCellSize), int32(totalCellSize),
					squareColor,
				)
				rl.DrawTexture(sprites[*board.At(x, y)], render_x, render_y, rl.White)
			}
		}
		rl.EndDrawing()
	}
}
