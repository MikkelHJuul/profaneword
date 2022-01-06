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
	studder     formatter = "studder"
	shuffle     formatter = "shuffle"
	horse       formatter = "horse"
	randomly    formatter = "randomly"
	random      formatter = "random"
)

var formatters = []string{
	string(l337), string(uberL337), string(sarcastic),
	string(scream), string(whisper), string(randomly),
	string(random), string(fatFingers), string(fastFingers),
	string(reverse), string(swear), string(studder),
	string(horse), string(shuffle),
}

type formatFunc func([]string, int) (int, profaneword.Formatter)

type plainFormatter func() profaneword.Formatter

func (p plainFormatter) formatF() formatFunc {
	return func(strings []string, i int) (int, profaneword.Formatter) {
		return i, p()
	}
}

var formatterFuncs map[formatter]formatFunc

func init() {
	formatterFuncs = map[formatter]formatFunc{
		sarcastic:   plainFormatter(profaneword.NewSarcasticFormatter).formatF(),
		l337:        plainFormatter(profaneword.L337Formatter).formatF(),
		uberL337:    plainFormatter(profaneword.Uber1337Formatter).formatF(),
		fatFingers:  plainFormatter(profaneword.NewFatFingerFormatter).formatF(),
		fastFingers: plainFormatter(profaneword.NewFastFingerFormatter).formatF(),
		scream:      plainFormatter(profaneword.NewUppercaseFormatter).formatF(),
		whisper:     plainFormatter(profaneword.NewLowercaseFormatter).formatF(),
		reverse:     plainFormatter(profaneword.NewWordReversingFormatter).formatF(),
		swear:       plainFormatter(profaneword.NewSwearFormatter).formatF(),
		studder:     plainFormatter(profaneword.NewStudderFormatter).formatF(),
		horse:       plainFormatter(profaneword.NewHorseFormatter).formatF(),
		shuffle:     plainFormatter(profaneword.NewShuffleFormatter).formatF(),
		randomly:    getRandomlyFormatter,
		random:      getRandomFormatter,
	}
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
	if formatterFunc, ok := formatterFuncs[formatter(args[i])]; ok {
		return formatterFunc(args, i)
	}
	return i, profaneword.UnitFormatter{}
}

func getRandomlyFormatter(args []string, i int) (int, profaneword.Formatter) {
	i++
	var wrappedFormatter profaneword.Formatter
	i, wrappedFormatter = getFormatter(args, i)
	randomlyFormatter := profaneword.NewRandomlyFormatter(wrappedFormatter)
	return i, randomlyFormatter
}

func getRandomFormatter(args []string, i int) (int, profaneword.Formatter) {
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
