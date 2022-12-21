package cmd

import (
	"github.com/abiosoft/ishell"
	"github.com/seaung/ipcsuite-go/internal/nmap"
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

		istype := nmap.Iscategorie(categoie)

		utils.New().Info("请您提供一个目标: ")
		hosts := c.ReadLine()
		utils.New().Info("请您提供一个端口或端口的范围: ")
		ports := c.ReadLine()

		nmap.Run(istype, hosts, ports)
	},
}

var customeNseCmd = &ishell.Cmd{
	Name:     "custome-nse",
	Help:     "自定义个NSE脚本分类",
	LongHelp: "这个命令提供了",
	Func: func(c *ishell.Context) {
		c.ShowPrompt(false)
		defer c.ShowPrompt(true)

		utils.New().Info("请您输入一个脚本的分类或脚本分类的组合:")
		categoies := c.ReadLine()

		utils.New().Info("请您提供一个目标: ")
		hosts := c.ReadLine()
		utils.New().Info("请您提供一个端口或端口的范围")
		ports := c.ReadLine()

		nmap.RunCutomeScript(categoies, hosts, ports)
	},
}

func init() {
	rootCmd.AddCmd(nseCmd)
	rootCmd.AddCmd(customeNseCmd)
}
