package main

import (
    rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
    screenWidth := int32(1280)
    screenHeight := int32(720)
    rl.InitWindow(screenWidth, screenHeight, "Go Dungeon Prototype")
    defer rl.CloseWindow()

    rl.SetTargetFPS(60)

    // Load dungeon models (OBJ)
	//character := rl.LoadModel("models/characters/obj/character-male-b.obj")
    corridor := rl.LoadModel("models/dungeon/obj/corridor.obj")
    corridorCorner := rl.LoadModel("models/dungeon/obj/corridor-corner.obj")
    corridorEnd := rl.LoadModel("models/dungeon/obj/corridor-end.obj")

    defer rl.UnloadModel(corridor)
    defer rl.UnloadModel(corridorCorner)
    defer rl.UnloadModel(corridorEnd)

    // Load character (static OBJ for now)
    playerModel := rl.LoadModel("assets/characters/obj/character_male_b.obj")
    defer rl.UnloadModel(playerModel)

    // Camera (first-person style)
    camera := rl.Camera3D{
        Position:   rl.NewVector3(0, 1.8, 4), // Slightly above ground
        Target:     rl.NewVector3(0, 1.8, 3),
        Up:         rl.NewVector3(0, 1, 0),
        Fovy:       60.0,
        Projection: rl.CameraPerspective,
    }
    rl.DisableCursor() // lock mouse for FPS control

    moveSpeed := float32(3.0)
    mouseSensitivity := float32(0.003)

    for !rl.WindowShouldClose() {
        dt := rl.GetFrameTime()

        // Mouse look
        delta := rl.GetMouseDelta()
        yaw := -delta.X * mouseSensitivity
        pitch := -delta.Y * mouseSensitivity

        forward := rl.Vector3Normalize(rl.Vector3Subtract(camera.Target, camera.Position))
        right := rl.Vector3Normalize(rl.Vector3CrossProduct(forward, rl.NewVector3(0, 1, 0)))

        // Apply yaw rotation (horizontal)
        forward = rl.Vector3RotateByAxisAngle(forward, rl.NewVector3(0, 1, 0), yaw)

        // Apply pitch rotation (vertical)
        forward = rl.Vector3RotateByAxisAngle(forward, right, pitch)

        camera.Target = rl.Vector3Add(camera.Position, forward)

        // Keyboard movement (WASD)
        movement := rl.NewVector3(0, 0, 0)
        if rl.IsKeyDown(rl.KeyW) {
            movement = rl.Vector3Add(movement, forward)
        }
        if rl.IsKeyDown(rl.KeyS) {
            movement = rl.Vector3Subtract(movement, forward)
        }
        if rl.IsKeyDown(rl.KeyA) {
            movement = rl.Vector3Subtract(movement, right)
        }
        if rl.IsKeyDown(rl.KeyD) {
            movement = rl.Vector3Add(movement, right)
        }

        if rl.Vector3Length(movement) > 0 {
            movement = rl.Vector3Normalize(movement)
            movement = rl.Vector3Scale(movement, moveSpeed*dt)
            camera.Position = rl.Vector3Add(camera.Position, movement)
            camera.Target = rl.Vector3Add(camera.Position, forward)
        }

        rl.BeginDrawing()
        rl.ClearBackground(rl.Black)

        rl.BeginMode3D(camera)

        // Place some dungeon tiles
        rl.DrawModel(corridor, rl.NewVector3(0, 0, 0), 1.0, rl.White)
        rl.DrawModel(corridorCorner, rl.NewVector3(4, 0, 0), 1.0, rl.White)
        rl.DrawModel(corridorEnd, rl.NewVector3(8, 0, 0), 1.0, rl.White)

        // Draw the player model as a test prop in the scene
        rl.DrawModel(playerModel, rl.NewVector3(2, 0, 2), 1.0, rl.White)

        rl.EndMode3D()

        rl.DrawFPS(10, 10)
        rl.EndDrawing()
    }
}
