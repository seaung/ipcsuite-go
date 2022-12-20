package cmd

import (
	"github.com/abiosoft/ishell"
)

var rootCmd = ishell.New()

func Executor() {
	rootCmd.SetPrompt("ipcsuite > ")
	rootCmd.Run()
}
