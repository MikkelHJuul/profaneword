package profaneword

import (
	"regexp"
	"strings"
)

//Formatter is a formatter, that formats the entire input text(string) and outputs the formatted text
type Formatter interface {
	Format(string) string
}

type CharFormatter interface {
	FormatRune(rune) []rune
}

type MultiFormatter struct {
	Formatters []Formatter
}

func (m *MultiFormatter) With(f Formatter) {
	if f != nil {
		m.Formatters = append(m.Formatters, f)
	}
}

func (m *MultiFormatter) Format(word string) string {
	for _, f := range m.Formatters {
		word = f.Format(word)
	}
	return word
}

type UnitFormatter struct{}

var _ Formatter = UnitFormatter{}

func (u UnitFormatter) Format(words string) string {
	return words
}

type CharFormatterDelegatingFormatter struct {
	CharFormatter
}

func (c *CharFormatterDelegatingFormatter) Format(word string) string {
	out := make([]rune, len(word))
	for _, r := range word {
		out = append(out, c.FormatRune(r)...)
	}
	return string(out)
}

type RandomlyFormattingFormatter struct {
	thresholdRandom
	Other Formatter
}

var _ Formatter = &RandomlyFormattingFormatter{
	thresholdRandom: thresholdRandom{},
	Other:           nil,
}

type PerWordFormattingFormatter struct {
	Other Formatter
}

func (rff *PerWordFormattingFormatter) Format(word string) string {
	words := strings.Split(word, ` `)
	ws := make([]string, len(words))
	for i, word := range words {
		ws[i] = rff.Other.Format(word)
	}
	return strings.Join(ws, ` `)
}

func (rff *RandomlyFormattingFormatter) Format(word string) string {
	if rff.Rand.Rand().Cmp(rff.Threshold) > 0 {
		return rff.Other.Format(word)
	}
	return word
}

func NewRandomlyFormatter(wrap Formatter) Formatter {
	random := &RandomlyFormattingFormatter{thresholdRandom: newFiftyFifty()}
	random.Other = wrap
	return &PerWordFormattingFormatter{random}
}

type TitleFormatter struct{}

var _ Formatter = TitleFormatter{}

func (t TitleFormatter) Format(word string) string {
	return strings.Title(word)
}

func RandomTitleFormatter() Formatter {
	randomFormatter := NewRandomlyFormatter(TitleFormatter{})
	return randomFormatter
}

var _ Formatter = &RandomlyFormattingFormatter{
	thresholdRandom: thresholdRandom{},
	Other:           nil,
}

type RegexReplacingFormatter struct {
	PatternMatcher *regexp.Regexp
	Replacement    string
}

func (rr *RegexReplacingFormatter) Format(word string) string {
	return rr.PatternMatcher.ReplaceAllString(word, rr.Replacement)
}

func DelimiterFormatterWith(repl string) Formatter {
	return &RegexReplacingFormatter{
		regexp.MustCompile(` `),
		repl,
	}
}

type ReversingFormatter struct{}

func (r *ReversingFormatter) Format(text string) string {
	l := len(text)
	reversed := make([]rune, l)
	for i, t := range text {
		reversed[l-i-1] = t
	}
	return string(reversed)
}

func NewWordReversingFormatter() Formatter {
	return &PerWordFormattingFormatter{&ReversingFormatter{}}
}
