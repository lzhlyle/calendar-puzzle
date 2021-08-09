package quick

import (
	"strconv"
	"testing"
	"time"
)

func TestCalc(t *testing.T) {
	// Fill directions
	for b := range Blocks {
		for d := 0; d < 4; d++ {
			if d > 0 {
				Blocks[b][d] = MoveToTopLeft(Rotation(Blocks[b][d-1]))
			}
		}
	}

	board := InitBoard(time.Now())

	// dfs
	res, err := Fill(board, Blocks)
	if err != nil {
		t.Logf("FAIL, err: %v", err)
	}
	Output(res)
}

func TestZipBoard(t *testing.T) {
	exp, act := "11110000000000000000000000000000000000010000001000000", strconv.FormatInt(zipBoard(emptyBoard), 2)
	if exp != act {
		t.Logf("\nexp:%s\nact:%s", exp, act)
		t.Fail()
	}
}
