package compression

import (
	"image"
)

type CompressionParams struct {
	Quality int
}

type CompressionAlgorithm interface {
	Name() string
	
	CompressImage(img image.Image, params CompressionParams) image.Image
}