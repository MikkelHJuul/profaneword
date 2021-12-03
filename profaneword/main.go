package main

import (
	"fmt"
	"github.com/MikkelHJuul/profaneword/profaneword/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
