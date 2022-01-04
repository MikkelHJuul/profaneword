package profaneword

import (
	"math/big"
	"strings"
	"unicode"
)

//NewRandomFormatter is a Returns a RandomlyFormattingCharFormatter that delegates to another CharFormatter at a rate of 50%
func NewRandomFormatter() *RandomlyFormattingCharFormatter {
	return &RandomlyFormattingCharFormatter{
		thresholdRandom: newFiftyFifty(),
	}
}

//RandomlyFormattingCharFormatter is a CharFormatter that delegates randomly to the embedded CharFormatter: Other
type RandomlyFormattingCharFormatter struct {
	thresholdRandom
	Other CharFormatter
}

func (rff *RandomlyFormattingCharFormatter) SetCharFormatter(wrap CharFormatter) {
	rff.Other = wrap
}

func (rff *RandomlyFormattingCharFormatter) GetCharFormatter() CharFormatter {
	return rff.Other
}

var _ WrappingCharFormatter = &RandomlyFormattingCharFormatter{}

//FormatRune formats a single rune or not given the Random with threshold
func (rff *RandomlyFormattingCharFormatter) FormatRune(r rune) []rune {
	if rff.Rand.Rand().Cmp(rff.Threshold) > 0 {
		return rff.Other.FormatRune(r)
	}
	return []rune{r}
}

var l337Map = map[rune][]rune{ //small subset of 1337 alphabet
	'A': {'4'},
	'B': {'8'},
	'E': {'3'},
	'G': {'6'},
	'I': {'1'},
	'L': {'1'},
	'O': {'0'},
	'S': {'5'},
	'T': {'7'},
	'Z': {'2'},
}

//curated list values from https://da.wikipedia.org/wiki/Leetspeak
var uberl337Map = map[rune][][]rune{
	'A': {{'4'}, {'/', '\\'}, {'@'}, {'/', '-', '\\'}},
	'B': {{'8'}, {'1', '3'}, {'|', '3'}, {'!', '3'}},
	'C': {{'['}, {'('}, {'<'}},
	'D': {{')'}, {'[', ')'}},
	'E': {{'3'}},
	'F': {{'|', '='}, {'|', '#'}},
	'G': {{'6'}, {'(', '_', '+'}},
	'H': {{'#'}, {']', '-', '['}, {'|', '-', '|'}},
	'I': {{'1'}, {'!'}, {'|'}},
	'J': {{'_', '|'}},
	'K': {{'|', '<'}},
	'L': {{'1'}, {'|', '_'}, {'|'}},
	'M': {{'|', 'v', '|'}, {'|', '\\', '/', '|'}},
	'N': {{'|', '\\', '|'}, {'|', 'V'}},
	'O': {{'0'}, {'(', ')'}},
	'P': {{'|', '>'}},
	'Q': {{'(', ')', '_'}},
	'R': {{'2'}, {'1', '2'}, {'|', '?'}},
	'S': {{'5'}, {'$'}, {'§'}, {'z'}, {'Z'}},
	'T': {{'7'}, {'+'}},
	'U': {{'(', '_', ')'}, {'|', '_', '|'}},
	'V': {{'\\', '/'}},
	'W': {{'\\', '/', '\\', '/'}, {'v', 'v'}, {'\'', '/', '/'}, {'\\', '\\', '\''}},
	'X': {{'>', '<'}, {'}', '{'}},
	'Y': {{'`', '/'}},
	'Z': {{'2'}, {'~', '/', '_'}},
}

//L337CharFormatter is a CharFormatter that formats by replacing
//the given rune by a slice if runes as stated in the internal map
type L337CharFormatter struct {
	uber1337 map[rune][]rune
}

var _ CharFormatter = L337CharFormatter{}

//FormatRune returns the rune-slice given in the internal map, or returns the input value
func (u L337CharFormatter) FormatRune(r rune) []rune {
	rKey := unicode.ToUpper(r)
	if leetVal, ok := u.uber1337[rKey]; ok {
		return leetVal
	}
	return []rune{r}
}

//Uber1337Formatter returns an initiated randomly chosen L337CharFormatter with the special uber1337-map
func Uber1337Formatter() Formatter {
	uber133Map := buildRandomMap(uberl337Map)
	return &CharFormatterDelegatingFormatter{
		&L337CharFormatter{
			uber133Map,
		},
	}
}

func buildRandomMap(m map[rune][][]rune) map[rune][]rune {
	randDev := CryptoRand{}
	var randomMap = make(map[rune][]rune, len(m))
	for k, v := range m {
		idx := randDev.RandMax(len(v))
		randomMap[k] = v[idx]
	}
	return randomMap
}

//L337Formatter returns a L337CharFormatter with a predefined mapping, the l337map
func L337Formatter() Formatter {
	return &CharFormatterDelegatingFormatter{
		&L337CharFormatter{
			l337Map,
		},
	}
}

