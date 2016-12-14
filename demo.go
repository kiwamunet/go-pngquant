package main

import (
	"crypto/rand"
	"encoding/binary"
	"io/ioutil"
	"log"
	"strconv"

	"time"

	"github.com/kiwamunet/go-pngquant/binding"
)

const (
	usePram = structPram
	srcPath = "testdata/demo.png"
	dirPath = "testdata/"
)

type formatParams int

const (
	slicePram formatParams = iota
	stringPram
	structPram
)

func main() {
	var b []byte
	var err error

	b, err = getImageData()
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()
	log.Println("ZopfliPng Starting .....")

	switch usePram {
	case slicePram:
		b, err = sliceParam(b)
		if err != nil {
			log.Fatal(err)
		}
	case stringPram:
		b, err = stringParam(b)
		if err != nil {
			log.Fatal(err)
		}
	case structPram:
		b, err = structParam(b)
		if err != nil {
			log.Fatal(err)
		}
	}

	// output test
	err = binding.OutputFile(b, dirPath+random()+".png")
	if err != nil {
		log.Fatal(err)
	}
	elapsed0 := time.Since(start)
	log.Printf("elapsed time: %.3f secs", elapsed0.Seconds())
}

func sliceParam(src []byte) ([]byte, error) {
	strings := []string{"Pngquant", "256", "--speed", "3", "--quality", "0-100"}
	return binding.Pngquant(strings, src)
}

func stringParam(src []byte) ([]byte, error) {
	string := "Pngquant 256 --speed 3 --quality 0-100"
	return binding.PngquantOneLine(string, src)
}

func structParam(src []byte) ([]byte, error) {
	st := binding.PngquantParams{
		NumColors:  256,
		Speed:      3,
		QualityMin: 0,
		QualityMax: 100,
	}
	return binding.PngquantStruct(st, src)
}

func getImageData() (b []byte, e error) {
	b, e = ioutil.ReadFile(srcPath)
	return
}

func random() string {
	var n uint64
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	return strconv.FormatUint(n, 36)
}
