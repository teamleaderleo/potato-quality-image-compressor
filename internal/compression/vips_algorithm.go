package compression

import (
	"bytes"
	"image"
	"image/png"

	"github.com/davidbyttow/govips/v2/vips"
)

type VipsAlgorithm struct{}

func NewVipsAlgorithm() *VipsAlgorithm {
	vips.Startup(nil)
	return &VipsAlgorithm{}
}

func (a *VipsAlgorithm) Name() string {
	return "vips"
}

func (a *VipsAlgorithm) CompressImage(img image.Image, params CompressionParams) image.Image {
	// Encode input image to buffer (WebP encoder expects raw image bytes)
	var inputBuf bytes.Buffer
	_ = png.Encode(&inputBuf, img) // fallback to PNG as common baseline

	// Load image with libvips
	vipsImage, err := vips.NewImageFromBuffer(inputBuf.Bytes())
	if err != nil {
		return img // fallback: return original
	}
	defer vipsImage.Close()

	// Determine scale factor from quality
	scale := float64(params.Quality) / 100.0
	if scale <= 0.0 {
		scale = 1.0
	}

	_ = vipsImage.Resize(scale, vips.KernelNearest)

	webpBuf, _, err := vipsImage.ExportWebp(&vips.WebpExportParams{
		Quality:        params.Quality,
		StripMetadata:  true,
	})
	if err != nil {
		return img
	}

	// Decode final image for Go pipeline (must return image.Image)
	finalImg, _, err := image.Decode(bytes.NewReader(webpBuf))
	if err != nil {
		return img
	}
	return finalImg
}
