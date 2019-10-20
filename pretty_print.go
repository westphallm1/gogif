package main

import (
	"image"
	"log"
	"strconv"
	"strings"
)

type ImageConvert struct {
	image          image.Image
	bounds         image.Rectangle
	width, height  int
	downscaled     []RGBA
	xScale, yScale int
}

func (convert *ImageConvert) set(x, y int, value RGBA) {
	if idx := x*convert.height + y; idx < len(convert.downscaled) {
		convert.downscaled[idx] = value
	}
}

func (convert *ImageConvert) samplePixels(x, y int) {
	startX := x * convert.xScale
	endX := startX + convert.xScale
	startY := y * convert.yScale
	endY := startY + convert.yScale
	if convert.bounds.Max.X < endX {
		endX = convert.bounds.Max.X
	}
	if convert.bounds.Max.Y < endY {
		endY = convert.bounds.Max.Y
	}
	nPixels := uint32((endX-startX)*(endY-startY)) * 0x101
	var r, g, b, a uint32
	for i := startX; i < endX; i++ {
		for j := startY; j < endY; j++ {
			r1, g1, b1, a1 := convert.image.At(i, j).RGBA()
			r += r1
			g += g1
			b += b1
			a += a1
		}
	}
	r /= nPixels
	g /= nPixels
	b /= nPixels
	a /= nPixels
	r -= r % 10
	g -= g % 10
	b -= b % 10
	convert.set(x, y, RGBA{int(r), int(g), int(b), int(a)})
}

func getScale(image image.Image, newWidth, newHight int) (int, int) {
	if image.Bounds().Min.X != 0 || image.Bounds().Min.Y != 0 {
		log.Fatal("Image doesn't start at 0!")
	}
	xScale := image.Bounds().Max.X / newWidth
	yScale := image.Bounds().Max.Y / newHight
	return xScale, yScale
}

// func makeBlankFrame(width, height int) ImageConvert {
// 	frame := ImageConvert{
// 		width:      width,
// 		height:     height,
// 		downscaled: make([]RGBA, width*height),
// 	}
// }
func downscaleImage(image image.Image, width, height, xScale, yScale int) ImageConvert {
	bounds := image.Bounds()
	convert := ImageConvert{
		image:      image,
		bounds:     bounds,
		width:      width,
		height:     height,
		xScale:     xScale,
		yScale:     yScale,
		downscaled: make([]RGBA, width*height),
	}
	// process the rows in parallel
	startX := bounds.Min.X / xScale
	startY := bounds.Min.Y / yScale
	endX := width - (width*xScale-bounds.Max.X)/xScale
	endY := height - (height*yScale-bounds.Max.Y)/yScale
	filledRows := make(chan bool)
	unfilledRows := make(chan bool)
	//zero out the rows (todo be more efficent)
	for x := 0; x < width; x++ {
		go func(x int) {
			for y := 0; y < height; y++ {
				convert.set(x, y, RGBA{0, 0, 0, 0})
			}
			unfilledRows <- true
		}(x)
	}
	for x := 0; x < width; x++ {
		<-unfilledRows
	}
	for x := startX; x < endX; x++ {
		go func(x int) {
			for y := startY; y < endY; y++ {
				convert.samplePixels(x, y)
			}
			filledRows <- true
		}(x)
	}
	for x := startX; x < endX; x++ {
		<-filledRows
	}
	return convert
}

func spliceImages(rgb1, rgb2 []RGBA) []RGBA {
	for i := range rgb1 {
		if rgb2[i].a > rgb1[i].a-20 {
			rgb1[i] = rgb2[i]
		}
	}
	return rgb1
}

func (convert *ImageConvert) printTo(sb *strings.Builder) {
	var lastRGB RGBA
	for j := 0; j < convert.height; j++ {
		for i := 0; i < convert.width; i++ {
			rgb := convert.downscaled[convert.height*i+j]
			if lastRGB.r != rgb.r || lastRGB.g != rgb.g || lastRGB.b != rgb.b {
				printFg24(sb, rgb)
			}
			sb.WriteRune(rgb.getRune())
		}
		sb.WriteRune('\n')
	}
}

func (convert *ImageConvert) toString() string {
	var sb strings.Builder
	convert.printTo(&sb)
	return sb.String()
}

type RGBA struct {
	r, g, b, a int
}

func (rgb *RGBA) printEscape(sb *strings.Builder) {
	sb.WriteString(strconv.FormatInt(int64(rgb.r), 10))
	sb.WriteRune(';')
	sb.WriteString(strconv.FormatInt(int64(rgb.g), 10))
	sb.WriteRune(';')
	sb.WriteString(strconv.FormatInt(int64(rgb.b), 10))
}

func (rgb *RGBA) getRune() rune {
	if rgb.a < 10 {
		return ' '
	}
	switch brightness := int(rgb.r) + int(rgb.g) + int(rgb.b); {
	case brightness < 100:
		return '*'
	case brightness < 200:
		return '!'
	case brightness < 300:
		return '('
	case brightness < 400:
		return '&'
	case brightness < 500:
		return '$'
	case brightness < 600:
		return '%'
	case brightness < 700:
		return '#'
	default:
		return '@'
	}
}

/*
 *Print the escape sequence for a 24-bit color foreground
 */
func printFg24(sb *strings.Builder, rgb RGBA) {
	sb.WriteString("\033[38;2;")
	rgb.printEscape(sb)
	sb.WriteRune('m')
}

/*
 * Print the escape sequence for a 24-bit color background
 */
func printBg24(sb *strings.Builder, rgb RGBA) {
	sb.WriteString("\033[48;2;")
	rgb.printEscape(sb)
	sb.WriteRune('m')
}

/*
 * Print the escape sequence for clearing the fore/background
 */
func printTerm(sb *strings.Builder) {
	sb.WriteString("\033[0m")
}

func moveCursor(sb *strings.Builder, x, y int) {
	sb.WriteString("\033[")
	sb.WriteString(strconv.FormatInt(int64(x), 10))
	sb.WriteRune(';')
	sb.WriteString(strconv.FormatInt(int64(y), 10))
	sb.WriteRune('H')
}

func clearScreen(sb *strings.Builder) {
	sb.WriteString("\033[2J")
}
