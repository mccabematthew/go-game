package main

import (
	"encoding/json"
	"fmt"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Prop is a placed object in the level
type Prop struct {
	Type  string  `json:"type"`
	X     float32 `json:"x"`
	Y     float32 `json:"y"`
	Z     float32 `json:"z"`
	Scale float32 `json:"scale"`
}

// LevelData holds level props and spawn
type LevelData struct {
	Props       []Prop `json:"props"`
	PlayerSpawn struct {
		X float32 `json:"x"`
		Y float32 `json:"y"`
		Z float32 `json:"z"`
	} `json:"player_spawn"`
}

func loadLevel(path string) (LevelData, error) {
	var lvl LevelData
	f, err := os.ReadFile(path)
	if err != nil {
		return lvl, err
	}
	err = json.Unmarshal(f, &lvl)
	return lvl, err
}

func main() {
	const screenW = 1280
	const screenH = 720
	rl.InitWindow(screenW, screenH, "Raylib-go 3D Level - Prototype")
	rl.SetTargetFPS(60)

	// --- Camera setup (simple first-person style)
	camera := rl.Camera3D{
		Position:   rl.NewVector3(0, 2, 10), // will be set by level spawn later
		Target:     rl.NewVector3(0, 1.5, 0),
		Up:         rl.NewVector3(0, 1, 0),
		Fovy:       60.0,
		Projection: rl.CameraPerspective,
	}

	// Load assets: try to load models if present; if not, we'll use primitives.
	// Put your Kenney models/textures into assets/models and assets/textures.
	var crateModel rl.Model
	var treeModel rl.Model
	var crateTexture rl.Texture2D
	cratesLoaded := false
	treesLoaded := false

	// Attempt to load model files; if fail, use fallback primitives.
	if _, err := os.Stat("assets/models/crate.obj"); err == nil {
		crateModel = rl.LoadModel("assets/models/crate.obj")
		if _, texErr := os.Stat("assets/textures/crate.png"); texErr == nil {
			crateTexture = rl.LoadTexture("assets/textures/crate.png")
			crateModel.Materials[0].Maps[rl.MAP_DIFFUSE].Texture = crateTexture
		}
		cratesLoaded = true
	}
	if _, err := os.Stat("assets/models/tree.obj"); err == nil {
		treeModel = rl.LoadModel("assets/models/tree.obj")
		if _, texErr := os.Stat("assets/textures/tree.png"); texErr == nil {
			tex := rl.LoadTexture("assets/textures/tree.png")
			treeModel.Materials[0].Maps[rl.MAP_DIFFUSE].Texture = tex
		}
		treesLoaded = true
	}

	// Load a ground texture if present
	var groundTex rl.Texture2D
	if _, err := os.Stat("assets/textures/ground.png"); err == nil {
		groundTex = rl.LoadTexture("assets/textures/ground.png")
	}

	// Level data: try loading JSON, otherwise default
	level, err := loadLevel("level/level1.json")
	if err != nil {
		// fallback in-code level
		level = LevelData{
			Props: []Prop{
				{Type: "crate", X: 2, Y: 0, Z: -4, Scale: 1},
				{Type: "crate", X: -2, Y: 0, Z: -6, Scale: 1},
				{Type: "tree", X: -3, Y: 0, Z: 6, Scale: 1.5},
			},
		}
		level.PlayerSpawn.X = 0
		level.PlayerSpawn.Y = 1.6
		level.PlayerSpawn.Z = 8
	}

	// Set camera to player spawn
	camera.Position = rl.NewVector3(level.PlayerSpawn.X, level.PlayerSpawn.Y, level.PlayerSpawn.Z)
	camera.Target = rl.NewVector3(level.PlayerSpawn.X, level.PlayerSpawn.Y, level.PlayerSpawn.Z-1)

	// Simple collision boxes for props (AABB)
	type AABB struct {
		Min rl.Vector3
		Max rl.Vector3
	}
	propBoxes := make([]AABB, 0, len(level.Props))

	for _, p := range level.Props {
		half := float32(0.5 * p.Scale)
		// simple cube box around prop center (works for crates; trees may need larger)
		min := rl.NewVector3(p.X-half, p.Y, p.Z-half)
		max := rl.NewVector3(p.X+half, p.Y+1.5*p.Scale, p.Z+half)
		propBoxes = append(propBoxes, AABB{Min: min, Max: max})
	}

	playerRadius := float32(0.3)
	playerHeight := float32(1.6)

	// Movement settings
	moveSpeed := float32(6)
	mouseSensitivity := float32(0.0035)

	// Hide cursor for mouselook
	rl.DisableCursor()

	// Main loop
	for !rl.WindowShouldClose() {
		// delta time
		dt := rl.GetFrameTime()

		// --- Mouse look: rotate camera's target using mouse delta
		mx := rl.GetMouseDelta().X
		my := rl.GetMouseDelta().Y

		// Camera forward vector
		forward := rl.Vector3Normalize(rl.Vector3Subtract(camera.Target, camera.Position))

		// Yaw: rotate around Y axis by mouse X
		yaw := mx * mouseSensitivity
		// Pitch: clamp vertical rotation
		pitch := -my * mouseSensitivity

		// Apply yaw (rotate forward vector around Up)
		forward = rl.Vector3RotateByQuaternion(forward, rl.QuaternionFromAxisAngle(rl.NewVector3(0, 1, 0), yaw))

		// Apply pitch: find right axis
		right := rl.Vector3Normalize(rl.Vector3CrossProduct(forward, rl.NewVector3(0, 1, 0)))
		forward = rl.Vector3Normalize(rl.Vector3RotateByQuaternion(forward, rl.QuaternionFromAxisAngle(right, pitch)))

		// Update camera target (keep small distance in front)
		camera.Target = rl.Vector3Add(camera.Position, rl.Vector3Scale(forward, 1.0))

		// Movement relative to forward/right
		moveDir := rl.NewVector3(0, 0, 0)
		if rl.IsKeyDown(rl.KeyW) {
			moveDir = rl.Vector3Add(moveDir, forward)
		}
		if rl.IsKeyDown(rl.KeyS) {
			moveDir = rl.Vector3Subtract(moveDir, forward)
		}
		if rl.IsKeyDown(rl.KeyA) {
			moveDir = rl.Vector3Subtract(moveDir, right)
		}
		if rl.IsKeyDown(rl.KeyD) {
			moveDir = rl.Vector3Add(moveDir, right)
		}

		// Normalize and apply speed
		if rl.Vector3Length(moveDir) > 0.001 {
			moveDir = rl.Vector3Normalize(moveDir)
			moveVec := rl.Vector3Scale(moveDir, moveSpeed*dt)

			// Attempt move with simple collision: check new position against propBoxes
			newPos := rl.Vector3Add(camera.Position, moveVec)

			// Simple ground clamp
			if newPos.Y < playerHeight {
				newPos.Y = playerHeight
			}

			// AABB collision: treat player as capsule approximated by AABB centered on player's feet
			playerMin := rl.NewVector3(newPos.X-playerRadius, newPos.Y-playerHeight, newPos.Z-playerRadius)
			playerMax := rl.NewVector3(newPos.X+playerRadius, newPos.Y, newPos.Z+playerRadius)

			collided := false
			for _, b := range propBoxes {
				if aabbIntersect(playerMin, playerMax, b.Min, b.Max) {
					collided = true
					break
				}
			}
			if !collided {
				camera.Position = newPos
				camera.Target = rl.Vector3Add(camera.Position, rl.Vector3Scale(forward, 1.0))
			}
		}

		// Quick reset spawn
		if rl.IsKeyPressed(rl.KeyR) {
			camera.Position = rl.NewVector3(level.PlayerSpawn.X, level.PlayerSpawn.Y, level.PlayerSpawn.Z)
			camera.Target = rl.NewVector3(level.PlayerSpawn.X, level.PlayerSpawn.Y, level.PlayerSpawn.Z-1)
		}

		// --- Drawing
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)
		// Draw ground plane
		if groundTex.ID != 0 {
			// textured quad
			rl.DrawPlane(rl.NewVector3(0, 0, 0), rl.NewVector2(50, 50), rl.NewColor(255, 255, 255, 255)) // placeholder
			// Note: raylib-go plane draw with texture map setup is more involved; consider using rl.DrawModel with a large box textured.
		} else {
			rl.DrawPlane(rl.NewVector3(0, 0, 0), rl.NewVector2(50, 50), rl.Green)
		}

		// Draw props from level
		for i, p := range level.Props {
			pos := rl.NewVector3(p.X, p.Y, p.Z)
			switch p.Type {
			case "crate":
				if cratesLoaded {
					rl.DrawModel(crateModel, pos, p.Scale, rl.White)
				} else {
					rl.DrawCube(pos, 1*p.Scale, 1*p.Scale, 1*p.Scale, rl.Brown)
					rl.DrawCubeWires(pos, 1*p.Scale, 1*p.Scale, 1*p.Scale, rl.Black)
				}
			case "tree":
				if treesLoaded {
					rl.DrawModel(treeModel, pos, p.Scale, rl.White)
				} else {
					// simple trunk + leaves
					trunkPos := rl.NewVector3(p.X, p.Y+0.5*p.Scale, p.Z)
					rl.DrawCylinder(trunkPos, 0.2*p.Scale, 1.0*p.Scale, rl.NewColor(100, 60, 20, 255))
					leafPos := rl.NewVector3(p.X, p.Y+1.2*p.Scale, p.Z)
					rl.DrawSphere(leafPos, 0.8*p.Scale, rl.DarkGreen)
				}
			default:
				rl.DrawSphere(pos, 0.5*p.Scale, rl.Purple)
			}

			// Optionally draw bounding boxes for debugging
			box := propBoxes[i]
			rl.DrawCubeWires(rl.NewVector3((box.Min.X+box.Max.X)/2, (box.Min.Y+box.Max.Y)/2, (box.Min.Z+box.Max.Z)/2),
				box.Max.X-box.Min.X, box.Max.Y-box.Min.Y, box.Max.Z-box.Min.Z, rl.Maroon)
		}

		// Draw a small marker for the player
		rl.DrawSphere(camera.Position, 0.1, rl.Yellow)

		rl.EndMode3D()

		rl.DrawText("WASD to move, mouse to look, R to reset", 10, 10, 20, rl.DarkGray)
		rl.DrawFPS(10, screenH-30)
		rl.EndDrawing()
	}

	// Cleanup
	if cratesLoaded {
		crateModel.Unload()
	}
	if treesLoaded {
		treeModel.Unload()
	}
	if groundTex.ID != 0 {
		rl.UnloadTexture(groundTex)
	}
	rl.CloseWindow()
}

// aabbIntersect checks intersection of 2 AABBs
func aabbIntersect(minA, maxA, minB, maxB rl.Vector3) bool {
	return (minA.X <= maxB.X && maxA.X >= minB.X) &&
		(minA.Y <= maxB.Y && maxA.Y >= minB.Y) &&
		(minA.Z <= maxB.Z && maxA.Z >= minB.Z)
}
