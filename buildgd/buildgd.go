package buildgd

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/irishgreencitrus/godot-build-cli-go/helper"
	"github.com/irishgreencitrus/godot-build-cli-go/variables"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Initialises building for every type of every version.
func BuildInitialiser(types []string, version []string) {
	for _, t := range types {
		for _, v := range version {
			if _, err := os.Stat("download/godot-" + v); os.IsNotExist(err) {
				fmt.Printf("Build directory not found for %s. Try downloading it using the download command. Tried to build type: %s\n", v, t)
				continue
			}
			BuildGodot(v, helper.GetFlagsFromType(t))
		}
	}
}

// A better entry point than BuildInitialiser if you are building a custom tool.
// Currently the actual building only supports linux/amd64, linux/arm, and linux/arm64
// May expand to other platforms in the near future.
// TIP: If you want to add another platform in a PR, add the method in here!
func BuildGodot(ver string, typ string) {
	switch runtime.GOOS {
	//case "windows":
	//	buildWithFlags(ver, strings.Fields("-j"+fmt.Sprint(runtime.NumCPU())+" platform=windows"))
	case "linux":
		if runtime.GOARCH == "amd64" {
			BuildWithFlags(ver, strings.Fields(typ+" -j"+fmt.Sprint(runtime.NumCPU())))
		} else if runtime.GOARCH == "arm64" {
			os.Setenv("CCFLAGS", "-mtune=cortex-a72 -mcpu=cortex-a72 -mfloat-abi=hard -mlittle-endian -munaligned-access -mfpu=neon-fp-armv8")
			BuildWithFlags(ver, strings.Fields("bits="+variables.ToolAnswers.Bits+" "+typ+" use_llvm=yes -j"+fmt.Sprint(runtime.NumCPU())))
		} else if runtime.GOARCH == "arm" {
			os.Setenv("CCFLAGS", "-mtune=cortex-a72 -mcpu=cortex-a72 -mfloat-abi=hard -mlittle-endian -munaligned-access -mfpu=neon-fp-armv8")
			BuildWithFlags(ver, strings.Fields(typ+" use_llvm=yes -j"+fmt.Sprint(runtime.NumCPU())))
		}
	}
}

// Probably the best entrypoint if you are using this as a module.
// Directly builds the versions using scons, following the output.
func BuildWithFlags(vers string, flags []string) {
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

// Experimental replacement for BuildInitialiser
