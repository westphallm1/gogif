package main

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

var TEST_IMAGE = "/home/mwestphall/Pictures/squidward.jpg"
var TEST_GIF = "/home/mwestphall/Pictures/squidward.gif"
var WIDTH = 120
var HEIGHT = WIDTH / 2

var printLock = sync.Mutex{}

func printSynch(sb *strings.Builder) {
	printLock.Lock()
	print(sb.String())
	printLock.Unlock()
}

/*
 * Continually process a gif, downsampling each frame to its
 * ASCII representation and printing it to the screen.
 *
 * TODO: Cache result
 */

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
	index := 1
	for {
		os.Stdin.Read(b)
		var sb strings.Builder
		switch b[0] {
		case '\n':
			index = 1
			moveCursor(&sb, 1, 1)
			clearLine(&sb)
		case '\x7F':
			if index -= 1; index < 1 {
				index = 1
			}
			moveCursor(&sb, 1, index)
			sb.WriteRune(' ')
		default:
			moveCursor(&sb, 1, index)
			sb.WriteString(string(b[0]))
			index++
		}
		printSynch(&sb)
	}
}

var gifs = make(map[int]chan struct{})

func main() {
	var sb strings.Builder
	clearScreen(&sb)
	printSynch(&sb)
	quit := make(chan struct{})
	height, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	for i := range os.Args[2:] {
		gifs[i] = make(chan struct{})
		agif := NewAsciiGif(os.Args[i+2], height*2, height, 1+i*(height*2+1), 2)
		go agif.printLoop(gifs[i])
	}
	go pollKeyStrokes()
	<-quit
}
