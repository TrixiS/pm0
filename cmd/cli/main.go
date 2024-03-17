package main

import (
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
				Name:      "start",
				Aliases:   []string{"s"},
				Usage:     "start a process",
				UsageText: "command",
				Args:      true,
				Action:    contextProvider.Wraps(commands.Start),
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
