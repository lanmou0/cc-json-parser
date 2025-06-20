package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"unicode"
	"unicode/utf8"
)

type BufferedScanner struct {
	scanner *bufio.Scanner
	buffer  rune
	hasPeek bool
}

func NewBufferedScanner(file *os.File) *BufferedScanner {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes)
	return &BufferedScanner{
		scanner: scanner}
}

func (bs *BufferedScanner) Scan() (rune, bool) {
	if bs.hasPeek {
		bs.hasPeek = false
		return bs.buffer, false
	}
	hasToken := bs.scanner.Scan()
	if bs.scanner.Err() != nil {
		fmt.Printf("error while scanning err: %s", bs.scanner.Err().Error())
	}
	ch, _ := utf8.DecodeRune(bs.scanner.Bytes())
	return ch, !hasToken
}

func (bs *BufferedScanner) ScanUntil(r rune) []rune {
	rb := make([]rune, 0)
	for {
		ch, done := bs.Scan()
		if done {
			break
		}

		rb = append(rb, ch)
		if ch == r {
			break
		}
	}

	return rb
}

func (bs *BufferedScanner) ScanN(n int) []rune {
	rb := make([]rune, 0)
	for n > 0 {
		ch, done := bs.Scan()
		if done {
			break
		}
		fmt.Printf("scan N %c\n", ch)

		rb = append(rb, ch)
		n -= 1
	}

	return rb
}

func (bs *BufferedScanner) ScanUntilExclude(r rune) []rune {
	rb := make([]rune, 0)
	for {
		ch, done := bs.Scan()
		if done {
			break
		}
		if ch == r {
			bs.Buffer(ch)
			break
		}
		rb = append(rb, ch)
	}

	return rb
}

func (bs *BufferedScanner) ScanUntilExcludeAll(r ...rune) []rune {
	rb := make([]rune, 0)
	for {
		ch, done := bs.Scan()
		if done {
			break
		}
		if slices.Contains(r, ch) {
			bs.Buffer(ch)
			break
		}
		rb = append(rb, ch)
	}

	return rb
}

func (bs *BufferedScanner) ScanNoSpace() rune {
	var r rune
	for {
		ch, done := bs.Scan()
		if done || !unicode.IsSpace(ch) {
			r = ch
			break
		}
	}

	return r
}

func (bs *BufferedScanner) Buffer(r rune) {
	bs.buffer = r
	bs.hasPeek = true
}

func (bs *BufferedScanner) PeekNoSpace() rune {
	var r rune
	for {
		ch, done := bs.Scan()
		if done || !unicode.IsSpace(ch) {
			bs.Buffer(ch)
			r = ch
			break
		}
	}

	return r
}

func (bs *BufferedScanner) Peek() rune {
	ch, _ := bs.Scan()
	bs.Buffer(ch)

	return ch
}
