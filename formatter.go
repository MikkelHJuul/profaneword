package profaneword

import (
	"math/big"
	"unicode"
)

type Formatter interface {
	Format(string) string
}

type MultiFormatter struct {
	formatters []Formatter
}

func (m *MultiFormatter) With(f Formatter) {
	m.formatters = append(m.formatters, f)
}

func (m MultiFormatter) Format(word string) string {
	for _, f := range m.formatters {
		word = f.Format(word)
	}
	return word
}

type SarcasticFormatter struct {
	Rand      RandomDevice
	Threshold *big.Rat
}

func (sf SarcasticFormatter) Format(word string) string {
	out := make([]rune, len(word))
	for i, r := range word {
		if sf.Rand.Rand().Cmp(sf.Threshold) > 0 {
			if unicode.IsUpper(r) {
				out[i] = unicode.ToLower(r)
			} else {
				out[i] = unicode.ToUpper(r)
			}
		} else {
			out[i] = r
		}
	}
	return string(out)
}

func NewSarcasticFormatter(device RandomDevice, threshold *big.Rat) SarcasticFormatter {
	if threshold == nil {
		threshold = big.NewRat(1, 2)
	}
	if device == nil { //pointer-nil??
		device = CryptoRand{}
	}
	return SarcasticFormatter{
		Rand:      device,
		Threshold: threshold,
	}
}

var l337Map = map[rune]rune{ //small subset of 1337 alphabet
	'A': '4',
	'B': '8',
	'E': '3',
	'G': '6',
	'I': '1',
	'L': '1',
	'O': '0',
	'S': '5',
	'T': '7',
	'Z': '2',
}

type L337Formatter struct{}

func (sf L337Formatter) Format(word string) string {
	out := make([]rune, len(word))
	for i, r := range word {
		rKey := unicode.ToUpper(r)
		if leetVal, ok := l337Map[rKey]; ok {
			out[i] = leetVal
		} else {
			out[i] = r
		}
	}
	return string(out)
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

type UberL337Formatter struct{}

func (_ UberL337Formatter) Format(word string) string {
	out := make([]rune, len(word))
	randDevice := CryptoRand{}
	for i, r := range word {
		rKey := unicode.ToUpper(r)
		if leetVal, ok := uberl337Map[rKey]; ok {
			randomIndex := randDevice.RandMax(int64(len(leetVal)) + 1)
			leetChars := append(leetVal, []rune{r})[randomIndex]
			out = append(out, leetChars...)
		} else {
			out[i] = r
		}
	}
	return string(out)
}
