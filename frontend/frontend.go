package frontend

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/irishgreencitrus/godot-build-cli-go/buildgd"
	"github.com/irishgreencitrus/godot-build-cli-go/helper"
	"github.com/irishgreencitrus/godot-build-cli-go/variables"
	"github.com/irishgreencitrus/godot-build-cli-go/web"
)

func PrintLogo() {
	logo := [...]string{

		"       `#-#`     ##-`                                                                          ",
		"       #oooo+````/oooo/                                                                        ",
		"       -oooooooooooooo/                                                                        ",
		" #/##`-/oooooooooooooooo+-`##/#                                                                ",
		"#oooooooooooooooooooooooooooooo/`        `#####`   `#####`   ######`     `#####`  ########`    ",
		"/oooooooooooooooooooooooooooooo+`      `-##---#` `-##--###-  -##-###-` `-##--###- #--###--`    ",
		" #oooooyyyoooooooooooosyyooooo/        ####      -##-  `-### -##` -##- -##-  `-###  `##-       ",
		" #oooyNyoodyoooyyooosdsosNhooo#       `###` ---- -###   -### -##` #### -###   -##-  `##-       ",
		" #ooodm.../moooMMsooho...hmooo#        -##-##-#- `##-#`####` -####-### `##-#`####`  `##-       ",
		" #oooohdyyhooooNNooooyyyhdoooo#         `-#####-  `-#####-`  -#####-#   `-#####-`   `##-       ",
		" -ssooooooooooooooooooooooooss/            ```       ```      ````         ```       ```       ",
		" -hhhdmdoooosddddddyooooyNdhhh+        Unofficial Building CLI                                 ",
		" `ooooyNhhhhmdoooohNhhhhmdoooo-                                                                ",
		"  -+ooooossssoooooossssoooooo#                                                                 ",
		"    -/+oooooooooooooooooo+/-`                                                                  ",
		"       `##-###////###--#`                                                                      \n",
	}
	for line := range logo {
		fmt.Println(logo[line])
	}
	fmt.Printf("Detected OS: %s\n", runtime.GOOS)
	fmt.Printf("Detected Architecture: %s\n\n", runtime.GOARCH)

}


func InteractiveMode() {
	input := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("> ")
		input, _ := input.ReadString('\n')
		input = strings.Replace(input, "\n", "", -1)
		input = strings.Replace(input, "\r", "", -1)
		command := strings.Split(input, " ")
		commandword := command[0]

		//fmt.Println("\n")
		//fmt.Println(command)
		switch commandword {
		case "exit":
			fmt.Println("Exiting...")
			os.Exit(0)
		case "versions":
			fmt.Println("Versions available:")
			for i := range variables.Versions {
				fmt.Println(variables.Versions[i])
			}
		case "download":
			if len(command) != 2 {
				fmt.Println("Usage: download <version>")
				continue
			}
			chosenversion := command[1]
			web.DownloadInitialiser(chosenversion)
			
		case "build":
			if len(command) != 2 {
				fmt.Println("Usage: build <version>")
				continue
			}
			buildver := command[1]
			buildgd.BuildInitialiser(buildver)
		case "move_built":
			if len(command) != 2 {
				fmt.Println("Usage: move_built <version>")
				continue
			}
			helper.MoveBuilt(command[1])
		case "cleanzips":
			helper.CleanZips(variables.Versions)
		case "type":
			if len(command) != 2 {
				fmt.Println("Usage: type <type>")
				continue
			}
			variables.CurrentType, variables.CurrentTypeFlag = helper.TypeInitialiser(command[1])
		default:
			fmt.Println("Command list:")
			fmt.Println("exit")
			fmt.Println("build")
			fmt.Println("versions")
			fmt.Println("move_built")
			fmt.Println("cleanzips")
			fmt.Println("type")
		}
	}
}