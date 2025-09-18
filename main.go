package main

import (
    "github.com/gen2brain/raylib-go/raylib"
    "github.com/JSONAlexander/go-game/game"
)

func main() {
    // --- Window & Camera Setup ---
    screenWidth := int32(1280)
    screenHeight := int32(720)
    rl.InitWindow(screenWidth, screenHeight, "Go Game - Single Player Loop")
    rl.SetTargetFPS(60)

    camera := rl.Camera3D{}
    camera.Position = rl.NewVector3(0, 10, 10) // Start above and behind player
    camera.Target = rl.NewVector3(0, 0, 0)
    camera.Up = rl.NewVector3(0, 1, 0)
    camera.Fovy = 60.0
    camera.Projection = rl.CameraPerspective

    // --- Game Objects ---
    world := game.NewWorld("assets/models/Modular Dungeon Pack - Jan 2018/obj/Carpet.obj")
    player := game.NewPlayer(
        "assets/models/characters/kenney_blocky-characters_20/glb/character-j.glb",
        rl.NewVector3(0, 0, 0),
    )

    // --- Main Game Loop ---
    for !rl.WindowShouldClose() {
        dt := rl.GetFrameTime()

        // Update
        player.Update(dt)
        camera.Target = player.Position
        camera.Position = rl.NewVector3(
            player.Position.X,
            player.Position.Y+10,
            player.Position.Z+10,
        )

        // Draw
        rl.BeginDrawing()
        rl.ClearBackground(rl.RayWhite)

        rl.BeginMode3D(camera)
        world.Draw()
        player.Draw()
        rl.EndMode3D()

        rl.DrawFPS(10, 10)
        rl.EndDrawing()
    }

    // --- Cleanup ---
    world.Unload()
    player.Unload()
    rl.CloseWindow()
}
