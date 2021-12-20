package profanities

import (
	"strings"
	"testing"
)

func TestSentencesCorrectLength(t *testing.T) {
	for l, items := range sentences {
		for _, words := range items {
			if len(words.words) != l {
				t.Errorf("Length of sentence: `%s` is not correct, should be %d, found %d", words.string, l, len(words.words))
			}
			c := strings.Count(words.string, "%s")
			if c != l {
				t.Errorf("number of string-formats incorrect in: `%s`, should be %d, found %d", words.string, l, c)
			}
		}
	}
}
