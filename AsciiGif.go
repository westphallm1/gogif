package main

import (
	"image"
	"image/gif"
	"log"
	"os"
	"strings"
	"time"
)

type AsciiGif struct {
	gif           *gif.GIF // the gif to print
	x, y          int      // the top left cursor position of the gif
	width, height int      // the width and height of the gif
	index         int      // index of the current frame
}

func (agif *AsciiGif) blankFrame() ImageConvert {
	return ImageConvert{
		width:      agif.width,
		height:     agif.height,
		downscaled: make([]RGBA, agif.width*agif.height),
	}
}

func (agif *AsciiGif) getScale() (int, int) {
	return getScale(agif.gif.Image[0], agif.width, agif.height)
}

func spliceImages(rgb1, rgb2 []RGBA) []RGBA {
	for i := range rgb1 {
		if rgb2[i].a < rgb1[i].a-20 {
			rgb2[i] = rgb1[i]
		}
	}
	return rgb2
}

func (agif *AsciiGif) currentFrame() image.Image {
	return agif.gif.Image[agif.index]
}

func (agif *AsciiGif) printLoop(quit chan struct{}) {
	gif := agif.gif
	xScale, yScale := agif.getScale()
	for {
		lastFrame := agif.blankFrame()
		for agif.index = range gif.Image {
			var sb strings.Builder
			printBold(&sb)
			printBg24(&sb, RGBA{0, 0, 0, 0})
			img := downscaleImage(
				agif.currentFrame(), agif.width, agif.height, xScale, yScale)
			img.downscaled = spliceImages(lastFrame.downscaled, img.downscaled)
			img.printTo(&sb, agif.x, agif.y)
			printTerm(&sb)
			printSynch(&sb)
			time.Sleep(time.Duration(gif.Delay[agif.index]*10) * time.Millisecond)
			if gif.Disposal[agif.index] == 1 {
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

func readGif(filePath string) *gif.GIF {
	reader, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	err = nil
	gif, err := gif.DecodeAll(reader)
	if err != nil {
		log.Fatal(err)
	}
	return gif
}

func NewAsciiGif(filePath string, width, height, x, y int) AsciiGif {
	return AsciiGif{
		gif:    readGif(filePath),
		x:      x,
		y:      y,
		width:  width,
		height: height,
	}
}
