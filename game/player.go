// game/player.go
package game

import rl "github.com/gen2brain/raylib-go/raylib"

type Player struct {
    Model    rl.Model
    Position rl.Vector3
    Speed    float32
}

// NewPlayer loads the player model and creates a new Player instance.
func NewPlayer(modelPath string, startPos rl.Vector3) Player {
    model := rl.LoadModel(modelPath)
    return Player{
        Model:    model,
        Position: startPos,
        Speed:    5.0,
    }
}

// Update handles player movement input.
func (p *Player) Update(dt float32) {
    // Movement in X/Z plane
    if rl.IsKeyDown(rl.KeyW) {
        p.Position.Z -= p.Speed * dt
    }
    if rl.IsKeyDown(rl.KeyS) {
        p.Position.Z += p.Speed * dt
    }
    if rl.IsKeyDown(rl.KeyA) {
        p.Position.X -= p.Speed * dt
    }
    if rl.IsKeyDown(rl.KeyD) {
        p.Position.X += p.Speed * dt
    }
}

// Draw renders the player model.
func (p *Player) Draw() {
    rl.DrawModel(p.Model, p.Position, 1.0, rl.White)
}

// Unload frees resources when the game closes.
func (p *Player) Unload() {
    rl.UnloadModel(p.Model)
}
