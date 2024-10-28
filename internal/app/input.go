package app

import (
	"github.com/charmbracelet/log"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

// Handles mouse and keyboard input. Modifies camera fields
func (self *App) HandleInput() bool {
	var keyboardSpeed float32 = 0.1
	var mouseSpeed float32 = 0.01
	up_dir := mgl32.Vec3{0, 1, 0}
	right_dir := self.camera.Direction.Cross(up_dir)

	shouldUpdate := false
	x, y := self.window.GetCursorPos()
	self.camera.YawDelta = (float32(x) - self.camera.MousePos[0]) * mouseSpeed
	self.camera.PitchDelta = (float32(y) - self.camera.MousePos[1]) * mouseSpeed

	if self.camera.YawDelta != 0 {
		self.camera.MousePos[0] = float32(x)
		shouldUpdate = true
	}

	if self.camera.PitchDelta != 0 {
		self.camera.MousePos[1] = float32(y)
		shouldUpdate = true
	}

	// Move forward
	if self.window.GetKey(glfw.KeyW) == glfw.Press || self.window.GetKey(glfw.KeyUp) == glfw.Press {
		self.camera.Pos = self.camera.Pos.Add(self.camera.Direction.Mul(keyboardSpeed))
		shouldUpdate = true
	}

	// Move back
	if self.window.GetKey(glfw.KeyS) == glfw.Press || self.window.GetKey(glfw.KeyDown) == glfw.Press {
		self.camera.Pos = self.camera.Pos.Add(self.camera.Direction.Mul(-keyboardSpeed))
		shouldUpdate = true
	}

	// Move right
	if self.window.GetKey(glfw.KeyD) == glfw.Press || self.window.GetKey(glfw.KeyRight) == glfw.Press {
		self.camera.Pos = self.camera.Pos.Add(right_dir.Mul(keyboardSpeed))
		shouldUpdate = true
	}

	// Move left
	if self.window.GetKey(glfw.KeyA) == glfw.Press || self.window.GetKey(glfw.KeyLeft) == glfw.Press {
		self.camera.Pos = self.camera.Pos.Add(right_dir.Mul(-keyboardSpeed))
		shouldUpdate = true
	}

	// Move up
	if self.window.GetKey(glfw.KeySpace) == glfw.Press {
		self.camera.Pos = self.camera.Pos.Add(up_dir.Mul(keyboardSpeed))
		shouldUpdate = true
	}

	// Move down
	if self.window.GetKey(glfw.KeyLeftShift) == glfw.Press {
		self.camera.Pos = self.camera.Pos.Add(up_dir.Mul(-keyboardSpeed))
		shouldUpdate = true
	}

	if self.window.GetKey(glfw.KeyR) == glfw.Press && !self.shaderReloading {
		self.shaderReloading = true
		err := reloadShaders(self)
		if err != nil {
			log.Errorf("Reload failed with: %v", err)
		}
	}

	if self.window.GetKey(glfw.KeyR) == glfw.Release {
		self.shaderReloading = false
	}

	return shouldUpdate
}
