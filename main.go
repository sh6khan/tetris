package main

import(
  "time"
  "runtime"
  "math/rand"

  "github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

// CONSTANTS
const(
  TIMER_PERIOD = 255 // milliseconds
  FIELD_WIDTH = 10
  FIELD_LENGTH = 20
  BLOCK_SIZE = 20

  // Window size for GL
	WINDOW_WIDTH = BLOCK_SIZE * FIELD_WIDTH
	WINDOW_HEIGHT = BLOCK_SIZE * FIELD_LENGTH
)

// Each peice in the game is defined by
// and array of 4 points, the points
// are defined by integer x and y values
type Tetromino [4]struct {
  x int
  y int
}

type Color struct{
  R int
  G int
  B int
}

// GLOBAL VARIABLES
var (
  posX = 0
  posY = 0

  field [FIELD_LENGTH+2][FIELD_WIDTH+2]int

  InitialTetros = [][][]int{
    [][]int {
        []int {0,0,1,0},
        []int {0,0,1,0},
        []int {0,0,1,1},
        []int {0,0,0,0},
    },


    [][]int {
        []int {0,0,0,0},
        []int {0,1,1,1},
        []int {0,0,1,0},
        []int {0,0,0,0},
    },

    [][]int {
        []int {0,0,0,0},
        []int {0,1,1,0},
        []int {0,1,1,0},
        []int {0,0,0,0},
    },

    [][]int {
        []int {0,0,0,0},
        []int {0,0,1,1},
        []int {0,1,1,0},
        []int {0,0,0,0},
    },

    [][]int {
        []int {0,0,1,0},
        []int {0,0,1,0},
        []int {0,0,1,0},
        []int {0,0,1,0},
    },
  }

  // There will only be one peice active at a time
  tetromino Tetromino
  matrixTetromino [][]int
  colorIdx int


  Colors = []Color{
		Color{0, 0, 0},
		Color{170, 0, 0},
		Color{192, 192, 192},
		Color{170, 0, 170},
		Color{0, 0, 170},
		Color{0, 170, 0},
		Color{170, 85, 0},
		Color{0, 170, 170},
	}
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

  startGame()

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
    rotateTetromino()
  	case glfw.KeyLeft:
  		lateralTranslation(-1)
  	case glfw.KeyRight:
  		lateralTranslation(1)
  	case glfw.KeyDown:
  		moveDown()
  }
}

func gameTicker() {
  ticker := time.NewTicker(time.Millisecond * TIMER_PERIOD)
  for _ = range ticker.C {
    tick()
  }
}

func tick() {
  moveDown()
}

// The deleteCompletedLine function will delete an entire line
// from the field if there is not single element with the value of
// 0. We the move everthing above that row down by done
func deleteCompletedLine() {
  for i, row := range field {
    if i == 0 || i == FIELD_LENGTH + 1 {
      continue
    }


    completedRow := true

    for _, cell := range row {
      if cell == 0 {
        completedRow = false
        //break
      }
    }

    // move everything down by one
    if completedRow == true {
      println(i)
      for x := i - 1; x >= 1; x-- {
        for y, _ := range field[x] {
          field[x+1][y] = field[x][y]
        }
      }
      return
    }
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
  setColor(colorIdx)
  for _, pos := range tetromino {
    drawBlock(posY + pos.y, posX + pos.x)
  }
}

func glVertex(x, y int) {
	gl.Vertex2i(int32(x), int32(y))
}

func drawField() {
  for i, row := range field {
    for j, cell := range row {
      if cell > 0 {
        setColor(cell - 1)
        drawBlock(i, j)
      }
    }
  }
}

func startGame() {
  initGrid()
  generateTetromino()
}

// move the tetromino peice left and right
func lateralTranslation(dx int) {
  for _, pos := range tetromino {
    if field[pos.y + posY][pos.x + posX + dx] != 0 {
      return
    }
  }

  posX += dx
}

func rotateTetromino() {
  // generate rotated matrix
  newTetrominoMatrix := rotateMatrix(matrixTetromino)
  newTetromino := makeTretroObject(newTetrominoMatrix)

  for _, pos := range newTetromino {
    if field[posY + pos.y][posX + pos.x] != 0 {
      return
    }
  }

  // assign the current tetromino to the new matrix
  tetromino = newTetromino
}

func moveDown() {
  for _, pos := range tetromino {
    if field[pos.y + posY + 1][pos.x + posX] != 0 {

      // the Game is over
      if posY < 2 {
				startGame()
				return
			}

      // leave the peice on the field
      placeTetro()
      deleteCompletedLine()
      generateTetromino()
      return
    }
  }

  posY++
}

func placeTetro() {
  for _, pos := range tetromino {
    field[posY + pos.y][posX + pos.x] = colorIdx + 1
  }
}

func setColor(idx int) {
  c := Colors[idx]
  gl.Color3ub(uint8(c.R), uint8(c.G), uint8(c.B))
}

// since all the pecies are represented as a 2D matrix
// we can rotate the matrix in order to rotate the peice
// counter clockwise
func rotateMatrix(mat [][]int) [][]int {
    for x := 0; x < 2; x++ {
        for y := x; y < 3-x; y++ {
            var temp int = mat[x][y];
            mat[x][y] = mat[y][3-x];
            mat[y][3-x] = mat[3-x][3-y];
            mat[3-x][3-y] = mat[3-y][x];
            mat[3-y][x] = temp;
        }
    }
    return mat
}

func makeTretroObject(matrix [][]int) (res Tetromino){
  x := 0
  for i, row := range matrix {
    for j, cell := range row {
      if cell == 1 {
        res[x].x = i
        res[x].y = j
        x++
      }
    }
  }

  return res
}

// Genearte a new Tetromino object to be placed into the field
func generateTetromino() {
  posY = 0
  posX = FIELD_WIDTH / 2
  colorIdx = rand.Intn(8)
  r := rand.Intn(5)
  matrixTetromino = InitialTetros[r]
  tetromino = makeTretroObject(matrixTetromino)
}

// initGrid will populate the values of the global grid object
// it initializes all values as 0 with the outer edges defined as -1
// the reason for this is to check if a left/right translation leads
// to an invalid move
func initGrid() {
  for i, row := range field {
    for j, _ := range row {
      field[i][j] = 0
      if j == 0 || i == 0 || i == FIELD_LENGTH+1 || j == FIELD_WIDTH+1 {
        field[i][j] = -1
      }
    }
  }
}
