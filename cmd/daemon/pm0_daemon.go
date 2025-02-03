package main

import (
	"log/slog"
	"net"
	"os"
	"path"
	"sync"

	"github.com/TrixiS/pm0/internal/daemon"
	"github.com/TrixiS/pm0/internal/daemon/pb"
	"github.com/TrixiS/pm0/internal/utils"
	"github.com/asdine/storm/v3"
	"google.golang.org/grpc"
)

const DaemonDBFilename = "pm0_daemon.db"

func main() {
	pm0Dirpath, err := utils.GetPM0Dirpath()

	if err != nil {
		panic(err)
	}

	logsDirpath := path.Join(pm0Dirpath, "logs")

	if err := os.MkdirAll(logsDirpath, 0777); err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", "localhost:7777")

	if err != nil {
		panic(err)
	}

	dbFilepath := path.Join(pm0Dirpath, DaemonDBFilename)
	dbFactory := func() *storm.DB {
		db, err := storm.Open(dbFilepath)

		if err != nil {
			panic(err)
		}

		return db
	}

	daemonServer := daemon.NewDaemonServer(
		daemon.DaemonServerOptions{LogsDirpath: logsDirpath, DBFactory: dbFactory},
	)

	db := dbFactory()
	var unitModels []daemon.UnitModel

	if err := db.All(&unitModels); err != nil {
		panic(err)
	}

	db.Close()

	wg := sync.WaitGroup{}
	wg.Add(len(unitModels))

	for _, unitModel := range unitModels {
		model := unitModel

		go func() {
			defer wg.Done()
			_, err := daemonServer.StartUnit(model)

			if err != nil {
				slog.Error("start unit", model.ID, "err", err)
			}
		}()
	}

	wg.Wait()

	grpcServer := grpc.NewServer()
	pb.RegisterProcessServiceServer(grpcServer, daemonServer)

	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
