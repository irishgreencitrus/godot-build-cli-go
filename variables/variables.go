package variables


const ALL_SELECTOR string = "all"
// Supported versions for building
var Versions = []string{
	"3.2.3-stable",
	"3.2.2-stable",
	"3.2.1-stable",
	"3.2-stable",
	"3.1.2-stable",
	"3.1.1-stable",
	"3.1-stable",
}
// Supported platforms for building
var Platforms = []string{
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
// Supported binaries or types for building
var Types = []string{
	"editor",
	"export",
	"headless",
	"server",
}
// Flags for CurrentType
var CurrentTypeFlag = ""
// One of the types in Types[]
var CurrentType = Types[0]
// More human readable names for the output build files. 
// Will be used in helper's move built methods in the future
var FriendlyNames = map[string]string{
	"godot.x11.opt.tools" : "editor",
	"godot.x11.opt" : "export",
	"godot_server.x11.opt" : "server",
	//? "godot_server.x11.opt.tools" : "headless",
	
	
}