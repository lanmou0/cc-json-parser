package main

import (
	"errors"
	"strconv"
	"unicode/utf8"
)

func exitOnError(err error, msg string) {
	if err != nil {
		Logger.Error(errorMsg)
		Logger.Error(msg)
		panic(err)
	}
}

func isStrNumber(text string) bool {
	_, err := strconv.Atoi(text)
	return err == nil
}

func getStrLast(str string) rune {
	r, _ := utf8.DecodeLastRuneInString(str)
	// FIXME
	if r != utf8.RuneError {
		exitOnError(errors.New("Invalid Rune"), "Invalid symbol")
	}

	return r
}
