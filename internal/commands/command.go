package commands

import (
	"github.com/buildkite/go-buildkite/v4"
	"github.com/rs/zerolog"
)

type Globals struct {
	Client  *buildkite.Client
	Version string
	Logger  zerolog.Logger
}
