package main

import rl "github.com/gen2brain/raylib-go/raylib"

const scale int32 = 4
const cell_size int32 = 16
const total_cell_size int32 = scale * cell_size

const w int32 = 8
const h int32 = 8

func loadSprite(filepath string) rl.Texture2D {
	image := rl.LoadImage(filepath)
	defer rl.UnloadImage(image)
	rl.ImageResizeNN(image, image.Width * scale, image.Height * scale)
	return rl.LoadTextureFromImage(image)
}

func main() {
	rl.InitWindow(w * total_cell_size, h * total_cell_size, "girvel's chess app")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	pawn := loadSprite("sprites/pawn.png")
	white := rl.GetColor(0xb5bdbaff)
	black := rl.GetColor(0x3a373dff)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
			for x := range w {
				for y := range h {
					var square_color rl.Color
					if (x + y) % 2 == 0 {
						square_color = white
					} else {
						square_color = black
					}
					rl.DrawRectangle(
						int32(x) * total_cell_size, int32(y) * total_cell_size,
						total_cell_size, total_cell_size,
						square_color,
					)
				}
			}
			rl.DrawTexture(pawn, 0, 0, rl.White)
		rl.EndDrawing()
	}
}
