// unit tests for image processing functions

package imageprocessing_test

import (
	imageprocessing "goroutines_pipeline/image_processing"
	"image"
	"image/color"
	"image/draw"
	"testing"
)

// createTestImage generates a simple test image for unit tests
func createTestImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{255, 0, 0, 255}}, image.Point{}, draw.Src)
	return img
}

// test ReadImage if it successfully loads an image
func TestResize(t *testing.T) {
	img := createTestImage()
	resized := imageprocessing.Resize(img)
	if resized.Bounds().Dx() != 500 || resized.Bounds().Dy() != 500 {
		t.Errorf("Expected resized image to be 500x500 but got %dx%d", resized.Bounds().Dx(), resized.Bounds().Dy())
	}
}

// test Grayscale if it successfully converts an image to grayscale
func TestGrayscale(t *testing.T) {
	img := createTestImage()
	gray := imageprocessing.Grayscale(img)

	// Compare original and new pixel values (the grayscale image should not be pure red)
	original := img.At(0, 0)
	newPixel := gray.At(0, 0)
	if original == newPixel {
		t.Errorf("Expected grayscale image pixel to change, but it did not")
	}
}
