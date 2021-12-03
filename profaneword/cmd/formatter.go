package cmd

import "github.com/MikkelHJuul/profaneword"

type formatter string

const (
	sarcastic formatter = "/s"
	l337 formatter = "1337"
)
var formatters = []string{string(l337), string(sarcastic)}

func formatterOf(args []string) profaneword.Formatter {
	mul := &profaneword.MultiFormatter{}
	for _, arg := range args {
		switch formatter(arg) {
		case sarcastic:
			mul.With(profaneword.NewSarcasticFormatter(nil, nil))
		case l337:
			mul.With(profaneword.L337Formatter{})
		default:
		}
	}
	return mul
}
