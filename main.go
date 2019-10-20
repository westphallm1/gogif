package main

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var TEST_IMAGE = "/home/mwestphall/Pictures/squidward.jpg"
var TEST_GIF = "/home/mwestphall/Pictures/squidward.gif"

var printQueue chan *strings.Builder

func printSync(sb *strings.Builder) {
	if printQueue == nil {
		printQueue = make(chan *strings.Builder)
		go func() {
			for {
				builder := <-printQueue
				print(builder.String())
			}
		}()
	}
	printQueue <- sb
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
	for {
		os.Stdin.Read(b)
		go searchBar.onKey(b)
	}
}

var runningGifs = make(map[string]chan struct{})
var searchBar = NewAsciiInput("Search: ", 1, 1, 50)

func showPreviews(giphys []GifResponse, height int) {
	for _, oldId := range runningGifs {
		close(oldId)
	}
	xIdx := 1
	var sb strings.Builder
	clearLines(&sb, 2, height+2)
	printSync(&sb)
	for _, giphy := range giphys[:3] {
		agif := NewAsciiGif(readGiphy(giphy.Id), 0, height, xIdx, 2)
		agif.scaleToHeight()
		xIdx += agif.width + 1
		runningGifs[giphy.Id] = make(chan struct{})
		go agif.printLoop(runningGifs[giphy.Id])

	}

}

func main() {
	var sb strings.Builder
	clearScreen(&sb)
	printSync(&sb)
	quit := make(chan struct{})
	height, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	searchBar.draw()
	searchBar.callback = func(text string) {
		giphys := getGiphyJSON(text)
		showPreviews(giphys, height)
	}
	go pollKeyStrokes()
	<-quit
}
