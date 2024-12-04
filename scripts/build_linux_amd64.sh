GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" cmd/daemon/pm0_daemon.go
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" cmd/cli/pm0.go
