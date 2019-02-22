package cmd

import (
	"fmt"

	"github.com/urfave/cli"
)

// Run command for run the script file
var Run = cli.Command{
	Name: "run",
	Usage: "Run the script",
	Description: "Run the provided script file",
	Action: runRun,
}

func runRun(c *cli.Context) error {
	fmt.Println("Run run run ...")

	return nil
}
