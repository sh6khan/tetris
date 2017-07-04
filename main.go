package main

import(
  "fmt"
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

  // by making this thread the current context
  // we are saying that all drawings will happen on this thread
  // since we dont have multiple windows, is this needed?
  window.MakeContextCurrent()

  // set the key press callback for everytime we get a keyboard input
  window.SetKeyCallback(keyPress)

  // intialize the Grid
  initGrid()
  t := InitialTetros[0]
  printTetromino(t)
  t = rotateMatrix(t)
  printTetromino(t)




  // set up the ticker wich will run the ticker every round
  //go gameTick()

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
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func keyPress(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
  fmt.Printf("key: %v, scancode: %d, action: %v\n", key, scancode, action)
}

func update() {
  fmt.Printf("%v", time.Now())
}

func gameTick() {
  ticker := time.NewTicker(time.Millisecond * TIMER_PERIOD)
  for t := range ticker.C {
    fmt.Println("Ticker at", t)
  }
}

// move the tetromino peice left and right
func lateralTranslation(dx int) {
  for _, pos := range tetromino {
    if field[pos.x + posX + dx][pos.y + posY] != 0 {
      return
    }
  }

  posX += dx
}

func moveDown() {
  for _, pos := range tetromino {
    if field[pos.x + posX][pos.y + posY + 1] != 0 {
      // leave the peice on the field
      placeTetro()
      return
    }
  }

  posY += 1
}

func placeTetro() {
  for _, pos := range tetromino {
    field[pos.x][pos.y] = 1
  }
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


// for testing
func printTetromino(matrix [][]int) {
  for _, row := range matrix {
    for _, cell := range row {
      fmt.Printf("%d ", cell)
    }
    fmt.Println("")
  }
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
  posX = 0
  posY = FIELD_WIDTH / 2
  r := rand.Intn(4)
  tetromino = makeTretroObject(InitialTetros[r])
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
