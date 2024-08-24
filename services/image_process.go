package services

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"

	"github.com/nfnt/resize"
)

func ReadImage(path string) image.Image {
	// Check if the file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("File does not exist at path: " + path)
	}

	inputFile, err := os.Open(path)
	if err != nil {
		panic("Error opening file: " + err.Error())
	}
	defer inputFile.Close()

	// Decode the image
	img, _, err := image.Decode(inputFile)
	if err != nil {
		panic("Error decoding image: " + err.Error())
	}
	return img
}

func WriteImage(path string, img image.Image) {
	outputFile, err := os.Create(path)
	if err != nil {
		panic("Error creating output file: " + err.Error())
	}
	defer outputFile.Close()

	// Encode the image to the new file
	err = jpeg.Encode(outputFile, img, nil)
	if err != nil {
		panic(err)
	}
}

func Grayscale(img image.Image) image.Image {
	// Create a new grayscale image with the same dimensions
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	// Convert the original image to grayscale
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalPixel := img.At(x, y)
			grayPixel := color.GrayModel.Convert(originalPixel).(color.Gray)
			grayImg.Set(x, y, grayPixel)
		}
	}
	return grayImg
}

func Resize(img image.Image) image.Image {
	newWidth := uint(600)
	newHeight := uint(600)
	resizedImg := resize.Resize(newWidth, newHeight, img, resize.Lanczos3)
	return resizedImg
}
