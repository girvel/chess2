package iosystem

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	chess2 "github.com/girvel/chess2/src"
)

const scale int = 6
const cellSize int = 16
const totalCellSize int = scale * cellSize
const windowSize = chess2.BoardSize * totalCellSize

var colorWhiteSquare rl.Color = rl.GetColor(0xedededff)
var colorBlackSquare rl.Color = rl.GetColor(0x3a373dff)
var colorWhitePiece rl.Color = rl.GetColor(0xedededff)
var colorBlackPiece rl.Color = rl.GetColor(0x544747ff)
var colorSelected rl.Color = rl.GetColor(0xcfa867ff)
var colorLastMoveDark rl.Color = rl.GetColor(0x5d863fff)
var colorLastMoveLight rl.Color = rl.GetColor(0x869d42ff)

var pieceSprites []rl.Texture2D
var moveSprite, winSprite, lossSprite rl.Texture2D

var selectedX, selectedY int
var potentialMoves []chess2.Move
var isSelected = false

func Init() {
	rl.InitWindow(int32(windowSize), int32(windowSize), "girvel's chess app")
	rl.SetTargetFPS(60)

	pieceSprites = []rl.Texture2D{
		loadSprite("sprites/none.png"),
		loadSprite("sprites/pawn.png"),
		loadSpriteColored("sprites/pawn.png"),
		loadSprite("sprites/knight.png"),
		loadSpriteColored("sprites/knight.png"),
		loadSprite("sprites/bishop.png"),
		loadSpriteColored("sprites/bishop.png"),
		loadSprite("sprites/rook.png"),
		loadSpriteColored("sprites/rook.png"),
		loadSprite("sprites/queen.png"),
		loadSpriteColored("sprites/queen.png"),
		loadSprite("sprites/king.png"),
		loadSpriteColored("sprites/king.png"),
	}

	moveSprite = loadSprite("sprites/move.png")
	winSprite = loadSprite("sprites/win.png")
	lossSprite = loadSprite("sprites/loss.png")
}

func Draw(board *chess2.Board) {
	rl.BeginDrawing()
	for x := range chess2.BoardSize {
		for y := range chess2.BoardSize {
			var squareColor rl.Color
			switch {
			case isSelected && x == selectedX && y == selectedY:
				squareColor = colorSelected
			case (x + y) % 2 == 0:
				if board.LastMove != nil && (
					x == board.LastMove.X1 && y == board.LastMove.Y1 ||
					x == board.LastMove.X2 && y == board.LastMove.Y2) {
					squareColor = colorLastMoveLight
				} else {
					squareColor = colorWhiteSquare
				}
			default:
				if board.LastMove != nil && (
					x == board.LastMove.X1 && y == board.LastMove.Y1 ||
					x == board.LastMove.X2 && y == board.LastMove.Y2) {
					squareColor = colorLastMoveDark
				} else {
					squareColor = colorBlackSquare
				}
			}

			renderX := int32(x * totalCellSize)
			renderY := int32(y * totalCellSize)
			rl.DrawRectangle(
				renderX, renderY,
				int32(totalCellSize), int32(totalCellSize),
				squareColor,
			)
			
			piece := *board.At(x, y)
			if piece != chess2.PieceNone {
				rl.DrawTexture(pieceSprites[piece], renderX, renderY, rl.White)
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
}

func ReadInput(board *chess2.Board) (*chess2.Move, bool) {
	shouldClose := rl.WindowShouldClose()
	if board.Winner != chess2.SideNone || board.Turn == chess2.SideBlack {
		return nil, shouldClose
	}

	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		x := int(rl.GetMouseX()) / totalCellSize
		y := int(rl.GetMouseY()) / totalCellSize

		if isSelected {
			isSelected = false
			move := chess2.NewMove(selectedX, selectedY, x, y)
			if board.IsMoveLegal(move) {
				return &move, shouldClose
			}
		} else {
			if board.At(x, y).Is(board.Turn) {
				selectedX = x
				selectedY = y
				isSelected = true
				potentialMoves = board.GetMoves(x, y)
			}
		}
	}

	return nil, shouldClose
}

func Deinit() {
	for _, sprite := range pieceSprites {
		rl.UnloadTexture(sprite)
	}

	rl.UnloadTexture(moveSprite)
	rl.UnloadTexture(winSprite)
	rl.UnloadTexture(lossSprite)
	rl.CloseWindow()
}

func loadSprite(filepath string) rl.Texture2D {
	image := rl.LoadImage(filepath)
	defer rl.UnloadImage(image)
	rl.ImageResizeNN(image, image.Width * int32(scale), image.Height * int32(scale))
	return rl.LoadTextureFromImage(image)
}

func loadSpriteColored(filepath string) rl.Texture2D {
	image := rl.LoadImage(filepath)
	defer rl.UnloadImage(image)
	rl.ImageColorReplace(image, colorWhitePiece, colorBlackPiece)
	rl.ImageResizeNN(image, image.Width * int32(scale), image.Height * int32(scale))
	return rl.LoadTextureFromImage(image)
}
