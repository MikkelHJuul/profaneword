package profaneword

import (
	"regexp"
	"strings"
	"unicode"
)

// Formatter is a formatter, that formats the entire input text(string) and outputs the formatted text
type Formatter interface {
	Format(string) string
}

// CharFormatter is a formatter that emits a slice of runes given a single rune.
type CharFormatter interface {
	FormatRune(rune) []rune
}

// WrappingCharFormatter is an interface of a CharFormatter that wraps another CharFormatter
type WrappingCharFormatter interface {
	CharFormatter
	SetCharFormatter(CharFormatter)
	GetCharFormatter() CharFormatter
}

// WrappingFormatter is an interface of a Formatter that wraps another Formatter
type WrappingFormatter interface {
	Formatter
	SetFormatter(Formatter)
	GetFormatter() Formatter
}

// MultiFormatter is a wrapper to handle formatting with multiple Formatters
type MultiFormatter struct {
	Formatters []Formatter
}

// With appends a given Formatter to the Formatters in MultiFormatter
func (m *MultiFormatter) With(f Formatter) {
	if f != nil {
		m.Formatters = append(m.Formatters, f)
	}
}

// Format delegates/reassigns the input string sequentially for all formatters
func (m *MultiFormatter) Format(word string) string {
	for _, f := range m.Formatters {
		word = f.Format(word)
	}
	return word
}

// UnitFormatter is a Noop Formatter and CharFormatter implementation
type UnitFormatter struct{}

var _ Formatter = UnitFormatter{}
var _ CharFormatter = UnitFormatter{}

// Format returns the input
func (UnitFormatter) Format(words string) string {
	return words
}

// FormatRune return the rune wrapped in a slice
func (UnitFormatter) FormatRune(r rune) []rune {
	return []rune{r}
}

// CharFormatterDelegatingFormatter is a Formatter that delegates each letter in the input to a CharFormatter
type CharFormatterDelegatingFormatter struct {
	CharFormatter
}

var _ Formatter = &CharFormatterDelegatingFormatter{}

// Format calls the wrapped CharFormatter's FormatRune-method
func (c *CharFormatterDelegatingFormatter) Format(word string) string {
	var out []rune
	for _, r := range word {
		out = append(out, c.FormatRune(r)...)
	}
	return string(out)
}

// SetCharFormatter sets the CharFormatter to be wrapped by CharFormatterDelegatingFormatter
func (c *CharFormatterDelegatingFormatter) SetCharFormatter(formatter CharFormatter) {
	c.CharFormatter = formatter
}

// GetCharFormatter returns the wrapped CharFormatter, which is delegated to
func (c *CharFormatterDelegatingFormatter) GetCharFormatter() CharFormatter {
	return c.CharFormatter
}

// PerWordFormattingFormatter delegates to another Formatter, but calls the Format method for each word
type PerWordFormattingFormatter struct {
	Other Formatter
}

// SetFormatter sets the Formatter that PerWordFormattingFormatter should wrap
func (p *PerWordFormattingFormatter) SetFormatter(formatter Formatter) {
	p.Other = formatter
}

// GetFormatter returns the wrapped formatter of PerWordFormattingFormatter
func (p *PerWordFormattingFormatter) GetFormatter() Formatter {
	return p.Other
}

// Format splits the text on space, and calls Other > Format for each word, then joining the text again.
func (p *PerWordFormattingFormatter) Format(word string) string {
	words := strings.Split(word, ` `)
	ws := make([]string, len(words))
	for i, word := range words {
		ws[i] = p.Other.Format(word)
	}
	return strings.Join(ws, ` `)
}

// RandomlyFormattingFormatter is a formatter that formats at a rate of 50%
type RandomlyFormattingFormatter struct {
	thresholdRandom
	Other Formatter
}

var _ Formatter = &RandomlyFormattingFormatter{}

// Format calls the Other > Format at a rate of 50%
func (rff *RandomlyFormattingFormatter) Format(word string) string {
	if rff.Rand.Rand().Cmp(rff.Threshold) > 0 {
		return rff.Other.Format(word)
	}
	return word
}

// SetFormatter sets the wrapped Formatter of RandomlyFormattingFormatter
func (rff *RandomlyFormattingFormatter) SetFormatter(wrap Formatter) {
	rff.Other = wrap
}

// GetFormatter returns the Formatter that RandomlyFormattingFormatter wraps
func (rff *RandomlyFormattingFormatter) GetFormatter() Formatter {
	return rff.Other
}

// NewRandomlyFormatter is a method that wraps a given Formatter,
// with the RandomlyFormattingFormatter that is wrapped with a PerWordFormattingFormatter
func NewRandomlyFormatter(wrap Formatter) Formatter {
	random := &RandomlyFormattingFormatter{thresholdRandom: newFiftyFifty()}
	random.SetFormatter(wrap)
	return &PerWordFormattingFormatter{random}
}

// TitleFormatter is a Formatter that Titles the given text
type TitleFormatter struct{}

var _ Formatter = TitleFormatter{}

// Format calls strings.Title on the given text
func (TitleFormatter) Format(word string) string {
	return strings.Title(word)
}

// RandomTitleFormatter returns a formatter that titles only every other time
func RandomTitleFormatter() Formatter {
	return NewRandomlyFormatter(TitleFormatter{})
}

var _ Formatter = &RandomlyFormattingFormatter{}

// RegexReplacingFormatter is a Formatter that performs a regex replace functionality on the given text
type RegexReplacingFormatter struct {
	// PatternMatcher is a regexp.Regexp, to match against the text
	PatternMatcher *regexp.Regexp
	// Replacement is whatever is to be replaced by the given PatternMatcher
	Replacement string
}

