package main

import (
	"fmt"
	"hash/fnv"
	"math"
	"math/rand"
)

const (
	alpha = 0.1
)

// Sort a float array and return the sorted array and the indices
func sort(arr []float64) ([]float64, []int) {
	indices := make([]int, len(arr))
	for i := 0; i < len(arr); i++ {
		indices[i] = i
	}

	newArr := make([]float64, len(arr))
	for i := 0; i < len(arr); i++ {
		newArr[i] = arr[i]
	}

	for i := 0; i < len(arr); i++ {
		var maxI = i
		for j := i; j < len(arr); j++ {
			if newArr[j] < newArr[maxI] {
				maxI = j
			}
		}
		newArr[i], newArr[maxI] = newArr[maxI], newArr[i]
		indices[i], indices[maxI] = indices[maxI], indices[i]
	}

	return newArr, indices
}

// Generate a seeded gaussian sequence (mean=0, stdev=1.0) with the given length
func genGaussianSeq(key string, length int) []float64 {
	h := fnv.New64a()
	h.Write([]byte(key))
	seed := int64(h.Sum64())

	g := rand.New(rand.NewSource(seed))

	arr := make([]float64, length)
	for i := 0; i < length; i++ {
		arr[i] = g.NormFloat64()
	}

	return arr
}

// Encodes the given sequence into the provided DCT, starting from the greatest DCT values
func encode(dct []float64, seq []float64) []float64 {
	// Sort the DCT and obtain the indices of the greatest values
	_, indices := sort(dct)

	// Create a new array for the encoded DCT
	newDct := make([]float64, len(dct))
	for i := 0; i < len(dct); i++ {
		newDct[i] = dct[i]
	}

	// Encode the sequence
	for i := 0; i < len(seq); i++ {
		newDct[indices[i]] = dct[indices[i]] + alpha*seq[i]
	}

	return newDct
}

// Extracts the watermark from a watermarked image, provided with the original
func extract(dctOriginal []float64, dctWatermarked []float64, wlen int) []float64 {
	// Sort the original DCT and obtain the indices of the greatest values (in order to locate the watermark)
	_, indices := sort(dctOriginal)

	// Undo the watermarking process to extract the watermark
	watermark := make([]float64, wlen)
	for i := 0; i < wlen; i++ {
		watermark[i] = (dctWatermarked[indices[i]] - dctOriginal[indices[i]]) / (alpha)
	}

	return watermark
}

// Calculate the similarity between the original watermark sequence and an extracted watermark sequence
func similarity(originalSeq []float64, watermarkSeq []float64) (float64, error) {
	top, err := dot(watermarkSeq, originalSeq)
	if err != nil {
		return 0.0, err
	}

	bottom, err := dot(watermarkSeq, watermarkSeq)
	if err != nil {
		return 0.0, err
	}

	if math.Sqrt(bottom) == 0 {
		return 0.0, nil
	}

	return math.Abs(top / math.Sqrt(bottom)), nil
}

// Calculate the dot product between two arrays
func dot(x, y []float64) (d float64, err error) {
	if len(x) != len(y) {
		return 0.0, fmt.Errorf("lengths of vectors must be equal")
	}

	for i := 0; i < len(x); i++ {
		d += x[i] * y[i]
	}

	return
}
