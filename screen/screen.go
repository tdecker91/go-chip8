package screen

import (
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

const SCREEN_WIDTH int = 640
const SCREEN_HEIGHT int = 320

type Screen struct {
	title string
	display Display
}

func (screen *Screen) Init(title string) {
	var err error

	err = glfw.Init()

	if err != nil {
		return err
	}

	glfw.OpenWindowHint(glfw.WindowNoResize, 1)
	err = glfw.OpenWindow(SCREEN_WIDTH, SCREEN_HEIGHT, 0, 0, 0, 0, 0, 0, glfw.Windowed)
	if err != nil {
		return err
	}

	glfw.SetWindowTitle(title)
	desktopMode := glfw.DesktopMode()
	glfw.SetWindowPos((desktopMode.W-SCREEN_WIDTH)/2, (desktopMode.H-SCREEN_HEIGHT)/2)

	gl.ClearColor(0.255, 0.255, 0.255, 0)

	glfw.SetKeyCallback(func(key, state int) {
		if state == glfw.KeyPress {
			//i.KeyHandler.KeyDown(key)
		} else {
			//i.KeyHandler.KeyUp(key)
		}
	})

	glfw.SetWindowCloseCallback(func() int {
		glfw.CloseWindow()
		glfw.Terminate()
		return 0
	})

	return nil
}