package helper

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"github.com/irishgreencitrus/godot-build-cli-go/variables"
)

func CleanZips(v []string){
	for _,i := range v {
		err := os.Remove(i+".zip")
		if err != nil{
			fmt.Println(i+".zip not found. Can't be removed.")
		} else {
			fmt.Println("Removed "+i+".zip")
		}
		
	}
}
func MoveBuilt(ver string){
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
func MoveInitialiser(vers string){
	if vers == variables.ALL_SELECTOR {
		fmt.Println("Moving all builds.")
		for _,i := range variables.Versions {
			MoveBuilt(i)
		}
	} else {
		MoveBuilt(vers)
	}
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
// Arguably the most useful helper function in this file, allowing checking of 
// strings in a list of strings. Used a lot when validating versions, for example
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
// This function returns (NewType, NewTypeFlags) if the type is supported
// otherwise it just returns the old type and flag.
func TypeInitialiser(v string) (string,string){
	if !StringInSlice(v,variables.Types){
		
		fmt.Println("Type invalid. Keeping type at last")
		return variables.CurrentType, variables.CurrentTypeFlag
	} else {
		return v, GetFlagsFromType(v)
	}
}

// Returns flags for a type based on t.
// Not recommended calling this with user input as it'll deliberatly error out if incorrect.
func GetFlagsFromType(t string) string{
	switch t {
	case "editor":
		return variables.EDITOR_FLAGS
	case "export":
		return variables.EXPORT_FLAGS
	case "headless":
		return variables.HEADLESS_FLAGS
	case "server":
		return variables.SERVER_FLAGS
	default:
		log.Panicln("Somehow called with incorrect value",t)
		return ""
	}
}