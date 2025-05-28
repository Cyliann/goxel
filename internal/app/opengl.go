package app

import (
	"fmt"
	"os"
	"strings"
	"unsafe"

	"Cyliann/goxel/internal/voxel_data"

	"github.com/charmbracelet/log"
	"github.com/go-gl/gl/v4.1-core/gl"
)

// initOpenGL initializes OpenGL and returns an intiialized program and a fragment shader.
func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Debugf("OpenGL version: %v", version)

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
	gl.UseProgram(prog)

	return prog
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

// compiles the given shader from path and returns it as a memory address.
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

		error := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(error))

		return 0, fmt.Errorf("failed to compile %v: %v", string(source), error)
	}

	return shader, nil
}

// Creates and sends an SSBO
func sendSSBO(flat []voxel_data.FlatNode) uint32 {
	var ssbo uint32
	gl.GenBuffers(1, &ssbo)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, ssbo)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, len(flat)*int(unsafe.Sizeof(flat[0])), gl.Ptr(flat), gl.STATIC_DRAW)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 0, ssbo)

	return ssbo
}
