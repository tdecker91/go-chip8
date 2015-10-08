package screen

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/tdecker91/go-chip8/chip8"
	"fmt"
)

const SCREEN_WIDTH int = 640
const SCREEN_HEIGHT int = 320

type Screen struct {
	title string
	chip *chip8.Chip8
}

func (screen *Screen) Init(title string, onCloseHandler func(), chip *chip8.Chip8) {
	var err error

	screen.chip = chip

	err = glfw.Init()

	if err != nil {
		panic(err)
	}

	window, err := glfw.CreateWindow(SCREEN_WIDTH, SCREEN_HEIGHT, title, nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	gl.ClearColor(255.0, 0.0, 0.0, 0.0);  //Set the cleared screen colour to black
	gl.Viewport(0, 0, int32(SCREEN_WIDTH), int32(SCREEN_HEIGHT));   //This sets up the viewport so that the coordinates (0, 0) are at the top left of the window

	for !window.ShouldClose() {
		// Do OpenGL stuff
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		window.SwapBuffers()
		glfw.PollEvents()
	}

	onCloseHandler()

}

func drawRect() {
	
}