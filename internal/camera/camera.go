package camera

import (
	"github.com/charmbracelet/log"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Camera struct {
	X      float32
	Y      float32
	Z      float32
	Pitch  float32
	Yaw    float32
	MouseX float64
	MouseY float64
}

func New(window *glfw.Window) Camera {
	// Lock cursor to screen
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	x, y := window.GetCursorPos()
	return Camera{0, 0, 2, 0, 0, x, y}
}

// Handles mouse and keyboard input. Modifies camera fields
func (self *Camera) HandleInput(window *glfw.Window) {

	x, y := window.GetCursorPos()
	self.Yaw = float32(x - self.MouseX)
	self.Pitch = float32(y - self.MouseY)

	if self.Yaw != 0 {
		self.MouseX = x
		log.Debug(self.Yaw)
	}

	if self.Pitch != 0 {
		self.MouseY = y
	}

	if window.GetKey(glfw.KeyW) == glfw.Press || window.GetKey(glfw.KeyUp) == glfw.Press {
		self.Z -= 0.1
	}

	if window.GetKey(glfw.KeyS) == glfw.Press || window.GetKey(glfw.KeyDown) == glfw.Press {
		self.Z += 0.1
	}

	if window.GetKey(glfw.KeyD) == glfw.Press || window.GetKey(glfw.KeyRight) == glfw.Press {
		self.X += 0.1
	}

	if window.GetKey(glfw.KeyA) == glfw.Press || window.GetKey(glfw.KeyLeft) == glfw.Press {
		self.X -= 0.1
	}

	if window.GetKey(glfw.KeySpace) == glfw.Press {
		self.Y += 0.1
	}

	if window.GetKey(glfw.KeyLeftShift) == glfw.Press {
		self.Y -= 0.1
	}
}
