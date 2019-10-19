package main

import (
	"strings"
	"strconv"
	"image"
	"os"
	"log"
)

type ImageChunk struct {
}

type ImageConvert struct {
	image * image.Image
	bounds  image.Rectangle
	width, height int
	downscaled [] RGB
	xScale, yScale int
}
func (convert * ImageConvert) samplePixels(x, y int) {
	startX := convert.bounds.Min.X + x * convert.xScale
	endX := startX + convert.xScale
	startY := convert.bounds.Min.Y +  y * convert.yScale
	endY := startY + convert.yScale
} 

func downscaleImage(filePath string, height int, width int) []RGB{
	reader, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	image, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	bounds := image.Bounds()
	convert := ImageConvert {
		image: &image, 
		bounds: bounds,
		width: width, 
		height: height, 
		downscaled: make([]RGB, width * height),
	}

	convert.xScale = (bounds.Max.X - bounds.Min.X) / width
	convert.yScale = (bounds.Max.Y - bounds.Min.Y) / width
	
	for i := 0; i < width; i ++ {
		for j := 0; j < height; j++ {
			convert.samplePixels(i, j)
		}
	}
	return convert.downscaled
}

type RGB struct {
	r, g, b byte
}

func (rgb * RGB) printEscape(sb * strings.Builder) {
	sb.WriteString(strconv.FormatInt(int64(rgb.r), 10))
	sb.WriteRune(';')
	sb.WriteString(strconv.FormatInt(int64(rgb.g), 10))
	sb.WriteRune(';')
	sb.WriteString(strconv.FormatInt(int64(rgb.b), 10))
}

func (rgb * RGB) getRune() rune {
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
func printFg24(sb * strings.Builder, rgb RGB) {
	sb.WriteString("\033[38;2;")
	rgb.printEscape(sb)
	sb.WriteRune('m')
}

/*
 * Print the escape sequence for a 24-bit color background
 */
func printBg24(sb * strings.Builder, rgb RGB) {
	sb.WriteString("\033[48;2;")
	rgb.printEscape(sb)
	sb.WriteRune('m')
}

/*
 * Print the escape sequence for clearing the fore/background
 */
func printTerm(sb * strings.Builder) {
	sb.WriteString("\033[0m")
}


