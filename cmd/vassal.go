package main

import (
	"github.com/alecthomas/kong"
	"github.com/sayden/counters/output"
)

type vassalCli struct {
	Templates  []string `help:"List of templates to embed in the Vassal module" short:"t"`
	OutputPath string   `help:"Path to the temp folder" short:"o"`
	output.VassalConfig
}

func (i *vassalCli) Run(ctx *kong.Context) error {
	if i.Templates != nil {
		return i.TemplatesToVassal()
	}

	return output.CSVToVassalFile(i.VassalConfig)
}

func (v *vassalCli) TemplatesToVassal() error {
	return output.VassalModule(v.OutputPath, v.Templates)
}
