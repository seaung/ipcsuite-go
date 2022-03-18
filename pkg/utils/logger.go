package utils

import (
	"fmt"

	"github.com/fatih/color"
)

type Logger struct{}

func New() *Logger {
	return &Logger{}
}

func (l *Logger) Warnning(message string) {
	highlight := color.New(color.FgYellow).SprintFunc()
	fmt.Println(highlight("[!]"), highlight(message))
}

func (l *Logger) Success(message string) {
	highlight := color.New(color.FgHiGreen).SprintFunc()
	fmt.Println(highlight("[+]"), highlight(message))
}

func (l *Logger) Info(message string) {
	highlight := color.New(color.FgHiBlue).SprintFunc()
	reset := color.New(color.FgWhite).SprintFunc()
	fmt.Println(highlight("[*]"), reset(message))
}

func (l *Logger) Errors(message string) {
	highlight := color.New(color.FgHiRed).SprintFunc()
	fmt.Println(highlight("[-]"), highlight(message))
}
