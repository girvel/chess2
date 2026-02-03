package chess2

import (
	"fmt"
)

const BoardW int = 8
const BoardH int = 8

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
	inner [BoardW * BoardH]Piece
	Turn Side
	LastMove Move
}

func EmptyBoard() *Board {
	var result Board
	result.Turn = SideWhite

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
	if x < 0 || y < 0 || x >= BoardW || y >= BoardW {
		panic(fmt.Sprintf("attempt to access (%d, %d)", x, y))
	}
	return &b.inner[x + y * BoardW]
}

type Move struct {
	X1, Y1, X2, Y2 int
}

func NewMove(x1, y1, x2, y2 int) Move {
	return Move{ X1: x1, Y1: y1, X2: x2, Y2: y2 }
}

func (b *Board) Move(move Move) {
	if b.WillBeEnPassant(move) {
		*b.At(move.X2, move.Y1) = PieceNone;
	}
	*b.At(move.X2, move.Y2) = *b.At(move.X1, move.Y1)
	*b.At(move.X1, move.Y1) = PieceNone
	b.Turn = 1 - b.Turn
	b.LastMove = move
}

func (b *Board) WillBeEnPassant(m Move) bool {
	direction := int(1 - 2 * b.Turn)
	centerline := int(4 - b.Turn)
	if m.Y1 != centerline ||
		m.Y2 != centerline + direction ||
		Abs(m.X2 - m.X1) != 1 ||
		b.LastMove.Y1 != centerline + 2 * direction ||
		b.LastMove.X1 != m.X2 ||
		b.LastMove.Y2 != centerline ||
		*b.At(m.X2, m.Y2) != PieceNone {
		return false
	}

	neighbor := *b.At(m.X2, m.Y1)
	if b.Turn == SideWhite {
		return neighbor == PieceBlackPawn
	} else {
		return neighbor == PieceWhitePawn
	}
}

func (b *Board) IsMoveLegal(m Move) bool {
	if (m.X1 < 0 || m.X2 < 0 || m.Y1 < 0 || m.Y2 < 0 ||
		m.X1 >= BoardW || m.X2 >= BoardW || m.Y1 >= BoardH || m.Y2 >= BoardH) {
		return false
	}

	source := *b.At(m.X1, m.Y1)
	dest := *b.At(m.X2, m.Y2)
	if !source.Is(b.Turn) {
		return false
	}

	switch source {
	case PieceWhitePawn, PieceBlackPawn:
		direction := int(1 - 2 * b.Turn)
		if m.X2 == m.X1 &&
			m.Y2 == m.Y1 + direction &&
			dest == PieceNone {
			return true
		}

		baseline := int(1 + b.Turn * 5)
		if m.Y1 == baseline &&
			m.Y2 == baseline + 2 * direction &&
			m.X1 == m.X2 &&
			dest == PieceNone &&
			*b.At(m.X2, m.Y1 + direction) == PieceNone {
			return true
		}

		other_side := 1 - b.Turn
		if (m.X2 == m.X1 - 1 || m.X2 == m.X1 + 1) && 
			m.Y2 == m.Y1 + direction &&
			dest.Is(other_side) {
			return true
		}

		return b.WillBeEnPassant(m)

	case PieceWhiteKnight, PieceBlackKnight:
		x := m.X2 - m.X1
		y := m.Y2 - m.Y1
		return Abs(y) == 3 - Abs(x) && x != 0 && y != 0

	case PieceWhiteBishop, PieceBlackBishop:
		if Abs(m.X2 - m.X1) != Abs(m.Y2 - m.Y1) {
			return false
		}

		dx := Sign(m.X2 - m.X1)
		dy := Sign(m.Y2 - m.Y1)
		x := m.X1
		y := m.Y1
		for {
			x += dx
			y += dy
			piece := *b.At(x, y)
			if x == m.X2 {
				return !piece.Is(b.Turn)
			}
			if piece != PieceNone {
				return false
			}
		}
	
	case PieceWhiteRook, PieceBlackRook:
		if (m.X2 - m.X1 != 0) == (m.Y2 - m.Y1 != 0) {
			return false
		}

		dx := Sign(m.X2 - m.X1)
		dy := Sign(m.Y2 - m.Y1)
		x := m.X1
		y := m.Y1
		for {
			x += dx
			y += dy
			piece := *b.At(x, y)
			if x == m.X2 && y == m.Y2 {
				return !piece.Is(b.Turn)
			}
			if piece != PieceNone {
				return false
			}
		}
	}
	return false
}

