package app

import (
	"Cyliann/goxel/internal/camera"
	"Cyliann/goxel/internal/graphics"
	"Cyliann/goxel/internal/voxel_data"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// New creates a new app. Calls initGlfw() and initOpenGL().
func New() App {
	var err error
	app := App{shaderReloading: false}

	app.window = graphics.InitGlfw()
	app.program = graphics.InitOpenGL()
	app.addCallbacks()
	app.camera = camera.New(app.window)

	width := app.window.GetMonitor().GetVideoMode().Width
	height := app.window.GetMonitor().GetVideoMode().Height

	app.renderTexture = graphics.CreateRenderTexture(width, height)
	app.framebuffer, err = graphics.CreateFramebufferWithTexture(app.renderTexture)
	if err != nil {
		panic("Failed to create framebuffer")
	}

	flat_nodes := voxel_data.GetVoxels()
	graphics.SendSSBO(flat_nodes)

	return app
}

type App struct {
	window          *glfw.Window
	program         uint32
	vao             uint32
	camera          camera.Camera
	shaderReloading bool
	renderTexture   uint32
	framebuffer     uint32
}

// App.Run is the main app loop. Polls events and calls App.draw()
func (self *App) Run() {
	timeStart := time.Now()
	self.camera.Update(self.window) // update for the first time to set up matrices
	for !self.window.ShouldClose() {
		// frameTime := time.Now()
		elapsedTime := float32(time.Since(timeStart))
		self.draw()
		glfw.PollEvents() // has to be after draw()
		shouldUpdate := self.HandleInput()
		if shouldUpdate {
			self.camera.Update(self.window)
		}
		self.updateUniforms(elapsedTime)
		// fmt.Print("\033[H\033[2J")
		// fmt.Printf("Frame time: %f", float32(time.Since(frameTime).Milliseconds()))
	}
}

// App.draw redraws frames
func (self *App) draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	width := self.window.GetMonitor().GetVideoMode().Width
	height := self.window.GetMonitor().GetVideoMode().Height

	graphics.RunCompute(self.program, self.renderTexture, width, height)

	graphics.BlitFramebuffer(self.framebuffer, width, height)
	self.window.SwapBuffers()
}

func (self *App) addCallbacks() {
	graphics.SetCallbacks(self.window, self.program, self.renderTexture, self.framebuffer)
}

func (self *App) updateUniforms(elapsedTime float32) {
	uTime := gl.GetUniformLocation(self.program, gl.Str("uTime\x00"))
	gl.Uniform1f(uTime, elapsedTime/1000000000)

	uPlayerPos := gl.GetUniformLocation(self.program, gl.Str("uPlayerPos\x00"))
	gl.Uniform3f(uPlayerPos, self.camera.Pos.X(), self.camera.Pos.Y(), self.camera.Pos.Z())

	uInvView := gl.GetUniformLocation(self.program, gl.Str("uInvView\x00"))
	gl.UniformMatrix4fv(uInvView, 1, false, &self.camera.InverseView[0]) // pass the pointer to the first element. The rest is calculated by opengl

	uInvProj := gl.GetUniformLocation(self.program, gl.Str("uInvProj\x00"))
	gl.UniformMatrix4fv(uInvProj, 1, false, &self.camera.InverseProj[0]) // pass the pointer to the first element. The rest is calculated by opengl
}

func (self *App) reloadShaders() error {
	computeShader, err := graphics.CompileShader("shaders/compute.glsl", gl.COMPUTE_SHADER)
	if err != nil {
		return (err)
	}

	prog := gl.CreateProgram()

	gl.AttachShader(prog, computeShader)
	gl.LinkProgram(prog)
	gl.UseProgram(prog)

	self.program = prog
	graphics.ForceSizeUpdate(self.window, self.program, self.renderTexture, self.framebuffer)
	log.Debug("Reloaded: ", "Program", self.program, "shader", computeShader)

	return nil
}

// App.Close is run at the end of the program. Terminates the GLFW window.
func (self *App) Close() {
	glfw.Terminate()
}
