package main

import (
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

var printQueue chan *strings.Builder

const N_PICS = 6

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

func showPreviews(giphys []GifResponse, height int, start int) {
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
	for _, giphy := range giphys[start : start+N_PICS] {
		gifReaders <- downloadGiphy(giphy.Id)
	}
}

func loadGifs(height int) {
	for {
		reader := <-gifReaders
		agif := NewAsciiGif(reader, 0, height, xIdx, 2)
		if len(gifs) == 0 {
			showBigGif(agif)
		}
		agif.scaleToHeight()
		xIdx += agif.width + 1
		gifs = append(gifs, agif)
		go agif.printLoop()
	}
}

var BigGif AsciiGif
var BigGifIdx = 0

func showBigGif(other AsciiGif) {
	if (BigGif != AsciiGif{}) {
		BigGif.pause <- struct{}{}
		BigGif = AsciiGif{}
	}
	BigGif = CopyAsciiGif(other)
	BigGif.x = 1
	BigGif.y = BigGif.height + 3
	BigGif.height *= 4
	BigGif.scaleToHeight()
	go BigGif.printLoop()
}

func switchBigGif() {
	if BigGifIdx < 0 {
		BigGifIdx = 0
	}
	if BigGifIdx >= N_PICS {
		BigGifIdx = N_PICS - 1
	}
	showBigGif(gifs[BigGifIdx])
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
		showPreviews(giphys, height, 0)
	}
	go loadGifs(height)
	go pollKeyStrokes()
	<-quit
}
