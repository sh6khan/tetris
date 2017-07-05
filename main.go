package main

import(
  "time"
  "runtime"

  "tetris/game"

  "github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

// CONSTANTS
const(
  TIMER_PERIOD = 255 // milliseconds
  BLOCK_SIZE = 20

  // Window size for GL
	WINDOW_WIDTH = BLOCK_SIZE * game.FIELD_WIDTH
	WINDOW_HEIGHT = BLOCK_SIZE * game.FIELD_LENGTH
)

func init() {
  // ensure that the only go routine that is running
  // on this thread is the main func
  runtime.LockOSThread()
}

func main() {
  if err := glfw.Init(); err != nil {
    panic(err)
  }

  if err := gl.Init(); err != nil {
		panic(err)
	}

  defer glfw.Terminate()

  window, err := glfw.CreateWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "Testing", nil, nil)
	if err != nil {
		panic(err)
	}

  game.StartGame()

  // by making this thread the current context
  // we are saying that all drawings will happen on this thread
  // since we dont have multiple windows, is this needed?
  window.MakeContextCurrent()

  // set the key press callback for everytime we get a keyboard input
  window.SetKeyCallback(keyPress)


  // set up the ticker wich will run the ticker every round
  go gameTicker()

  // Init OpenGL
	gl.Ortho(0, WINDOW_WIDTH, WINDOW_HEIGHT, 0, -1, 1)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	//gl.ClearColor(20, 20, 0, 0)
	gl.ClearColor(255, 255, 255, 0)
	gl.LineWidth(1)
	gl.Color3f(1, 0, 0)

  for !window.ShouldClose() {
		// Do OpenGL stuff.
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
    drawScene()
    window.SwapBuffers()
		glfw.PollEvents()
	}
}

func keyPress(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
  switch key {
  case glfw.KeyUp:
  // we dont want to trigger on the edges of a key press
  if action != glfw.Press {
    return
  }
  game.RotateTetromino()
	case glfw.KeyLeft:
		game.LateralTranslation(-1)
	case glfw.KeyRight:
		game.LateralTranslation(1)
	case glfw.KeyDown:
		game.MoveDown()
  }
}

func gameTicker() {
  ticker := time.NewTicker(time.Millisecond * TIMER_PERIOD)
  for _ = range ticker.C {
    game.Tick()
  }
}


func drawScene() {
  drawTetromino()
	drawField()
}

func drawBlock(i int, j int) {
  i--
  j--
  gl.Begin(gl.POLYGON)
	glVertex(j * BLOCK_SIZE, i * BLOCK_SIZE)
	glVertex((j+1) * BLOCK_SIZE - 1, i * BLOCK_SIZE)
	glVertex((j+1) * BLOCK_SIZE - 1, (i+1) * BLOCK_SIZE -1)
	glVertex(j * BLOCK_SIZE, (i+1) * BLOCK_SIZE - 1)
	gl.End()
}

func drawTetromino() {
  setColor(game.GetTetro().ColorIdx)
  for _, pos := range game.GetTetro().Pos {
    drawBlock(game.GetY() + pos.GetY(), game.GetX() + pos.GetX())
  }
}

func glVertex(x, y int) {
	gl.Vertex2i(int32(x), int32(y))
}

func drawField() {
  for i, row := range game.GetField() {
    for j, cell := range row {
      if cell > 0 {
        setColor(cell - 1)
        drawBlock(i, j)
      }
    }
  }
}


func setColor(idx int) {
  c := game.Colors[idx]
  gl.Color3ub(uint8(c.R), uint8(c.G), uint8(c.B))
}
