package cmd

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"

	"github.com/zeuxisoo/go-skriplang/pkg/logger"
)

// Run command for run the script file
var Run = cli.Command{
	Name: "run",
	Usage: "Run the script",
	Description: "Run the provided script file",
	Action: runRun,
}

func runRun(c *cli.Context) error {
	filePath := c.Args().Get(0)

	if len(strings.TrimSpace(filePath)) <= 0 {
		logger.Fatal("Please enter the script file path")
	}

	fmt.Println("Run run run ...")

	return nil
}
