package configs

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

var (
	// Version is the version of the SKM.
	Version = "dev"
	// BuildDate is the date the SKM was built.
	BuildDate = ""
	// Commit is the commit hash the SKM was built from.
	Commit = ""
	// GoVersion is the version of Go used to build the SKM.
	GoVersion = runtime.Version()
	// OS is the operating system the SKM was built for.
	OS = runtime.GOOS
	// Arch is the architecture the SKM was built for.
	Arch = runtime.GOARCH
	// Compiler is the compiler used to build the SKM.
	Compiler = runtime.Compiler
	// Distribution is the distribution channel of the SKM e.g. apt, brew, snap, etc.
	Distribution = "Direct"
)

func init() {
	initVersion()
}

func initVersion() {
	if bi, isAvailable := debug.ReadBuildInfo(); isAvailable {
		if bi.Main.Version != "" {
			Version = bi.Main.Version
		}
		if Commit == "" {
			Commit = fmt.Sprintf("(unknown, sum=%s)", bi.Main.Sum)
		}
		if BuildDate == "" {
			BuildDate = "(unknown)"
		}
	}
}
