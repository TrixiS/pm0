package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/TrixiS/pm0/internal/cli/command_context"
	"github.com/TrixiS/pm0/internal/cli/commands"
	"github.com/TrixiS/pm0/internal/daemon/pb"
	"github.com/TrixiS/pm0/internal/utils"
	"github.com/asdine/storm/v3"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const CLIClientDBFilename = "pm0_cli.db"

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
			) // TODO: get host from somewhere (env/config)

			if err != nil {
				log.Fatalf("grpc dial: %v", err)
			}

			defer conn.Close()

			client := pb.NewProcessServiceClient(conn)
			return f(client)
		},
	}

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name: "start",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "name",
						Required: true,
					},
				},
				Usage:     "start a unit",
				UsageText: "command",
				Args:      true,
				Action:    contextProvider.Wraps(commands.Start),
			},
			{
				Name:    "list",
				Aliases: []string{"ls"},
				Usage:   "list units",
				Args:    false,
				Action:  contextProvider.Wraps(commands.List),
			},
			{
				Name:   "stop",
				Usage:  "stop a unit",
				Args:   true,
				Action: contextProvider.Wraps(commands.Stop),
			},
			{
				Name:   "restart",
				Usage:  "restart a unit",
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
				Usage:  "show unit logfile contents",
				Args:   true,
				Action: contextProvider.Wraps(commands.Logs),
			},
			{
				Name:   "delete",
				Usage:  "delete units",
				Args:   true,
				Action: contextProvider.Wraps(commands.Delete),
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
