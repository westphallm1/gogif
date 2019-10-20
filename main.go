package main

import (
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var TEST_IMAGE = "/home/mwestphall/Pictures/squidward.jpg"
var TEST_GIF = "/home/mwestphall/Pictures/squidward.gif"
var WIDTH = 40
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
func printGif(fileName string, x, y int, quit chan struct{}) {
	gif := readGif(fileName)
	xScale, yScale := getScale(gif.Image[0], WIDTH, HEIGHT)
	for {
		lastFrame := makeBlankFrame(WIDTH, HEIGHT)
		for i := range gif.Image {
			var sb strings.Builder
			moveCursor(&sb, x, y)
			printBg24(&sb, RGBA{0, 0, 0, 0})
			img := downscaleImage(gif.Image[i], WIDTH, HEIGHT, xScale, yScale)
			img.downscaled = spliceImages(lastFrame.downscaled, img.downscaled)
			img.printTo(&sb, x, y)
			printTerm(&sb)
			printSynch(&sb)
			time.Sleep(time.Duration(gif.Delay[i]*10) * time.Millisecond)
			if gif.Disposal[i] == 1 {
				lastFrame = img
			}
		}
		select {
		case <-quit:
			return
		default:
		}
	}
}

// https://stackoverflow.com/a/17278730
func pollKeyStrokes() {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	// restore the echoing state when exiting
	defer exec.Command("stty", "-F", "/dev/tty", "echo").Run()
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
	for i := range os.Args[1:] {
		gifs[i] = make(chan struct{})
		go printGif(os.Args[i+1], i*(WIDTH+2), 2, gifs[i])
	}
	go pollKeyStrokes()
	<-quit
}