// Format replaces all instances of the given regex PatternMatcher with the Replacement
func (rr *RegexReplacingFormatter) Format(word string) string {
	return rr.PatternMatcher.ReplaceAllString(word, rr.Replacement)
}

// DelimiterFormatterWith replaces all spaces with the given string
func DelimiterFormatterWith(repl string) Formatter {
	return &RegexReplacingFormatter{
		regexp.MustCompile(` `),
		repl,
	}
}

// ReversingFormatter reverses the string
type ReversingFormatter struct{}

// Format reverses the input text
func (ReversingFormatter) Format(text string) string {
	l := len(text)
	reversed := make([]rune, l)
	for i, t := range text {
		reversed[l-i-1] = t
	}
	return string(reversed)
}

// NewWordReversingFormatter returns a ReversingFormatter that reverses each words in a group,
// and not the entire text as one
func NewWordReversingFormatter() Formatter {
	return &PerWordFormattingFormatter{ReversingFormatter{}}
}

type swearFormatter struct {
	CharFormatter
}

// SetCharFormatter sets the formatter to be used by swearFormatter
func (s *swearFormatter) SetCharFormatter(wrap CharFormatter) {
	s.CharFormatter = wrap
}

// GetCharFormatter return the swearFormatter CharFormatter
func (s *swearFormatter) GetCharFormatter() CharFormatter {
	return s.CharFormatter
}

var _ Formatter = &swearFormatter{swearCharFormatter{}}
var _ WrappingCharFormatter = &swearFormatter{swearCharFormatter{}}

// Format of swearFormatter will return the input string if the input does not start with a letter,
// for all other cases it replaces using the wrapped CharFormatter until it meets a non-letter character
func (s *swearFormatter) Format(word string) string {
	if !unicode.IsLetter(rune(word[0])) {
		return word
	}
	var runes []rune
	var i int
	var c rune
	for i, c = range word {
		if unicode.IsLetter(c) {
			runes = append(runes, s.FormatRune(c)...)
		} else {
			break
		}
	}
	suffix := ""
	if i < len(word)-1 {
		suffix = word[i:]
	}
	return string(runes) + `!` + suffix
}

// NewSwearFormatter reuturns a Formatter that replaces each character in a word with cartoonish swear
func NewSwearFormatter() Formatter {
	return &PerWordFormattingFormatter{&swearFormatter{&swearCharFormatter{
		CryptoRand{},
	}}}
}

// StudderFormatter is a formatter that writes text that appear like it's studdering
type StudderFormatter struct {
	RandomDevice
}

// Format will return the first character of a word plus '-' up to 4 times, and the word
func (s StudderFormatter) Format(word string) string {
	if !unicode.IsLetter(rune(word[0])) {
		return word
	}
	numStudder := s.RandMax(4)
	sb := strings.Builder{}
	for i := 0; i < numStudder; i++ {
		sb.WriteString(word[:1] + "-")
	}
	return sb.String() + word
}

var _ Formatter = StudderFormatter{}

// NewStudderFormatter returns a PerWordFormattingFormatter that wraps a StudderFormatter
func NewStudderFormatter() Formatter {
	return &PerWordFormattingFormatter{StudderFormatter{CryptoRand{}}}
}

// HorseFormatter returns horse-related banter for each call
type HorseFormatter struct {
	RandomDevice
}

var horsewords = []string{
	"horse",
	"horsey",
	"Horsey-Mc-Horseface",
	"mare",
	"stallion",
	"mustang",
	"ride",
	"riding",
	"saddle",
	"saddlebags",
	"blinders",
	"colt",
	"filly",
	"bronco",
	"foal",
	"gilding",
	"steed",
	"pony",
	"cavalry",
	"cavalier",
	"jenny",
	"equine",
	"jock",
	"jockey",
	"horseshoe",
	"horseback",
	"yearling",
	"dressage",
	"strider",
	"trot",
	"trotting",
	"bridle",
	"stirrup",
	"spurs",
	"bareback",
	"bareback",
	"reins",
	"broodmare",
	"damsire",
	"gallop",
	"mule",
	"whip",
	"quirt",
	"stud",
	"zebroid",
}

// Format returns a random horse-related word
func (s HorseFormatter) Format(_ string) string {
	idx := s.RandMax(len(horsewords))
	return horsewords[idx]
}

var _ Formatter = HorseFormatter{}

// NewHorseFormatter returns a PerWordFormattingFormatter wrapping a HorseFormatter
func NewHorseFormatter() Formatter {
	return &PerWordFormattingFormatter{HorseFormatter{CryptoRand{}}}
}

// ShuffleFormatter shuffles the given string
type ShuffleFormatter struct {
	RandomDevice
}

// Format shuffles the characters of an input string and returns a new shuffled string
func (s ShuffleFormatter) Format(text string) string {
	word := []byte(text)
	wlen := len(word)
	shuffledWord := make([]byte, wlen)
	for i := 0; i < wlen; i++ {
		idx := s.RandMax(wlen - i)
		shuffledWord[i] = word[idx]
		word[idx] = word[wlen-1-i]
	}
	return string(shuffledWord)
}

var _ Formatter = ShuffleFormatter{}

// NewShuffleFormatter returns a PerWordFormattingFormatter wrapping a ShuffleFormatter
func NewShuffleFormatter() Formatter {
	return &PerWordFormattingFormatter{ShuffleFormatter{CryptoRand{}}}
}
