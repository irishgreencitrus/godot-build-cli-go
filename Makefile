all:
	GOOS=linux GOARCH=arm64 go build -o bin/godot-build-cli-go-linux.arm.64 -ldflags="-s -w" main.go
	GOOS=linux GOARCH=arm go build -o bin/godot-build-cli-go-linux.arm -ldflags="-s -w" main.go 
	GOOS=linux GOARCH=amd64 go build -o bin/godot-build-cli-go-linux.amd64 -ldflags="-s -w" main.go
