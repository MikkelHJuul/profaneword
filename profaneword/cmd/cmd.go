package cmd

import (
	"bufio"
	"github.com/spf13/cobra"
	"io"
	"os"
)

// usageTpl is a copy-paste of the underlying UsageTemplate in order to inject Args (and [..args])
const usageTpl = `Usage:{{if .Runnable}}
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

Args:
	1337  output formatted as 1337-speak
	/s    sARcaSTiC OUtpUt


Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`


var (
	profaneCmd = &cobra.Command{
		Use: "profaneword",
		Short: "A generator for profane passwords as requested by u/gatestone",
		Long:
`profaneword is a CLI library and tool for generating obscene/profane passwords. 
It's probably not particularly safe to use, as these passwords will be easy to brute force; 
if an attacker knows you use this generator. But hey, it's just for fun.`,
		Args: cobra.OnlyValidArgs,
		ValidArgs: formatters,
		Run: profaneWords,
	}

	obscure = &cobra.Command{
		Use: "obscure",
		Short: "apply formatters on std in",
		Long: "obscure applies formatters on std in thus you can format any text, or format an output given by profaneword",
		Args: cobra.OnlyValidArgs,
		ValidArgs: formatters,
		Run: obscureFunc,
	}

	version = &cobra.Command{
		Use:   "version",
		Short: "print the version and exit",
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.Print("0.1.1")
		},
	}
)

func profaneWords(cmd *cobra.Command, args []string) {
	numWords := numWordsFrom(cmd)
	text := GetWords(numWords)
	formatter := formatterOf(args)
	cmd.Println(formatter.Format(text))
}


func numWordsFrom(cmd *cobra.Command) int { //errs?
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
	formatter := formatterOf(args)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			cmd.Println(err)
			os.Exit(1)
		}
		cmd.Print(formatter.Format(text))
	}
}


// Execute executes the root command.
func Execute() error {
	return profaneCmd.Execute()
}

func init() {
	profaneCmd.AddCommand(version)
	profaneCmd.AddCommand(obscure)

	profaneCmd.PersistentFlags().Int16P("extensiveness","e",2, "how long (number of words) the password should be")
	profaneCmd.PersistentFlags().Bool("extend",false, "lengthen the output (extensiveness+1)")
	profaneCmd.PersistentFlags().Bool("EXTEND",false, "lengthen the output further (extensiveness+3)")

	profaneCmd.SetUsageTemplate(usageTpl)
}

