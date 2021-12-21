package profanities

import (
	"fmt"
	"github.com/MikkelHJuul/profaneword"
	"strings"
)

//Sentence is a linked-list of formattable structures, each with a format string,
//a word (type) that would fit there, and a pointer to the next part of the Sentence
type Sentence struct {
	format string
	word   Word
	next   *Sentence
}

func (s *Sentence) getPart(word string) string {
	return fmt.Sprintf(s.format, word)
}

//Sentencer is any object that can return a string using a Sentence
type Sentencer interface {
	Sentence(*Sentence) string
}

//ProfanitySentencer is a type that implements Sentencer, while integrating the profanities database.
//ProfanitySentencer is configurable with dissallowed words.
type ProfanitySentencer struct {
	profaneword.RandomDevice
	dissallowedWord Word
}

var _ Sentencer = &ProfanitySentencer{}
var _ SentenceFetcher = &ProfanitySentencer{}

//Sentence implements the Sentencer interface, using randomized text from the profanities database.
func (pw *ProfanitySentencer) Sentence(sentence *Sentence) string {
	builder := strings.Builder{}
	for s := sentence; s != nil; s = s.next {
		text := pw.getRandomText(s.word, pw.dissallowedWord)
		builder.WriteString(s.getPart(text))
	}
	return builder.String()
}

func (pw *ProfanitySentencer) getRandomText(word, dissallowedWord Word) string {
	w := ""
	rLen := len(radix)
	for w == "" {
		base := radix[pw.RandMax(rLen)]
		wrds := base.GetOfSingle(word, dissallowedWord)
		if len(wrds) > 0 {
			w = wrds[pw.RandMax(len(wrds))]
		}
	}
	return w
}

//NewProfanitySentencer returns a ProfanitySentencer with the default configuration,
//passing a dissallowedWord to the Sentencer, and using a profaneword.CryptoRand
func NewProfanitySentencer(dissallowedWord Word) ProfanitySentencer {
	return ProfanitySentencer{profaneword.CryptoRand{}, dissallowedWord}
}

//SentenceFetcher is an interface for an object that returns a Sentence of a given length.
type SentenceFetcher interface {
	GetSentence(length int) *Sentence
}

//GetSentence implements SentenceFetcher for ProfanitySentencer.
//GetSentence builds a sentence of arbitrary length by using the internal
//flatSentence, recursively calling the internal map of flatSentence, and compiling a Sentence from it
func (pw *ProfanitySentencer) GetSentence(numWords int) *Sentence {
	sents, ok := sentences[numWords]
	if !ok {
		sentence := pw.GetSentence(5)
		cpy := sentence
		for cpy.next != nil {
			cpy = cpy.next
		}
		nextSent := pw.GetSentence(numWords - 5)
		nextSent.format = ` ` + nextSent.format
		cpy.next = nextSent
		return sentence
	}
	sentence := sents[pw.RandMax(len(sents))]
	return buildSentence(sentence)
}

//will not panic! test: TestSentencesCorrectLength validates the length!
func buildSentence(flat flatSentence) *Sentence {
	var sentence *Sentence = nil
	subst := `%s`
	indexPS := strings.LastIndex(flat.string, subst)
	var currFlat = flat
	for {
		format, next := currFlat.string[indexPS:], currFlat.string[:indexPS]
		sentence = &Sentence{
			format: format,
			word:   currFlat.words[len(currFlat.words)-1],
			next:   sentence,
		}

		indexPS = strings.LastIndex(next, subst)
		if indexPS == -1 {
			sentence.format = next + sentence.format
			break
		}
		currFlat = flatSentence{
			string: next,
			words:  currFlat.words[:len(currFlat.words)-1],
		}
	}
	return sentence
}

const EFE = EXCL | FILLER | END

type flatSentence struct {
	string
	words []Word
}

