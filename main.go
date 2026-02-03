package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	chess2 "github.com/girvel/chess2/src"
)

const scale int = 4
const cellSize int = 16
const totalCellSize int = scale * cellSize

var colorWhite rl.Color = rl.GetColor(0xedededff)
var colorBlack rl.Color = rl.GetColor(0x3a373dff)
var colorSelected rl.Color = rl.GetColor(0xcfa867ff)

func LoadWhitePiece(filepath string) rl.Texture2D {
	image := rl.LoadImage(filepath)
	defer rl.UnloadImage(image)
	rl.ImageResizeNN(image, image.Width * int32(scale), image.Height * int32(scale))
	return rl.LoadTextureFromImage(image)
}

func LoadBlackPiece(filepath string) rl.Texture2D {
	image := rl.LoadImage(filepath)
	defer rl.UnloadImage(image)
	rl.ImageColorReplace(image, colorWhite, colorBlack)
	rl.ImageResizeNN(image, image.Width * int32(scale), image.Height * int32(scale))
	return rl.LoadTextureFromImage(image)
}

func main() {
	rl.InitWindow(int32(chess2.BoardW * totalCellSize), int32(chess2.BoardH * totalCellSize), "girvel's chess app")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	board := chess2.EmptyBoard()

	sprites := []rl.Texture2D{
		LoadWhitePiece("sprites/none.png"),
		LoadWhitePiece("sprites/pawn.png"),
		LoadBlackPiece("sprites/pawn.png"),
		LoadWhitePiece("sprites/knight.png"),
		LoadBlackPiece("sprites/knight.png"),
		LoadWhitePiece("sprites/bishop.png"),
		LoadBlackPiece("sprites/bishop.png"),
		LoadWhitePiece("sprites/rook.png"),
		LoadBlackPiece("sprites/rook.png"),
		LoadWhitePiece("sprites/queen.png"),
		LoadBlackPiece("sprites/queen.png"),
		LoadWhitePiece("sprites/king.png"),
		LoadBlackPiece("sprites/king.png"),
	}

	var selectedX, selectedY int
	isSelected := false

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		for x := range chess2.BoardW {
			for y := range chess2.BoardH {
				var squareColor rl.Color
				switch {
				case isSelected && x == selectedX && y == selectedY:
					squareColor = colorSelected
				case (x + y) % 2 == 0:
					squareColor = colorWhite
				default:
					squareColor = colorBlack
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
		rl.EndDrawing()

		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			x := int(rl.GetMouseX()) / totalCellSize
			y := int(rl.GetMouseY()) / totalCellSize

			if isSelected {
				move := chess2.NewMove(selectedX, selectedY, x, y)
				if board.IsMoveLegal(move) {
					board.Move(move)
				}
				isSelected = false
			} else {
				if board.At(x, y).Is(board.Turn) {
					selectedX = x
					selectedY = y
					isSelected = true
				}
			}
		}
	}

	for _, sprite := range sprites {
		rl.UnloadTexture(sprite)
	}
}
