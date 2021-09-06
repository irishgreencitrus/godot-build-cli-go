package frontend

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/irishgreencitrus/godot-build-cli-go/v2/buildgd"
	"github.com/irishgreencitrus/godot-build-cli-go/v2/helper"
	"github.com/irishgreencitrus/godot-build-cli-go/v2/variables"
	"github.com/irishgreencitrus/godot-build-cli-go/v2/web"
	"runtime"
)

// Prints the Godot Logo as ascii art
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

// New version of the old interactive mode.
// About 100x more intuative than the old one.
func SurveyMode() {
	survey.Ask(variables.ToolQuestions, &variables.ToolAnswers,
		survey.WithIcons(
			func(is *survey.IconSet) {

				is.Error.Text = "✗"
				is.Error.Format = "red+hb"

				is.Question.Text = "⁇"
				is.Question.Format = "white"

				is.SelectFocus.Text = "»"
				is.SelectFocus.Format = "white+hb"

				is.MarkedOption.Text = "✓"
				is.MarkedOption.Format = "green"

				is.UnmarkedOption.Text = "✗"
				is.UnmarkedOption.Format = "red"
			},
		),
	)
	//fmt.Println(variables.ToolAnswers)
	//fmt.Println(variables.ToolQuestions)

	web.DownloadInitialiser(variables.ToolAnswers.DownloadVer)
	if variables.ToolAnswers.RemoveZips {
		helper.CleanZips(variables.ToolAnswers.DownloadVer)
	}
	buildgd.BuildInitialiser(variables.ToolAnswers.BinaryTypes, variables.ToolAnswers.BuildVer)
	if variables.ToolAnswers.MoveBuilt {
		helper.MoveInitialiser(variables.Versions)
	}
	if variables.ToolAnswers.RenameFriendly {
		helper.RenameBuilt()
	}
}
