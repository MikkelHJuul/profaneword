package cmd

import (
	"fmt"
	"testing"
)

type s string

func (s s) String() string {
	return string(s)
}

func Test_getWords(t *testing.T) {

	wfuns := [][]fmt.Stringer{
		{s("a")},
		{s("b")},
		{s("c")},
		{s("d")},
		{s("e")},
	}

	tests := []struct {
		name string
		words int
		want string
	}{
		{"returns d", 4, "d"},
		{"returns a", 1, "a"},
		{"two", 6, "e a"},
		{"more", 14, "e e d"},
		{"even more", 20, "e e e e"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getWords(tt.words, wfuns); got != tt.want {
				t.Errorf("getWords() = %v, want %v", got, tt.want)
			}
		})
	}
}