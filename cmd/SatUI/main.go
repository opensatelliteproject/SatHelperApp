package main

import (
	"github.com/alecthomas/kong"
	"github.com/opensatelliteproject/SatHelperApp"
)

var cli struct {
	UI     struct{} `cmd help:"User Interface for SatUI" default:"1"`
	Server struct{} `cmd help:"SatHelperApp Headless mode"`
}

func main() {
	var err error
	ctx := kong.Parse(&cli,
		kong.Name("SatUI "+SatHelperApp.GetVersion()),
		kong.Description("GOES Satellite SDR Receiver"))

	satlog = satlog.WithCustomWriter(ctx.Stdout)

	switch ctx.Command() {
	case "ui":
		ctx.Printf("UI Mode")
		err = startUI(ctx)
	case "server":
		ctx.Printf("Server Mode")
		err = startServer(ctx)
	}

	ctx.FatalIfErrorf(err)
}
