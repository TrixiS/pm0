package cli

import (
	"fmt"

	"github.com/fatih/color"
)

var MainColorFunc = color.New(color.FgBlue, color.Bold).SprintfFunc()
var TableHeaderColorFunc = color.New(color.FgBlue, color.Underline).SprintfFunc()

func Printf(format string, args ...interface{}) {
	fmt.Printf(MainColorFunc("[PM0] ")+format+"\n", args...)
}
