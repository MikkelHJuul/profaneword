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
	scream      formatter = "SCREAM"
	whisper     formatter = "whisper"
	reverse     formatter = "esrever"
	swear       formatter = "swear"
	randomly    formatter = "randomly"
	random      formatter = "random"
)

var formatters = []string{
	string(l337), string(uberL337), string(sarcastic),
	string(scream), string(whisper), string(randomly),
	string(random), string(fatFingers), string(fastFingers),
	string(reverse), string(swear),
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
	case scream:
		return i, &profaneword.CharFormatterDelegatingFormatter{CharFormatter: profaneword.UppercaseCharFormatter{}}
	case whisper:
		return i, &profaneword.CharFormatterDelegatingFormatter{CharFormatter: profaneword.LowercaseCharFormatter{}}
	case reverse:
		return i, profaneword.NewWordReversingFormatter()
	case swear:
		return i, profaneword.NewSwearFormatter()
	case randomly:
		i++
		var wrappedFormatter profaneword.Formatter
		i, wrappedFormatter = getFormatter(args, i)
		randomlyFormatter := profaneword.NewRandomlyFormatter(wrappedFormatter)
		return i, randomlyFormatter
	case random:
		i++
		var wrappedFormatter profaneword.Formatter
		i, wrappedFormatter = getFormatter(args, i)
		charFormatter, ok := wrappedFormatter.(profaneword.CharFormatter)
		if !ok {
			if delegating, isType := wrappedFormatter.(profaneword.WrappingFormatter); isType {
				if charFormatter, ok = delegating.GetFormatter().(profaneword.CharFormatter); ok {
					formatterToWrap := wrapRandom(charFormatter)
					delegating.SetFormatter(formatterToWrap)
					return i, delegating
				}
			}
			return i, wrappedFormatter
		}
		wrapped := wrapRandom(charFormatter)
		if delegating, isType := wrappedFormatter.(profaneword.WrappingFormatter); isType {
			delegating.SetFormatter(wrapped)
			return i, delegating
		}
		return i, wrapped
	}
	return i, profaneword.UnitFormatter{}
}

func wrapRandom(charFormatter profaneword.CharFormatter) profaneword.Formatter {
	randomFormatter := profaneword.NewRandomFormatter()
	if delegating, isType := charFormatter.(profaneword.WrappingCharFormatter); isType {
		randomFormatter.SetCharFormatter(delegating.GetCharFormatter()) // we don't want a circular reference
		delegating.SetCharFormatter(randomFormatter)
		return delegating.(profaneword.Formatter) // this is OKAY, as this is given already
	}
	randomFormatter.SetCharFormatter(charFormatter)
	return &profaneword.CharFormatterDelegatingFormatter{CharFormatter: randomFormatter}
}
