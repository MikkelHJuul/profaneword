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

type SarcasticFormatter struct{
	Rand RandomDevice
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
		Rand: device,
		Threshold: threshold,
	}
}

var l337Map = map[rune]rune {  //small subset of 1337 alphabet
	'a': '4',
	'A': '4',
	'b': '8',
	'B': '8',
	'e': '3',
	'E': '3',
	'g': '6',
	'G': '6',
	'i': '1',
	'I': '1',
	'l': '1',
	'L': '1',
	'o': '0',
	'O': '0',
	's': '5',
	'S': '5',
	't': '7',
	'T': '7',
	'z': '2',
	'Z': '2',
}

type L337Formatter struct {}

func (sf L337Formatter) Format(word string) string {
	out := make([]rune, len(word))
	for i, r := range word {
		if leetVal, ok := l337Map[r]; ok {
			out[i] = leetVal
		} else {
			out[i] = r
		}
	}
	return string(out)
}
