package cmd

import (
	"fmt"
	"github.com/MikkelHJuul/profaneword/profanities"
	"strings"
)

func GetWords(words int) string {
	return getWords(words, funs)
}

func getWords(words int, wfuns [][]fmt.Stringer) string{
	numFuns := len(funs)
	var stringers []fmt.Stringer
	for ; words > 0 ; words -= numFuns{
		idx := words
		if words > numFuns{
			idx = numFuns
		}
		stringers = append(stringers, wfuns[idx - 1]...)
	}

	var text strings.Builder

	stringers, last := stringers[:len(stringers)-1], stringers[len(stringers)-1]
	for _, f := range stringers {
		text.WriteString(f.String())
		text.WriteRune(' ')
	}
	text.WriteString(last.String())

	return text.String()
}

var funs = [][]fmt.Stringer{
	{profanities.NounStringer{}},
	{profanities.VerbStringer{}, profanities.NounStringer{}},
	{profanities.VerbStringer{}, profanities.SwearStringer{}, profanities.NounStringer{}},
	{profanities.VerbStringer{}, profanities.SwearStringer{}, profanities.VerbStringer{}, profanities.NounStringer{}},
	{profanities.VerbStringer{}, profanities.NounStringer{}, profanities.SwearStringer{}, profanities.VerbStringer{}, profanities.NounStringer{}},
}