package main

import (
	"fmt"
	"log"
	"math"
	"os"
)

const (
	WATERMARK_LENGTH = 1000
	THRESHOLD        = 6
)

func main() {
	checkArgs()

	mode := os.Args[1]

	if mode == "-e" {
		performEncode()
	} else if mode == "-d" {
		performDecode()
	} else {
		fmt.Printf("Unknown flag: %s\n", mode)
	}
}

// Perform encoding of an image
func performEncode() {
	infile := os.Args[2]
	key := os.Args[3]
	outfile := os.Args[4]

	// Load image
	img, err := loadImageRGBA(infile)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(2)
	}

	// Extract 8-bit pixel components and convert to 32-bit pixel (preprocessing for DCT)
	p := uint8ArrayToUint32Array(img.Pix)
	predct := make([]float64, len(p))
	for i := 0; i < len(p); i++ {
		predct[i] = float64(p[i])
	}

	// Perform 1-D DCT-II on the pixel array
	log.Println("Calculating DCT...")
	dct := applyDCT1D(predct)

	// Generate a watermark sequence from the key and encode into DCT values
	log.Println("Encoding...")
	seq := genGaussianSeq(key, WATERMARK_LENGTH)
	newDct := encode(dct, seq)

	// Invert the DCT process to obtain (possibly modified) 32-bit pixels
	log.Println("Applying IDCT...")
	idct := applyIDCT1D(newDct)
	log.Println("Done")

	// Convert 32-bit pixel format back to 8-bit pixel format
	for i := 0; i < len(idct); i++ {
		p[i] = uint32(idct[i])
	}
	img.Pix = uint32ArrayToUint8Array(p)

	// Save image
	err = saveImageRGBA(img, outfile)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(2)
	}
}

// Decode and determine match of watermark
func performDecode() {
	originalImage := os.Args[2]
	key := os.Args[3]
	watermarkedImage := os.Args[4]

	// Load images
	oimg, err := loadImageRGBA(originalImage)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(2)
	}

	wimg, err := loadImageRGBA(watermarkedImage)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(2)
	}

	// Extract 8-bit pixel components and convert to 32-bit pixel (preprocessing for DCT)
	op := uint8ArrayToUint32Array(oimg.Pix)
	opredct := make([]float64, len(op))
	for i := 0; i < len(op); i++ {
		opredct[i] = float64(op[i])
	}

	wp := uint8ArrayToUint32Array(wimg.Pix)
	wpredct := make([]float64, len(wp))
	for i := 0; i < len(wp); i++ {
		wpredct[i] = float64(wp[i])
	}

	// Perform 1-D DCT-II on the pixel array
	log.Println("Calculating DCT (1/2)...")
	odct := applyDCT1D(opredct)

	log.Println("Calculating DCT (2/2)...")
	wdct := applyDCT1D(wpredct)

	// Generate true (original) watermark from key
	log.Println("Generating original watermark...")
	seq := genGaussianSeq(key, WATERMARK_LENGTH)

	// Extract watermark
	log.Println("Extracting watermark...")
	watermark := extract(odct, wdct, WATERMARK_LENGTH)

	// Calculate similarity
	log.Println("Calculating similarity...")
	sim, err := similarity(seq, watermark)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(3)
	}

	match := false
	if sim > THRESHOLD {
		match = true
	}

	log.Println("Done")
	fmt.Println()
	fmt.Printf("Similarity: %v (Match: %v)\n", sim, match)
}

// Convert float64 to uint32, with rounding and overflow prevention
func float64ToUint32(val float64) uint32 {
	if val <= 0.0 {
		return uint32(0)
	}

	if val >= math.MaxInt32 {
		return uint32(math.MaxInt32)
	}

	r := val - float64(uint32(val))
	if r < 0.5 {
		return uint32(val)
	}
	return uint32(val) + 1
}

func checkArgs() {
	if len(os.Args) != 5 {
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Printf("Encode new watermark: %s -e [input image] [key] [output image]\n", os.Args[0])
	fmt.Printf("Decode watermark: %s -d [original image] [key] [watermarked image]\n", os.Args[0])
}
