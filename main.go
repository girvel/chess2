package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func sign(x int) int {
	switch {
	case x == 0: return 0
	case x < 0: return -1
	default: return 1
	}
}

const scale int = 4
const cellSize int = 16
const totalCellSize int = scale * cellSize

const w int = 8
const h int = 8

var colorWhite rl.Color = rl.GetColor(0xedededff)
var colorBlack rl.Color = rl.GetColor(0x3a373dff)
var colorSelected rl.Color = rl.GetColor(0xcfa867ff)

type Side int

const SideBlack Side = 0
const SideWhite Side = 1

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

func (p Piece) Is(side Side) bool {
	return p != PieceNone && int(p) % 2 == int(side)
}

type Board struct {
	inner [w * h]Piece
	turn Side
	lastMove Move
}

func EmptyBoard() *Board {
	var result Board
	result.turn = SideWhite

	*result.At(0, 0) = PieceBlackRook
	*result.At(1, 0) = PieceBlackKnight
	*result.At(2, 0) = PieceBlackBishop
	*result.At(3, 0) = PieceBlackQueen
	*result.At(4, 0) = PieceBlackKing
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
	*result.At(3, 7) = PieceWhiteQueen
	*result.At(4, 7) = PieceWhiteKing
	*result.At(5, 7) = PieceWhiteBishop
	*result.At(6, 7) = PieceWhiteKnight
	*result.At(7, 7) = PieceWhiteRook

	return &result
}

func (b *Board) At(x, y int) *Piece {
	if x < 0 || y < 0 || x >= w || y >= w {
		panic(fmt.Sprintf("attempt to access (%d, %d)", x, y))
	}
	return &b.inner[x + y * w]
}

type Move struct {
	x1, y1, x2, y2 int
}

func (b *Board) Move(move Move) {
	if b.WillBeEnPassant(move) {
		*b.At(move.x2, move.y1) = PieceNone;
	}
	*b.At(move.x2, move.y2) = *b.At(move.x1, move.y1)
	*b.At(move.x1, move.y1) = PieceNone
	b.turn = 1 - b.turn
	b.lastMove = move
}

func (b *Board) WillBeEnPassant(m Move) bool {
	direction := int(1 - 2 * b.turn)
	centerline := int(4 - b.turn)
	if m.y1 != centerline ||
		m.y2 != centerline + direction ||
		abs(m.x2 - m.x1) != 1 ||
		b.lastMove.y1 != centerline + 2 * direction ||
		b.lastMove.x1 != m.x2 ||
		b.lastMove.y2 != centerline ||
		*b.At(m.x2, m.y2) != PieceNone {
		return false
	}

	neighbor := *b.At(m.x2, m.y1)
	if b.turn == SideWhite {
		return neighbor == PieceBlackPawn
	} else {
		return neighbor == PieceWhitePawn
	}
}

func (b *Board) IsMoveLegal(m Move) bool {
	if (m.x1 < 0 || m.x2 < 0 || m.y1 < 0 || m.y2 < 0 ||
		m.x1 >= w || m.x2 >= w || m.y1 >= h || m.y2 >= h) {
		return false
	}

	source := *b.At(m.x1, m.y1)
	dest := *b.At(m.x2, m.y2)
	if !source.Is(b.turn) {
		return false
	}

	switch source {
	case PieceWhitePawn, PieceBlackPawn:
		direction := int(1 - 2 * b.turn)
		if m.x2 == m.x1 &&
			m.y2 == m.y1 + direction &&
			dest == PieceNone {
			return true
		}

		baseline := int(1 + b.turn * 5)
		if m.y1 == baseline &&
			m.y2 == baseline + 2 * direction &&
			m.x1 == m.x2 &&
			dest == PieceNone &&
			*b.At(m.x2, m.y1 + direction) == PieceNone {
			return true
		}

		other_side := 1 - b.turn
		if (m.x2 == m.x1 - 1 || m.x2 == m.x1 + 1) && 
			m.y2 == m.y1 + direction &&
			dest.Is(other_side) {
			return true
		}

		return b.WillBeEnPassant(m)

	case PieceWhiteKnight, PieceBlackKnight:
		x := m.x2 - m.x1
		y := m.y2 - m.y1
		return abs(y) == 3 - abs(x) && x != 0 && y != 0

	case PieceWhiteBishop, PieceBlackBishop:
		if abs(m.x2 - m.x1) != abs(m.y2 - m.y1) {
			return false
		}

		dx := sign(m.x2 - m.x1)
		dy := sign(m.y2 - m.y1)
		x := m.x1
		y := m.y1
		for {
			x += dx
			y += dy
			piece := *b.At(x, y)
			if x == m.x2 {
				return !piece.Is(b.turn)
			}
			if piece != PieceNone {
				return false
			}
		}
	
	case PieceWhiteRook, PieceBlackRook:
		if (m.x2 - m.x1 != 0) == (m.y2 - m.y1 != 0) {
			return false
		}

		dx := sign(m.x2 - m.x1)
		dy := sign(m.y2 - m.y1)
		x := m.x1
		y := m.y1
		for {
			x += dx
			y += dy
			piece := *b.At(x, y)
			if x == m.x2 && y == m.y2 {
				return !piece.Is(b.turn)
			}
			if piece != PieceNone {
				return false
			}
		}
	}
	return false
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
	rl.ImageColorReplace(image, colorWhite, colorBlack)
	rl.ImageResizeNN(image, image.Width * int32(scale), image.Height * int32(scale))
	return rl.LoadTextureFromImage(image)
}

func main() {
	rl.InitWindow(int32(w * totalCellSize), int32(h * totalCellSize), "girvel's chess app")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	board := EmptyBoard()

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
		for x := range w {
			for y := range h {
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
				if piece != PieceNone {
					rl.DrawTexture(sprites[piece], render_x, render_y, rl.White)
				}
			}
		}
		rl.EndDrawing()

		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			x := int(rl.GetMouseX()) / totalCellSize
			y := int(rl.GetMouseY()) / totalCellSize

			if isSelected {
				move := Move{selectedX, selectedY, x, y}
				if board.IsMoveLegal(move) {
					board.Move(move)
				}
				isSelected = false
			} else {
				if *board.At(x, y) != PieceNone {
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
