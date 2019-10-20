package main

import (
	"image"
	"image/gif"
	"io"
	"log"
	"strings"
	"time"
)

type AsciiGif struct {
	gif           *gif.GIF // the gif to print
	x, y          int      // the top left cursor position of the gif
	width, height int      // the width and height of the gif
	index         int      // index of the current frame
	pause         chan struct{}
	play          chan struct{}
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

func (agif *AsciiGif) printLoop() {
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
			printSync(&sb)
			time.Sleep(time.Duration(gif.Delay[agif.index]*10) * time.Millisecond)
			if gif.Disposal[agif.index] == 1 {
				lastFrame = img
			}
			select {
			case <-agif.pause:
				<-agif.play
			default:
			}
		}
	}

}

func readGif(reader io.Reader) *gif.GIF {
	gif, err := gif.DecodeAll(reader)
	if err != nil {
		log.Fatal(err)
	}
	return gif
}

func NewAsciiGif(reader io.Reader, width, height, x, y int) AsciiGif {
	return AsciiGif{
		gif:    readGif(reader),
		x:      x,
		y:      y,
		width:  width,
		height: height,
		pause:  make(chan struct{}),
		play:   make(chan struct{}),
	}
}

func CopyAsciiGif(other AsciiGif) AsciiGif {
	return AsciiGif{
		gif:    other.gif,
		x:      other.x,
		y:      other.y,
		width:  other.width,
		height: other.height,
		pause:  make(chan struct{}),
		play:   make(chan struct{}),
	}
}

func (agif *AsciiGif) scaleToHeight() {
	bounds := agif.gif.Image[0].Bounds()
	agif.width = 2 * agif.height * bounds.Max.X / bounds.Max.Y
}
