package web

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/irishgreencitrus/godot-build-cli-go/helper"
)
// Downloads the chosen Godot version source as on the github release.
// Although this is only intended to work with the supported versions, you might get lucky with newer
// or older versions on the github. However I have only tested the versions in variables.go
func DownloadVersion(chosenversion string) {
	url := fmt.Sprintf("https://github.com/godotengine/godot/archive/%s.zip", chosenversion)
	path := fmt.Sprintf("%s.zip", chosenversion)
	fmt.Println("Downloading", path)
	DownloadFile(path, url)
	fmt.Println("Unzipping", path)
	helper.Unzip(path, "download")
	fmt.Println("Successfully Unzipped", path)
}
// Initialises the downloading for godot version. This is more recommended than DownloadVersion as it 
// checks against supported version and also has support for downloading every version
// by simply calling DownloadInitialiser("all")
func DownloadInitialiser(versions []string){
	var wg sync.WaitGroup
	for _, v := range versions{
		fmt.Printf("Downloading %s to downloads/\n",v)
		wg.Add(1)
		go func(v string) {
			DownloadVersion(v)
			wg.Done()
		}(v)
	}
	wg.Wait()
}
// This function is stolen from
// https://golangcode.com/download-a-file-from-a-url/
// It writes to the file as it downloads it to prevent it from loading the entire file into memory.
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
