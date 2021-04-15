package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/irishgreencitrus/godot-build-cli-go/buildgd"
	"github.com/irishgreencitrus/godot-build-cli-go/frontend"
	"github.com/irishgreencitrus/godot-build-cli-go/helper"
	"github.com/irishgreencitrus/godot-build-cli-go/variables"
	"github.com/irishgreencitrus/godot-build-cli-go/web"
)

func main() {
	args := os.Args[1:]

	downloadFlag := flag.Int("download", 0, "Downloads specified version")
	moveFlag := flag.Int("move", 0, "Moves specified builds to an easier to access location")
	buildFlag := flag.Int("build", 0, "Builds specified version")
	typeFlag := flag.Int("type", 1, "Chooses the type to build")
	shouldPrintVersion := flag.Bool("V", false, "Prints available versions")
	shouldPrintPlatform := flag.Bool("P", false, "Prints available platforms")
	shouldPrintTypes := flag.Bool("T", false, "Prints possible types")
	shouldRemoveZips := flag.Bool("Z", false, "Removes version zip files")
	shouldRenameFriendly := flag.Bool("R", false, "Renames builds to more readable names")
	flag.Parse()
	if *shouldPrintVersion {
		fmt.Println("Available Versions:")
		fmt.Println(strings.Join(variables.Versions, "\n"))
	}
	if *shouldPrintPlatform {
		fmt.Println("Available Platforms")
		fmt.Println(strings.Join(variables.Platforms, "\n"))
	}
	if *shouldPrintTypes {
		fmt.Println("Available types")
		fmt.Println(strings.Join(variables.Types, "\n"))
	}
	if *downloadFlag != 0 {
		c := helper.ListWithBitFilter(variables.Versions, byte(*downloadFlag))
		fmt.Println("Downloading versions:", c)
		web.DownloadInitialiser(c)
	}

	if *buildFlag != 0 {
		t := helper.ListWithBitFilter(variables.Types, byte(*typeFlag))
		fmt.Println("Types specified:", t)
		c := helper.ListWithBitFilter(variables.Versions, byte(*buildFlag))
		fmt.Println("Build Versions Specified:", c)
		buildgd.BuildInitialiser(t, c)
	}
	if *shouldRemoveZips {
		helper.CleanZips(variables.Versions)
	}
	if *moveFlag != 0 {
		c := helper.ListWithBitFilter(variables.Versions, byte(*moveFlag))
		helper.MoveInitialiser(c)
	}
	if *shouldRenameFriendly {
		helper.RenameBuilt()
	}
	if len(args) == 0 {
		frontend.PrintLogo()
		frontend.SurveyMode()
	}
}
