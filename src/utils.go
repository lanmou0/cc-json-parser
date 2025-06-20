package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func exitOnError(err error, msg string) {
	if err != nil {
		Logger.Error(errorMsg)
		Logger.Error(msg)
		os.Exit(-1)
	}
}

func isStrNumber(text string) bool {
	_, err := strconv.Atoi(text)
	return err == nil
}

func trimSpace(rs []rune) []rune {
	return []rune(strings.TrimSpace(string(rs)))
}

func dump(data JsonValue) {
	b, _ := json.MarshalIndent(data, "", "  ")
	fmt.Print(string(b))
}
