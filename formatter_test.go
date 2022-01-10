package profaneword

import (
	"fmt"
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

type countRandomDevice struct {
	counter int
}

func (n *countRandomDevice) Rand() *big.Rat {
	return big.NewRat(1, 1)
}

func (n *countRandomDevice) RandMax(max int) int {
	if n.counter == max {
		n.counter = 0
		return max
	}
	n.counter++
	return n.counter - 1
}

var _ RandomDevice = &countRandomDevice{}

func TestHorseFormatter_Format(t *testing.T) {
	words := make([]string, len(horsewords))
	h := HorseFormatter{&countRandomDevice{}}
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

var _ RandomDevice = &countRandomDevice{}

type maxRandomDevice struct{}

func (maxRandomDevice) Rand() *big.Rat {
	panic("implement me")
}

func (maxRandomDevice) RandMax(max int) int {
	return max - 1
}

var _ RandomDevice = maxRandomDevice{}

func TestShuffleFormatter_Format_NonRandom(t *testing.T) {
	in := "SOMETHING"
	h := ShuffleFormatter{&countRandomDevice{}}
	got := h.Format(in)
	if in != got {
		t.Errorf("non-random shuffle should be No-op, expected: %s, got %s", in, got)
	}
	sh := ShuffleFormatter{maxRandomDevice{}}
	got = sh.Format(in)
	if got != "GNIHTEMOS" { // S gets picked first before a s
		t.Errorf("non-random max-shuffle should return GNIHTEMOS (reverse), got %s", got)
	}
}

func TestStudderFormatter_Format(t *testing.T) {
	in := "zero one two three four zero"
	s := PerWordFormattingFormatter{StudderFormatter{&countRandomDevice{}}}
	got := s.Format(in)
	expected := "zero o-one t-t-two t-t-t-three f-f-f-f-four zero"
	if got != expected {
		t.Errorf("non random studder should be non-random: expected %s, got %s", expected, got)
	}
	ss := StudderFormatter{&countRandomDevice{}}
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
	if got != "asdd!" {
		t.Errorf("Expected swearFormatter to add exclamation")
	}
	in = "asd-asd"
	got = sf.Format(in)
	if got != "asd!-asd" {
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
	if r.Format("dsa") != "asd" {
		t.Errorf("expected ReversingFormatter to reverse the string")
	}
}

func TestTitleFormatter_Format(t *testing.T) {
	tf := TitleFormatter{}
	if tf.Format("asd") != "Asd" {
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

type countingFormatter struct {
	int
}

func (c *countingFormatter) Format(_ string) string {
	defer func() { c.int++ }()
	return fmt.Sprintf(`%d`, c.int)
}

func TestPerWordFormattingFormatter_Format(t *testing.T) {
	pw := PerWordFormattingFormatter{}
	pw.SetFormatter(&countingFormatter{})
	got := pw.Format("any three words")
	if got != "0 1 2" {
		t.Errorf("unexpected return from PerWordFormattingFormatter")
	}
}

type appendingFormatter string

func (a appendingFormatter) Format(word string) string {
	return word + string(a)
}

func TestMultiFormatter_Format(t *testing.T) {
	one, two, three := appendingFormatter("one"), appendingFormatter("two"), appendingFormatter("three")
	mf := MultiFormatter{}
	mf.With(one)
	mf.With(two)
	mf.With(three)
	if mf.Format("") != "onetwothree" {
		t.Errorf("unexpected formatting of MultiFormatter")
	}
}
