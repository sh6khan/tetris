package game

import (
  "math/rand"
  "tetris/util"
)

const (
  FIELD_WIDTH = 10
  FIELD_LENGTH = 20
)

type Tetromino struct {
  Pos [4]Position
  ColorIdx int
  matrix [][]int
}

// Each peice in the game is defined by
// and array of 4 points, the points
// are defined by integer x and y values
type Position struct {
  x int
  y int
}

func (p *Position) GetX() int {
  return p.x
}

func (p *Position) GetY() int {
  return p.y
}

type Color struct{
  R int
  G int
  B int
}

// GLOBAL VARIABLES
var (
  PosX = 0
  PosY = 0

  field [FIELD_LENGTH+2][FIELD_WIDTH+2]int

  InitialTetros = [][][]int{
    [][]int {
        []int {0,0,1,0},
        []int {0,0,1,0},
        []int {0,0,1,1},
        []int {0,0,0,0},
    },

    [][]int {
        []int {0,0,1,0},
        []int {0,0,1,0},
        []int {0,1,1,0},
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
        []int {0,0,0,0},
        []int {1,1,0,0},
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
  Tetro Tetromino

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


func Tick() {
  MoveDown()
}

func StartGame() {
  initGrid()
  generateTetromino()
}

func GetX() int {
  return PosX
}

func GetY() int {
  return PosY
}

func GetField() [FIELD_LENGTH+2][FIELD_WIDTH+2]int {
  return field
}

func GetTetro() Tetromino {
  return Tetro
}


// move the Tetro peice left and right
func LateralTranslation(dx int) {
  for _, Pos := range Tetro.Pos {
    if field[Pos.y + PosY][Pos.x + PosX + dx] != 0 {
      return
    }
  }

  PosX += dx
}

func RotateTetromino() {
  // generate rotated matrix
  newTetrominoMatrix := util.RotateMatrix(Tetro.matrix)
  newTetromino := makeTretroObject(newTetrominoMatrix)

  for _, Pos := range newTetromino.Pos {
    if field[PosY + Pos.y][PosX + Pos.x] != 0 {
      return
    }
  }

  // assign the current Tetro to the new matrix
  newTetromino.ColorIdx = Tetro.ColorIdx
  Tetro = newTetromino
}

func MoveDown() {
  for _, Pos := range Tetro.Pos {
    if field[Pos.y + PosY + 1][Pos.x + PosX] != 0 {

      // the Game is over
      if PosY < 2 {
				StartGame()
				return
			}

      // leave the peice on the field
      placeTetro()
      deleteCompletedLine()
      generateTetromino()
      return
    }
  }

  PosY++
}

func makeTretroObject(matrix [][]int) (res Tetromino) {
  x := 0
  for i, row := range matrix {
    for j, cell := range row {
      if cell == 1 {
        res.Pos[x].x = i
        res.Pos[x].y = j
        x++
      }
    }
  }

  res.ColorIdx = rand.Intn(8)
  res.matrix = matrix

  return res
}

func placeTetro() {
  for _, Pos := range Tetro.Pos {
    field[PosY + Pos.y][PosX + Pos.x] = Tetro.ColorIdx + 1
  }
}

// Genearte a new Tetromino object to be placed into the field
func generateTetromino() {
  PosY = 0
  PosX = FIELD_WIDTH / 2
  r := rand.Intn(len(InitialTetros))
  Tetro = makeTretroObject(InitialTetros[r])
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
        break
      }
    }

    // move everything down by one
    if completedRow == true {
      for x := i - 1; x >= 1; x-- {
        for y, _ := range field[x] {
          field[x+1][y] = field[x][y]
        }
      }
      return
    }
  }
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
