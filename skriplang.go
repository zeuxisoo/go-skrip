package main

import (
	"os"

	"github.com/urfave/cli"

	"github.com/zeuxisoo/go-skrip/cmd"
)

const AppVersion = "0.1.0"

func main() {
	app := cli.NewApp()
	app.Name = "Skrip"
	app.Usage = "This is a skrip language usage"
	app.Version = AppVersion
	app.Commands = []cli.Command{
		cmd.Run,
	}

	app.Run(os.Args)
}
