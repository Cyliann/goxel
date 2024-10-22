package camera

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
