package commands

import (
	"fmt"
	"runtime"

	"github.com/buildkite/go-buildkite/v4"
	"github.com/rs/zerolog"
)

type Globals struct {
	Client  *buildkite.Client
	Version string
	Logger  zerolog.Logger
}

func UserAgent(version string) string {
	os := runtime.GOOS
	arch := runtime.GOARCH

	return fmt.Sprintf("buildkite-mcp-server/%s (%s; %s)", version, os, arch)
}
