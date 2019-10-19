package main

import "strings"

var TEST_IMAGE = "/home/mwestphall/Pictures/squidward.jpg"

func main() {
	var sb strings.Builder
	printBg24(&sb, RGB{0, 0, 0})
	// for i := 0; i < 77; i ++ {
	// 	g := i*510/76
	// 	if g > 255 {
	// 		g = 510 - g
	// 	}
	// 	rgb := RGB {
	// 		byte(255 - (i*255)/76),
	// 		byte(g),
	// 		byte(i*255/76),
	// 	}
	// 	printFg24(&sb, rgb)
	// 	sb.WriteRune(rgb.getRune())
	// }
	// printTerm(&sb)
	// str := sb.String()
	// println(str)
	// println(len(str))
	img := downscaleImage(readImage(TEST_IMAGE), 80, 40)
	img.printTo(&sb)
	printTerm(&sb)
	println(sb.String())
}
