package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	vertexShaderSource = `
    #version 460
    in vec3 vp;

    void main() {
        gl_Position = vec4(vp, 1.0);
    }
` + "\x00"

	fragmentShaderSource = `
    #version 460
    out vec4 frag_color;

		uniform vec2 uSize;

    void main() {
  			vec2 color = gl_FragCoord.xy/uSize;
        frag_color = vec4(color.xy, pow(1 - color.x * color.y, 2), 1);
    }
` + "\x00"
)

var (
	triangle = []float32{
		1, 1, 0, // top right
		1, -1, 0, // top left
		-1, 1, 0, // bottom right
	}
	width  float32 = 1000
	height float32 = 1000
)

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()

	program := initOpenGL(window)

	vao := makeVao(triangle)

	for !window.ShouldClose() {
		draw(vao, window, program)
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

	window, err := glfw.CreateWindow(int(width), int(height), "Goxel engine", nil, nil)
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
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()

	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)

	window.SetSizeCallback(windowResizeWrapper(prog))

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

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
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
func windowResizeWrapper(program uint32) func(*glfw.Window, int, int) {
	return func(window *glfw.Window, width int, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
		uSize := gl.GetUniformLocation(program, gl.Str("uSize"+"\x00"))
		gl.Uniform2f(uSize, float32(width), float32(height))
	}
}
