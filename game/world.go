// game/world.go
package game

import rl "github.com/gen2brain/raylib-go/raylib"

type World struct {
    RoomModel rl.Model
    Position  rl.Vector3
}

// NewWorld loads the dungeon room.
func NewWorld(roomPath string) World {
    return World{
        RoomModel: rl.LoadModel(roomPath),
        Position:  rl.NewVector3(0, 0, 0),
    }
}

// Draw renders the room.
func (w *World) Draw() {
    rl.DrawModel(w.RoomModel, w.Position, 1.0, rl.White)
}

// Unload frees world resources.
func (w *World) Unload() {
    rl.UnloadModel(w.RoomModel)
}
