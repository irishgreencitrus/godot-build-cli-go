# Godot Build CLI

- [Godot Build CLI](#godot-build-cli)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Usage](#usage)
## Prerequisites

In order to use the build CLI at all, you should install the dependencies for compiling godot [here](https://docs.godotengine.org/en/stable/development/compiling/compiling_for_x11.html#distro-specific-one-liners)

Also if your OS is not in this list, consider it useless for now. However downloading should work no matter the OS. 
```
linux/amd64
linux/arm
linux/arm64
```
> Want to add support for more OSes?
> Make an issue or a PR and i'll address it in a later update
## Installation

If you have go installed you can use this one liner.
```bash
$ go install github.com/irishgreencitrus/godot-build-cli-go
```
Then you can run the tool with
```bash
$ ~/go/bin/godot-build-cli-go
```
> Add gopath to your path to run it with godot-build-cli-go


If you don't have go installed you can fetch the latest release from the releases page, but I recommend using the method above!

## Usage

```sh
$ ./go-build-cli-go -h
Usage of ./go-build-cli-go:
  -P	Prints available platforms
  -V	Prints available versions
  -Z	Removes version zip files
  -build string
    	Builds specified version
  -download string
    	Downloads specified version
  -move string
    	Moves specified builds to an easier to access location
```

Example of downloading every version supported building it moving the binaries and removing the zips

```
$ ./go-build-cli-go -download all -build all -Z -move all
```

Example of listing the versions
```
$ ./go-build-cli-go -V
```

Example of downloading and building 3.2.3-stable

```
$ ./go-build-cli-go -download 3.2.3-stable -build 3.2.3-stable
```

> Not enough info? Open an issue asking what you want to be documented, or a feature request





