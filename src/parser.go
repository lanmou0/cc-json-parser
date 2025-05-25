package main

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"unicode"
)

var hasError = false
var errorMsg = ""
var line = 0
var offset = 0
var loop = 0

func logLoop() {
	loop += 1
	if loop > 150 {
		panic("too much looping")
	}
}

func parseJsonObject(scanner *BufferedScanner) JsonObject {
	store := make(JsonObject)
	logLoop()
	for {
		ch, done := scanner.Scan()

		Logger.Debug(fmt.Sprintf("in parse json object loop: %#U, done: %t", ch, done))
		if done {
			break
		}
		if hasError {
			exitOnError(errors.New("Malformed Json File"), "malformed Json File")
		}
		Logger.Debug(fmt.Sprintf("tscan %#U", ch))
		if unicode.IsSpace(ch) {
			if ch == LINE_BREAK {
				//FIXME move to BufferedScanner Interface
				offset = 0
				line += 1
			}
			continue
		}
		if R_CURLY == ch {
			return store
		}
		if QUOTE == ch {
			scanner.Buffer(ch)
			Logger.Debug(fmt.Sprintf("in parse json object quote handler: %#U", ch))
			key := parseJsonKey(scanner)
			value := parseJsonValue(scanner)
			if key == "" && value == nil {
				break
			}
			store[key] = value
		}
	}
	hasError = true
	errorMsg = "Malformed Object"

	return nil
}

func parseJsonArray(scanner *BufferedScanner) JsonArray {
	array := make(JsonArray, 0)

	logLoop()
	for {
		ch, done := scanner.Scan()
		if done {
			break
		}
		if hasError {
			exitOnError(errors.New("Malformed Json File"), "malformed Json File")
		}
		Logger.Debug(fmt.Sprintf("arscan %#U", ch))
		if unicode.IsSpace(ch) {
			continue
		}
		if R_BRAC == ch {
			return JsonArray(array)
		}

		value := parseJsonValue(scanner)
		array = append(array, value)
		Logger.Debug(fmt.Sprintf("array %v", value))
		Logger.Debug(fmt.Sprintf("array %v", array))
	}

	hasError = true
	errorMsg = "Malformed Array"
	return nil
}

func parseJsonPrimitive(scanner *BufferedScanner, primitive string) JsonValue {
	logLoop()
	token := scanner.ScanUntilExcludeAll(COMMA, R_BRAC, R_CURLY)
	token = trimSpace(token)

	if string(token) != primitive {
		hasError = true
		errorMsg = "Unrecognized Symbol: " + string(token)
	}

	return JsonValue(string(token))
}

func parseJsonNumber(scanner *BufferedScanner) JsonValue {
	token := scanner.ScanUntilExcludeAll(COMMA, R_BRAC, R_CURLY)
	logLoop()

	num := string(token)
	if !isStrNumber(num) {
		//FIXME add a error struct and methods
		hasError = true
		errorMsg = "Unrecognized Symbol: " + num
	}
	return string(num)
}

func parseJsonValue(scanner *BufferedScanner) JsonValue {
	for {
		ch, done := scanner.Scan()
		if done {
			break
		}
		if hasError {
			exitOnError(errors.New("Malformed Json File"), "malformed Json File")
		}
		logLoop()
		offset += 1
		if unicode.IsSpace(ch) {
			continue
		}
		Logger.Debug(fmt.Sprintf("kvscan %#U", ch))
		if COMMA == ch {
			break
		}
		scanner.Buffer(ch)
		if L_CURLY == ch {
			return parseJsonObject(scanner)
		}
		if L_BRAC == ch {
			return parseJsonArray(scanner)
		}
		if QUOTE == ch {
			return parseJsonString(scanner)
		}
		if 't' == ch {
			return parseJsonPrimitive(scanner, PRIM_TRUE)
		}
		if 'f' == ch {
			return parseJsonPrimitive(scanner, PRIM_FALSE)
		}
		if 'n' == ch {
			return parseJsonPrimitive(scanner, PRIM_NULL)
		}
		if unicode.IsDigit(ch) {
			return parseJsonNumber(scanner)
		}
	}
	//FIXME return error if we arrive here
	return nil
}

func parseJsonKey(scanner *BufferedScanner) JsonString {
	b := scanner.ScanUntilExclude(COLON)
	//for formatting purpose
	key := string(b[1 : len(b)-1])
	if b[0] != QUOTE && b[len(b)-1] != QUOTE {
		hasError = true
		errorMsg = "Malformed key" + key
	}
	colon, done := scanner.Scan()
	if done || COLON != colon {
		hasError = true
		errorMsg = "Missing token: ':'"
	}

	return JsonString(key)
}

func parseJsonString(scanner *BufferedScanner) JsonString {
	token := scanner.ScanUntilExcludeAll(COMMA, R_BRAC, R_CURLY)

	token = trimSpace(token)
	//for formatting purpose
	key := string(token[1 : len(token)-1])
	fmt.Printf("parse json string %s\n", string(token))
	if token[0] != QUOTE && token[len(token)-1] != QUOTE {
		hasError = true
		errorMsg = "Malformed key" + key
	}
	delim := scanner.Peek()
	if !slices.Contains([]rune{COMMA, R_BRAC, R_CURLY}, delim) {
		hasError = true
		errorMsg = "Missing token ',', ']' or '}'"
	}

	return JsonString(key)
}

func ParseJson(file *os.File) (JsonValue, error) {
	scanner := NewBufferedScanner(file)
	Logger.Debug(fmt.Sprint("in parse json"))
	ch := scanner.Peek()
	switch ch {
	case L_BRAC:
		return parseJsonArray(scanner), nil
	case L_CURLY:
		return parseJsonObject(scanner), nil
	}
	return nil, errors.New("Malformed Json File")
}
