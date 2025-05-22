package main

import (
	"bufio"
	"os"
	"unicode/utf8"
)

type BufferedScanner struct {
	scanner *bufio.Scanner
	buffer  RuneBuffer
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
		r, hp := bs.buffer.Pop()
		bs.hasPeek = hp
		return r, false
	}
	done := bs.scanner.Scan()
	ch, _ := utf8.DecodeRune(bs.scanner.Bytes())
	return ch, done
}

func (bs *BufferedScanner) PeekOne() rune {

	return bs.Peek(1).Get()[0]
}

func (bs *BufferedScanner) Peek(n int) *RuneBuffer {
	if bs.hasPeek {
		return &bs.buffer
	}
	for bs.scanner.Scan() && n > 0 {
		ch, _ := utf8.DecodeRune(bs.scanner.Bytes())
		bs.buffer.Push(ch)
		n -= 1
	}
	bs.hasPeek = true

	return &bs.buffer
}

func (bs *BufferedScanner) PeekUntil(r rune) *RuneBuffer {
	if bs.hasPeek {
		return &bs.buffer
	}
	for bs.scanner.Scan() {
		ch, _ := utf8.DecodeRune(bs.scanner.Bytes())
		bs.buffer.Push(ch)
		if ch == r {
			break
		}
	}
	//FIXME error if not found
	bs.hasPeek = true

	return &bs.buffer
}

func (bs *BufferedScanner) PeekUntilExclude(r rune) *RuneBuffer {
	if bs.hasPeek {
		return &bs.buffer
	}
	for bs.scanner.Scan() {
		ch, _ := utf8.DecodeRune(bs.scanner.Bytes())
		if ch == r {
			break
		}
		bs.buffer.Push(ch)
	}
	//FIXME error if not found
	bs.hasPeek = true

	return &bs.buffer
}
