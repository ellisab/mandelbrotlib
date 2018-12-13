package mandelbrotlib

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func BenchmarkGenMandelbrotCmplx128z1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f, err := ioutil.TempFile("", "1024*1024*z1.png")
		defer os.Remove(f.Name())
		if err != nil {
			log.Fatalf("%s\n", err)
		}
		GenMandelbrotCmplx128(f, 1)
	}
}

func BenchmarkGenMandelbrotCmplx128z10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f, err := ioutil.TempFile("", "1024x1024*z10.png")
		defer os.Remove(f.Name())
		if err != nil {
			log.Fatalf("%s\n", err)
		}
		GenMandelbrotCmplx128(f, 10)
	}
}

func BenchmarkGenMandelbrotCmplx128z30(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f, err := ioutil.TempFile("", "1024x1024*z30")
		defer os.Remove(f.Name())
		if err != nil {
			log.Fatalf("%s\n", err)
		}
		GenMandelbrotCmplx128(f, 30)
	}
}

// using math.big datatypes is very slow (expected)

func BenchmarkGenMandelbrotBigFLoatz1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f, err := ioutil.TempFile("", "1024*1024*z1.png")
		defer os.Remove(f.Name())
		if err != nil {
			log.Fatalf("%s\n", err)
		}
		GenMandelbrotBigFloat(f, 1)
	}
}

func BenchmarkGenMandelbrotBigFloatz3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f, err := ioutil.TempFile("", "1024x1024*z3.png")
		defer os.Remove(f.Name())
		if err != nil {
			log.Fatalf("%s\n", err)
		}
		GenMandelbrotBigFloat(f, 3)
	}
}
