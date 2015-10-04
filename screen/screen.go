package screen

import (
	//"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

const SCREEN_WIDTH int = 640
const SCREEN_HEIGHT int = 320

type Screen struct {
	title string
}

func (screen *Screen) Init(title string) {
	var err error

	err = glfw.Init()

	if err != nil {
		panic(err)
	}

	window, err := glfw.CreateWindow(SCREEN_WIDTH, SCREEN_HEIGHT, title, nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	for !window.ShouldClose() {
		// Do OpenGL stuff
		window.SwapBuffers()
		glfw.PollEvents()
	}

}