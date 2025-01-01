package main

import (
	"flag"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/log"
)

var logger = log.New(os.Stderr)

var Cli struct {
	Assets           AssetsOutput     `cmd:"" help:"Generate images of some short, using either counters or cards, from a JSON file"`
	Json             JsonOutput       `cmd:"" help:"Generate a JSON of some short, by transforming another JSON as input"`
	Vassal           vassalCli        `cmd:"" help:"Create a vassal module for testing"`
	GenerateTemplate GenerateTemplate `cmd:"" help:"Generates a new counter template file with default values"`
	CheckTemplate    CheckTemplate    `cmd:"" help:"Check if a JSON file is a valid counter template"`
	Tilemap          Tilemap          `cmd:"" help:"Generate a tilemap"`
}

func main() {
	now := time.Now()
	defer func(now time.Time) {
		logger.Infof("Execution time: %v", time.Since(now))
	}(now)

	flag.Parse()

	logger.SetReportTimestamp(false)
	logger.SetReportCaller(false)
	logger.SetLevel(log.DebugLevel)

	ctx := kong.Parse(&Cli)

	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
