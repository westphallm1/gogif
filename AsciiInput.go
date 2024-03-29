package main

import (
	"os"
	"os/exec"
	"strings"
)

type AsciiInput struct {
	prompt       string
	x, y, length int
	chars        []byte
	index        int
	callback     func(string)
}

func (ainput *AsciiInput) draw() {
	var sb strings.Builder
	moveCursor(&sb, ainput.x, ainput.y)
	sb.WriteString(ainput.prompt)
	printUnderline(&sb)
	sb.WriteString(strings.Repeat(" ", ainput.length))
	printTerm(&sb)
	printSync(&sb)
}

func (ainput *AsciiInput) text() string {
	return string(ainput.chars[:ainput.index])
}

func (ainput *AsciiInput) onKey(key []byte) {
	var sb strings.Builder
	start := len(ainput.prompt) + 1
	switch key[0] {
	case '\n':
		if ainput.index == 0 {
			go switchBigGif(1)
		} else if ainput.callback != nil {
			go ainput.callback(ainput.text())
		}
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
		if ainput.index >= ainput.length {
			return
		}
		moveCursor(&sb, 1, start+ainput.index)
		printUnderline(&sb)
		sb.WriteString(string(key[0]))
		printTerm(&sb)
		ainput.chars[ainput.index] = key[0]
		ainput.index++
	}
	printSync(&sb)
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

// TODO this correctly
func isArrow(key []byte) bool {
	if key[0] == '\033' {
		os.Stdin.Read(key)
		os.Stdin.Read(key)
		if key[0] == 'C' {
			go switchBigGif(1)
		} else if key[0] == 'D' {
			go switchBigGif(-1)
		}
		return true
	}
	return false
}

// https://stackoverflow.com/a/17278730
func pollKeyStrokes() {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	exec.Command("tput", "civis").Run()
	// restore the echoing state when exiting
	defer exec.Command("stty", "-F", "/dev/tty", "echo").Run()
	defer exec.Command("tput", "cvvis").Run()
	var b []byte = make([]byte, 1)
	for {
		os.Stdin.Read(b)
		if !isArrow(b) {
			searchBar.onKey(b)
		}
	}
}
