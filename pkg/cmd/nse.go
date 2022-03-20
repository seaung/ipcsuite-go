package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/Ullaakut/nmap/v2"
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

		runCutomeScript(categoies, hosts, ports)
	},
}

func init() {
	rootCmd.AddCmd(nseCmd)
	rootCmd.AddCmd(customeNseCmd)
}

func runCutomeScript(categoies, hosts, ports string) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)

	defer cancel()

	scanner, err := nmap.NewScanner(
		nmap.WithTargets(hosts),
		nmap.WithPorts(ports),
		nmap.WithContext(ctx),
		nmap.WithScripts(categoies),
	)

	if err != nil {
		utils.New().Errors(fmt.Sprintf("Unable to create nmap scanner : %v\n", err))
	}

	result, _, err := scanner.Run()

	if err != nil {
		utils.New().Errors(fmt.Sprintf("Unable to run nmap scan : %v\n", err))
	}

	for _, host := range result.Hosts {
		if len(host.Ports) == 0 || len(host.Addresses) == 0 {
			continue
		}
		utils.New().Warnning(fmt.Sprintf("Host %q: \n", host.Addresses[0]))

		for _, script := range host.HostScripts {
			utils.New().Warnning(fmt.Sprintf("Script ID      : %s\n", script.ID))
			utils.New().Warnning(fmt.Sprintf("Script Element : %s\n", script.Elements[0].Value))
			utils.New().Warnning(fmt.Sprintf("Script Table   : %s\n", script.Tables[0].Key))
			utils.New().Warnning(fmt.Sprintf("Script Output  : %s\n", script.Output))
		}
	}
	utils.New().Success(fmt.Sprintf("Nmap done: %d hosts up scanned in %3f seconds\n", len(result.Hosts), result.Stats.Finished.Elapsed))
}

func run(categoies, hosts, ports string) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)

	defer cancel()

	scanner, err := nmap.NewScanner(
		nmap.WithTargets(hosts),
		nmap.WithPorts(ports),
		nmap.WithContext(ctx),
		nmap.WithScripts(categoies),
	)

	if err != nil {
		utils.New().Errors(fmt.Sprintf("Unable to create nmap scanner : %v\n", err))
	}

	result, _, err := scanner.Run()

	if err != nil {
		utils.New().Errors(fmt.Sprintf("Unable to run nmap scan : %v\n", err))
	}

	for _, host := range result.Hosts {
		if len(host.Ports) == 0 || len(host.Addresses) == 0 {
			continue
		}
		utils.New().Warnning(fmt.Sprintf("Host %q: \n", host.Addresses[0]))

		for _, script := range host.HostScripts {
			utils.New().Warnning(fmt.Sprintf("Script ID      : %s\n", script.ID))
			utils.New().Warnning(fmt.Sprintf("Script Element : %s\n", script.Elements[0].Value))
			utils.New().Warnning(fmt.Sprintf("Script Table   : %s\n", script.Tables[0].Key))
			utils.New().Warnning(fmt.Sprintf("Script Output  : %s\n", script.Output))
		}
	}
	utils.New().Success(fmt.Sprintf("Nmap done: %d hosts up scanned in %3f seconds\n", len(result.Hosts), result.Stats.Finished.Elapsed))
}

func iscategorie(categoies int) string {
	switch categoies {
	default:
		return "Unknown"
	case 0:
		return "auth"
	case 1:
		return "broadcast"
	case 2:
		return "brute"
	case 3:
		return "default"
	case 4:
		return "discovery"
	case 5:
		return "dos"
	case 6:
		return "exploit"
	case 7:
		return "external"
	case 8:
		return "fuzzer"
	case 9:
		return "intrusive"
	case 10:
		return "malware"
	case 11:
		return "safe"
	case 12:
		return "version"
	case 13:
		return "vuln"
	}
}
