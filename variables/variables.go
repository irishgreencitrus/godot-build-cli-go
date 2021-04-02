package variables


const ALL_SELECTOR string = "all"
var Versions = []string{
	"3.2.3-stable",
	"3.2.2-stable",
	"3.2.1-stable",
	"3.2-stable",
	"3.1.2-stable",
	"3.1.1-stable",
	"3.1-stable",
}
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
var Types = []string{
	"editor",
	"export",
	"headless",
	"server",
}

var CurrentTypeFlag = ""
var CurrentType = Types[0]