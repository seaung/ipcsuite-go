package modules

import (
	"context"
	"fmt"
	"time"

	"github.com/Ullaakut/nmap"
)

func RunNse(categories string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	scanner, err := nmap.NewScanner(
		nmap.WithScripts(categories),
		nmap.WithContext(ctx),
	)

	if err != nil {
		fmt.Println(err)
	}

	result, _, err := scanner.Run()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(result)
}
