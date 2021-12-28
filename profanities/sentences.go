package profanities

import (
	"fmt"
	"github.com/MikkelHJuul/profaneword"
	"strings"
)

type sentnc struct {
	format string
	word   Word
}

//Sentence is a linked-list of formattable structures, each with a format string,
//a word (type) that would fit there, and a pointer to the next part of the Sentence
type Sentence struct {
	sentnc
	next *Sentence
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

	randSent := func() (*Sentence, sent) {
		idx := pw.RandMax(len(sentences))
		s := sentences[idx]
		return &Sentence{sentnc: s.sentnc}, s
	}
	var cur *Sentence
	var s sent
	pos := not_last
	for pos&not_last != 0 {
		cur, s = randSent()
		pos = s.sentPos
	}
	if strings.HasSuffix(cur.format, " ") {
		cur.format = cur.format[:len(cur.format)-1]
	}
	numWords--
	for ; numWords > 0; numWords-- {
		prev := cur
		cur, s = randSent()
		cur.next = prev
	}
	return cur
}

const efe = EXCL | FILLER | END
const all Word = 255

type sentPos uint8

const (
	not_last sentPos = 16
)

type sent struct {
	sentnc
	sentPos
}

var sentences = [...]sent{
	{sentnc: sentnc{format: `%s `, word: all}},
	{sentnc: sentnc{format: `%s-fucker `, word: efe}},
	{sentnc: sentnc{format: `%s-fucker! `, word: efe}},
	{sentnc: sentnc{format: `sex-%s `, word: efe}},
	{sentnc: sentnc{format: `sex-%s! `, word: efe}},
	{sentnc: sentnc{format: `%s! `, word: all}},
	{sentnc: sentnc{format: `%s? `, word: all}},
	{sentnc: sentnc{format: `%s?! `, word: all}},
	{sentnc: sentnc{format: `%s!? `, word: all}},
	{sentnc: sentnc{format: `%s!! `, word: efe}},
	{sentnc: sentnc{format: `%s...NOT! `, word: efe}},
	{sentnc: sentnc{format: `%s-`, word: all}, sentPos: not_last},
	{sentnc: sentnc{format: `%s vs `, word: all}, sentPos: not_last},
	{sentnc: sentnc{format: `%s vs. `, word: all}, sentPos: not_last},
	{sentnc: sentnc{format: `%s, `, word: all}, sentPos: not_last},
	{sentnc: sentnc{format: `%s: `, word: all}, sentPos: not_last},
	{sentnc: sentnc{format: `%s; `, word: all}, sentPos: not_last},
	{sentnc: sentnc{format: `%s - `, word: all}, sentPos: not_last},
	{sentnc: sentnc{format: `%s -> `, word: all}, sentPos: not_last},
	{sentnc: sentnc{format: `%s = `, word: all}, sentPos: not_last},
	{sentnc: sentnc{format: `%s == `, word: all}, sentPos: not_last},
	{sentnc: sentnc{format: `%s is `, word: all}, sentPos: not_last},
	{sentnc: sentnc{format: `%s or `, word: all}, sentPos: not_last},
	{sentnc: sentnc{format: `%s of `, word: all}, sentPos: not_last},
	{sentnc: sentnc{format: `%s `, word: efe}},
}
