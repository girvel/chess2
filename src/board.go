package chess2

import (
	"fmt"
)

const BoardSize int = 8

type Side int

const SideBlack Side = 0
const SideWhite Side = 1
const SideNone Side = -1

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

func (p Piece) Side() Side {
	switch p {
		case PieceNone: return SideNone
		default: return Side(p % 2)
	}
}

type Board struct {
	inner [BoardSize * BoardSize]Piece
	Turn Side
	LastMove Move
	A1Moved, A8Moved, E1Moved, E8Moved, H1Moved, H8Moved bool
	Winner Side
}

func EmptyBoard() *Board {
	var result Board
	result.Turn = SideWhite
	result.Winner = SideNone

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
	if x < 0 || y < 0 || x >= BoardSize || y >= BoardSize {
		panic(fmt.Sprintf("attempt to access (%d, %d)", x, y))
	}
	return &b.inner[x + y * BoardSize]
}

type Move struct {
	X1, Y1, X2, Y2 int
}

func NewMove(x1, y1, x2, y2 int) Move {
	return Move{ X1: x1, Y1: y1, X2: x2, Y2: y2 }
}

func (b *Board) Move(move Move) {
	switch {
	case b.WillBeEnPassant(move):
		*b.At(move.X2, move.Y1) = PieceNone;

	case b.WillBeCastle(move):
		direction := Sign(move.X2 - move.X1)
		var rookX int
		if direction < 0 {
			rookX = 0
		} else {
			rookX = 7
		}
		*b.At(move.X2 - direction, move.Y2) = *b.At(rookX, move.Y2);
		*b.At(rookX, move.Y2) = PieceNone;
	}

	source := b.At(move.X1, move.Y1)
	dest := b.At(move.X2, move.Y2)

	switch *dest {
	case PieceWhiteKing: b.Winner = SideBlack
	case PieceBlackKing: b.Winner = SideWhite
	}

	*dest = *source
	*source = PieceNone
	b.Turn = 1 - b.Turn
	b.LastMove = move

	switch {
	case move.X1 == 0 && move.Y1 == 0: b.A8Moved = true;
	case move.X1 == 4 && move.Y1 == 0: b.E8Moved = true;
	case move.X1 == 7 && move.Y1 == 0: b.H8Moved = true;
	case move.X1 == 0 && move.Y1 == 7: b.A1Moved = true;
	case move.X1 == 4 && move.Y1 == 7: b.E1Moved = true;
	case move.X1 == 7 && move.Y1 == 7: b.H1Moved = true;
	}
}

