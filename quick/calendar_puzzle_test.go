package quick

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"
)

// 棋盘可通过 [8][7]int8 表示
// -1: empty, 20: date, 99: forbid, [0, 9]: block
const (
	Empty  = -1
	Date   = 20
	Forbid = 99
)

var emptyBoard = [8][7]int8{
	{Empty, Empty, Empty, Empty, Empty, Empty, Forbid},
	{Empty, Empty, Empty, Empty, Empty, Empty, Forbid},
	{Empty, Empty, Empty, Empty, Empty, Empty, Empty},
	{Empty, Empty, Empty, Empty, Empty, Empty, Empty},
	{Empty, Empty, Empty, Empty, Empty, Empty, Empty},
	{Empty, Empty, Empty, Empty, Empty, Empty, Empty},
	{Empty, Empty, Empty, Empty, Empty, Empty, Empty},
	{Forbid, Forbid, Forbid, Forbid, Empty, Empty, Empty},
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

// 不需要在图形的基础上移动，压缩后右移即可
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
	res[month/6][month%6] = Date
	res[day/7+2][day%7] = Date
	res[week/4+6][week%4+3+week/4] = Date

	return res
}

func fill(board [8][7]int8, blocks [10][4][4][4]int8) ([8][7]int8, error) {
	// zip
	zBoard := zipBoard(board)
	zBlocks := [10][4]int64{}
	for i := range blocks {
		for d := range blocks[i] {
			zBlocks[i][d] = zipBlock(blocks[i][d])
		}
	}

	// prepare dfs
	res, ok := tryFill(board, zBoard, zBlocks, 0)
	if !ok {
		return [8][7]int8{}, errors.New("can not fill")
	}
	return res, nil
}

func tryFill(curr [8][7]int8, zipCurr int64, zBlocks [10][4]int64, i int) ([8][7]int8, bool) {
	// i: block index
	if i >= 10 {
		return curr, true
	}

	//outputUnzipBoard(zipCurr)

	// 选方向
	for d := range zBlocks[i] {
		zBlock := zBlocks[i][d]

		//outputUnzipBlock(zBlock)

		// 选位置
		// 在棋盘上从左上角开始尽量拼进去，直到可放入：在二进制从右向左滑动
		for loc := 0; loc < 64-16; loc++ {
			var move int64 = int64(zBlock << loc)
			tryMove := move

			// 尝试放入：检查移动后的每一位
			for tryMove > 0 {
				//outputBinaryString(tryMove)
				//outputBinaryString(zipCurr)
				// 看最低位的 1 是否有占用
				lowest := tryMove & -tryMove
				if zipCurr == zipCurr|lowest {
					// 已被占用
					break
				}
				// 消掉最低位的 1
				tryMove &= tryMove - 1
			}

			// 不可放入
			if tryMove != 0 {
				continue
			}

			// copy to next curr
			nextCurr := [8][7]int8{}
			for ni := range curr {
				for nj := range curr[ni] {
					nextCurr[ni][nj] = curr[ni][nj]
				}
			}
			// 将 blocks[i][d] 放 curr 的 loc 位置上，标记为 blocks 的索引 i
			out := false
			locI, locJ := loc/7, loc%7
			for bi := range blocks[i][d] {
				for bj := range blocks[i][d][bi] {
					if blocks[i][d][bi][bj] > 0 {
						// 检查是否超出棋盘（二进制时不直观）
						if locI+bi >= 8 || locJ+bj >= 7 {
							out = true
							break
						}
						nextCurr[locI+bi][locJ+bj] = int8(i)
					}
				}
				if out {
					break
				}
			}
			if out {
				continue
			}

			//output(nextCurr)
			res, ok := tryFill(nextCurr, zipCurr|move, zBlocks, i+1)
			if ok {
				return res, true
			}
		} // loc
	} // direction

	// 这个 block，任何方向、任何格子都放不进去，剪枝
	return curr, false
}

// 也要压缩成 int64，与棋盘统一，方便判断「是否可放入棋盘」
func zipBlock(block [4][4]int8) int64 {
	var res int64
	for i := range block {
		for j := range block[i] {
			if block[i][j] > 0 {
				res |= 1 << (7*i + j)
			}
		}
	}
	return res
}

func zipBoard(board [8][7]int8) int64 {
	var res int64
	for i := range board {
		for j := range board[i] {
			if board[i][j] > Empty {
				res |= 1 << (len(board[i])*i + j)
			}
		}
	}
	return res
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
	res, err := fill(board, blocks)
	if err != nil {
		t.Logf("FAIL, err: %v", err)
	}
	output(res)
}

func outputBinaryString(val int64) {
	fmt.Println(strconv.FormatInt(val, 2))
	fmt.Println()
}

func outputUnzipBoard(val int64) {
	var board [8][7]int8
	for i := 0; i < 8; i++ {
		for j := 0; j < 7; j++ {
			board[i][j] = int8(val & 1)
			val >>= 1
		}
	}
	output(board)
}

func output(board [8][7]int8) {
	for _, b := range board {
		for _, row := range b {
			fmt.Printf("%v\t", row)
		}
		fmt.Println()
	}
	fmt.Println()
}

func outputUnzipBlock(val int64) {
	var block [8][7]int8
	for i := 0; i < 8; i++ {
		for j := 0; j < 7; j++ {
			block[i][j] = int8(val & 1)
			val >>= 1
		}
	}
	output(block)
}

func TestZipBoard(t *testing.T) {
	exp, act := "11110000000000000000000000000000000000010000001000000", strconv.FormatInt(zipBoard(emptyBoard), 2)
	if exp != act {
		t.Logf("\nexp:%s\nact:%s", exp, act)
		t.Fail()
	}
}
