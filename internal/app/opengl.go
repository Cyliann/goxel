package app

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/go-gl/gl/v4.1-core/gl"
)

const WORLD_SIZE = 32

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

// Creates a 3D texture
func createTexture() uint32 {
	data := createWorldMap()
	var textureID uint32
	gl.GenTextures(1, &textureID)
	gl.BindTexture(gl.TEXTURE_3D, textureID)
	gl.TexImage3D(gl.TEXTURE_3D, 0, gl.RED, WORLD_SIZE, WORLD_SIZE, WORLD_SIZE, 0, gl.RED, gl.FLOAT, gl.Ptr(&data[0]))
	gl.TexParameteri(gl.TEXTURE_3D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_3D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_3D, gl.TEXTURE_WRAP_R, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_3D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_3D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.BindTexture(gl.TEXTURE_3D, 0)

	return textureID
}

// Passes the 3D texture to the shader
func sendTexture(textureID uint32, program uint32) {
	gl.BindTexture(gl.TEXTURE_3D, textureID)
	textureUniformLocation := gl.GetUniformLocation(program, gl.Str("voxelMap\x00"))
	gl.Uniform1i(textureUniformLocation, 0)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_3D, textureID)
}

func createWorldMap() [WORLD_SIZE * WORLD_SIZE * WORLD_SIZE]float32 {
	var data [WORLD_SIZE * WORLD_SIZE * WORLD_SIZE]float32
	radius := 15
	for x := range WORLD_SIZE {
		for y := range WORLD_SIZE {
			for z := range WORLD_SIZE {
				i := x + WORLD_SIZE*y + WORLD_SIZE*WORLD_SIZE*z
				if math.Pow(float64(x-radius), 2.)+math.Pow(float64(y-radius), 2)+math.Pow(float64(z-radius), 2) < float64(radius*radius) {
					data[i] = 1
				} else {
					data[i] = 0
				}
			}
		}
	}

	return data
}