//TODO - programmatically?!
var sentences = map[int][]flatSentence{
	1: {
		{`%s`, []Word{EFE}},
		{`%s-fucker`, []Word{EFE}},
		{`sex-%s`, []Word{EFE}},
		{`%s!`, []Word{EFE}},
		{`%s?`, []Word{EFE}},
		{`%s?!`, []Word{EFE}},
		{`%s!?`, []Word{EFE}},
		{`%s!!`, []Word{EFE}},
		{`%s...NOT!`, []Word{EFE}},
	},
	2: {
		{`%s %s`, []Word{EFE | START, EFE}},
		{`%s sex-%s`, []Word{EFE | START, EFE}},
		{`%s %s-fucker`, []Word{EFE | START, EFE}},
		{`%s-%s`, []Word{EFE | START, EFE}},
		{`%s vs %s`, []Word{EFE, EFE}},
		{`%s vs. %s`, []Word{EFE, EFE}},
		{`%s, %s`, []Word{EFE | START, EFE}},
		{`%s: %s`, []Word{EFE | START, EFE}},
		{`%s; %s`, []Word{EFE | START, EFE}},
		{`%s - %s`, []Word{EFE | START, EFE}},
		{`%s -> %s`, []Word{EXCL | END, EXCL | END}},
		{`%s = %s`, []Word{EXCL | END, EXCL | END}},
		{`%s == %s`, []Word{EXCL | END, EXCL | END}},
		{`%s is %s`, []Word{EXCL | END, EXCL | END}},
		{`%s or %s`, []Word{EXCL | END, EXCL | END}},
		{`%s! %s`, []Word{EXCL | END, FILLER | END}},
		{`%s-fucker! %s`, []Word{EXCL | END, FILLER | END}},
		{`%s?! %s`, []Word{EXCL | END, FILLER | END}},
		{`%s? %s!`, []Word{DEFAULT, EXCL | END}},
		{`%s? sex-%s!`, []Word{DEFAULT, EXCL | END}},
		{`%s? %s-fucker!`, []Word{DEFAULT, EXCL | END}},
		{`%s? %s`, []Word{DEFAULT, EXCL | END}},
	},
	3: {
		{`%s %s %s`, []Word{EFE | START, FILLER, EFE}},
		{`%s %s sex-%s`, []Word{EFE | START, FILLER, EFE}},
		{`%s %s %s-fucker`, []Word{EFE | START, FILLER, EFE}},
		{`%s %s %s!`, []Word{EFE | START, FILLER, EFE}},
		{`%s-%s %s!`, []Word{DEFAULT, FILLER, EFE}},
		{`%s-%s sex-%s!`, []Word{DEFAULT, FILLER, EFE}},
		{`%s-%s %s-fucker!`, []Word{DEFAULT, FILLER, EFE}},
		{`%s vs %s %s`, []Word{EFE, FILLER, EFE}},
		{`%s %s vs. %s`, []Word{DEFAULT, EFE, EFE}},
		{`%s %s or %s`, []Word{DEFAULT, EFE, EFE}},
		{`%s, %s %s`, []Word{EFE | START, DEFAULT, END}},
		{`%s %s, %s`, []Word{EFE | START, FILLER, END}},
		{`%s %s, %s!`, []Word{EFE | START, FILLER, END}},
		{`%s %s, sex-%s!`, []Word{EFE | START, FILLER, END}},
		{`%s %s, %s-fucker!`, []Word{EFE | START, FILLER, END}},
		{`%s: %s %s`, []Word{EFE | START, DEFAULT, EFE}},
		{`%s %s: %s`, []Word{EFE | START, END, EFE}},
		{`%s; %s %s`, []Word{EFE | START, DEFAULT, EFE}},
		{`%s - %s %s`, []Word{EFE | START, DEFAULT, EFE}},
		{`%s %s - %s`, []Word{EFE | START, END, EFE}},
		{`%s -> %s %s`, []Word{EXCL | END, FILLER, EXCL | END}},
		{`%s %s! %s`, []Word{DEFAULT, EXCL | END, FILLER | END}},
		{`%s? %s %s!`, []Word{START | FILLER, FILLER, EXCL | END}},
		{`%s? %s %s`, []Word{START | FILLER, FILLER, EXCL | END}},
		{`%s %s? %s!`, []Word{START | FILLER, FILLER, EXCL | END}},
		{`%s %s? %s`, []Word{START | FILLER, FILLER, EXCL | END}},
	},
	4: {
		{`%s %s %s %s`, []Word{EFE | START, FILLER, FILLER, EFE}},
		{`%s %s %s sex-%s`, []Word{EFE | START, FILLER, FILLER, EFE}},
		{`%s %s %s %s-fucker`, []Word{EFE | START, FILLER, FILLER, EFE}},
		{`%s %s %s %s!`, []Word{EFE | START, FILLER, FILLER, EFE}},
		{`%s vs %s %s %s`, []Word{EFE, DEFAULT, FILLER, EFE}},
		{`%s vs. %s %s %s`, []Word{EFE, DEFAULT, FILLER, EFE}},
		{`%s %s %s vs %s`, []Word{DEFAULT, FILLER, EFE, EFE}},
		{`%s %s %s vs %s-fucker`, []Word{DEFAULT, FILLER, EFE, EFE}},
		{`%s %s %s vs. %s`, []Word{DEFAULT, FILLER, EFE, EFE}},
		{`%s %s %s vs. sex-%s`, []Word{DEFAULT, FILLER, EFE, EFE}},
		{`%s %s vs %s %s`, []Word{DEFAULT, EFE, DEFAULT, EFE}},
		{`%s %s vs. %s %s`, []Word{DEFAULT, EFE, DEFAULT, EFE}},
		{`%s, %s %s %s`, []Word{EFE | START, DEFAULT, FILLER, END}},
		{`%s %s %s, %s`, []Word{EFE | START, FILLER, FILLER, END}},
		{`%s %s %s, %s!`, []Word{EFE | START, FILLER, FILLER, END}},
		{`%s %s %s, sex-%s!`, []Word{EFE | START, FILLER, FILLER, END}},
		{`%s %s %s, %s-fucker!`, []Word{EFE | START, FILLER, FILLER, END}},
		{`%s: %s %s %s`, []Word{EFE | START, DEFAULT, FILLER, EFE}},
		{`%s %s: %s %s`, []Word{DEFAULT, EFE, DEFAULT, EFE}},
		{`%s %s %s: %s`, []Word{EFE | START, FILLER, END, EFE}},
		{`%s; %s %s %s`, []Word{EFE | START, DEFAULT, FILLER, EFE}},
		{`%s? %s %s %s`, []Word{EXCL | END, DEFAULT, FILLER, EFE}},
		{`%s - %s %s %s`, []Word{EFE | START, DEFAULT, FILLER, EFE}},
		{`%s %s! %s %s`, []Word{DEFAULT, EFE, DEFAULT, EFE}},
		{`%s? %s %s %s!`, []Word{DEFAULT, DEFAULT, FILLER, EXCL | END}},
		{`%s %s? %s %s!`, []Word{START | FILLER, FILLER, DEFAULT, EXCL | END}},
		{`%s %s-fucker? %s %s!`, []Word{START | FILLER, FILLER, DEFAULT, EXCL | END}},
		{`%s %s? %s %s`, []Word{START | FILLER, FILLER, DEFAULT, EXCL | END}},
	},
	5: {
		{`%s %s %s %s %s`, []Word{EFE | START, FILLER, FILLER, FILLER, EFE}},
		{`%s %s %s %s sex-%s`, []Word{EFE | START, FILLER, FILLER, FILLER, EFE}},
		{`%s %s %s %s %s-fucker`, []Word{EFE | START, FILLER, FILLER, FILLER, EFE}},
		{`%s %s %s %s %s!`, []Word{EFE | START, FILLER, FILLER, FILLER, EFE}},
		{`%s vs %s %s %s %s`, []Word{EFE, DEFAULT, FILLER, FILLER, EFE}},
		{`%s vs. %s %s %s %s`, []Word{EFE, DEFAULT, FILLER, FILLER, EFE}},
		{`%s %s %s %s vs %s`, []Word{DEFAULT, FILLER, FILLER, EFE, EFE}},
		{`%s %s %s %s vs sex-%s`, []Word{DEFAULT, FILLER, FILLER, EFE, EFE}},
		{`%s %s %s %s vs %s-fucker`, []Word{DEFAULT, FILLER, FILLER, EFE, EFE}},
		{`%s %s %s %s vs. %s`, []Word{DEFAULT, FILLER, FILLER, EFE, EFE}},
		{`%s %s %s vs %s %s`, []Word{DEFAULT, FILLER, EFE, DEFAULT, EFE}},
		{`%s %s %s vs. %s %s`, []Word{DEFAULT, FILLER, EFE, DEFAULT, EFE}},
		{`%s, %s %s %s %s`, []Word{EFE | START, DEFAULT, FILLER, FILLER, END}},
		{`%s %s %s, %s %s`, []Word{EFE | START, FILLER, FILLER, FILLER, END}},
		{`%s %s %s, %s %s!`, []Word{EFE | START, FILLER, FILLER, FILLER, END}},
		{`%s %s %s, %s %s-fucker!`, []Word{EFE | START, FILLER, FILLER, FILLER, END}},
		{`%s %s %s, %s sex-%s!`, []Word{EFE | START, FILLER, FILLER, FILLER, END}},
		{`%s: %s %s %s %s`, []Word{EFE | START, DEFAULT, FILLER, FILLER, EFE}},
		{`%s %s: %s %s %s`, []Word{DEFAULT, EFE, DEFAULT, FILLER, EFE}},
		{`%s %s %s %s: %s`, []Word{EFE | START, FILLER, FILLER, END, EFE}},
		{`%s %s %s %s: %s-fucker`, []Word{EFE | START, FILLER, FILLER, END, EFE}},
		{`%s; %s %s %s %s`, []Word{EFE | START, DEFAULT, FILLER, FILLER, EFE}},
		{`%s? %s %s %s %s`, []Word{EXCL | END, DEFAULT, FILLER, FILLER, EFE}},
		{`%s %s? %s %s %s`, []Word{EXCL | END, DEFAULT, FILLER, FILLER, EFE}},
		{`%s-%s? %s %s %s`, []Word{EXCL | END, DEFAULT, FILLER, FILLER, EFE}},
		{`%s - %s %s %s %s`, []Word{EFE | START, DEFAULT, FILLER, FILLER, EFE}},
		{`%s %s - %s %s %s`, []Word{EFE | START, DEFAULT, FILLER, FILLER, EFE}},
		{`%s %s %s! %s %s`, []Word{DEFAULT, FILLER, EFE, DEFAULT, EFE}},
		{`%s %s? %s %s %s!`, []Word{DEFAULT, END, DEFAULT, FILLER, EXCL | END}},
		{`%s %s %s? %s %s!`, []Word{START | FILLER, FILLER, END, DEFAULT, EXCL | END}},
		{`%s %s %s? %s sex-%s!`, []Word{START | FILLER, FILLER, END, DEFAULT, EXCL | END}},
		{`%s %s %s? %s %s-fucker!`, []Word{START | FILLER, FILLER, END, DEFAULT, EXCL | END}},
		{`%s %s %s? %s %s`, []Word{START | FILLER, FILLER, END, DEFAULT, EXCL | END}},
	},
}
