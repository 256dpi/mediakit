package main

import (
	"flag"
	"os"

	"github.com/256dpi/xo"
	"github.com/kr/pretty"

	"github.com/256dpi/mediakit"
)

func main() {
	// parse flags
	flag.Parse()

	// open file
	file, err := os.Open(flag.Arg(0))
	xo.PanicIf(err)
	defer file.Close()

	// analyze file
	report, err := mediakit.Analyze(nil, file)
	xo.PanicIf(err)

	// print report
	pretty.Println(report)
}
