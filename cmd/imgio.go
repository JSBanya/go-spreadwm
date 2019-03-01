package main

import (
	"fmt"
	"golang.org/x/image/bmp"
	"image"
	"image/png"
	"os"
	"path/filepath"
)

// Load an image from a file
// Compatible with PNG and BMP formats
func loadImageRGBA(path string) (*image.RGBA, error) {
	ext := filepath.Ext(path)
	if ext != ".png" && ext != ".bmp" {
		return nil, fmt.Errorf("Unsupported file format: %s", ext)
	}

	// Open the image
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// Create decoder
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	// Covert type to RGBA
	var imgRGBA *image.RGBA
	switch img.(type) {
	case *image.RGBA:
		{
			imgRGBA = img.(*image.RGBA)
		}
	case *image.NRGBA:
		{
			imgNRGBA := img.(*image.NRGBA)
			imgRGBA = image.NewRGBA(imgNRGBA.Bounds())
			imgRGBA.Pix = imgNRGBA.Pix
			imgRGBA.Stride = imgNRGBA.Stride
		}
	default:
		return nil, fmt.Errorf("Unsupported image format")
	}

	return imgRGBA, nil
}

// Save an RGBA image to file
// Compatible with PNG and BMP formats
func saveImageRGBA(img *image.RGBA, path string) error {
	ext := filepath.Ext(path)
	if ext != ".png" && ext != ".bmp" {
		return fmt.Errorf("Unsupported file format: %s", ext)
	}

	// Create file
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write image
	err = nil
	switch ext {
	case ".png":
		err = png.Encode(f, img)
	case ".bmp":
		err = bmp.Encode(f, img)
	default:
		err = fmt.Errorf("Unknown or unsupported image format: %s", ext)
	}

	return err
}

// Convert uint8 array (in R, G, B, A order) to a 32 bit pixel array
func uint8ArrayToUint32Array(pixels []uint8) []uint32 {
	newPixels := make([]uint32, len(pixels)/4)

	for i := 0; i < len(pixels); i += 4 {
		newPixels[i/4] = (uint32(pixels[i]) << 24) | (uint32(pixels[i+1]) << 16) | (uint32(pixels[i+2]) << 8) | uint32(pixels[i+3])
	}

	return newPixels
}

// Convert a 32 bit pixel array to a uint8 array (in R, G, B, A order)
func uint32ArrayToUint8Array(pixels []uint32) []uint8 {
	newPixels := make([]uint8, len(pixels)*4)

	for i := 0; i < len(pixels); i++ {
		newPixels[i*4] = uint8(pixels[i] >> 24)
		newPixels[i*4+1] = uint8(pixels[i] >> 16)
		newPixels[i*4+2] = uint8(pixels[i] >> 8)
		newPixels[i*4+3] = uint8(pixels[i])
	}

	return newPixels
}
