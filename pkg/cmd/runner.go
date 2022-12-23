package cmd

import (
	"fmt"

	"github.com/abiosoft/ishell"
	"github.com/seaung/ipcsuite-go/pkg/utils"
)

var runnerCmd = &ishell.Cmd{
	Name:     "audit",
	Help:     "审计漏洞",
	LongHelp: "审计网络摄像头漏洞",
	Func: func(c *ishell.Context) {
		c.ShowPrompt(false)
		defer c.ShowPrompt(true)

		choice := c.MultiChoice([]string{
			"single",
			"multiple",
		}, "请您选择一个分类")

		istype := choiceCategorie(choice)

		utils.New().Info("请您提供一个目标:")
		target := c.ReadLine()

		utils.New().Info(fmt.Sprintf("您选择的分类是: %s 审计的目标是: %s", istype, target))

	},
}

func init() {
	rootCmd.AddCmd(runnerCmd)
}

func audit(istype int) {
	// process
}

func choiceCategorie(choice int) string {
	switch choice {
	default:
		return "Unknown"
	case 0:
		return "single"
	case 1:
		return "multiple"
	}
}