// TODO split detection & validation
func (b *Board) WillBeEnPassant(m Move) bool {
	source := *b.At(m.X1, m.Y1)
	if source != PieceWhitePawn && source != PieceBlackPawn {
		return false
	}

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

func (b *Board) WillBeCastle(m Move) bool {
	// NEXT check attack on passed square
	var backline int
	this_side := b.Turn
	opposite_side := 1 - b.Turn
	direction := Sign(m.X2 - m.X1)
	if b.Turn == SideWhite {
		if b.H1Moved || b.E1Moved {
			return false
		}
		backline = 7
	} else {
		if b.H8Moved || b.E8Moved {
			return false
		}
		backline = 0
	}

	if (m != NewMove(4, backline, 4 + 2 * direction, backline) ||
		*b.At(4 + direction, backline) != PieceNone ||
		*b.At(4 + 2 * direction, backline) != PieceNone) {
		return false
	}

	b.Turn = opposite_side
	defer func() {b.Turn = this_side}()
	for x := range BoardSize {
		for y := range BoardSize {
			piece := *b.At(x, y)
			if piece.Is(opposite_side) && b.IsMoveLegal(NewMove(x, y, 4 + direction, backline)) {
				return false
			}
		}
	}

	return true
}

// TODO should it check for turn?
func (b *Board) IsMoveLegal(m Move) bool {
	if b.Winner != SideNone {
		return false
	}

	if (m.X1 < 0 || m.X2 < 0 || m.Y1 < 0 || m.Y2 < 0 ||
		m.X1 >= BoardSize || m.X2 >= BoardSize || m.Y1 >= BoardSize || m.Y2 >= BoardSize) {
		return false
	}

	if b.WillBeEnPassant(m) ||
		b.WillBeCastle(m) {
		return true
	}

	source := *b.At(m.X1, m.Y1)
	dest := *b.At(m.X2, m.Y2)
	if !source.Is(b.Turn) || dest.Is(b.Turn) {
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

		if (m.X2 == m.X1 - 1 || m.X2 == m.X1 + 1) && 
			m.Y2 == m.Y1 + direction {
			return true
		}

		return false
	}

	ox := m.X2 - m.X1
	oy := m.Y2 - m.Y1

	switch source {
	case PieceWhiteKnight, PieceBlackKnight:
		return Abs(oy) == 3 - Abs(ox) && ox != 0 && oy != 0

	case PieceWhiteBishop, PieceBlackBishop:
		if Abs(ox) != Abs(oy) {
			return false
		}
	
	case PieceWhiteRook, PieceBlackRook:
		if (ox != 0) == (oy != 0) {
			return false
		}
	
	case PieceWhiteKing, PieceBlackKing:
		return Abs(ox) * Abs(oy) <= 1
	
	case PieceWhiteQueen, PieceBlackQueen:
		if Abs(ox) != Abs(oy) && (ox != 0) == (oy != 0) {
			return false
		}
	}

	dx := Sign(ox)
	dy := Sign(oy)
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

func (b *Board) GetMoves(x, y int) []Move {
	var potential []Move
	source := *b.At(x, y)
	switch source {
	case PieceWhitePawn, PieceBlackPawn:
		direction := int(1 - source.Side() * 2)
		potential = append(potential, NewMove(x, y, x, y + direction))
		potential = append(potential, NewMove(x, y, x, y + 2 * direction))
		potential = append(potential, NewMove(x, y, x + 1, y + direction))
		potential = append(potential, NewMove(x, y, x - 1, y + direction))

	case PieceWhiteKnight, PieceBlackKnight:
		potential = append(potential, NewMove(x, y, x + 1, y + 2))
		potential = append(potential, NewMove(x, y, x + 1, y - 2))
		potential = append(potential, NewMove(x, y, x - 1, y + 2))
		potential = append(potential, NewMove(x, y, x - 1, y - 2))
		potential = append(potential, NewMove(x, y, x + 2, y + 1))
		potential = append(potential, NewMove(x, y, x + 2, y - 1))
		potential = append(potential, NewMove(x, y, x - 2, y + 1))
		potential = append(potential, NewMove(x, y, x - 2, y - 1))

	case PieceWhiteBishop, PieceBlackBishop:
		for v := range BoardSize {
			potential = append(potential, NewMove(x, y, v, v - x + y))
			potential = append(potential, NewMove(x, y, v, -v + x + y))
		}

	case PieceWhiteRook, PieceBlackRook:
		for v := range BoardSize {
			potential = append(potential, NewMove(x, y, v, y))
			potential = append(potential, NewMove(x, y, x, v))
		}
	
	case PieceWhiteQueen, PieceBlackQueen:
		for v := range BoardSize {
			potential = append(potential, NewMove(x, y, v, y))
			potential = append(potential, NewMove(x, y, x, v))
			potential = append(potential, NewMove(x, y, v, v - x + y))
			potential = append(potential, NewMove(x, y, v, -v + x + y))
		}
	
	case PieceWhiteKing, PieceBlackKing:
		potential = append(potential, NewMove(x, y, x + 1, y + 1))
		potential = append(potential, NewMove(x, y, x + 1, y))
		potential = append(potential, NewMove(x, y, x + 2, y))
		potential = append(potential, NewMove(x, y, x + 1, y - 1))
		potential = append(potential, NewMove(x, y, x, y + 1))
		potential = append(potential, NewMove(x, y, x, y - 1))
		potential = append(potential, NewMove(x, y, x - 1, y + 1))
		potential = append(potential, NewMove(x, y, x - 1, y))
		potential = append(potential, NewMove(x, y, x - 2, y))
		potential = append(potential, NewMove(x, y, x - 1, y - 1))
	}

	var result []Move = make([]Move, 0, len(potential))
	for _, m := range potential {
		if b.IsMoveLegal(m) {
			result = append(result, m)
		}
	}

	return result
}

