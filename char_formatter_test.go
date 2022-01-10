package profaneword

import (
	"math/big"
	"reflect"
	"testing"
)

func TestL337CharFormatter_FormatRune(t *testing.T) {
	abFormatter := L337CharFormatter{map[rune][]rune{'a': {'b'}}}
	got := abFormatter.FormatRune('a')
	if len(got) != 1 && got[0] != 'b' {
		t.Errorf("L337CharFormatter misformatted, expected: {'b'} got: %v", got)
	}
	got = abFormatter.FormatRune('d')
	if len(got) != 1 && got[0] != 'd' {
		t.Errorf("L337CharFormatter misformatted, expected: {'d'} got: %v", got)
	}
}

type zeroRandomDevice int

func (z zeroRandomDevice) Rand() *big.Rat {
	return big.NewRat(int64(z), 1)
}

func (zeroRandomDevice) RandMax(_ int) int {
	return 0
}

var _ RandomDevice = zeroRandomDevice(0)

func TestFastFingerCharFormatter_FormatRune(t *testing.T) {
	ff := FastFingerCharFormatter{zeroRandomDevice(2)}
	if got := ff.FormatRune('a'); len(got) != 1 && got[0] != 'a' {
		t.Errorf("FastFingerCharFormatter did return unit-slice as expected")
	}
	ff.RandomDevice = zeroRandomDevice(0)
	if got := ff.FormatRune('a'); len(got) != 0 {
		t.Errorf("FastFingerCharFormatter did return empty slice, as expected")
	}
}

func TestFatFingerCharFormatter_FormatRune(t *testing.T) {
	fatf := FatFingerCharFormatter{zeroRandomDevice(0)}
	if got := fatf.FormatRune('a'); len(got) != 1 && got[0] != 'a' {
		t.Errorf("FatFingerCharFormatter did return unit-slice as expected, got: %v", got)
	}
	fatf.RandomDevice = zeroRandomDevice(0)
	if got := fatf.FormatRune('a'); len(got) != 4 && reflect.DeepEqual([]rune{'a', 'q', 'q', 'a'}, got) {
		t.Errorf("FatFingerCharFormatter did return aqqa, as expected, got: %v", got)
	}
}

func TestL337Formatter(t *testing.T) {
	l := L337Formatter()
	if l.Format("asd") != "45d" {
		t.Errorf("L337formatter did not format as expected")
	}
}
