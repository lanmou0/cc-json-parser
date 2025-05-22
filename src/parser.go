package main

import (
	"errors"
	"fmt"
	"os"
	"unicode"
)

var hasError = false
var errorMsg = ""
var line = 0
var offset = 0

func parseJsonArray(scanner *BufferedScanner) JsonArray {
	array := make(JsonArray, 0)

	for {
		ch, done := scanner.Scan()
		if done {
			break
		}
		if hasError {
			exitOnError(errors.New("Malformed Json File"), "malformed Json File")
		}
		Logger.Debug(fmt.Sprintf("arscan %#U\n", ch))
		if unicode.IsSpace(ch) {
			continue
		}

		if R_BRAC == ch {
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

func parseJsonPrimitive(scanner *BufferedScanner, primitive string) JsonValue {
	ch := scanner.PeekUntil(getStrLast(primitive))
	if string(ch.Get()) != primitive {
		hasError = true
		errorMsg = "Unrecognized Symbol: " + string(ch.Get())
	}

	return JsonValue(string(ch.Get()))
}

func parseJsonNumber(scanner *BufferedScanner) JsonValue {
	buf := scanner.PeekUntilExclude(COMMA)

	num := string(buf.Get())
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
		offset += 1
		Logger.Debug(fmt.Sprintf("kvscan %#U\n", ch))
		if unicode.IsSpace(ch) {
			continue
		}
		if COMMA == ch {
			break
		}
		if QUOTE == ch {
			scanner.buffer.Push(ch)
			return parseJsonString(scanner)
		}
		if L_CURLY == ch {
			scanner.buffer.Push(ch)
			return parseJsonObject(scanner)
		}
		if L_BRAC == ch {
			scanner.buffer.Push(ch)
			return parseJsonArray(scanner)
		}
		if 't' == ch {
			scanner.buffer.Push(ch)
			return parseJsonPrimitive(scanner, PRIM_TRUE)
		}
		if 'f' == ch {
			scanner.buffer.Push(ch)
			return parseJsonPrimitive(scanner, PRIM_FALSE)
		}
		if 'n' == ch {
			scanner.buffer.Push(ch)
			return parseJsonPrimitive(scanner, PRIM_NULL)
		}
		if unicode.IsDigit(ch) {
			scanner.buffer.Push(ch)
			return parseJsonNumber(scanner)
		}
	}
	//FIXME return error if we arrive here
	return nil
}

func parseJsonString(scanner *BufferedScanner) JsonString {
	b := scanner.PeekUntil(QUOTE)
	scanner.buffer.PopAll()

	return JsonString(string(b.Get()))
}

func parseJsonObject(scanner *BufferedScanner) JsonObject {
	store := make(JsonObject)
	for {
		ch, done := scanner.Scan()
		if done {
			break
		}
		if hasError {
			exitOnError(errors.New("Malformed Json File"), "malformed Json File")
		}
		Logger.Debug(fmt.Sprintf("tscan %#U\n", ch))
		if unicode.IsSpace(ch) {
			if ch == '\n' {
				//FIXME move to BufferedScanner Interface
				offset = 0
				line += 1
			}
			continue
		}
		if QUOTE == ch {
			key := parseJsonString(scanner)
			value := parseJsonValue(scanner)
			if key == "" && value == nil {
				break
			}
			store[key] = value
		}
		}
	}

	return store
}

func ParseJson(file *os.File) (JsonValue, error) {
	scanner := NewBufferedScanner(file)
	ch := scanner.PeekOne()
	switch ch {
	case L_BRAC:
		return parseJsonArray(scanner), nil
	case L_CURLY:
		return parseJsonObject(scanner), nil
	}
	return nil, errors.New("Malformed Json File")
}
