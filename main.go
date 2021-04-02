package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/irishgreencitrus/godot-build-cli-go/buildgd"
	"github.com/irishgreencitrus/godot-build-cli-go/helper"
	"github.com/irishgreencitrus/godot-build-cli-go/variables"
	"github.com/irishgreencitrus/godot-build-cli-go/web"
	"github.com/irishgreencitrus/godot-build-cli-go/frontend"
)

func main() {
	args := os.Args[1:]
	helper.TypeInitialiser(variables.CurrentType)

	downloadFlag := flag.String("download","","Downloads specified version")
	moveFlag := flag.String("move","","Moves specified builds to an easier to access location")
	buildFlag := flag.String("build","","Builds specified version")
	typeFlag := flag.String("type","","Chooses the type to build")
	shouldPrintVersion := flag.Bool("V",false,"Prints available versions")
	shouldPrintPlatform := flag.Bool("P",false,"Prints available platforms")
	shouldPrintTypes := flag.Bool("T",false,"Prints possible types")
	shouldRemoveZips := flag.Bool("Z",false,"Removes version zip files")
	
	flag.Parse()
	if *shouldPrintVersion {
		fmt.Println("Available Versions:")
		fmt.Println(strings.Join(variables.Versions,"\n"))
	}
	if *shouldPrintPlatform {
		fmt.Println("Available Platforms")
		fmt.Println(strings.Join(variables.Platforms,"\n"))
	}
	if *shouldPrintTypes {
		fmt.Println("Available types")
		fmt.Println(strings.Join(variables.Types,"\n"))
	}
	if *downloadFlag != "" {
		fmt.Println("Download Version Specified:",*downloadFlag)
		web.DownloadInitialiser(*downloadFlag)
	}
	if *typeFlag != ""{
		fmt.Println("Type specified:",*typeFlag)
		helper.TypeInitialiser(*typeFlag)
	}

	if *buildFlag != ""{
		fmt.Println("Build Version Specified:", *buildFlag)
		buildgd.BuildInitialiser(*buildFlag)
	}
	if *shouldRemoveZips {
		helper.CleanZips(variables.Versions)
	}
	if *moveFlag != ""{
		helper.MoveInitialiser(*moveFlag)
	}
	if len(args) == 0 {
		frontend.PrintLogo()
		frontend.InteractiveMode()
	}
}







