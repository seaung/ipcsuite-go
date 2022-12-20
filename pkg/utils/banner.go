package utils

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var (
	author  string = "seaung"
	version string = "1.0.0"
)

func ShowBanner() {
	name := fmt.Sprintf("ipcsuite-go (v.%s)", version)
	banner := `
    _                       _ __                       
   (_)___  ____________  __(_) /____        ____ _____ 
  / / __ \/ ___/ ___/ / / / / __/ _ \______/ __ '/ __ \
 / / /_/ / /__(__  ) /_/ / / /_/  __/_____/ /_/ / /_/ /
/_/ .___/\___/____/\__,_/_/\__/\___/      \__, /\____/ 
 /_/                                     /____/        

	`
	all_lines := strings.Split(banner, "\n")
	w := len(all_lines[1])

	fmt.Println(banner)
	color.Yellow(fmt.Sprintf("%[1]*s", -w, fmt.Sprintf("%[1]*s", (w+len(name))/2, name)))
	color.Cyan(fmt.Sprintf("%[1]*s", -w, fmt.Sprintf("%[1]*s", (w+len(author))/2, author)))
	fmt.Println()
}
