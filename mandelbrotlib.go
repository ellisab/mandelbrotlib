// Package mandelbrotlib emits a PNG image of the Mandelbrot fractal.
package mandelbrotlib

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"math/big"
	"math/cmplx"
)

var (
	// Width of the mandelbrot image
	Width = 1024
	// Height of the mandelbrot image
	Height = 1024
	Delta  = 0.3 / float64(Width)
)

const (
	xmin, ymin, xmax, ymax = -2, -2, +2, +2
)

// GenMandelbrotCmplx128 creates a Mandelbrot image using complex128 type
func GenMandelbrotCmplx128(w io.Writer, zoom uint8) {
	img := image.NewRGBA(image.Rect(0, 0, Width, Height))
	for py := 0; py < Height; py++ {
		y := float64(py)/float64(Height)*(ymax-ymin) + ymin
		for px := 0; px < Width; px++ {
			x := float64(px)/float64(Width)*(xmax-xmin) + xmin
			// Image point(px, py) represents complex value z.
			img.Set(px, py, supersamplingCmplx128(x, y, zoom))
		}
	}
	png.Encode(w, img) // NOTE: ignoring errors
}

func mandelbrotCmplx128(z complex128) uint8 {
	const iterations = 200

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return n
		}
	}
	return 0
}

func supersamplingCmplx128(a, b float64, zoom uint8) color.Color {
	var p, x uint8
	if zoom == 1 {
		p = mandelbrotCmplx128(complex(a, b))
	}
	zo := zoom
	for zoom > 0 {
		p1 := mandelbrotCmplx128(complex(a+Delta, b+Delta))
		p2 := mandelbrotCmplx128(complex(a+Delta, b-Delta))
		p3 := mandelbrotCmplx128(complex(a-Delta, b-Delta))
		p4 := mandelbrotCmplx128(complex(a-Delta, b+Delta))
		x = (p1 + p2 + p3 + p4) / 4
		p += x
		zoom--
	}
	p /= zo

	switch {
	case p != 0 && p < 63:
		return color.RGBA{p * 4, 0, 0, 240}
	case p > 63 && p < 126:
		return color.RGBA{p * 4, (p - 63) * 4, 0, 240}
	case p > 126:
		return color.RGBA{p * 4, p * 4, (p - 126) * 4, 240}
	}
	return color.Black
}

// GenMandelbrotBigFloat creates a Mandelbrot image using math/big Float type
func GenMandelbrotBigFloat(w io.Writer, zoom uint8) {
	deltabf := big.NewFloat(4 / 1024)
	img := image.NewRGBA(image.Rect(0, 0, Width, Height))
	for py := 0; py < Height; py++ {
		y := float64(py)/float64(Height)*(ymax-ymin) + ymin
		by := big.NewFloat(y)
		for px := 0; px < Width; px++ {
			x := float64(px)/float64(Width)*(xmax-xmin) + xmin
			bx := big.NewFloat(x)
			// Image point (px, py) represents complex value z.
			img.Set(px, py, supersamplingBigFloat(bx, by, deltabf, zoom))
		}
	}
	png.Encode(w, img) // NOTE: ignoring errors
}

func mandelbrotBigFloat(r, i *big.Float) uint8 {
	const iterations = 200

	vR, vI := &big.Float{}, &big.Float{}
	abs := big.NewFloat(0)
	for n := uint8(0); n < iterations; n++ {
		vRn, vIn := big.NewFloat(0), big.NewFloat(0)
		vRn.Mul(vR, vR).Sub(vRn, (&big.Float{}).Mul(vI, vI)).Add(vRn, r)
		vIn.Mul(vR, vI).Mul(vIn, big.NewFloat(2)).Add(vIn, i)
		vR, vI = vRn, vIn
		abs.Mul(vR, vR).Add(abs, (&big.Float{}).Mul(vI, vI))
		if abs.Cmp(big.NewFloat(4)) == 1 {
			return n
		}
		abs.Set(big.NewFloat(0))
	}
	return 0
}

func supersamplingBigFloat(a, b, Delta *big.Float, zoom uint8) color.Color {
	var p, x uint8
	if zoom == 1 {
		p = mandelbrotBigFloat(a, b)
	}
	zero := big.Float{}
	suma, sumb := big.NewFloat(0), big.NewFloat(0)
	zo := zoom
	for zoom > 0 {
		p1 := mandelbrotBigFloat(suma.Add(a, Delta), sumb.Add(b, Delta))
		suma.Set(&zero)
		sumb.Set(&zero)
		p2 := mandelbrotBigFloat(suma.Add(a, Delta), sumb.Sub(b, Delta))
		suma.Set(&zero)
		sumb.Set(&zero)
		p3 := mandelbrotBigFloat(suma.Sub(a, Delta), sumb.Sub(b, Delta))
		suma.Set(&zero)
		sumb.Set(&zero)
		p4 := mandelbrotBigFloat(suma.Sub(a, Delta), sumb.Add(b, Delta))
		suma.Set(&zero)
		sumb.Set(&zero)
		x = (p1 + p2 + p3 + p4) / 4
		p += x
		zoom--
	}
	p /= zo
	switch {
	case p != 0 && p < 63:
		return color.RGBA{p * 4, 0, 0, 240}
	case p > 63 && p < 126:
		return color.RGBA{p * 4, (p - 63) * 4, 0, 240}
	case p > 126:
		return color.RGBA{p * 4, p * 4, (p - 126) * 4, 240}
	}
	return color.Black
}
