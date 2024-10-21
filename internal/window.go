package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	// . "github.com/tylerwince/godbg"
)

var (
	triangle = []float32{
		3, -1, 0, // top left
		-1, 3, 0, // bottom right
		-1, -1, 0, // bottom left
	}
)

type App struct{}

func (self *App) Create() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()

	program := initOpenGL(window)

	vao := makeVao(triangle)

	timeStart := time.Now()
	for !window.ShouldClose() {
		// frameTime := time.Now()
		uTime := gl.GetUniformLocation(program, gl.Str("uTime\x00"))
		elapsedTime := float32(time.Since(timeStart))
		gl.Uniform1f(uTime, elapsedTime/1000000000)

		draw(vao, window, program)
		// fmt.Print("\033[H\033[2J")
		// fmt.Printf("Frame time: %f", float32(time.Since(frameTime).Milliseconds()))
	}
}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.ScaleToMonitor, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
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
	fmt.Printf("Monitor: %v, Selected Resolution: %dx%d @ %dHz\n", monitor.GetName(), bestMode.Width, bestMode.Height, bestMode.RefreshRate)

	window, err := glfw.CreateWindow(bestMode.Width, bestMode.Height, "Goxel engine", monitor, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL(window *glfw.Window) uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}

	// version := gl.GoStr(gl.GetString(gl.VERSION))
	// log.Println("OpenGL version", version)

	vertexShader, err := compileShader("shaders/vert.glsl", gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := compileShader("shaders/frag.glsl", gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()

	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)

	scale_x, scale_y := float32(1), float32(1) //window.GetMonitor().GetContentScale()
	window.SetSizeCallback(windowResizeCallback(prog, scale_x, scale_y))
	fbheight, fbwidth := window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(fbheight), int32(fbwidth))

	return prog
}

func draw(vao uint32, window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle)/3))
	gl.UseProgram(program)

	glfw.PollEvents()
	window.SwapBuffers()
}

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}

func compileShader(path string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	source, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("Failed to fetch shader: %s", err)
	}

	csources, free := gl.Strs(string(source) + "\x00")
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

// Returns a closure, so I can pass parameters to it (eg. program)
func windowResizeCallback(program uint32, scale_x float32, scale_y float32) func(*glfw.Window, int, int) {
	return func(window *glfw.Window, width int, height int) {
		gl.Viewport(0, 0, int32(float32(width)*scale_x), int32(float32(height)*scale_y))
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
