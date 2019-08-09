# Go Skrip

## Development

Using the go module by default

    export GO111MODULE=on

Install the test package

    go get github.com/smartystreets/goconvey

## Testing

Run all test case

    go test ./... -count

Run all test case in specify the package

    go test ./parser -count=1

Run selected test case in specify the package

    go test ./evaluator -run "TestForEachHashExpression"
