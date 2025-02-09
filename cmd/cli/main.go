package main

import (
	"log"
	"os"
	"path"

	pm0 "github.com/TrixiS/pm0/internal/cli"
	"github.com/TrixiS/pm0/internal/cli/command"
	"github.com/TrixiS/pm0/internal/cli/commands"
	"github.com/TrixiS/pm0/internal/daemon/pb"
	"github.com/TrixiS/pm0/internal/utils"
	"github.com/asdine/storm/v3"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	cliClientDBFilename     = "pm0_cli.db"
	maxRecvMessageSizeBytes = 8 * 1024 * 1024
)

func main() {
	pm0Dirpath, err := utils.GetPM0Dirpath()

	if err != nil {
		log.Fatal(err)
	}

	dbFilepath := path.Join(pm0Dirpath, cliClientDBFilename)

	contextProvider := command.ContextProvider{
		DBFactory: func() *storm.DB {
			db, err := storm.Open(dbFilepath)

			if err != nil {
				log.Fatalf("db open: %v", err)
			}

			return db
		},
		WithClient: func(f func(pb.ProcessServiceClient) error) error {
			conn, err := grpc.NewClient(
				"localhost:7777",
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxRecvMessageSizeBytes)),
			) // TODO: get address from somewhere (env/args)

			if err != nil {
				log.Fatalf("grpc dial: %v", err)
			}

			defer conn.Close()

			client := pb.NewProcessServiceClient(conn)
			return f(client)
		},
	}

	app := &cli.App{
		Name:  "pm0",
		Usage: "CLI client for PM0 daemon",
		Commands: []*cli.Command{
			{
				Name: "start",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "name",
						Required: false,
					},
					&cli.StringFlag{
						Name:     "cwd",
						Required: false,
					},
					&cli.StringSliceFlag{
						Name:     "env",
						Required: false,
						Aliases:  []string{"e"},
					},
				},
				Usage:     "Start a unit",
				UsageText: "command",
				Args:      true,
				Action:    contextProvider.Wraps(commands.Start),
			},
			{
				Name:    "list",
				Aliases: []string{"ls"},
				Usage:   "List units",
				Args:    false,
				Action:  contextProvider.Wraps(commands.List),
			},
			{
				Name:   "stop",
				Usage:  "Stop a unit",
				Args:   true,
				Action: contextProvider.Wraps(commands.Stop),
				Subcommands: []*cli.Command{
					createAllSubcommand(contextProvider.Wraps(commands.StopAll)),
				},
			},
			{
				Name:   "restart",
				Usage:  "Restart a unit",
				Args:   true,
				Action: contextProvider.Wraps(commands.Restart),
				Subcommands: []*cli.Command{
					createAllSubcommand(contextProvider.Wraps(commands.RestartAll)),
				},
			},
			{
				Name: "logs",
				Flags: []cli.Flag{
					&cli.Uint64Flag{
						Name:     "lines",
						Required: false,
					},
					&cli.BoolFlag{
						Name:     "follow",
						Required: false,
						Aliases:  []string{"f"},
					},
				},
				Usage:  "Show unit logfile contents",
				Args:   true,
				Action: contextProvider.Wraps(commands.Logs),
				Subcommands: []*cli.Command{
					{
						Name:   "clear",
						Usage:  "Clear unit log file",
						Args:   true,
						Action: contextProvider.Wraps(commands.LogsClear),
					},
				},
			},
			{
				Name:    "delete",
				Usage:   "Delete units",
				Aliases: []string{"rm"},
				Args:    true,
				Action:  contextProvider.Wraps(commands.Delete),
				Subcommands: []*cli.Command{
					createAllSubcommand(contextProvider.Wraps(commands.DeleteAll)),
				},
			},
			{
				Name:   "show",
				Usage:  "Show unit info",
				Args:   true,
				Action: contextProvider.Wraps(commands.Show),
			},
			{
				Name:   "setup",
				Action: contextProvider.Wraps(commands.Setup),
			},
			{
				Name: "update",
				Flags: []cli.Flag{
					&cli.Uint64Flag{
						Name:     "id",
						Required: true,
					},
					&cli.StringFlag{
						Name: "name",
					},
					&cli.StringSliceFlag{
						Name:    "env",
						Aliases: []string{"e"},
					},
				},
				Action: contextProvider.Wraps(commands.Update),
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		pm0.Printf(err.Error())
	}
}

var allSubCommandFlags = []cli.Flag{
	&cli.Uint64SliceFlag{
		Name:     "except",
		Aliases:  []string{"e"},
		Required: false,
	},
}

func createAllSubcommand(action cli.ActionFunc) *cli.Command {
	return &cli.Command{
		Name:   "all",
		Flags:  allSubCommandFlags,
		Args:   false,
		Action: action,
	}
}
