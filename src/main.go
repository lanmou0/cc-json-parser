package main

import (
	"fmt"
	"os"
)

var Logger = NewLogger(false)

func main() {
	args := os.Args[1:]

	if len(args) > 1 {
		fmt.Println("Usage ccjp /path/to/file.json")
		return
	}

	filePath := args[0]
	file, err := os.Open(filePath)
	exitOnError(err, fmt.Sprintf("failed to open file: %s", filePath))
	defer file.Close()

	output, err := ParseJson(file)
	if err != nil {
		fmt.Printf("error parsing file: %s", err.Error())
	}
	dump(output)
}
