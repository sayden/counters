package main

import (
	"github.com/alecthomas/kong"
	"github.com/sayden/counters/output"
)

type vassalCli struct {
	Templates  []string `help:"List of templates to embed in the Vassal module" short:"t"`
	OutputPath string   `help:"Path for the temp folder" short:"o"`
}

func (v *vassalCli) Run(ctx *kong.Context) error {
	return output.VassalModule(v.OutputPath, v.Templates)
}
