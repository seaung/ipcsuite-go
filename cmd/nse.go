package main

import (
	"github.com/abiosoft/ishell"
	"github.com/seaung/ipcsuite-go/pkg/utils"
)

var nseCmd = &ishell.Cmd{
	Name:     "nse",
	Help:     "运行NSE脚本",
	LongHelp: "根据用户提供的NSE脚本分类运行特定的脚本",
	Func: func(c *ishell.Context) {
		c.ShowPrompt(false)
		defer c.ShowPrompt(true)

		categoie := c.MultiChoice([]string{
			"auth",
			"broadcast",
			"brute",
			"default",
			"discovery",
			"dos",
			"exploit",
			"external",
			"fuzzer",
			"intrusive",
			"malware",
			"safe",
			"version",
			"vuln",
		}, "请您选择一个分类: ")

		istype := iscategorie(categoie)

		utils.New().Info("请您提供一个目标: ")
		hosts := c.ReadLine()
		utils.New().Info("请您提供一个端口或端口的范围: ")
		ports := c.ReadLine()

		run(istype, hosts, ports)
	},
}

func run(categoies, hosts, ports string) {
}

func iscategorie(categoies int) string {
	switch categoies {
	default:
		return ""
	case 0:
		return ""
	case 1:
		return ""
	case 2:
		return ""
	case 3:
		return ""
	}
}
