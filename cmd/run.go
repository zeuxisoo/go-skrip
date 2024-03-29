package cmd

import (
	"io/ioutil"
	"strings"

	"github.com/urfave/cli"

	"github.com/zeuxisoo/go-skrip/evaluator"
	"github.com/zeuxisoo/go-skrip/lexer"
	"github.com/zeuxisoo/go-skrip/object"
	"github.com/zeuxisoo/go-skrip/parser"
	"github.com/zeuxisoo/go-skrip/pkg/logger"
)

// Run command for run the script file
var Run = cli.Command{
	Name:        "run",
	Usage:       "Run the script",
	Description: "Run the provided script file",
	Action:      runRun,
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

	theLexer := lexer.NewLexer(string(contentBytes))
	theParser := parser.NewParser(theLexer)
	theProgram := theParser.Parse()

	if len(theParser.Errors()) > 0 {
		for _, message := range theParser.Errors() {
			logger.Error(message)
		}
	} else {
		theEnvironment := object.NewEnvironment()
		theEvaluator := evaluator.Eval(theProgram, theEnvironment)

		if theEvaluator == nil {
			// nothing todo
		}
	}

	return nil
}
