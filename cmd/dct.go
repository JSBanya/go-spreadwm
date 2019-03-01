package main

import (
	"math"
)

/*
* Equation implementations based off https://www.mathworks.com/help/images/discrete-cosine-transform.html
 */

// 1-D DCT-II
func applyDCT1D(arr []float64) []float64 {
	dct := make([]float64, len(arr))

	for k := 0; k < len(arr); k++ {
		c := math.Sqrt(2.0 / float64(len(arr)))
		if k == 0 {
			c = 1.0 / math.Sqrt(float64(len(arr)))
		}

		sum := 0.0
		for n := 0; n < len(arr); n++ {
			sum += arr[n] * math.Cos(math.Pi*(2.0*float64(n)+1.0)*float64(k)/(2.0*float64(len(arr))))
		}
		dct[k] = c * sum
	}

	return dct
}

// 1-D DCT-III
func applyIDCT1D(arr []float64) []float64 {
	idct := make([]float64, len(arr))

	for k := 0; k < len(arr); k++ {
		sum := 0.0
		for n := 0; n < len(arr); n++ {
			c := math.Sqrt(2.0 / float64(len(arr)))
			if n == 0 {
				c = 1.0 / math.Sqrt(float64(len(arr)))
			}

			sum += c * arr[n] * math.Cos(math.Pi*(2.0*float64(k)+1.0)*float64(n)/(2.0*float64(len(arr))))
		}
		idct[k] = sum
	}

	return idct
}

// 2-D DCT-II
func applyDCT2D(mat [][]float64) [][]float64 {
	height := len(mat)
	if height == 0 {
		return nil
	}
	width := len(mat[0])

	dct := make([][]float64, height)
	for i := 0; i < height; i++ {
		dct[i] = make([]float64, width)
	}

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {

			ci := math.Sqrt(2.0 / float64(height))
			if i == 0 {
				ci = 1.0 / math.Sqrt(float64(height))
			}

			cj := math.Sqrt(2.0 / float64(width))
			if j == 0 {
				cj = 1.0 / math.Sqrt(float64(width))
			}

			var sum float64
			for k := 0; k < height; k++ {
				for l := 0; l < width; l++ {
					sum += mat[k][l] * math.Cos(math.Pi*(2.0*float64(k)+1.0)*float64(i)/(2.0*float64(height))) * math.Cos(math.Pi*(2.0*float64(l)+1.0)*float64(j)/(2.0*float64(width)))
				}
			}

			dct[i][j] = ci * cj * sum
		}
	}

	return dct
}

// 2-D DCT-III
func applyIDCT2D(mat [][]float64) [][]float64 {
	height := len(mat)
	if height == 0 {
		return nil
	}
	width := len(mat[0])

	idct := make([][]float64, height)
	for i := 0; i < height; i++ {
		idct[i] = make([]float64, width)
	}

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			var sum float64
			for k := 0; k < height; k++ {
				for l := 0; l < width; l++ {

					ck := math.Sqrt(2.0 / float64(height))
					if k == 0 {
						ck = 1.0 / math.Sqrt(float64(height))
					}

					cl := math.Sqrt(2.0 / float64(width))
					if k == 0 {
						cl = 1.0 / math.Sqrt(float64(width))
					}

					sum += ck * cl * mat[k][l] * math.Cos(math.Pi*(2.0*float64(i)+1.0)*float64(k)/(2.0*float64(height))) * math.Cos(math.Pi*(2.0*float64(j)+1.0)*float64(l)/(2.0*float64(width)))
				}
			}

			idct[i][j] = sum
		}
	}

	return idct
}
