package profaneword

import "testing"

func TestShuffleFormatter_Format(t *testing.T) {
	tests := []string{
		"asd",
		"hello",
		"world",
		"a",
		"",
	}
	s := ShuffleFormatter{CryptoRand{}}
	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			if got := s.Format(tt); len(got) != len(tt) {
				t.Errorf("length mismatch: %v, %v", got, tt)
			}
		})
	}
}
