package main

import (
	"flag"
	"os"
	"runtime/pprof"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/log"
)

var logger = log.New(os.Stderr)

var Cli struct {
	// Asset option is used to generate assets directly. This is usually counters images or Card sheets.
	Assets AssetsOutput `cmd:"" help:"Generate images of some short, using either counters or cards, from a JSON file"`

	// JSON uses a JSON input to generate another JSON output.
	Json JsonOutput `cmd:"" help:"Generate a JSON of some short, by transforming another JSON as input"`

	// Vassal is used to generate a Vassal module for testing purposes.
	Vassal vassal `cmd:"" help:"Create a vassal module for testing. It searches for the 'template.xml' in the same folder"` //FIXME

	GenerateTemplate GenerateTemplate `cmd:"" help:"Generates a new counter template file with default values"`

	CheckTemplate CheckTemplate `cmd:"" help:"Check if a JSON file is a valid counter template"`
}

func main() {
	flag.Parse()
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	logger.SetReportTimestamp(false)
	logger.SetReportCaller(false)
	logger.SetLevel(log.DebugLevel)

	ctx := kong.Parse(&Cli)

	err = ctx.Run()
	ctx.FatalIfErrorf(err)

	logger.Info("Done")
}
