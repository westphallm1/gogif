package main

import "strings"

type AsciiInput struct {
	prompt       string
	x, y, length int
	chars        []byte
	index        int
}

func (ainput *AsciiInput) draw() {
	var sb strings.Builder
	moveCursor(&sb, ainput.x, ainput.y)
	sb.WriteString(ainput.prompt)
	printUnderline(&sb)
	sb.WriteString(strings.Repeat(" ", ainput.length))
	printTerm(&sb)
	printSynch(&sb)
}
func (ainput *AsciiInput) onKey(key []byte) {
	var sb strings.Builder
	start := len(ainput.prompt) + 1
	switch key[0] {
	case '\n':
		ainput.index = 0
		ainput.draw()
	case '\x7F':
		if ainput.index -= 1; ainput.index < 0 {
			ainput.index = 0
		}
		moveCursor(&sb, 1, start+ainput.index)
		printUnderline(&sb)
		sb.WriteRune(' ')
		printTerm(&sb)
	default:
		moveCursor(&sb, 1, start+ainput.index)
		printUnderline(&sb)
		sb.WriteString(string(key[0]))
		printTerm(&sb)
		ainput.index++
	}
	printSynch(&sb)
}

func NewAsciiInput(prompt string, x, y, length int) AsciiInput {
	return AsciiInput{
		prompt: prompt,
		x:      x,
		y:      y,
		length: length,
		index:  0,
		chars:  make([]byte, length),
	}
}
