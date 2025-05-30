package graphics

import (
	"github.com/charmbracelet/log"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// InitGlfw initializes glfw and returns a Window to use.
func InitGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.ScaleToMonitor, glfw.True)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// Get the primary monitor
	monitor := glfw.GetPrimaryMonitor()
	videoMode := monitor.GetVideoMode()
	if videoMode == nil {
		panic("Failed to get monitor video mode")
	}

	// Get the available video modes and set the monitor to its native resolution
	modes := monitor.GetVideoModes()
	bestMode := findBestMode(modes, videoMode.Width, videoMode.Height, videoMode.RefreshRate)

	// Print selected mode information (for debugging)
	log.Debugf("Monitor: %v, Selected Resolution: %dx%d @ %dHz\n", monitor.GetName(), bestMode.Width, bestMode.Height, bestMode.RefreshRate)

	window, err := glfw.CreateWindow(bestMode.Width, bestMode.Height, "Goxel engine", monitor, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

// Returns a closure, so you can pass parameters to it (eg. program)
func windowResizeCallback(program, renderTexture, framebuffer uint32, scale_x, scale_y float32) func(*glfw.Window, int, int) {
	return func(window *glfw.Window, width int, height int) {
		ResizeTexture(framebuffer, renderTexture, width, height)
		uSize := gl.GetUniformLocation(program, gl.Str("uSize"+"\x00"))
		gl.Uniform2f(uSize, float32(width)*scale_x, float32(height)*scale_y)
	}
}

// Utility function to find the best video mode based on resolution and refresh rate
func findBestMode(modes []*glfw.VidMode, targetWidth, targetHeight, targetRefreshRate int) *glfw.VidMode {
	var bestMode *glfw.VidMode
	for _, mode := range modes {
		// Try to match the target resolution and refresh rate
		if mode.Width == targetWidth && mode.Height == targetHeight && mode.RefreshRate == targetRefreshRate {
			bestMode = mode
			break
		}
	}
	// If we didn't find an exact match, fallback to the first available mode (shouldn't happen normally)
	if bestMode == nil && len(modes) > 0 {
		bestMode = modes[0]
	}
	return bestMode
}

// Manualy sets uSize uniform
func ForceSizeUpdate(window *glfw.Window, program, renderTexture, framebuffer uint32) {
	scale_x, scale_y := window.GetMonitor().GetContentScale()
	width := window.GetMonitor().GetVideoMode().Width
	height := window.GetMonitor().GetVideoMode().Height

	windowResizeCallback(program, renderTexture, framebuffer, scale_x, scale_y)(window, width-1, height-1)
}

func SetCallbacks(window *glfw.Window, program, renderTexture, framebuffer uint32) {
	scale_x, scale_y := window.GetMonitor().GetContentScale()
	window.SetSizeCallback(windowResizeCallback(program, renderTexture, framebuffer, scale_x, scale_y))
}
