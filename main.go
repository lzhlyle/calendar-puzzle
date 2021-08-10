package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/lzhlyle/calendar-puzzle/quick"
)

func main() {
	// fill directions
	for b := range quick.Blocks {
		for d := 0; d < 4; d++ {
			if d > 0 {
				quick.Blocks[b][d] = quick.MoveToTopLeft(quick.Rotation(quick.Blocks[b][d-1]))
			}
		}
	}

	date := time.Now()
	if len(os.Args) == 4 {
		y, errY := strconv.Atoi(os.Args[1])
		m, errM := strconv.Atoi(os.Args[2])
		d, errD := strconv.Atoi(os.Args[3])
		if errY == nil && errM == nil && errD == nil {
			date = time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local)
		}
	}
	fmt.Println(date.Format("2006-01-02"))
	board := quick.InitBoard(date)

	// dfs
	res, err := quick.Fill(board, quick.Blocks)
	if err != nil {
		fmt.Printf("FAIL, err: %v", err)
	}
	quick.Output(res)
}
