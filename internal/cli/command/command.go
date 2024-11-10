package command

import (
	"github.com/TrixiS/pm0/internal/daemon/pb"
	"github.com/asdine/storm/v3"
	"github.com/urfave/cli/v2"
)

type Context struct {
	Provider ContextProvider
	CLI      *cli.Context
}

type CommandFunc func(*Context) error

type ContextProvider struct {
	DBFactory  func() *storm.DB
	WithClient func(f func(pb.ProcessServiceClient) error) error
}

func (provider ContextProvider) Wraps(f CommandFunc) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		commandContext := Context{
			Provider: provider,
			CLI:      ctx,
		}

		return f(&commandContext)
	}
}
