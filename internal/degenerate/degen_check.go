package degenerate

import (
	"fmt"
	"os"
)

func DegenCheck(boardName string) {
	var boards []string = []string{"s", "hc", "hm", "h", "e", "u", "d", "y", "t", "hr", "gif", "aco", "r"}
	for _, board := range boards {
		if boardName == board {
			fmt.Println("Stop being a degenerate.")
			os.Exit(0)
		}
	}
}
