package camera

import (
	"math"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	Pos         mgl32.Vec3
	Direction   mgl32.Vec3
	MousePos    mgl32.Vec2
	pitchDelta  float32
	yawDelta    float32
	Fov         float32
	NearClip    float32
	FarClip     float32
	InverseView mgl32.Mat4
	InverseProj mgl32.Mat4
}

func New(window *glfw.Window) Camera {
	// Lock cursor to screen
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	x, y := window.GetCursorPos()
	return Camera{
		mgl32.Vec3{0, 0, 2},                // pos
		mgl32.Vec3{0, 0, -1},               // direction
		mgl32.Vec2{float32(x), float32(y)}, // mous pos
		0,                                  // pitch
		0,                                  // yaw
		45,                                 // fov
		0,                                  // nearclip
		0,                                  // farclip
		mgl32.Ident2().Mat4(),
		mgl32.Ident2().Mat4(),
	}
}

// Handles mouse and keyboard input. Modifies camera fields
func (self *Camera) HandleInput(window *glfw.Window) bool {
	var speed float32 = 0.2
	up_dir := mgl32.Vec3{0, 1, 0}
	right_dir := self.Direction.Cross(up_dir)

	shouldUpdate := false
	x, y := window.GetCursorPos()
	self.yawDelta = float32(x) - self.MousePos[0]
	self.pitchDelta = float32(y) - self.MousePos[1]

	if self.yawDelta != 0 {
		self.MousePos[0] = float32(x)
		shouldUpdate = true
	}

	if self.pitchDelta != 0 {
		self.MousePos[1] = float32(y)
		shouldUpdate = true
	}

	// Move forward
	if window.GetKey(glfw.KeyW) == glfw.Press || window.GetKey(glfw.KeyUp) == glfw.Press {
		self.Pos = self.Pos.Add(self.Direction.Mul(speed))
		shouldUpdate = true
	}

	// Move back
	if window.GetKey(glfw.KeyS) == glfw.Press || window.GetKey(glfw.KeyDown) == glfw.Press {
		self.Pos = self.Pos.Add(self.Direction.Mul(-speed))
		shouldUpdate = true
	}

	// Move right
	if window.GetKey(glfw.KeyD) == glfw.Press || window.GetKey(glfw.KeyRight) == glfw.Press {
		self.Pos = self.Pos.Add(right_dir.Mul(speed))
		shouldUpdate = true
	}

	// Move left
	if window.GetKey(glfw.KeyA) == glfw.Press || window.GetKey(glfw.KeyLeft) == glfw.Press {
		self.Pos = self.Pos.Add(right_dir.Mul(-speed))
		shouldUpdate = true
	}

	// Move up
	if window.GetKey(glfw.KeySpace) == glfw.Press {
		self.Pos = self.Pos.Add(up_dir.Mul(speed))
		shouldUpdate = true
	}

	// Move down
	if window.GetKey(glfw.KeyLeftShift) == glfw.Press {
		self.Pos = self.Pos.Add(up_dir.Mul(-speed))
		shouldUpdate = true
	}

	return shouldUpdate
}

func (self *Camera) Update() {
	up_dir := mgl32.Vec3{0, 1, 0}
	right_dir := self.Direction.Cross(up_dir)
	var speed float32 = 0.01

	// Create quaternion for pitch (rotation around the right axis)
	quatPitch := AngleAxis(self.pitchDelta*speed, right_dir)

	// Create quaternion for yaw (rotation around the up axis, which is (0, 1, 0))
	quatYaw := AngleAxis(self.yawDelta*speed, mgl32.Vec3{0, 1, 0})

	// Normalize the resulting quaternion
	quat := quatPitch.Mul(quatYaw).Normalize()

	// Rotate the forward direction using the quaternion
	self.Direction = quat.Rotate(self.Direction)

	self.recalculateView()
}

func (self *Camera) recalculateView() {
	// For some reason mgl32 expects 9 floats instead of 3 vecs, and go doesn't support tuple unpacking
	self.InverseView = mgl32.LookAt(
		self.Pos[0],
		self.Pos[1],
		self.Pos[2],
		self.Pos.Add(self.Direction)[0],
		self.Pos.Add(self.Direction)[1],
		self.Pos.Add(self.Direction)[2],
		0, // up_dir x
		1, // up_dir y
		0, // up_dir z
	).Inv()
}

func (self *Camera) recalculateProjection(window *glfw.Window) {
	self.InverseProj = mgl32.Perspective(
		mgl32.DegToRad(self.Fov),
		float32(window.GetMonitor().GetVideoMode().Width/window.GetMonitor().GetVideoMode().Height), // aspect ratio
		self.NearClip,
		self.FarClip,
	).Inv()
}

func AngleAxis(angle float32, axis mgl32.Vec3) mgl32.Quat {
	halfAngle := angle / 2.0
	s := float32(math.Sin(float64(halfAngle)))
	return mgl32.Quat{
		W: float32(math.Cos(float64(halfAngle))),
		V: axis.Mul(s), // axis * sin(angle / 2)
	}
}
