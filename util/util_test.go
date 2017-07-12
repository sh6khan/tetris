package util


import "testing"


func TestRotation(t *testing.T) {
  input := [1]struct{
    input [][]int
    expectedOutput [][]int
  } {
    // test case 1
    {
      [][]int{
  			[]int{0, 0, 0, 1},
  			[]int{0, 0, 0, 1},
  			[]int{0, 0, 0, 1},
  			[]int{0, 0, 0, 1},
  		},
      [][]int{
  			[]int{1, 1, 1, 1},
  			[]int{0, 0, 0, 0},
  			[]int{0, 0, 0, 0},
  			[]int{0, 0, 0, 0},
  		},
    },
  }

  for _, testCase := range input {
    res := RotateMatrix(testCase.input)
    comapareMatrix(res, testCase.expectedOutput, t)
  }
}

func BenchmarkRotation(b *testing.B) {
  input := [][]int {
    []int{0, 0, 0, 0},
    []int{0, 0, 1, 1},
    []int{0, 1, 1, 0},
    []int{0, 0, 0, 0},
  }

  for n := 0; n < b.N; n++ {
    RotateMatrix(input)
  }
}

func comapareMatrix(a [][]int, b [][]int, t *testing.T) {
  for i, row := range a {
    for j, cell := range row {
      if b[i][j] != cell {
        t.Fatalf("at index: %d %d, expectedOutput: %d, but got: %d", i, j, b[i][j], cell)
      }
    }
  }
}
