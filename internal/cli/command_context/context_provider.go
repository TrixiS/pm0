package command_context

import (
	"github.com/TrixiS/pm0/internal/daemon/pb"
	"github.com/asdine/storm/v3"
	"github.com/urfave/cli/v2"
)

type CommandContextProvider struct {
	DBFactory  func() *storm.DB
	WithClient func(f func(pb.ProcessServiceClient) error) error
}

func (provider *CommandContextProvider) Wraps(f CommandFunc) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		commandContext := CommandContext{
			Provider:   provider,
			CLIContext: ctx,
		}

		return f(&commandContext)
	}
}
