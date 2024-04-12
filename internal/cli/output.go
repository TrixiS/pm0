package cli

import (
	"errors"
	"fmt"
	"time"

	"github.com/TrixiS/pm0/internal/daemon"
	"github.com/jedib0t/go-pretty/v6/text"
)

const tableNoneString = "None"

var pm0OutputPrefix = text.FgHiCyan.Sprint("[PM0] ")
var runningStatusString = text.FgGreen.Sprint("Running")
var exitedStatusString = text.FgWhite.Sprint("Exited")
var failedStatusString = text.FgRed.Sprint("Failed")
var stoppedStatusString = text.FgYellow.Sprint("Stopped")

var ErrNoIdent = errors.New("provide at least one unit identifier (id or name)")

func Printf(format string, args ...interface{}) {
	fmt.Printf(pm0OutputPrefix+format+"\n", args...)
}

func FormatNillablePointer[T any](ptr *T) string {
	if ptr == nil {
		return tableNoneString
	}

	return fmt.Sprintf("%v", *ptr)
}

func FormatUnitUptime(startedAt int64, status daemon.UnitStatus) string {
	if status != daemon.RUNNING {
		return tableNoneString
	}

	startedAtTime := time.Unix(startedAt, 0)
	now := time.Now()
	diff := now.Sub(startedAtTime)
	return diff.Round(time.Second).String()
}

func FormatUnitStatus(unitStatus daemon.UnitStatus) string {
	switch unitStatus {
	case daemon.RUNNING:
		return runningStatusString
	case daemon.EXITED:
		return exitedStatusString
	case daemon.FAILED:
		return failedStatusString
	case daemon.STOPPED:
		return stoppedStatusString
	default:
		return "Unknown"
	}
}
