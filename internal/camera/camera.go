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
	PitchDelta  float32
	YawDelta    float32
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
		mgl32.Vec3{64, 64, -64},            // pos
		mgl32.Vec3{0, 0, 1},                // direction
		mgl32.Vec2{float32(x), float32(y)}, // mous pos
		0,                                  // pitchDelta
		0,                                  // yawDelta
		60,                                 // fov
		1,                                  // nearclip
		100,                                // farclip
		mgl32.Ident2().Mat4(),
		mgl32.Ident2().Mat4(),
	}
}

func (self *Camera) Update(window *glfw.Window) {
	up_dir := mgl32.Vec3{0, 1, 0}
	right_dir := self.Direction.Cross(up_dir)

	// Create quaternion for pitch (rotation around the right axis)
	quatPitch := angleAxis(-self.PitchDelta, right_dir)

	// Create quaternion for yaw (rotation around the up axis, which is (0, 1, 0))
	quatYaw := angleAxis(-self.YawDelta, mgl32.Vec3{0, 1, 0})

	// Normalize the resulting quaternion
	quat := quatPitch.Mul(quatYaw).Normalize()

	// Rotate the forward direction using the quaternion
	self.Direction = quat.Rotate(self.Direction)

	self.recalculateView()
	self.recalculateProjection(window)
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

func angleAxis(angle float32, axis mgl32.Vec3) mgl32.Quat {
	halfAngle := angle / 2.0
	s := float32(math.Sin(float64(halfAngle)))
	return mgl32.Quat{
		W: float32(math.Cos(float64(halfAngle))),
		V: axis.Mul(s), // axis * sin(angle / 2)
	}
}
