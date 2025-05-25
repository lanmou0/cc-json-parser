package main

import "fmt"

type RuneBuffer struct {
	buffer []rune
}

func (rb RuneBuffer) Get() []rune {
	return rb.buffer
}

func (rb *RuneBuffer) Push(r rune) {
	if rb.buffer == nil {
		rb.buffer = make([]rune, 0)
	}
	rb.buffer = append(rb.buffer, r)
}

func (rb *RuneBuffer) Pop() (rune, bool) {
	if rb.isEmpty() {
		return 0, false
	}
	r := rb.buffer[0]

	s := make([]rune, len(rb.buffer)-1)
	s = append(s, rb.buffer[1:]...)
	rb.buffer = s
	Logger.Debug(fmt.Sprintf("buffer: %s, %s, %c, %t\n", string(rb.Get()), string(s), r, rb.isEmpty()))
	return r, true
}

func (rb *RuneBuffer) PopAll() []rune {
	r := make([]rune, len(rb.buffer))
	r = append(r, rb.buffer[:]...)
	rb.buffer = make([]rune, 0)
	return r
}

func (rb RuneBuffer) isEmpty() bool {
	return rb.Len() == 0
}

func (rb RuneBuffer) Len() int {
	return len(rb.buffer)
}
