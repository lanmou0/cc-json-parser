package main

import "strconv"

func exitOnError(err error, msg string) {
	if err != nil {
		Logger.Error(errorMsg)
		Logger.Error(msg)
		panic(err)
	}
}

func isNumber(text string) bool {
	_, err := strconv.Atoi(text)
	return err == nil
}
