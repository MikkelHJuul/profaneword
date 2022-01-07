package profaneword

import (
	"math/big"
	"testing"
)

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

type nonRandomDevice struct {
	counter int
}

func (n *nonRandomDevice) Rand() *big.Rat {
	return big.NewRat(1, 1)
}

func (n *nonRandomDevice) RandMax(max int) int {
	if n.counter == max {
		n.counter = 0
		return max
	}
	n.counter++
	return n.counter - 1
}

var _ RandomDevice = &nonRandomDevice{}

func TestHorseFormatter_Format(t *testing.T) {
	words := make([]string, len(horsewords))
	h := HorseFormatter{&nonRandomDevice{}}
	for i := 0; i < len(horsewords); i++ {
		got := h.Format("")
		words[i] = got
	}
	for i, w := range words {
		if horsewords[i] != w {
			t.Errorf("Elements should be the same, expected %s, got %s", horsewords[i], w)
		}
	}
}

var _ RandomDevice = &nonRandomDevice{}

func TestShuffleFormatter_Format_NonRandom(t *testing.T) {
	in := "SOMETHING"
	h := ShuffleFormatter{&nonRandomDevice{}}
	got := h.Format(in)
	if in != got {
		t.Errorf("non-random shuffle should be No-op, expected: %s, got %s", in, got)
	}
}

func TestStudderFormatter_Format(t *testing.T) {
	in := "zero one two three four zero"
	s := PerWordFormattingFormatter{StudderFormatter{&nonRandomDevice{}}}
	got := s.Format(in)
	expected := "zero o-one t-t-two t-t-t-three f-f-f-f-four zero"
	if got != expected {
		t.Errorf("non random studder should be non-random: expected %s, got %s", expected, got)
	}
	ss := StudderFormatter{&nonRandomDevice{}}
	got = ss.Format(in)
	if got != in {
		t.Errorf("non random studder should be non-random: expected %s, got %s", expected, got)
	}
}

func TestSwearFormatter_Format(t *testing.T) {
	sf := swearFormatter{UnitFormatter{}}
	in := "!asd "
	if in != sf.Format(in) {
		t.Errorf("Expected No-op when the first letter is not a unicode letter")
	}
	in = "asdd"
	got := sf.Format(in)
	if "asdd!" != got {
		t.Errorf("Expected swearFormatter to add exclamation")
	}
	in = "asd-asd"
	got = sf.Format(in)
	if "asd!-asd" != got {
		t.Errorf("Expected swearFormatter to add exclamation correctly")
	}
}

func TestNewSwearFormatter(t *testing.T) {
	input := "ASDASD"
	sf := NewSwearFormatter()
	got := sf.Format(input)
	if len(got) != len(input)+1 {
		t.Errorf("expected exactly one exclamation to be added, input: %s, got: %s", input, got)
	}
}

func TestReversingFormatter_Format(t *testing.T) {
	r := ReversingFormatter{}
	if "asd" != r.Format("dsa") {
		t.Errorf("expected ReversingFormatter to reverse the string")
	}
}

func TestTitleFormatter_Format(t *testing.T) {
	tf := TitleFormatter{}
	if "Asd" != tf.Format("asd") {
		t.Errorf("incorrect titling")
	}
}

type cachingCharFormatter struct {
	text []rune
}

func (c *cachingCharFormatter) FormatRune(r rune) []rune {
	c.text = append(c.text, r)
	return []rune{r}
}

var _ CharFormatter = &cachingCharFormatter{}

func TestCharFormatterDelegatingFormatter_Format(t *testing.T) {
	del := CharFormatterDelegatingFormatter{}
	del.SetCharFormatter(&cachingCharFormatter{})
	del.Format("asd")
	got := string(del.GetCharFormatter().(*cachingCharFormatter).text)
	if got != "asd" {
		t.Errorf("CharFormatterDelegatingFormatter did not pass text correctly")
	}
}
