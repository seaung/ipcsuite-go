package nmap

import (
	"context"
	"fmt"
	"time"

	"github.com/Ullaakut/nmap/v2"
	"github.com/seaung/ipcsuite-go/pkg/utils"
)

func RunCutomeScript(categoies, hosts, ports string) {

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

func Run(categoies, hosts, ports string) {

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

func Iscategorie(categoies int) string {
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
