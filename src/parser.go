package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

var curlyCount = 0
var bracCount = 0
var hasError = false
var errorMsg = ""
var line = 0
var offset = 0
var runeBuf []rune
var breakScan = false

type JsonValue any
type JsonString string
type JsonArray []JsonValue
type Store map[JsonString]JsonValue

func parseJsonArray(scanner *BufferedScanner) JsonArray {
	array := make([]JsonValue, 0)

	for {
		ch, done := scanner.Scan()
		if done {
			break
		}
		Logger.Debug(fmt.Sprintf("arscan %#U\n", ch))
		if unicode.IsSpace(ch) {
			continue
		}

		if R_BRAC == ch {
			bracCount -= 1
			breakScan = true
			break
		}

		if COMMA == ch {
			continue
		}

		value := parseJsonValue(scanner)
		array = append(array, value)
		Logger.Debug(fmt.Sprintf("arscan %v\n", array))
	}

	return JsonArray(array)
}

func parseJsonValue(scanner *BufferedScanner) JsonValue {
	var vbuilder strings.Builder
	internal_store := make(Store)
	var array []Store
	qcount := 0
	for scanner.Scan() {
		offset += 1
		text, _ := utf8.DecodeRune(scanner.Bytes())
		Logger.Debug(fmt.Sprintf("kvscan %#U\n", text))
		if unicode.IsSpace(text) {
			continue
		}
		if COMMA == text {
			break
		}
		if QUOTE == text {
			qcount += 1
		}
		if L_CURLY == text {
			runeBuf = text
			scan(scanner, internal_store)
		}
		if L_BRAC == text {
			bracCount += 1
			array = parseArray(scanner)
		}
		if R_CURLY == text {
			breakScan = true
			curlyCount -= 1
			break
		}
		vbuilder.WriteRune(text)
	}

	value := vbuilder.String()
	if len(internal_store) != 0 {
		return internal_store
	}
	if len(array) != 0 {
		return array
	}
	if qcount == 0 {
		isPrim := PRIM_NULL == value || PRIM_FALSE == value || PRIM_TRUE == value
		isNum := isNumber(value)

		if !isPrim && !isNum {
			hasError = true
			errorMsg = fmt.Sprintf("%d:%d|unrecognized value", line, offset)
			return nil
		}
	} else if qcount != 2 {
		hasError = true
		errorMsg = fmt.Sprintf("%d:%d|mistmatched quotes", line, offset)
		return nil
	}
	return value
}

func parseJsonString(scanner *BufferedScanner) JsonString {
	var kbuilder strings.Builder
	qfound := false
	for scanner.Scan() {
		offset += 1
		text, _ := utf8.DecodeRune(scanner.Bytes())
		Logger.Debug(fmt.Sprintf("kvscan %#U\n", text))
		if QUOTE == text {
			qfound = true
			continue
		}
		if COLON == text {
			if !qfound {
				hasError = true
				errorMsg = fmt.Sprintf("%d:%d|Missing colon", line, offset)
				return ""
			}
			break
		}

		if COMMA == text {
			hasError = true
			errorMsg = fmt.Sprintf("%d:%d|Missing comma", line, offset)
			return ""
		}

		if R_CURLY == text {
			hasError = true
			errorMsg = fmt.Sprintf("%d:%d|Missmatched curly brace", line, offset)
			return ""
		}
		kbuilder.WriteRune(text)
	}
	key := kbuilder.String()

	return JsonString(key)
}

func tokenize(ch rune, scanner *BufferedScanner, store Store) {
	if unicode.IsSpace(ch) {
		if ch == '\n' {
			offset = 0
			line += 1
		}
		return
	}
	switch ch {
	case L_CURLY:
		curlyCount += 1
	case R_CURLY:
		breakScan = true
		curlyCount -= 1
	case QUOTE:
		key := parseJsonString(scanner)
		value := parseJsonValue(scanner)
		if key == "" && value == nil {
			break
		}
		store[key] = value
	}
}

func parseJsonObject(scanner *BufferedScanner, store Store) {
	if hasError {
		exitOnError(errors.New("malformed Json file"), "malformed Json file")
	}
	Logger.Debug(fmt.Sprintf("tscan %#U\n", ch))
	tokenize(ch, scanner, store)
	if breakScan {
		breakScan = false
		break
	}
}

func ParseJson(file *os.File) {
	scanner := NewBufferedScanner(file)
	var store = make(Store)
	rb := scanner.Peek(1)
	ch := rb.buffer[0]
	switch ch {
	case L_BRAC:
		parseJsonArray(scanner)
	case L_CURLY:
		parseJsonObject(scanner, store)
	}

	Logger.Info(fmt.Sprintf("store %v\n", store))
	Logger.Debug(fmt.Sprintf("errors %d, %d\n", bracCount, curlyCount))

	if bracCount+curlyCount != 0 {
		exitOnError(errors.New("malformed Json file"), "malformed Json file")
	}
}
