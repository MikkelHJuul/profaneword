package profanities

import (
	"fmt"
	"testing"
)

func TestRadixDatabase(t *testing.T) {
	countMap := make(map[Word]int, 7)
	for _, r := range radix {
		wrds := r.GetWords([]Word{START, FILLER, END, EXCL, SPLIT, MISSPELL, POSITIVE})
		for k, v := range wrds {
			countMap[k] = countMap[k] + len(v)
		}
	}
	fmt.Println(countMap)
}
