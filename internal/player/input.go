package player

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	// . "github.com/tylerwince/godbg"
)

// Returns a closure; so you can pass parameters to it (eg. program)
func KeyCallback(camera *Camera) glfw.KeyCallback {
	return func(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press || action == glfw.Repeat {
			switch key {
			case glfw.KeyW:
				camera.Z -= 0.1
			case glfw.KeyUp:
				camera.Z -= 0.1

			case glfw.KeyS:
				camera.Z += 0.1
			case glfw.KeyDown:
				camera.Z += 0.1

			case glfw.KeyD:
				camera.X += 0.1
			case glfw.KeyRight:
				camera.X += 0.1

			case glfw.KeyA:
				camera.X -= 0.1
			case glfw.KeyLeft:
				camera.X -= 0.1

			case glfw.KeySpace:
				camera.Y += 0.1

			case glfw.KeyLeftShift:
				camera.Y -= 0.1
			}
		}
	}
}
