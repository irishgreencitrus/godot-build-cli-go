# Godot Build CLI
[![Go Reference](https://pkg.go.dev/badge/github.com/irishgreencitrus/godot-build-cli-go.svg)](https://pkg.go.dev/github.com/irishgreencitrus/godot-build-cli-go)
- [Godot Build CLI](#godot-build-cli)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Usage](#usage)
    - [An important change!](#an-important-change)
      - [Version Guide](#version-guide)
      - [Type Guide](#type-guide)
    - [Help and common usage](#help-and-common-usage)
## Prerequisites

In order to use the build CLI at all, you should install the dependencies for compiling godot [here](https://docs.godotengine.org/en/stable/development/compiling/compiling_for_x11.html#distro-specific-one-liners)
>### IMPORTANT NOTICE!!
>If you see an error like this,
>```
>scons: building terminated because of errors.
>sh: 1: clang++: not found
>```
>you don't have clang installed! Install it using your package manager before opening an issue

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
$ go install github.com/irishgreencitrus/godot-build-cli-go@latest
```
Then you can run the tool with
```bash
$ ~/go/bin/godot-build-cli-go
```
> Add gopath to your path to run it with godot-build-cli-go


If you don't have go installed you can fetch the latest release from the releases page, but I recommend using the method above!

## Usage
### An important change!
In the last commit you could only build/download/move one version or all of them.
In the latest commit you can now build any permutation of versions using this guide. One means download/build/move that version and zero means don't
#### Version Guide
```
3.2.3-stable  - 1
3.2.2-stable  - 1
3.2.1-stable  - 1
3.2-stable    - 1
3.1.2-stable  - 1
3.1.1-stable  - 1
3.1-stable    - 1
```
1111111 in decimal is 127 - so to build every version you need to put the flag 127. Bits are in this order! That means read from top to bottom. Here's a guide for types
#### Type Guide
```
editor    - 1
export    - 1
headless  - 0
server    - 1
```
So to build every type of binary *except* headless just use the flag `-type 13` as 1101 in decimal is 13.

### Help and common usage

```sh
$ ./godot-build-cli-go -h
Usage of ./godot-build-cli-go:
  -P	Prints available platforms
  -T	Prints possible types
  -V	Prints available versions
  -Z	Removes version zip files
  -build int
    	Builds specified version
  -download int
    	Downloads specified version
  -move int
    	Moves specified builds to an easier to access location
  -type int
    	Chooses the type to build (default 1)
```

Example of downloading every version supported building it moving the binaries and removing the zips

```
$ ./go-build-cli-go -download 127 -build 127 -Z -move all
```

Example of listing the versions
```
$ ./go-build-cli-go -V
```

Example of downloading and building 3.2.3-stable

```
$ ./go-build-cli-go -download 1 -build 1
```

> Not enough info? Open an issue asking what you want to be documented, or a feature request





