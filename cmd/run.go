package cmd

import (
	"fmt"
	"strings"
	"io/ioutil"

	"github.com/urfave/cli"

	"github.com/zeuxisoo/go-skriplang/pkg/logger"
	"github.com/zeuxisoo/go-skriplang/lexer"
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

	contentBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		logger.Error("Cannot open the script file")
		logger.Fatal("%v", err)
	}

	fmt.Println("Run run run ...")
	fmt.Println(string(contentBytes))

	fmt.Println(lexer.NewLexer(string(contentBytes)))

	return nil
}
