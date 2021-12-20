package cmd

import (
	"github.com/MikkelHJuul/profaneword"
)

type formatter string

const (
	sarcastic   formatter = "/s"
	l337        formatter = "1337"
	uberL337    formatter = "uber1337"
	fatFingers  formatter = "fat"
	fastFingers formatter = "fst"
	SCREAM      formatter = "SCREAM"
	whisper     formatter = "whisper"
	randomly    formatter = "randomly"
	random      formatter = "random"
)

var formatters = []string{
	string(l337), string(uberL337), string(sarcastic),
	string(SCREAM), string(whisper), string(randomly),
	string(random), string(fatFingers), string(fastFingers),
}

func formatterOf(args []string, formatters ...profaneword.Formatter) profaneword.Formatter {
	mulF := &profaneword.MultiFormatter{Formatters: formatters}
	for i := 0; i < len(args); i++ {
		var formatter profaneword.Formatter
		i, formatter = getFormatter(args, i)
		mulF.With(formatter)
	}
	return mulF
}

func getFormatter(args []string, i int) (int, profaneword.Formatter) {
	if i == len(args) {
		return i, profaneword.UnitFormatter{}
	}
	arg := args[i]
	switch formatter(arg) {
	case sarcastic:
		return i, profaneword.NewSarcasticFormatter()
	case l337:
		return i, profaneword.L337Formatter()
	case uberL337:
		return i, profaneword.Uber1337Formatter()
	case fatFingers:
		return i, profaneword.NewFatFingerFormatter()
	case fastFingers:
		return i, profaneword.NewFastFingerFormatter()
	case SCREAM:
		return i, &profaneword.CharFormatterDelegatingFormatter{CharFormatter: profaneword.UppercaseCharFormatter{}}
	case whisper:
		return i, &profaneword.CharFormatterDelegatingFormatter{CharFormatter: profaneword.LowercaseCharFormatter{}}
	case randomly:
		i++
		var wrappedFormatter profaneword.Formatter
		i, wrappedFormatter = getFormatter(args, i)
		randomlyFormatter := profaneword.NewRandomlyFormatter()
		randomlyFormatter.Other = wrappedFormatter
		return i, randomlyFormatter
	case random:
		i++
		var wrappedFormatter profaneword.Formatter
		i, wrappedFormatter = getFormatter(args, i)
		charFormatter, ok := wrappedFormatter.(profaneword.CharFormatter)
		if !ok {
			//this should only happen with random randomly, which is already validated against...
			//all other formatters are CharFormatters atm.
			return i, profaneword.UnitFormatter{}
		}
		randomFormatter := profaneword.NewRandomFormatter()
		randomFormatter.Other = charFormatter
		return i, &profaneword.CharFormatterDelegatingFormatter{CharFormatter: randomFormatter}
	}
	return i, profaneword.UnitFormatter{}
}
