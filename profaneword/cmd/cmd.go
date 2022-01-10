package cmd

import (
	"bufio"
	"github.com/MikkelHJuul/profaneword"
	"github.com/MikkelHJuul/profaneword/profanities"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
)

// usageTpl is a copy-paste of the underlying UsageTemplate in order to inject Args (and [..args])
const (
	usageTpl = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}} [..args]
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Args[formatters]:
  1337            output formatted as 1337-speak
  uber1337        output formatted with an extended 1337 alphabet
  fat             output some t3xt wifth fat fringers
  fst             otput sme tet writen wit haste
  esrever         desrever tuptuo, per word [random does not apply]
  shuffle         tuoput si ffudlehs
  SCREAM          OUTPUT IS UPPERCASE
  whisper         output is lowercased
  swear           output cartoonish #%$@!! 
  studder         o-o-output s-s-s-studdering t-text [random does not apply]
  horse           just output horse-related words in stead[very unsafe] [random does not apply]
  /s              sARcaSTiC OUtpUt
	
  randomly        the next formatter is applied only randomly (per word basis) threshold is 50:50
  random          the next formatter is applied only randomly (per character basis) threshold is 50:50
                  both "random" and "randomly" are chainable onto themselves, 
                  though "randomly" must be before "random"


Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
	alternateDelimiters = ".-/_:$%^+=!@'`,|<>\"~\\?*&"
	RAND                = "RAND"
)

var (
	profaneCmd = &cobra.Command{
		Use:       "profaneword",
		Short:     "A generator for profane passwords as requested by u/gatestone",
		Long:      `profaneword is a program for generating obscene/profane passwords.`,
		Args:      cobra.OnlyValidArgs,
		ValidArgs: formatters,
		Run:       profaneWords,
		PreRun:    validateArgs,
	}

	obscure = &cobra.Command{
		Use:       "obscure",
		Short:     "apply formatters on std in",
		Long:      "obscure applies formatters on stdin thus you can format any text, or post-format an output given by profaneword",
		Args:      cobra.OnlyValidArgs,
		ValidArgs: formatters,
		Run:       obscureFunc,
		PreRun:    validateArgs,
	}

	version = &cobra.Command{
		Use:   "version",
		Short: "print the version and exit",
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.Print("0.8.4")
		},
	}
)

func errUseEnd(cmd *cobra.Command, t string) {
	cmd.PrintErrln(t)
	cmd.Println()
	_ = cmd.Usage()
	os.Exit(1)
}

func validateArgs(cmd *cobra.Command, args []string) {
	for i, arg := range args {
		if formatter(arg) == random {
			if i == len(args)-1 {
				errUseEnd(cmd, `"random" cannot be used without a formatter`)
			}
			if formatter(args[i+1]) == randomly {
				errUseEnd(cmd, `"random" cannot appear before "randomly"`)
			}
		}
		if formatter(arg) == randomly {
			if i == len(args)-1 {
				errUseEnd(cmd, `"randomly" cannot be used without a formatter`)
			}
		}
	}
}

func profaneWords(cmd *cobra.Command, args []string) {
	numWords := numWordsFrom(cmd)
	delim := getDelimiter(cmd)
	var disallowW = disallowedWords(cmd)
	sentencer := profanities.NewProfanitySentencer(disallowW)
	sentence := sentencer.GetSentence(numWords)
	text := sentencer.Sentence(sentence)
	formatter := formatterOf(args, profaneword.RandomTitleFormatter(), profaneword.DelimiterFormatterWith(delim))
	cmd.Println(formatter.Format(text))
}

func disallowedWords(cmd *cobra.Command) (disallowed profanities.Word) {
	no, _ := cmd.PersistentFlags().GetString("no")
	for _, nope := range strings.Split(no, "|") {
		switch nope {
		case "MISSPELL":
			disallowed |= profanities.MISSPELL
		case "POSITIVE":
			disallowed |= profanities.POSITIVE
		case "":
		default:
			errUseEnd(cmd, "unknown disallowed word: "+nope)
		}
	}
	return
}

func getDelimiter(cmd *cobra.Command) (delim string) {
	delim, _ = cmd.PersistentFlags().GetString("delimiter")
	if delim == RAND {
		idx := profaneword.CryptoRand{}.RandMax(len(alternateDelimiters))
		delim = string(alternateDelimiters[idx])
	}
	return
}

func numWordsFrom(cmd *cobra.Command) int {
	pflags := cmd.PersistentFlags()
	ext, _ := pflags.GetInt16("extensiveness")
	if extend, _ := pflags.GetBool("extend"); extend {
		ext += 1
	}
	if extend, _ := pflags.GetBool("EXTEND"); extend {
		ext += 3
	}
	return int(ext)
}

func obscureFunc(cmd *cobra.Command, args []string) {
	reader := bufio.NewReader(os.Stdin)
	delim := getDelimiter(cmd.Root())
	formatter := formatterOf(args, profaneword.DelimiterFormatterWith(delim))
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if text != "" && text != "\n" {
					cmd.Println(formatter.Format(text))
				}
				break
			}
			cmd.Println(err)
			os.Exit(1)
		}
		text = text[:len(text)-1] // remove the newline to not print two newlines
		cmd.Println(formatter.Format(text))
	}
}

// Execute executes the root command.
func Execute() error {
	return profaneCmd.Execute()
}

func init() {
	profaneCmd.AddCommand(version)
	profaneCmd.AddCommand(obscure)

	profaneCmd.PersistentFlags().Int16P("extensiveness", "e", 2, "how long (number of words) the password should be. Default is 2")
	profaneCmd.PersistentFlags().Bool("extend", false, "lengthen the output (extensiveness+1)")
	profaneCmd.PersistentFlags().Bool("EXTEND", false, "lengthen the output further (extensiveness+3)")

	profaneCmd.PersistentFlags().StringP("delimiter", "d", " ", "a specific delimiter to use, or '"+RAND+"' for a randomly chosen one from: '"+alternateDelimiters+"'")

	profaneCmd.PersistentFlags().String("no", "", "exclude types of words: can be MISSPELL, POSITIVE or a '|' separated text of those")

	profaneCmd.SetUsageTemplate(usageTpl)
}
