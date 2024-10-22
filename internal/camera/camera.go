package camera

import "github.com/go-gl/glfw/v3.3/glfw"

type Camera struct {
	X     float32
	Y     float32
	Z     float32
	Pitch float32
	Yaw   float32
}

func New() Camera {
	return Camera{0, 0, 2, 0, 0}
}

// Returns a closure; so you can pass parameters to it (eg. program)
func HandleInput(camera *Camera) {
	if glfw.GetCurrentContext().GetKey(glfw.KeyW) == glfw.Press || glfw.GetCurrentContext().GetKey(glfw.KeyUp) == glfw.Press {
		camera.Z -= 0.1
	}

	if glfw.GetCurrentContext().GetKey(glfw.KeyS) == glfw.Press || glfw.GetCurrentContext().GetKey(glfw.KeyDown) == glfw.Press {
		camera.Z += 0.1
	}

	if glfw.GetCurrentContext().GetKey(glfw.KeyD) == glfw.Press || glfw.GetCurrentContext().GetKey(glfw.KeyRight) == glfw.Press {
		camera.X += 0.1
	}

	if glfw.GetCurrentContext().GetKey(glfw.KeyA) == glfw.Press || glfw.GetCurrentContext().GetKey(glfw.KeyLeft) == glfw.Press {
		camera.X -= 0.1
	}

	if glfw.GetCurrentContext().GetKey(glfw.KeySpace) == glfw.Press {
		camera.Y += 0.1
	}

	if glfw.GetCurrentContext().GetKey(glfw.KeyLeftShift) == glfw.Press {
		camera.Y -= 0.1
	}
}
