package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/urfave/cli"

	"github.com/zeuxisoo/go-skrip/evaluator"
	"github.com/zeuxisoo/go-skrip/lexer"
	"github.com/zeuxisoo/go-skrip/object"
	"github.com/zeuxisoo/go-skrip/parser"
	"github.com/zeuxisoo/go-skrip/pkg/logger"
)

var Cli = cli.Command{
	Name:        "cli",
	Usage:       "Start console mode",
	Description: "Start console mode to execute code",
	Action:      runCli,
}

var LivePrefixState struct {
	LivePrefix string
	IsEnable   bool
}

var keywords = []string{
	"func", "let", "true", "false", "if", "else",
	"return", "for", "in", "nil", "break", "continue",
}

var code = ""
var leftBraceCount = 0
var rightBraceCount = 0

var env *object.Environment

func init() {
	env = object.NewEnvironment()
}

func runCode(code string) {
	theLexer := lexer.NewLexer(code)
	theParser := parser.NewParser(theLexer)
	theProgram := theParser.Parse()

	if len(theParser.Errors()) > 0 {
		for _, message := range theParser.Errors() {
			logger.Error(message)
		}
	} else {
		theEvaluator := evaluator.Eval(theProgram, env)

		if theEvaluator != nil {
			if theEvaluator.Type() == object.ERROR_OBJECT {
				fmt.Println("Error")
			}

			if theEvaluator.Type() == object.NIL_OBJECT {
				fmt.Println()
			}
		} else {
			fmt.Println()
		}
	}
}

func executor(line string) {
	if strings.ToLower(strings.TrimSpace(line)) == "exit" {
		os.Exit(0)
	}

	// If the line end with {, it should be multiple line
	// 1. concat each line
	// 2. increase the left brace count
	// 2. change live prefix to ..
	if strings.HasSuffix(line, "{") {
		code = code + line
		leftBraceCount = leftBraceCount + 1

		LivePrefixState.LivePrefix = ".. "
		LivePrefixState.IsEnable = true

		return
	}

	// If the line end with }, it should be end of the multiple line
	// 1. increase the right brace count
	// 2. compare the left and right brace is or not equals
	// 3. if equals
	//		1. disable the multiple mode
	//		2. reset left and right brace count
	//		3. run all saved code
	//		4. clean all saved code
	if strings.HasSuffix(line, "}") && leftBraceCount > 0 {
		rightBraceCount = rightBraceCount + 1

		if leftBraceCount == rightBraceCount {
			code = code + line

			LivePrefixState.LivePrefix = ">> "
			LivePrefixState.IsEnable = true

			leftBraceCount = 0
			rightBraceCount = 0

			runCode(code)

			code = ""
		}

		return
	}

	// Single line code
	if leftBraceCount == 0 && rightBraceCount == 0 {
		runCode(line)
	} else {
		code = code + line
	}
}

func completer(document prompt.Document) []prompt.Suggest {
	suggests := []prompt.Suggest{
		// cli keywords
		prompt.Suggest{Text: "exit"},
	}

	// program keywords
	for _, keyword := range keywords {
		suggests = append(suggests, prompt.Suggest{
			Text: keyword,
		})
	}

	return prompt.FilterContains(suggests, document.GetWordBeforeCursor(), true)
}

func changeLivePrefixState() (string, bool) {
	return LivePrefixState.LivePrefix, LivePrefixState.IsEnable
}

func handleExit() {
	switch v := recover().(type) {
	case nil:
		return
	case int:
		os.Exit(int(v))
	default:
		fmt.Println(v)
	}
}

func runCli(c *cli.Context) error {
	defer handleExit()

	controlCKeyBinding := prompt.KeyBind{
		Key: prompt.ControlC,
		Fn: func(*prompt.Buffer) {
			panic(int(0))
		},
	}

	thePrompt := prompt.New(
		executor, completer,
		prompt.OptionPrefix(">> "),
		prompt.OptionTitle("skrip repl"),
		prompt.OptionLivePrefix(changeLivePrefixState),
		prompt.OptionAddKeyBind(controlCKeyBinding),
	)

	thePrompt.Run()

	return nil
}
