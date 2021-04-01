package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

)
const ALL_SELECTOR string = "all"
var versions = []string{
	"3.2.3-stable",
	"3.2.2-stable",
	"3.2.1-stable",
	"3.2-stable",
	"3.1.2-stable",
	"3.1.1-stable",
	"3.1-stable",
}
var types = []string{
	"editor",
	"export",
	"headless",
	"server",
}

var typ = types[0]
var typFlag = ""
var platforms = []string{
	"linux/amd64",
	"linux/arm",
	"linux/arm64",
}

const (
	EDITOR_FLAGS = "platform=x11 target=release_debug tools=yes"
	EXPORT_FLAGS = "platform=x11 target=release tools=no"
	HEADLESS_FLAGS = "platform=server target=release_debug tools=yes"
	SERVER_FLAGS = "platform=server target=release tools=no"
)
func main() {
	args := os.Args[1:]
	typeInitialiser(typ)
	//runtimeName := os.Args[0]

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
		fmt.Println(strings.Join(versions,"\n"))
	}
	if *shouldPrintPlatform {
		fmt.Println("Available Platforms")
		fmt.Println(strings.Join(platforms,"\n"))
	}
	if *shouldPrintTypes {
		fmt.Println("Available types")
		fmt.Println(strings.Join(types,"\n"))
	}
	if *downloadFlag != "" {
		fmt.Println("Download Version Specified:",*downloadFlag)
		downloadInitialiser(*downloadFlag)
	}
	if *typeFlag != ""{
		fmt.Println("Type specified:",*typeFlag)
		typeInitialiser(*typeFlag)
	}

	if *buildFlag != ""{
		fmt.Println("Build Version Specified:", *buildFlag)
		buildInitialiser(*buildFlag)
	}
	if *shouldRemoveZips {
		cleanZips()
	}
	if *moveFlag != ""{
		moveInitialiser(*moveFlag)
	}
	if len(args) == 0 {
		printLogo()
		interactiveMode()
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
			downloadInitialiser(chosenversion)
			
		case "build":
			if len(command) != 2 {
				fmt.Println("Usage: build <version>")
				continue
			}
			buildver := command[1]
			buildInitialiser(buildver)
		case "move_built":
			if len(command) != 2 {
				fmt.Println("Usage: move_built <version>")
				continue
			}
			moveBuilt(command[1])
		case "cleanzips":
			cleanZips()
		case "type":
			if len(command) != 2 {
				fmt.Println("Usage: type <type>")
				continue
			}
			typeInitialiser(command[1])
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
func downloadInitialiser(version string){
	if version == ALL_SELECTOR {
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
	} else if !stringInSlice(version, versions) {
		fmt.Println("Godot version not found or supported. To check versions type versions.")
	} else {
		downloadVersion(version)
		//fmt.Println("Unzipped:\n" + strings.Join(files, "\n"))
	}
}
func typeInitialiser(v string){
	if !stringInSlice(v,types){
		fmt.Println("Type invalid. Keeping type at",typ)
	} else {
		typ = v
		typFlag = getFlagsFromType(v)
	}
}
func getFlagsFromType(t string) string{
	switch t {
	case "editor":
		return EDITOR_FLAGS
	case "export":
		return EXPORT_FLAGS
	case "headless":
		return HEADLESS_FLAGS
	case "server":
		return SERVER_FLAGS
	default:
		log.Panicln("Somehow called with incorrect value",t)
		return ""
	}
}
func buildInitialiser(version string){
	if !stringInSlice(version, versions) && version != ALL_SELECTOR {
		fmt.Println("Godot version not found or supported. To check versions type versions.")
	} else {

		godotdir := fmt.Sprintf("download/godot-%s", version)

		if version == ALL_SELECTOR {
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
			fmt.Println("Building", version)
			buildGodot(version)
		}
	}
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
		buildWithFlags(ver, strings.Fields("-j"+fmt.Sprint(runtime.NumCPU())+" platform=windows"))
	case "linux":
		if runtime.GOARCH == "amd64" {
			buildWithFlags(ver, strings.Fields(typFlag+" -j"+fmt.Sprint(runtime.NumCPU())))
		} else if runtime.GOARCH == "arm64" {
			os.Setenv("CCFLAGS", "-mtune=cortex-a72 -mcpu=cortex-a72 -mfloat-abi=hard -mlittle-endian -munaligned-access -mfpu=neon-fp-armv8")
			buildWithFlags(ver, strings.Fields(typFlag+" use_llvm=yes -j"+fmt.Sprint(runtime.NumCPU())))
		} else if runtime.GOARCH == "arm"{
			os.Setenv("CCFLAGS", "-mtune=cortex-a72 -mcpu=cortex-a72 -mfloat-abi=hard -mlittle-endian -munaligned-access -mfpu=neon-fp-armv8")
			buildWithFlags(ver, strings.Fields(typFlag+" use_llvm=yes -j"+fmt.Sprint(runtime.NumCPU())))
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
func moveBuilt(ver string){
	files,err := ioutil.ReadDir("download/godot-"+ver+"/bin")
	if errors.Is(err, fs.ErrNotExist){
		_,er := ioutil.ReadDir("download/godot-"+ver)
		if errors.Is(er,fs.ErrNotExist){
			fmt.Println("This version doesn't exist in downloads")
		} else {
			fmt.Println("This version hasn't been built yet. Build it using build <version>.")
		}
		
	}
	for _,f := range files{
		fmt.Println("Moving: " + f.Name())
		if _,er := os.Stat("builds"); os.IsNotExist(er){
			os.Mkdir("builds",0755)
		}
		err := os.Rename("download/godot-"+ver+"/bin/"+f.Name(),"builds/"+ver+"."+f.Name())
		if err != nil {
			fmt.Println(err)
		}
	}	
}
func moveInitialiser(vers string){
	if vers == ALL_SELECTOR {
		fmt.Println("Moving all builds.")
		for _,i := range versions {
			moveBuilt(i)
		}
	} else {
		moveBuilt(vers)
	}
}
func cleanZips(){
	for _,i := range versions {
		err := os.Remove(i+".zip")
		if err != nil{
			fmt.Println(i+".zip not found. Can't be removed.")
		} else {
			fmt.Println("Removed "+i+".zip")
		}
		
	}
}
