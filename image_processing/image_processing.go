package imageprocessing

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"

	"github.com/nfnt/resize"
)

func ReadImage(path string) image.Image {
	inputFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	// Decode the image
	img, _, err := image.Decode(inputFile)
	if err != nil {
		fmt.Println(path)
		panic(err)
	}
	return img
}

func WriteImage(path string, img image.Image) {
	outputFile, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	// Encode the image to the new file
	err = jpeg.Encode(outputFile, img, nil)
	if err != nil {
		panic(err)
	}
}

func Grayscale(img image.Image) image.Image {
	// Create a new grayscale image
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	// Convert each pixel to grayscale
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalPixel := img.At(x, y)
			grayPixel := color.GrayModel.Convert(originalPixel)
			grayImg.Set(x, y, grayPixel)
		}
	}
	return grayImg
}

// Resize resizes the image to a maximum width or height of 500 pixels
// while maintaining the aspect ratio.
func Resize(img image.Image) image.Image {
	originalBounds := img.Bounds()
	width := originalBounds.Dx()
	height := originalBounds.Dy()

	var newWidth, newHeight uint

	if width >= height {
		newWidth = 500
		newHeight = uint((float64(height) / float64(width)) * 500)
	} else {
		newHeight = 500
		newWidth = uint((float64(width) / float64(height)) * 500)
	}

	resized := resize.Resize(newWidth, newHeight, img, resize.Lanczos3)
	return resized
}
