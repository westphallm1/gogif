package main

import (
	"io"
	"log"
	"os"
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

var gifs []AsciiGif
var searchBar = NewAsciiInput("Search: ", 1, 1, 50)
var gifReaders = make(chan io.Reader)
var xIdx = 1

func showPreviews(giphys []GifResponse, height int) {
	for len(gifs) > 0 {
		last := len(gifs) - 1
		gifs[last].pause <- struct{}{}
		gifs[last] = AsciiGif{}
		gifs = gifs[:last]
	}
	var sb strings.Builder
	clearLines(&sb, 2, height+2)
	printSync(&sb)
	xIdx = 1
	for _, giphy := range giphys[:3] {
		gifReaders <- downloadGiphy(giphy.Id)
	}
}

func loadGifs(height int) {
	for {
		reader := <-gifReaders
		agif := NewAsciiGif(reader, 0, height, xIdx, 2)
		agif.scaleToHeight()
		xIdx += agif.width + 1
		gifs = append(gifs, agif)
		go agif.printLoop()
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
		moveCursor(&sb, 70, 1)
		sb.WriteString("Searching for: ")
		sb.WriteString(text)
		giphys := getGiphyJSON(text)
		showPreviews(giphys, height)
	}
	go loadGifs(height)
	go pollKeyStrokes()
	<-quit
}
