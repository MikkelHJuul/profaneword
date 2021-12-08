package cmd

import "github.com/MikkelHJuul/profaneword"

type formatter string

const (
	sarcastic formatter = "/s"
	l337      formatter = "1337"
	uberL337  formatter = "uber1337"
)

var formatters = []string{string(l337), string(uberL337), string(sarcastic)}

func formatterOf(args []string) profaneword.Formatter {
	mulF := &profaneword.MultiFormatter{}
	for _, arg := range args {
		switch formatter(arg) {
		case sarcastic:
			mulF.With(profaneword.NewSarcasticFormatter(nil, nil))
		case l337:
			mulF.With(profaneword.L337Formatter{})
		case uberL337:
			mulF.With(profaneword.UberL337Formatter{})
		default:
		}
	}
	return mulF
}
