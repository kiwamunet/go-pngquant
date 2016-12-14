package binding

/*
#cgo pkg-config: libpng
#cgo LDFLAGS: -L${SRCDIR}/../vendor -lpngquant
#include "../vendor/pngquant.h"
#include <stdio.h>
#include <stdlib.h>
#include "fmemopen.h"
*/
import "C"
import (
	"bytes"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

const pngquant = "pngquant"

type PngquantParams struct {
	NumColors  int
	Speed      int
	QualityMin int
	QualityMax int
}

func PngquantStruct(p PngquantParams, src []byte) (b []byte, err error) {
	cmds := []string{pngquant}

	if p.NumColors > 0 && p.NumColors <= 256 {
		cmds = appendSlice(cmds, strconv.Itoa(p.NumColors))
	}

	if p.Speed > 0 && p.Speed <= 10 {
		cmds = appendSlice(cmds, "--speed")
		cmds = appendSlice(cmds, strconv.Itoa(p.Speed))
	}

	if p.QualityMin >= 0 && p.QualityMin <= 100 && p.QualityMax >= 0 && p.QualityMax <= 100 {
		if p.QualityMin <= p.QualityMax {
			cmds = appendSlice(cmds, "--quality")
			cmds = appendSlice(cmds, strconv.Itoa(p.QualityMin)+"-"+strconv.Itoa(p.QualityMax))
		}
	}
	cmds = appendSlice(cmds, "-")
	return Pngquant(cmds, src)
}

func PngquantOneLine(str string, src []byte) (b []byte, err error) {
	return Pngquant(strings.Split(strings.TrimSpace(str), " "), src)
}

func Pngquant(cmds []string, src []byte) (b []byte, err error) {
	cmds = checkCommands(cmds)

	argc := C.int(len(cmds))
	argv := make([]*C.char, int(argc))
	for i, arg := range cmds {
		argv[i] = C.CString(arg)
	}

	rb := C.CString("rb")
	defer C.free(unsafe.Pointer(rb))

	// TODO: fmemopenを使わずにReadCaptureWithCGoにする
	C.stdin = C.fmemopen(unsafe.Pointer(&src[0]), C.size_t(len(src)), (*C.char)(unsafe.Pointer(rb)))

	b, err = WriteCaptureWithCGo(func() {
		res := C.pngquant(argc, (**C.char)(&argv[0]))
		if res != 0 {
			// TODO: ここの処理をちゃんとする
			log.Println(res)
		}
	})
	return b, err
}

func OutputFile(b []byte, path string) (e error) {
	file, err := os.Create(path)
	if err != nil {
		e = err
	}
	defer file.Close()
	file.Write(b)
	return
}

func WriteCaptureWithCGo(call func()) ([]byte, error) {
	originalStdErr, originalStdOut := os.Stderr, os.Stdout
	originalCStdErr, originalCStdOut := C.stderr, C.stdout
	defer func() {
		os.Stderr, os.Stdout = originalStdErr, originalStdOut
		C.stderr, C.stdout = originalCStdErr, originalCStdOut
	}()

	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	cw := C.CString("w")
	defer C.free(unsafe.Pointer(cw))

	f := C.fdopen((C.int)(w.Fd()), cw)

	os.Stderr, os.Stdout = w, w
	C.stderr, C.stdout = f, f

	out := make(chan []byte)
	go func() {
		var b bytes.Buffer

		_, err := io.Copy(&b, r)
		if err != nil {
			// TODO: panic 使わない
			panic(err)
		}

		out <- b.Bytes()
	}()

	call()

	C.fflush(f)

	err = w.Close()
	if err != nil {
		return nil, err
	}

	return <-out, err
}

func checkCommands(src []string) (dst []string) {
	dst = src
	if len(src) == 0 {
		return
	}

	if strings.ToLower(src[0]) != pngquant {
		dst = insertSlice(src, 0, pngquant)
	}

	for _, v := range dst {
		if v == "-" {
			return
		}
	}
	dst = appendSlice(dst, "-")
	return
}

func appendSlice(src []string, v string) []string {
	t := make([]string, len(src), len(src)+1)
	copy(t, src)
	return append(t, v)
}

func insertSlice(slice []string, position int, value string) []string {
	if position > len(slice) {
		return slice
	}

	newSlice := make([]string, position+1, len(slice)+1)
	copy(newSlice, slice[:position])
	newSlice[position] = value
	return append(newSlice, slice[position:]...)
}