var keyboard = [5]string{
	"1234567890-",
	"qwertyuiop[",
	"asdfghjkl;",
	"zxcvbnm,.",
	"       ",
}

func getNeighbourgChars(r rune) string {
	for j, line := range keyboard[1:3] {
		jIdx := j + 1
		idx := strings.Index(line, string(r))
		if idx != -1 {
			idxlower := idx - 1
			if idxlower == -1 {
				idxlower = idx
			}
			var lower string
			if r != 'z' && r != 'x' {
				lower = keyboard[jIdx+1][idxlower:idx]
			}
			return keyboard[jIdx-1][idx:idx+1] + keyboard[jIdx][idxlower:idx+1] + lower
		}
	}
	return string(r)
}

//FatFingerCharFormatter formats the text/rune as if it was types with fat fingers
type FatFingerCharFormatter struct {
	RandomDevice
}

var _ CharFormatter = FatFingerCharFormatter{}

//FormatRune returns the slice of runes by finding the neighboring characters (keyboard) and
//returns a random set of characters from within that sequence. it may return the rune itself up to four times
func (ff FatFingerCharFormatter) FormatRune(r rune) []rune {
	if ff.Rand().Cmp(big.NewRat(1, 6)) < 0 {
		var outRunes []rune
		for len(outRunes) == 0 {
			if ff.Rand().Cmp(big.NewRat(1, 6)) < 0 {
				outRunes = append(outRunes, r)
			}
			newChars := getNeighbourgChars(r)
			if ff.Rand().Cmp(big.NewRat(2, 5)) < 0 {
				outRunes = append(outRunes, rune(newChars[ff.RandMax(len(newChars))]))
			}
			if ff.Rand().Cmp(big.NewRat(1, 12)) < 0 {
				outRunes = append(outRunes, rune(newChars[ff.RandMax(len(newChars))]))
			}
			if ff.Rand().Cmp(big.NewRat(1, 7)) < 0 {
				outRunes = append(outRunes, r)
			}
		}
		return outRunes
	}
	return []rune{r}
}

//NewFatFingerFormatter wraps the FatFingerCharFormatter in a CharFormatterDelegatingFormatter to produce a Formatter
func NewFatFingerFormatter() Formatter {
	return &CharFormatterDelegatingFormatter{CharFormatter: FatFingerCharFormatter{CryptoRand{}}}
}

//FastFingerCharFormatter formats as if written with haste, skipping characters at random
type FastFingerCharFormatter struct {
	RandomDevice
}

var _ CharFormatter = FastFingerCharFormatter{}

//FormatRune at a rate of 1/6 randomly skip a rune
func (ff FastFingerCharFormatter) FormatRune(r rune) []rune {
	if ff.Rand().Cmp(big.NewRat(1, 6)) < 0 {
		return []rune{}
	}
	return []rune{r}
}

//NewFastFingerFormatter returns an initiated FastFingerCharFormatter wrapped in a CharFormatterDelegatingFormatter to produce a Formatter
func NewFastFingerFormatter() Formatter {
	return &CharFormatterDelegatingFormatter{CharFormatter: FastFingerCharFormatter{CryptoRand{}}}
}

//UppercaseCharFormatter formats uppercase
type UppercaseCharFormatter struct{}

var _ CharFormatter = UppercaseCharFormatter{}

//FormatRune uppercases the rune
func (UppercaseCharFormatter) FormatRune(r rune) []rune {
	return []rune{unicode.ToUpper(r)}
}

//LowercaseCharFormatter formats lowercase
type LowercaseCharFormatter struct{}

var _ CharFormatter = LowercaseCharFormatter{}

//FormatRune the rune, but lowercase
func (LowercaseCharFormatter) FormatRune(r rune) []rune {
	return []rune{unicode.ToLower(r)}
}

//SwitchCaseCharFormatter switch the case
type SwitchCaseCharFormatter struct{}

var _ CharFormatter = SwitchCaseCharFormatter{}

//FormatRune switches case of the rune
func (SwitchCaseCharFormatter) FormatRune(r rune) []rune {
	if unicode.IsUpper(r) {
		return []rune{unicode.ToLower(r)}
	} else {
		return []rune{unicode.ToUpper(r)}
	}
}

//NewSarcasticFormatter returns a CharFormatterDelegatingFormatter that wraps a
//RandomlyFormattingCharFormatter that randomly delegates to SwitchCaseCharFormatter
func NewSarcasticFormatter() Formatter {
	randomFormatter := NewRandomFormatter()
	randomFormatter.Other = &SwitchCaseCharFormatter{}
	return &CharFormatterDelegatingFormatter{randomFormatter}
}

type swearCharFormatter struct {
	RandomDevice
}

var _ CharFormatter = swearCharFormatter{
	CryptoRand{},
}

func (s swearCharFormatter) FormatRune(_ rune) []rune {
	letters := `#&$@%+*"`
	idx := s.RandMax(len(letters))
	return []rune{rune(letters[idx])}
}
