package main

import (
	"image"
	"image/jpeg"
	"log"
	"os"
	"strconv"
	"strings"
)

type ImageConvert struct {
	image          image.Image
	bounds         image.Rectangle
	width, height  int
	downscaled     []RGB
	xScale, yScale int
}

func (convert *ImageConvert) set(x, y int, value RGB) {
	convert.downscaled[x*convert.width+y] = value
}

func (convert *ImageConvert) samplePixels(x, y int) {
	startX := convert.bounds.Min.X + x*convert.xScale
	endX := startX + convert.xScale
	startY := convert.bounds.Min.Y + y*convert.yScale
	endY := startY + convert.yScale
	nPixels := uint32((endX - startX) * (endY - startY)) * 0x101
	if convert.bounds.Max.X < endX {
		endX = convert.bounds.Max.X
	}
	if convert.bounds.Max.Y < endY {
		endY = convert.bounds.Max.Y
	}
	var r, g, b uint32
	for i := startX; i < endX; i++ {
		for j := startY; j < endY; j++ {
			r1, g1, b1, _ := convert.image.At(i, j).RGBA()
			r += r1
			g += g1
			b += b1
		}
	}
	r /= nPixels
	g /= nPixels
	b /= nPixels
	convert.set(x, y, RGB{byte(r), byte(g), byte(b)})
}

func downscaleImage(image image.Image, height int, width int) ImageConvert {
	bounds := image.Bounds()
	convert := ImageConvert{
		image:      image,
		bounds:     bounds,
		width:      width,
		height:     height,
		downscaled: make([]RGB, width*height),
	}

	convert.xScale = (bounds.Max.X - bounds.Min.X) / width
	convert.yScale = (bounds.Max.Y - bounds.Min.Y) / height

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			convert.samplePixels(i, j)
		}
	}
	return convert
}

func (convert *ImageConvert) printTo(sb *strings.Builder) {
	for j := 0; j < convert.height; j++ {
		for i := 0; i < convert.width; i++ {
			rgb := convert.downscaled[convert.width *i + j]
			printFg24(sb, rgb)
			sb.WriteRune(rgb.getRune())
		}
		sb.WriteRune('\n')
	}
	for i := 0; i < convert.width*convert.height; i++ {
	}
}

func (convert *ImageConvert) toString() string {
	var sb strings.Builder
	convert.printTo(&sb)
	return sb.String()
}

func readImage(filePath string) image.Image {
	reader, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	image, err := jpeg.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	return image
}

type RGB struct {
	r, g, b byte
}

func (rgb *RGB) printEscape(sb *strings.Builder) {
	sb.WriteString(strconv.FormatInt(int64(rgb.r), 10))
	sb.WriteRune(';')
	sb.WriteString(strconv.FormatInt(int64(rgb.g), 10))
	sb.WriteRune(';')
	sb.WriteString(strconv.FormatInt(int64(rgb.b), 10))
}

func (rgb *RGB) getRune() rune {
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
func printFg24(sb *strings.Builder, rgb RGB) {
	sb.WriteString("\033[38;2;")
	rgb.printEscape(sb)
	sb.WriteRune('m')
}

/*
 * Print the escape sequence for a 24-bit color background
 */
func printBg24(sb *strings.Builder, rgb RGB) {
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
