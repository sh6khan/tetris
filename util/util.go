package util

func RotateMatrix(mat [][]int) [][]int {
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
