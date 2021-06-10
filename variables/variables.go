package variables

import "github.com/AlecAivazis/survey/v2"

const ALL_SELECTOR string = "all"

// Supported versions for building
var Versions = []string{
	"3.3.2-stable",
	"3.3.1-stable",
	"3.3-stable",
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
	EDITOR_FLAGS   = "platform=x11 target=release_debug tools=yes"
	EXPORT_FLAGS   = "platform=x11 target=release tools=no"
	HEADLESS_FLAGS = "platform=server target=release_debug tools=yes"
	SERVER_FLAGS   = "platform=server target=release tools=no"
)
var Bits = []string{
	"32",
	"64",
}
// Supported binaries or types for building
var Types = []string{
	"editor",
	"export",
	"headless",
	"server",
}

// More human readable names for the output build files.
// Will be used in helper's move built methods in the future
// Example "3.2.3-stable.godot_server.x11.opt.64.llvm" -> "3.2.3-stable.server.64.llvm"
var FriendlyNames = map[string]string{
	"godot.x11.opt.tools":        "editor",
	"godot.x11.opt":              "export",
	"godot_server.x11.opt":       "server",
	"godot_server.x11.opt.tools": "headless",
}
var ToolQuestions = []*survey.Question{
	{
		Name: "downloadver",
		Prompt: &survey.MultiSelect{
			Message: "Choose versions to download",
			Options: Versions,
		},
	},
	{
		Name: "buildver",
		Prompt: &survey.MultiSelect{
			Message: "Choose versions to build",
			Options: Versions,
		},
	},
	{
		Name: "binarytypes",
		Prompt: &survey.MultiSelect{
			Message: "Choose binary types to build",
			Options: Types,
			Default: []string{Types[0]},
		},
	},
	{
		Name: "bits",
		Prompt: &survey.Select{
			Message: "64bit or 32bit?",
			Options: Bits,
			Default: Bits[0],
		},
	},
	{
		Name: "removezips",
		Prompt: &survey.Confirm{
			Message: "Remove downloaded zip files?",
			Default: false,
		},
	},
	{
		Name: "movebuilt",
		Prompt: &survey.Confirm{
			Message: "Move built binaries to build/ ?",
			Default: false,
		},
	},
	{
		Name: "renamefriendly",
		Prompt: &survey.Confirm{
			Message: "Rename builds to a more friendly name?",
			Default: false,
		},
	},
}

type ToolAnswerType struct {
	DownloadVer    []string
	BuildVer       []string
	BinaryTypes    []string
	Bits           string
	RemoveZips     bool
	MoveBuilt      bool
	RenameFriendly bool
}

var ToolAnswers ToolAnswerType
