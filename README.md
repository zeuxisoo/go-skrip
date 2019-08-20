# Go Skrip

## Usage

Execute from file

    skrip run main.sk

Execute from inline code

    skrip eval 'print("1234")'

    skrip eval 'let a="this is a test";print(a)'

## Development

Using the go module by default

    export GO111MODULE=on

Install the test package

    go get github.com/smartystreets/goconvey

Run the command to eval file

    go run *.go run ./test.sk

Run the command to eval inline code

    go run *.go eval 'print("1234")'

    go run *.go eval 'let a="this is a test";print(a)'

## Testing

Run all test case

    go test ./... -count

Run all test case in specify the package

    go test ./parser -count=1

Run selected test case in specify the package

    go test ./evaluator -run "TestForEachHashExpression"
