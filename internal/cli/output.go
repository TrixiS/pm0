package cli

import (
	"fmt"
	"time"

	"github.com/TrixiS/pm0/internal/daemon"
	"github.com/jedib0t/go-pretty/v6/text"
)

const tableNoneString = "None"

var (
	pm0OutputPrefix     = text.FgHiCyan.Sprint("[PM0] ")
	runningStatusString = text.FgGreen.Sprint("Running")
	exitedStatusString  = text.FgWhite.Sprint("Exited")
	failedStatusString  = text.FgRed.Sprint("Failed")
	stoppedStatusString = text.FgYellow.Sprint("Stopped")
)

func Printf(format string, args ...any) {
	fmt.Printf(pm0OutputPrefix+format+"\n", args...)
}

func FormatUnitUptime(startedAt int64, status daemon.UnitStatus) string {
	if status != daemon.UnitStatusRunning {
		return tableNoneString
	}

	startedAtTime := time.Unix(startedAt, 0)
	now := time.Now()
	diff := now.Sub(startedAtTime)
	return diff.Round(time.Second).String()
}

func FormatUnitStatus(unitStatus daemon.UnitStatus) string {
	switch unitStatus {
	case daemon.UnitStatusRunning:
		return runningStatusString
	case daemon.UnitStatusExited:
		return exitedStatusString
	case daemon.UnitStatusFailed:
		return failedStatusString
	case daemon.UnitStatusStopped:
		return stoppedStatusString
	default:
		return "Unknown"
	}
}
