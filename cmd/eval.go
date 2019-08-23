package cmd

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"

	"github.com/zeuxisoo/go-skrip/evaluator"
	"github.com/zeuxisoo/go-skrip/lexer"
	"github.com/zeuxisoo/go-skrip/object"
	"github.com/zeuxisoo/go-skrip/parser"
	"github.com/zeuxisoo/go-skrip/pkg/logger"
)

var Eval = cli.Command{
	Name:        "eval",
	Usage:       "Eval the inline code",
	Description: "Eval the provided inline code",
	Action:      runEval,
}

func runEval(c *cli.Context) error {
	code := c.Args().Get(0)

	cleanCode := strings.TrimSpace(code)

	if len(cleanCode) <= 0 {
		logger.Fatal("Please enter the code to eval")
	}

	theLexer := lexer.NewLexer(cleanCode)
	theParser := parser.NewParser(theLexer)
	theProgram := theParser.Parse()

	if len(theParser.Errors()) > 0 {
		for _, message := range theParser.Errors() {
			logger.Error(message)
		}
	} else {
		theEnvironment := object.NewEnvironment()
		theEvaluator := evaluator.Eval(theProgram, theEnvironment)

		if theEvaluator != nil {
			if theEvaluator.Type() == object.ERROR_OBJECT {
				fmt.Println("Error")
			}
		}
	}

	return nil
}
