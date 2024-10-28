package main

import (
	"log"
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
		log.Fatal(err)
	}

	logsDirpath := path.Join(pm0Dirpath, "logs")

	if err := os.MkdirAll(logsDirpath, 0777); err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", "localhost:7777")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	dbFilepath := path.Join(pm0Dirpath, DaemonDBFilename)
	dbFactory := func() *storm.DB {
		db, err := storm.Open(dbFilepath)

		if err != nil {
			log.Fatalf("db open: %v", err)
		}

		return db
	}

	daemonServer := daemon.NewDaemonServer(
		daemon.DaemonServerOptions{LogsDirpath: logsDirpath, DBFactory: dbFactory},
	)

	db := dbFactory()
	var unitModels []daemon.UnitModel

	if err := db.All(&unitModels); err != nil {
		log.Fatal(err)
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
				log.Println("start unit", model.ID, err.Error())
			}
		}()
	}

	wg.Wait()

	grpcServer := grpc.NewServer()
	pb.RegisterProcessServiceServer(grpcServer, daemonServer)
	log.Fatal(grpcServer.Serve(lis))
}
