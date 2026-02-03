package main

import rl "github.com/gen2brain/raylib-go/raylib"

func main() {
	rl.InitWindow(512, 512, "girvel's chess app")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		rl.DrawText("Hello, raylib!", 10, 10, 20, rl.LightGray)
		rl.EndDrawing()
	}
}
