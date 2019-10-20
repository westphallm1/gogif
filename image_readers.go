package main

import (
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
)

func readJpeg(filePath string) image.Image {
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

func readPng(filePath string) image.Image {
	reader, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	image, err := png.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	return image
}
