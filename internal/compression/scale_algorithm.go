package compression

import (
	"image"

	"golang.org/x/image/draw"
)

type ScaleAlgorithm struct{}

func NewScaleAlgorithm() *ScaleAlgorithm {
	return &ScaleAlgorithm{}
}

func (a *ScaleAlgorithm) Name() string {
	return "scale"
}

func (a *ScaleAlgorithm) CompressImage(img image.Image, params CompressionParams) image.Image {
	// If quality is 100, return the original image
	if params.Quality == 100 {
		return img
	}

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	
	// Calculate new dimensions based on quality
	scaleFactor := float64(params.Quality) / 100.0
	newWidth := int(float64(width) * scaleFactor)
	newHeight := int(float64(height) * scaleFactor)

	// Ensure minimum dimensions
	if newWidth < 10 {
		newWidth = 10
	}
	if newHeight < 10 {
		newHeight = 10
	}

	newImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.ApproxBiLinear.Scale(newImg, newImg.Bounds(), img, img.Bounds(), draw.Over, nil)

	return newImg
}