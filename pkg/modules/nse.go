package modules

import (
	"context"
	"fmt"
	"time"

	"github.com/Ullaakut/nmap"
)

func RunNseScripts(categories string, ports string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	scanner, err := nmap.NewScanner(
		nmap.WithScripts(categories),
		nmap.WithPorts(ports),
		nmap.WithContext(ctx),
	)

	if err != nil {
		fmt.Println(err)
	}

	result, _, err := scanner.Run()
	if err != nil {
		fmt.Println(err)
	}

	for _, r := range result.Hosts {
		if len(r.Ports) == 0 || len(r.Addresses) == 0 {
			continue
		}

		fmt.Printf("The Target Host : %s\n", r.Addresses[0])

		for _, port := range r.Ports {
			fmt.Printf("Protocol : %s\n", port.Protocol)
		}
	}

	for _, r := range result.PostScripts {
		fmt.Printf("ID : %s - Output - %s\n", r.ID, r.Output)

		for _, t := range r.Tables {
			fmt.Printf("Key : %s - Tables : %s - Element : %s\n", t.Key, t.Tables[0], t.Elements[0])
		}
	}
}
