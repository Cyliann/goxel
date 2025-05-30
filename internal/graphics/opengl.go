package graphics

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"unsafe"

	"Cyliann/goxel/internal/voxel_data"

	"github.com/charmbracelet/log"
	"github.com/go-gl/gl/v4.1-core/gl"
)

// InitOpenGL initializes OpenGL and returns an intiialized program and a fragment shader.
func InitOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Debugf("OpenGL version: %v", version)

	computeShader, err := CompileShader("shaders/compute.glsl", gl.COMPUTE_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()

	gl.AttachShader(prog, computeShader)
	gl.LinkProgram(prog)
	gl.UseProgram(prog)

	return prog
}

// compiles the given shader from path and returns it as a memory address.
func CompileShader(path string, shaderType uint32) (uint32, error) {
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
func SendSSBO(flat []voxel_data.FlatNode) uint32 {
	var ssbo uint32
	gl.GenBuffers(1, &ssbo)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, ssbo)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, len(flat)*int(unsafe.Sizeof(flat[0])), gl.Ptr(flat), gl.STATIC_DRAW)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 1, ssbo)

	return ssbo
}

func CreateRenderTexture(width, height int) uint32 {
	var texture uint32
	gl.CreateTextures(gl.TEXTURE_2D, 1, &texture)
	gl.TextureStorage2D(texture, 1, gl.RGBA32F, int32(width), int32(height))
	gl.TextureParameteri(texture, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TextureParameteri(texture, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	gl.TextureParameteri(texture, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TextureParameteri(texture, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	return texture
}

func CreateFramebufferWithTexture(texture uint32) (uint32, error) {
	var framebuffer uint32

	gl.CreateFramebuffers(1, &framebuffer)
	err := attachTextureToFramebuffer(texture, framebuffer)
	if err != nil {
		gl.DeleteFramebuffers(1, &framebuffer)
		return 0, errors.New("Failed to create framebuffer with texture")
	}

	return framebuffer, nil
}

func attachTextureToFramebuffer(texture, framebuffer uint32) error {
	gl.NamedFramebufferTexture(framebuffer, gl.COLOR_ATTACHMENT0, texture, 0)
	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		return errors.New("Framebuffer not complete")
	}
	return nil
}

func ResizeTexture(framebuffer, texture uint32, width, height int) {
	gl.DeleteTextures(1, &texture)
	newTexture := CreateRenderTexture(width, height)
	attachTextureToFramebuffer(newTexture, framebuffer)
}

func RunCompute(computeShader, renderTexture uint32, width, height int) {
	gl.UseProgram(computeShader)
	gl.BindImageTexture(0, renderTexture, 0, false, 0, gl.WRITE_ONLY, gl.RGBA32F)

	const workGroupSizeX uint32 = 16
	const workGroupSizeY uint32 = 16

	numGroupsX := (uint32(width) + workGroupSizeX - 1) / workGroupSizeX
	numGroupsY := (uint32(height) + workGroupSizeY - 1) / workGroupSizeY

	gl.DispatchCompute(numGroupsX, numGroupsY, 1)
	gl.MemoryBarrier(gl.SHADER_IMAGE_ACCESS_BARRIER_BIT)
}

func BlitFramebuffer(framebuffer uint32, width, height int) {
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, framebuffer)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0) // swapchain

	gl.BlitFramebuffer(0, 0, int32(width), int32(height), 0, 0, int32(width), int32(height), gl.COLOR_BUFFER_BIT, gl.NEAREST)
}
