package cmd

import (
	"fmt"
	"strings"
	"io/ioutil"

	"github.com/urfave/cli"

	"github.com/zeuxisoo/go-skrip/pkg/logger"
	"github.com/zeuxisoo/go-skrip/lexer"
	"github.com/zeuxisoo/go-skrip/parser"
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

	theLexer   := lexer.NewLexer(string(contentBytes))
	theParser  := parser.NewParser(theLexer)
	theProgram := theParser.Parse()

	if len(theParser.Errors()) > 0 {
		for _, message := range theParser.Errors() {
			logger.Error(message)
		}
	}else{
		fmt.Println(theProgram.Statements)
	}

	return nil
}
