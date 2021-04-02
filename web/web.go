package web

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/irishgreencitrus/godot-build-cli-go/helper"
	"github.com/irishgreencitrus/godot-build-cli-go/variables"
)

func DownloadVersion(chosenversion string) {
	url := fmt.Sprintf("https://github.com/godotengine/godot/archive/%s.zip", chosenversion)
	path := fmt.Sprintf("%s.zip", chosenversion)
	fmt.Println("Downloading", path)
	DownloadFile(path, url)
	fmt.Println("Unzipping", path)
	helper.Unzip(path, "download")
	fmt.Println("Successfully Unzipped", path)
}
func DownloadInitialiser(version string){
	if version == variables.ALL_SELECTOR {
		fmt.Println("Downloading every version.")
		var wg sync.WaitGroup
		for ver := range variables.Versions {
			wg.Add(1)
			go func(v int) {
				DownloadVersion(variables.Versions[v])
				wg.Done()
			}(ver)

		}
		wg.Wait()
		fmt.Println("All versions downloaded and extracted to /downloads directory")
	} else if !helper.StringInSlice(version, variables.Versions) {
		fmt.Println("Godot version not found or supported. To check versions type versions.")
	} else {
		DownloadVersion(version)
		//fmt.Println("Unzipped:\n" + strings.Join(files, "\n"))
	}
}
// This function is stolen from
// https://golangcode.com/download-a-file-from-a-url/
// It writes to the file as it downloads it to prevent it from loading the entire file into memory
func DownloadFile(filepath string, url string) error {

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
