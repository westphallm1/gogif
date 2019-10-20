package main

import (
	"os"
	"strings"
	"time"
)

var TEST_IMAGE = "/home/mwestphall/Pictures/squidward.jpg"
var TEST_GIF = "/home/mwestphall/Pictures/squidward.gif"
var WIDTH = 120
var HEIGHT = WIDTH / 2

func main() {
	gif := readGif(os.Args[1])
	xScale, yScale := getScale(gif.Image[0], WIDTH, HEIGHT)
	var lastFrame []RGBA
	for {
		for i := range gif.Image {
			var sb strings.Builder
			clearScreen(&sb)
			printBg24(&sb, RGBA{0, 0, 0, 0})
			moveCursor(&sb, 0, 0)
			img := downscaleImage(gif.Image[i], WIDTH, HEIGHT, xScale, yScale)
			if i > 0 {
				img.downscaled = spliceImages(lastFrame, img.downscaled)
			}
			lastFrame = img.downscaled
			img.printTo(&sb)
			printTerm(&sb)
			println(sb.String())
			time.Sleep(100 * time.Millisecond)
		}
	}
}
