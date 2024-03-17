package command_context

import (
	"github.com/urfave/cli/v2"
)

type CommandContext struct {
	Provider   *CommandContextProvider
	CLIContext *cli.Context
}

type CommandFunc func(*CommandContext) error
