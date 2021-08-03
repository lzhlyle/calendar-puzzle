package quick

import (
	"testing"
	"time"
)

// 棋盘可通过 [8][7]int8 表示
// 0: empty, 1: occupied, 9: cannot
var emptyBoard = [8][7]int8{
	{0, 0, 0, 0, 0, 0, 9},
	{0, 0, 0, 0, 0, 0, 9},
	{},
	{},
	{},
	{},
	{},
	{9, 9, 9, 9, 0, 0, 0},
}

// 4*4 可表示任意形状
// 10 种形状
// [10][4][4][4]bool 可表示任意形状的 4 个旋转结果，用于枚举
// [10:shape][4:direction][4:rows][4:cells]bool

type Sharp int8

const (
	B23 Sharp = 0 + iota
	C23
	L14
	L23
	L24
	L33
	T33
	Z23
	Z24
	Z33
)

type Direction int8

const (
	North Direction = 0 + iota
	East
	South
	West
)

var blocks = [10][4][4][4]int8{
	B23: {
		North: {
			{1},
			{1, 1},
			{1, 1},
		},
	},
	C23: {
		North: {
			{1, 1},
			{1},
			{1, 1},
		},
	},
	L14: {
		North: {
			{1, 1, 1, 1},
		},
	},
	L23: {
		North: {
			{1},
			{1, 1, 1},
		},
	},
	L24: {
		North: {
			{1},
			{1, 1, 1, 1},
		},
	},
	L33: {
		North: {
			{1},
			{1},
			{1, 1, 1},
		},
	},
	T33: {
		North: {
			{1},
			{1, 1, 1},
			{1},
		},
	},
	Z23: {
		North: {
			{0, 1, 1},
			{1, 1},
		},
	},
	Z24: {
		North: {
			{1, 1},
			{0, 1, 1, 1},
		},
	},
	Z33: {
		North: {
			{1, 1},
			{0, 1},
			{0, 1, 1},
		},
	},
}

func Rotation(curr [4][4]int8) [4][4]int8 {
	// clone
	res := [4][4]int8{}
	for i := range curr {
		for j := range curr[i] {
			res[i][j] = curr[i][j]
		}
	}

	// \-
	for i := range res {
		for j := range res[i] {
			if i > j {
				res[i][j], res[j][i] = res[j][i], res[i][j]
			}
		}
	}
	for i := range res {
		res[i][0], res[i][1], res[i][2], res[i][3] = res[i][3], res[i][2], res[i][1], res[i][0]
	}
	return res
}

func MoveToTopLeft(res [4][4]int8) [4][4]int8 {
	rowSum, colSum := [4]int8{}, [4]int8{}
	for i := range res {
		for j := 0; j < 4; j++ {
			rowSum[i] += res[i][j]
			colSum[j] += res[i][j]
		}
	}
	topEmpty, topContinue, leftEmpty, leftContinue := 0, true, 0, true
	for idx := 0; idx < 4; idx++ {
		if rowSum[idx] == 0 && topContinue {
			topEmpty++
		} else {
			topContinue = false
		}
		if colSum[idx] == 0 && leftContinue {
			leftEmpty++
		} else {
			leftContinue = false
		}
	}
	if topEmpty > 0 {
		for i := topEmpty; i < 4; i++ {
			for j := 0; j < 4; j++ {
				res[i-topEmpty][j], res[i][j] = res[i][j], 0
			}
		}
	}
	if leftEmpty > 0 {
		for j := leftEmpty; j < 4; j++ {
			for i := 0; i < 4; i++ {
				res[i][j-leftEmpty], res[i][j] = res[i][j], 0
			}
		}
	}
	return res
}

func initBoard(date time.Time) [8][7]int8 {
	month := date.Month() - 1 // [0, 11]
	day := date.Day() - 1     // [0, 30]
	week := date.Weekday()    // [0, 6]

	res := [8][7]int8{}
	for i := range emptyBoard {
		for j := range emptyBoard[i] {
			res[i][j] = emptyBoard[i][j]
		}
	}

	// current date should 2
	res[month/6][month%6] = 2
	res[day/7+2][day%7] = 2
	res[week/4+6][week%4+3+week/4] = 2

	return res
}

func fill(board [8][7]int8, blocks [10][4][4][4]int8, curr [][4][4]int8, visited map[Sharp]bool) ([][4][4]int8, bool) {
	if len(visited) == 10 {
		return curr, true
	}

	// todo @lzh 必须压缩，否则不好做「见缝插针」

	return nil, false
}

func TestCalc(t *testing.T) {
	// fill directions
	for b := range blocks {
		for d := 0; d < 4; d++ {
			if d > 0 {
				blocks[b][d] = MoveToTopLeft(Rotation(blocks[b][d-1]))
			}
		}
	}

	board := initBoard(time.Now())

	// dfs
	if res, ok := fill(board, blocks, [][4][4]int8{}, map[Sharp]bool{}); ok {
		// output
		for _, b := range res {
			for _, row := range b {
				t.Logf("%v", row)
			}
			t.Log()
		}
	}

	// output
	//for _, dir2Block := range blocks {
	//	for d := 0; d < 4; d++ {
	//		t.Logf("dir: %d", d)
	//		for _, row := range dir2Block[d] {
	//			t.Logf("%v", row)
	//		}
	//		t.Log()
	//	}
	//	t.Log("======")
	//}
	//
	//for i := range board {
	//	t.Logf("%v", board[i])
	//}
}
