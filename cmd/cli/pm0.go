package main

import (
	"log"
	"os"
	"path"

	pm0 "github.com/TrixiS/pm0/internal/cli"
	"github.com/TrixiS/pm0/internal/cli/command_context"
	"github.com/TrixiS/pm0/internal/cli/commands"
	"github.com/TrixiS/pm0/internal/daemon/pb"
	"github.com/TrixiS/pm0/internal/utils"
	"github.com/asdine/storm/v3"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TODO?: preserve stopped units

const CLIClientDBFilename = "pm0_cli.db"
const MaxRecvMessageSizeBytes = 8 * 8 * 1024 * 1024 // 8 MB

func main() {
	pm0Dirpath, err := utils.GetPM0Dirpath()

	if err != nil {
		log.Fatal(err)
	}

	dbFilepath := path.Join(pm0Dirpath, CLIClientDBFilename)

	contextProvider := &command_context.CommandContextProvider{
		DBFactory: func() *storm.DB {
			db, err := storm.Open(dbFilepath)

			if err != nil {
				log.Fatalf("db open: %v", err)
			}

			return db
		},
		WithClient: func(f func(pb.ProcessServiceClient) error) error {
			conn, err := grpc.Dial(
				"localhost:7777",
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MaxRecvMessageSizeBytes)),
			) // TODO: get host from somewhere (env/config)

			if err != nil {
				log.Fatalf("grpc dial: %v", err)
			}

			defer conn.Close()

			client := pb.NewProcessServiceClient(conn)
			return f(client)
		},
	}

	exceptFlag := &cli.StringSliceFlag{
		Name:     "except",
		Aliases:  []string{"e"},
		Required: false,
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
				Flags:  []cli.Flag{exceptFlag},
				Args:   true,
				Action: contextProvider.Wraps(commands.Stop),
			},
			{
				Name:   "restart",
				Usage:  "Restart a unit",
				Flags:  []cli.Flag{exceptFlag},
				Args:   true,
				Action: contextProvider.Wraps(commands.Restart),
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
				Flags:   []cli.Flag{exceptFlag},
				Args:    true,
				Action:  contextProvider.Wraps(commands.Delete),
			},
			{
				Name:   "show",
				Usage:  "Show unit info",
				Args:   true,
				Action: contextProvider.Wraps(commands.Show),
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		pm0.Printf(err.Error())
	}
}
