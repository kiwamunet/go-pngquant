package pngquant

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestPngquant(t *testing.T) {
	src, err := ioutil.ReadFile("testdata/demo.png")
	if err != nil {
		t.Fatalf("%v", err)
	}

	st := PngquantParams{
		NumColors:  256,
		Speed:      3,
		QualityMin: 0,
		QualityMax: 100,
	}

	img, err := PngquantStruct(st, src)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if err := ioutil.WriteFile("testdata/demo-result.png", img, os.ModePerm); err != nil {
		t.Fatalf("%v", err)
	}
}
