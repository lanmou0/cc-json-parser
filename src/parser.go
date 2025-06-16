package main

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"unicode"
)

var hasError = false
var errorMsg = ""
var line = 0
var offset = 0
var loop = 0

func logLoop() {
	loop += 1
	Logger.Debug(fmt.Sprintf("in loop n: %d", loop))
	if loop > 150 {
		panic("too much looping")
	}
}

func parseJsonObject(scanner *BufferedScanner) JsonObject {
	var store JsonObject
	logLoop()
	for {
		ch, done := scanner.Scan()

		Logger.Debug(fmt.Sprintf("in parse json object loop: %#U, done: %t", ch, done))
		Logger.Debug(fmt.Sprintf("Object: current object: %v", store))
		if done {
			break
		}
		if hasError {
			exitOnError(errors.New("Malformed Json File"), "malformed Json File")
		}
		if unicode.IsSpace(ch) {
			if ch == LINE_BREAK {
				//FIXME move to BufferedScanner Interface
				offset = 0
				line += 1
			}
			continue
		}
		if L_CURLY == ch {
			Logger.Debug(fmt.Sprintf("Object: start new object"))
			store = make(JsonObject)
			continue
		}
		if R_CURLY == ch {
			Logger.Debug(fmt.Sprintf("Object: parsed new object: %v", store))
			return store
		}
		if QUOTE == ch {
			scanner.Buffer(ch)
			key := parseJsonKey(scanner)
			value := parseJsonValue(scanner)
			Logger.Debug(fmt.Sprintf("in parse json object quote handler: %v", value))
			if key == "" && value == nil {
				break
			}

			delim := scanner.ScanNoSpace()
			Logger.Debug(fmt.Sprintf("obj scan delim %#U", delim))

			if !slices.Contains([]rune{COMMA, R_BRAC, R_CURLY}, delim) {
				hasError = true
				errorMsg = "Object: Missing token ',', ']' or '}'."
				return nil
			}

			if R_CURLY == delim {
				scanner.Buffer(delim)
			}
			store[key] = value
		}
	}
	hasError = true
	errorMsg = "Malformed Object"

	return nil
}

func parseJsonArray(scanner *BufferedScanner) JsonArray {
	var store JsonArray
	logLoop()
	ch, done := scanner.Scan()
	if done {
		return nil
	}
	if L_BRAC == ch {
		Logger.Debug(fmt.Sprintf("Array: created new array: %v", store))
		store = make(JsonArray, 0)
	}
	for {
		ch, done := scanner.Scan()
		if done {
			break
		}
		if hasError {
			exitOnError(errors.New("Malformed Json File"), "malformed Json File")
		}
		if unicode.IsSpace(ch) {
			continue
		}
		Logger.Debug(fmt.Sprintf("Array: scanned %#U", ch))
		if R_BRAC == ch {
			Logger.Debug(fmt.Sprintf("Array: parsed new array: %v", store))
			return JsonArray(store)
		}

		scanner.Buffer(ch)
		value := parseJsonValue(scanner)
		delim := scanner.ScanNoSpace()

		if !slices.Contains([]rune{COMMA, R_BRAC, R_CURLY}, delim) {
			hasError = true
			errorMsg = "Array: Missing token ',', ']' or '}'."
			return nil
		}

		if R_BRAC == delim {
			scanner.Buffer(delim)
		}

		store = append(store, value)
		Logger.Debug(fmt.Sprintf("array %v", value))
		Logger.Debug(fmt.Sprintf("array %v", store))
	}

	hasError = true
	errorMsg = "Malformed Array"
	return nil
}

func parseJsonValue(scanner *BufferedScanner) JsonValue {
	logLoop()
	for {
		ch, done := scanner.Scan()
		if done {
			break
		}
		if hasError {
			exitOnError(errors.New("Malformed Json File"), "malformed Json File")
		}
		offset += 1
		if unicode.IsSpace(ch) {
			continue
		}
		Logger.Debug(fmt.Sprintf("kvscan %#U", ch))
		scanner.Buffer(ch)
		fmt.Printf("json value %c\n", ch)
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

	return nil
}

func parseJsonKey(scanner *BufferedScanner) JsonString {
	token := scanner.ScanUntilExclude(COLON)
	token = trimSpace(token)

	if token[0] != QUOTE && token[len(token)-1] != QUOTE {
		hasError = true
		errorMsg = "Malformed key" + string(token)
	}

	//for formatting purpose
	key := string(token[1 : len(token)-1])
	if strings.Contains(key, string(QUOTE)) {
		hasError = true
		errorMsg = "Malformed key" + string(token)
		return ""
	}

	colon, done := scanner.Scan()
	if done || COLON != colon {
		hasError = true
		errorMsg = "Missing token: ':'"
		return ""
	}

	return JsonString(key)
}

func parseJsonString(scanner *BufferedScanner) JsonString {
	token := scanner.ScanUntilExcludeAll(COMMA, R_BRAC, R_CURLY)
	token = trimSpace(token)

	fmt.Printf("parse json string %s\n", string(token))
	if token[0] != QUOTE && token[len(token)-1] != QUOTE {
		hasError = true
		errorMsg = "Malformed key" + string(token)
	}

	key := string(token[1 : len(token)-1])
	if strings.Contains(key, string(QUOTE)) {
		hasError = true
		errorMsg = "Malformed key" + string(token)
		return ""
	}
	delim := scanner.Peek()
	if !slices.Contains([]rune{COMMA, R_BRAC, R_CURLY}, delim) {
		hasError = true
		errorMsg = "String: Missing token ',', ']' or '}'. Token: " + key
		return ""
	}

	return JsonString(key)
}

func parseJsonPrimitive(scanner *BufferedScanner, primitive string) JsonValue {
	token := scanner.ScanN(len(primitive))
	token = trimSpace(token)
	fmt.Printf("parse json prim %s\n", string(token))

	if string(token) != primitive {
		hasError = true
		errorMsg = "Unrecognized Symbol: " + string(token)
		return nil
	}

	delim := scanner.PeekNoSpace()
	fmt.Printf("parse json prim delim %s\n", string(delim))

	if !slices.Contains([]rune{COMMA, R_BRAC, R_CURLY}, delim) {
		hasError = true
		errorMsg = "Primitive: Missing token ',', ']' or '}'. Token: " + string(token)
		return nil
	}
	return JsonValue(string(token))
}

func parseJsonNumber(scanner *BufferedScanner) JsonValue {
	token := scanner.ScanUntilExcludeAll(COMMA, R_BRAC, R_CURLY)

	token = trimSpace(token)
	num := string(token)
	if !isStrNumber(num) {
		//FIXME add a error struct and methods
		hasError = true
		errorMsg = "Unrecognized Symbol: " + num
		return nil
	}

	delim := scanner.Peek()
	if !slices.Contains([]rune{COMMA, R_BRAC, R_CURLY}, delim) {
		hasError = true
		errorMsg = "Missing token ',', ']' or '}'"
		return nil
	}

	return num
}

func ParseJson(file *os.File) (JsonValue, error) {
	scanner := NewBufferedScanner(file)
	Logger.Debug(fmt.Sprintf("in parse json"))
	value := parseJsonValue(scanner)
	if hasError {
		return nil, errors.New(errorMsg)
	}
	return value, nil
}
