package main

type RuneBuffer struct {
	buffer []rune
}

func (rb *RuneBuffer) Push(r rune) {
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
	return r, true
}

func (rb RuneBuffer) isEmpty() bool {
	return rb.Len() == 0
}

func (rb RuneBuffer) Len() int {
	return len(rb.buffer)
}
