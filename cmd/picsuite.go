package main

import (
	"github.com/seaung/ipcsuite-go/pkg/cmd"
	"github.com/seaung/ipcsuite-go/pkg/utils"
)

func main() {
	utils.CheckSudo()
	utils.ShowBanner()

	cmd.Executor()
}
