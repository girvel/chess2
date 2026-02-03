package main

import rl "github.com/gen2brain/raylib-go/raylib"

const scale int32 = 4

func loadSprite(filepath string) rl.Texture2D {
	image := rl.LoadImage(filepath)
	defer rl.UnloadImage(image)
	rl.ImageResizeNN(image, image.Width * scale, image.Height * scale)
	return rl.LoadTextureFromImage(image)
}

func main() {
	rl.InitWindow(512, 512, "girvel's chess app")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	pawn := loadSprite("sprites/pawn.png")

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		rl.DrawTexture(pawn, 0, 0, rl.White)
		rl.EndDrawing()
	}
}
