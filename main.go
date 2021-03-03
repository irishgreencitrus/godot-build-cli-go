package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var versions = []string{
	"3.2.3-stable",
	"3.2.2-stable",
	"3.2.1-stable",
	"3.2-stable",
	"3.1.2-stable",
	"3.1.1-stable",
	"3.1-stable",
	"2.1.6-stable",
}

func main() {
	args := os.Args[1:]
	printLogo()
	if len(args) == 0 {
		fmt.Println("No arguments specified. Continuing in interactive mode.")
		interactiveMode()
	} else if stringInSlice("--download-only", args) {
		fmt.Println("Download Only Specified. Only downloading source and not building them.")
	}

}

func printLogo() {
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

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
func interactiveMode() {
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
			for i := range versions {
				fmt.Println(versions[i])
			}
		case "download":
			if len(command) != 2 {
				fmt.Println("Usage: download <version>")
				continue
			}
			chosenversion := command[1]
			if chosenversion == "*" {
				fmt.Println("Downloading every version.")
				var wg sync.WaitGroup
				for ver := range versions {
					wg.Add(1)
					go func(v int) {
						downloadVersion(versions[v])
						wg.Done()
					}(ver)

				}
				wg.Wait()
				fmt.Println("All versions downloaded and extracted to /downloads directory")
			} else if !stringInSlice(chosenversion, versions) {
				fmt.Println("Godot version not found or supported. To check versions type versions.")
				continue
			} else {
				downloadVersion(chosenversion)
				//fmt.Println("Unzipped:\n" + strings.Join(files, "\n"))
			}
		case "build":
			if len(command) != 2 {
				fmt.Println("Usage: build <version>")
				continue
			}
			buildver := command[1]
			if !stringInSlice(buildver, versions) && buildver != "*" {
				fmt.Println("Godot version not found or supported. To check versions type versions.")
				continue
			} else {

				godotdir := fmt.Sprintf("download/godot-%s", buildver)

				if buildver == "*" {
					fmt.Println("Building all versions")
					for i := range versions {
						if _, err := os.Stat("download/godot-" + versions[i]); os.IsNotExist(err) {
							fmt.Printf("Build directory not found for %s. Try downloading it using the download command\n", versions[i])
							continue
						}
						fmt.Println(versions[i])
						buildGodot(versions[i])
					}
				} else if _, err := os.Stat(godotdir); os.IsNotExist(err) {
					fmt.Println("Build directory not found. Try downloading it using the download command")
				} else {
					fmt.Println("Building", buildver)
					buildGodot(buildver)
				}
			}
		default:
			fmt.Println("Command list:")
			fmt.Println("exit")
			fmt.Println("build")
			fmt.Println("versions")
		}
	}
}

// This function is stolen from
// https://golangcode.com/download-a-file-from-a-url/
// It writes to the file as it downloads it to prevent it from loading the entire file into memory
func downloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// This function is stolen from
// https://golangcode.com/unzip-files-in-go/
// Unzips the file keeping directory structure
func Unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}

	return filenames, nil
}

func downloadVersion(chosenversion string) {
	url := fmt.Sprintf("https://github.com/godotengine/godot/archive/%s.zip", chosenversion)
	path := fmt.Sprintf("%s.zip", chosenversion)
	fmt.Println("Downloading", path)
	downloadFile(path, url)
	fmt.Println("Unzipping", path)
	Unzip(path, "download")
	fmt.Println("Successfully Unzipped", path)
}

func buildGodot(ver string) {
	switch runtime.GOOS {
	case "windows":
		buildWithFlags(ver, strings.Fields("-j8 platform=windows tools=no disable_3d=yes"))
	case "linux":
		if runtime.GOARCH == "amd64" {
			buildWithFlags(ver, strings.Fields("-j8 platform=x11"))
		} else if runtime.GOARCH == "arm64" {
			os.Setenv("CCFLAGS", "-mtune=cortex-a72 -mcpu=cortex-a72 -mfloat-abi=hard -mlittle-endian -munaligned-access -mfpu=neon-fp-armv8")
			buildWithFlags(ver, strings.Fields("platform=server target=release tools=no use_llvm=yes -j4"))
		}
	}
}

func buildWithFlags(vers string, flags []string) {

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
