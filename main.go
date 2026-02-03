package main

import rl "github.com/gen2brain/raylib-go/raylib"

const scale int = 4
const cellSize int = 16
const totalCellSize int = scale * cellSize

const w int = 8
const h int = 8

var white rl.Color = rl.GetColor(0xedededff)
var black rl.Color = rl.GetColor(0x3a373dff)

type Piece int
const (
	PieceNone Piece = iota
	PieceWhitePawn
	PieceBlackPawn
	PieceWhiteKnight
	PieceBlackKnight
	PieceWhiteBishop
	PieceBlackBishop
	PieceWhiteRook
	PieceBlackRook
	PieceWhiteQueen
	PieceBlackQueen
	PieceWhiteKing
	PieceBlackKing
)

type Board struct {
	inner [w * h]Piece
}

func EmptyBoard() *Board {
	var result Board

	*result.At(0, 0) = PieceBlackRook
	*result.At(1, 0) = PieceBlackKnight
	*result.At(2, 0) = PieceBlackBishop
	*result.At(3, 0) = PieceBlackKing
	*result.At(4, 0) = PieceBlackQueen
	*result.At(5, 0) = PieceBlackBishop
	*result.At(6, 0) = PieceBlackKnight
	*result.At(7, 0) = PieceBlackRook

	*result.At(0, 1) = PieceBlackPawn
	*result.At(1, 1) = PieceBlackPawn
	*result.At(2, 1) = PieceBlackPawn
	*result.At(3, 1) = PieceBlackPawn
	*result.At(4, 1) = PieceBlackPawn
	*result.At(5, 1) = PieceBlackPawn
	*result.At(6, 1) = PieceBlackPawn
	*result.At(7, 1) = PieceBlackPawn

	*result.At(0, 6) = PieceWhitePawn
	*result.At(1, 6) = PieceWhitePawn
	*result.At(2, 6) = PieceWhitePawn
	*result.At(3, 6) = PieceWhitePawn
	*result.At(4, 6) = PieceWhitePawn
	*result.At(5, 6) = PieceWhitePawn
	*result.At(6, 6) = PieceWhitePawn
	*result.At(7, 6) = PieceWhitePawn

	*result.At(0, 7) = PieceWhiteRook
	*result.At(1, 7) = PieceWhiteKnight
	*result.At(2, 7) = PieceWhiteBishop
	*result.At(3, 7) = PieceWhiteKing
	*result.At(4, 7) = PieceWhiteQueen
	*result.At(5, 7) = PieceWhiteBishop
	*result.At(6, 7) = PieceWhiteKnight
	*result.At(7, 7) = PieceWhiteRook

	return &result
}

func (b *Board) At(x, y int) *Piece {
	return &b.inner[x + y * w]
}

func LoadWhitePiece(filepath string) rl.Texture2D {
	image := rl.LoadImage(filepath)
	defer rl.UnloadImage(image)
	rl.ImageResizeNN(image, image.Width * int32(scale), image.Height * int32(scale))
	return rl.LoadTextureFromImage(image)
}

func LoadBlackPiece(filepath string) rl.Texture2D {
	image := rl.LoadImage(filepath)
	defer rl.UnloadImage(image)
	rl.ImageColorReplace(image, white, black)
	rl.ImageResizeNN(image, image.Width * int32(scale), image.Height * int32(scale))
	return rl.LoadTextureFromImage(image)
}

func main() {
	rl.InitWindow(int32(w * totalCellSize), int32(h * totalCellSize), "girvel's chess app")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	board := EmptyBoard()

	for !rl.WindowShouldClose() {
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
				
				piece := *board.At(x, y)
				if piece != PieceNone {
					rl.DrawTexture(sprites[piece], render_x, render_y, rl.White)
				}
			}
		}
		rl.EndDrawing()

		for _, sprite := range sprites {
			rl.UnloadTexture(sprite)
		}
	}
}
