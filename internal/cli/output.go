package cli

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/text"
)

var pm0OutputPrefix = text.FgHiCyan.Sprint("[PM0] ")

func Printf(format string, args ...interface{}) {
	fmt.Printf(pm0OutputPrefix+format+"\n", args...)
}
