package cmd

import (
	"fmt"

	"github.com/abiosoft/ishell"
	"github.com/seaung/ipcsuite-go/internal/ipcs"
	"github.com/seaung/ipcsuite-go/pkg/utils"
)

var runnerCmd = &ishell.Cmd{
	Name:     "audit",
	Help:     "审计漏洞",
	LongHelp: "审计网络摄像头漏洞",
}

var singleCmd = &ishell.Cmd{
	Name:     "sin",
	Help:     "单个poc进行检测",
	LongHelp: "使用单个poc进行检测",
	Func: func(c *ishell.Context) {
		c.ShowPrompt(false)
		defer c.ShowPrompt(true)

		utils.New().Info("请您输入单个poc的路径及文件名称")

		poc := c.ReadLine()

		utils.New().Warnning(fmt.Sprintf("您选择的poc是 : %s", poc))

		utils.New().Info("请您输入一个目标链接: e.g. https://localhost:8000/")

		target := c.ReadLine()

		utils.New().Warnning(fmt.Sprintf("您需要检测的目标链接为: %s", target))

		ipcs.RunSiglePoc(poc, target)
	},
}

var multipleCmd = &ishell.Cmd{
	Name:     "multi",
	Help:     "多个poc进行检测",
	LongHelp: "使用多个poc进行检测",
	Func: func(c *ishell.Context) {
		c.ShowPrompt(false)
		defer c.ShowPrompt(true)

		utils.New().Info("请您输入poc的路径")

		poc := c.ReadLine()

		utils.New().Warnning(fmt.Sprintf("您提供的poc路径是 : %s", poc))

		utils.New().Info("请您输入一个目标链接: e.g. https://localhost:8000/")

		target := c.ReadLine()

		utils.New().Warnning(fmt.Sprintf("您需要检测的目标链接为: %s", target))

		ipcs.RunMultiplePocs(poc, target)
	},
}

func init() {
	rootCmd.AddCmd(runnerCmd)
	runnerCmd.AddCmd(singleCmd)
	runnerCmd.AddCmd(multipleCmd)
}
