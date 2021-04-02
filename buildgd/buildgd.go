package buildgd

import (
	"github.com/irishgreencitrus/godot-build-cli-go/variables"
	"github.com/irishgreencitrus/godot-build-cli-go/helper"
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)
func BuildInitialiser(version string){
	if !helper.StringInSlice(version, variables.Versions) && version != variables.ALL_SELECTOR {
		fmt.Println("Godot version not found or supported. To check versions type versions.")
	} else {

		godotdir := fmt.Sprintf("download/godot-%s", version)

		if version == variables.ALL_SELECTOR {
			fmt.Println("Building all versions")
			for i := range variables.Versions {
				if _, err := os.Stat("download/godot-" + variables.Versions[i]); os.IsNotExist(err) {
					fmt.Printf("Build directory not found for %s. Try downloading it using the download command\n", variables.Versions[i])
					continue
				}
				fmt.Println(variables.Versions[i])
				BuildGodot(variables.Versions[i])
			}
		} else if _, err := os.Stat(godotdir); os.IsNotExist(err) {
			fmt.Println("Build directory not found. Try downloading it using the download command")
		} else {
			fmt.Println("Building", version)
			BuildGodot(version)
		}
	}
}
func BuildGodot(ver string) {
	switch runtime.GOOS {
	//case "windows":
	//	buildWithFlags(ver, strings.Fields("-j"+fmt.Sprint(runtime.NumCPU())+" platform=windows"))
	case "linux":
		if runtime.GOARCH == "amd64" {
			BuildWithFlags(ver, strings.Fields(variables.CurrentTypeFlag+" -j"+fmt.Sprint(runtime.NumCPU())))
		} else if runtime.GOARCH == "arm64" {
			os.Setenv("CCFLAGS", "-mtune=cortex-a72 -mcpu=cortex-a72 -mfloat-abi=hard -mlittle-endian -munaligned-access -mfpu=neon-fp-armv8")
			BuildWithFlags(ver, strings.Fields(variables.CurrentTypeFlag+" use_llvm=yes -j"+fmt.Sprint(runtime.NumCPU())))
		} else if runtime.GOARCH == "arm"{
			os.Setenv("CCFLAGS", "-mtune=cortex-a72 -mcpu=cortex-a72 -mfloat-abi=hard -mlittle-endian -munaligned-access -mfpu=neon-fp-armv8")
			BuildWithFlags(ver, strings.Fields(variables.CurrentTypeFlag+" use_llvm=yes -j"+fmt.Sprint(runtime.NumCPU())))
		}
	}
}
func BuildWithFlags(vers string, flags []string) {

	//for x := range flags{
	//	fmt.Println(flags[x])
	//}
	cmd := exec.Command("scons", flags...)
	cmd.Dir = fmt.Sprintf("download/godot-%s", vers)
	var errb bytes.Buffer
	cmd.Stderr = &errb
	out, _ := cmd.StdoutPipe()
	cmd.Start()
	scanner := bufio.NewScanner(out)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	cmd.Wait()
	fmt.Println(errb.String())
}