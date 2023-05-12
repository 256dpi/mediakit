package main

import (
	"flag"
	"os"

	"github.com/256dpi/mediakit"
	"github.com/kr/pretty"
)

func main() {
	// parse flags
	flag.Parse()

	// open file
	file, err := os.Open(flag.Arg(0))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// analyze file
	report, err := mediakit.Analyze(nil, file)
	if err != nil {
		panic(err)
	}

	// print report
	pretty.Println(report)
}
