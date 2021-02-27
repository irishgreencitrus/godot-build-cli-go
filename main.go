package main

import (
	"fmt"
	"os"
	"runtime"
	"bufio"
	"strings"
	"net/http"
	"io"
	"archive/zip"
	"path/filepath"
	"sync"
	//"os/exec"
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
	fmt.Printf("Detected OS: %s\n",runtime.GOOS)
	fmt.Printf("Detected Architecture: %s\n\n",runtime.GOARCH)

}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
func interactiveMode(){
	input := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("> ")
		input, _ := input.ReadString('\n')
		input = strings.Replace(input, "\n", "", -1)
		input = strings.Replace(input, "\r", "", -1)
		command := strings.Split(input," ")
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
				if len(command) != 2{
					fmt.Println("Usage: download <version>")
					continue
				}
				chosenversion := command[1]
				if chosenversion == "*" {
					fmt.Println("Downloading every version.")
					var wg sync.WaitGroup
					for ver := range versions {
						wg.Add(1)
						go downloadVersionRoutine(versions[ver],&wg)
					}
					wg.Wait()
					fmt.Println("All versions downloaded and extracted to /downloads directory")
				} else if !stringInSlice(chosenversion,versions) {
					fmt.Println("Godot version not found or supported. To check versions type versions.")
					continue
				} else {
					downloadVersion(chosenversion)
					//fmt.Println("Unzipped:\n" + strings.Join(files, "\n"))
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

func downloadVersion(chosenversion string){
	url := fmt.Sprintf("https://github.com/godotengine/godot/archive/%s.zip",chosenversion)
	path := fmt.Sprintf("%s.zip",chosenversion)
	fmt.Println("Downloading",path)
	downloadFile(path, url)
	fmt.Println("Unzipping",path)
	Unzip(path,"download")
	fmt.Println("Successfully Unzipped",path)
}
func downloadVersionRoutine(chosenversion string, wg *sync.WaitGroup){
	//defer 
	url := fmt.Sprintf("https://github.com/godotengine/godot/archive/%s.zip",chosenversion)
	path := fmt.Sprintf("%s.zip",chosenversion)
	fmt.Println("Downloading",path)
	downloadFile(path, url)
	fmt.Println("Unzipping",path)
	Unzip(path,"download")
	fmt.Println("Successfully Unzipped",path)
	wg.Done()
}