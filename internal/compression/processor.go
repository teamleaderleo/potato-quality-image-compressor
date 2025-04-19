package compression

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/chai2010/webp"
)

// ImageProcessor handles the common image processing operations
type ImageProcessor struct {
	algorithms     map[string]CompressionAlgorithm
	defaultAlgorithm CompressionAlgorithm
}

// NewImageProcessor creates a new ImageProcessor
func NewImageProcessor() *ImageProcessor {
	processor := &ImageProcessor{
		algorithms: make(map[string]CompressionAlgorithm),
	}

	// Register default (scale) algorithm
	scaleAlgorithm := NewScaleAlgorithm()
	processor.RegisterAlgorithm(scaleAlgorithm)
	processor.defaultAlgorithm = scaleAlgorithm

	// // Register libvips algorithm
	// vipsAlgorithm := NewVipsAlgorithm()
	// processor.RegisterAlgorithm(vipsAlgorithm)

	return processor
}

// RegisterAlgorithm registers a compression algorithm
func (p *ImageProcessor) RegisterAlgorithm(algorithm CompressionAlgorithm) {
	p.algorithms[algorithm.Name()] = algorithm
}

// SetDefaultAlgorithm sets the default algorithm to use
func (p *ImageProcessor) SetDefaultAlgorithm(name string) bool {
	if algorithm, exists := p.algorithms[name]; exists {
		p.defaultAlgorithm = algorithm
		return true
	}
	return false
}

// GetAlgorithm returns the algorithm with the given name
func (p *ImageProcessor) GetAlgorithm(name string) (CompressionAlgorithm, bool) {
	algorithm, exists := p.algorithms[name]
	return algorithm, exists
}

// GetDefaultAlgorithm returns the default algorithm
func (p *ImageProcessor) GetDefaultAlgorithm() CompressionAlgorithm {
	return p.defaultAlgorithm
}

// ProcessImage handles the complete process: decoding, compressing, and encoding
func (p *ImageProcessor) ProcessImage(input io.Reader, format string, quality int, algorithm CompressionAlgorithm) ([]byte, error) {
	// Decode the image
	img, _, err := image.Decode(input)
	if err != nil {
		return nil, fmt.Errorf("error decoding image: %v", err)
	}

	// Compress the image using the algorithm
	params := CompressionParams{
		Quality: quality,
	}
	compressedImg := algorithm.CompressImage(img, params)
	
	// Encode the image to the requested format
	return p.EncodeToFormat(compressedImg, format, quality)
}

// EncodeToFormat encodes an image to the specified format with the given quality
func (p *ImageProcessor) EncodeToFormat(img image.Image, format string, quality int) ([]byte, error) {
	var buf bytes.Buffer
	var err error

	switch format {
	case "webp":
		err = webp.Encode(&buf, img, &webp.Options{Quality: float32(quality)})
	case "jpeg", "jpg":
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
	case "png":
		err = png.Encode(&buf, img)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	if err != nil {
		return nil, fmt.Errorf("error encoding image to %s: %v", format, err)
	}

	return buf.Bytes(), nil
}