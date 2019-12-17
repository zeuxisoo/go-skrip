# Go Skrip

## Usage

Execute from file

    skrip run main.sk

Execute from inline code

    skrip eval 'print("1234")'

    skrip eval 'let a="this is a test";print(a)'

## Syntax

Define variable

    let name = "foo";
    let age  = 8;
    let boy  = true;
    let tail = 17.2;

Define array

    let array1 = [1,2,3];
    let array2 = ["foo", "bar", 1, 2.2];

    println(array1[1]);

Define hash

    let hash1 = {"a": 1, "b": 2};
    let hash2 = {1: "a", "b": 3.2};

    println(hash1["a"]);

If statement

    if (name == "foo") {
        println("hello");
    }else if (name == "bar") {
        println("hi");
    }else{
        println("bye");
    }

For statement

    for item in [1,2,3.2,"foo","bar"] {
        println(item)
    }

    for item in 1..3 {
        println(item)
    }

    for key, item in {"a": "foo", "b": "bar"} {
       println(key + " => " + item)
    }

Forever loop statement

    let x = 1;

    for {
        x = x + 1;

        println(x);

        if (x == 10) {
            break;
        }
    }

Function

    func name(first, last) {
        return first + " " + last;
    }

    println(name("tom", "cat"));

    let name2 = func(first, last) {
        return last + " " + first;
    };

    println(name2("tom", "cat"));

Syntax sugar

    let cat = {};
    cat.name   = "tom";
    cat.gender = "???";
    println(cat.name + "," + cat.gender);

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
