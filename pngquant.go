package pngquant

/*
#cgo pkg-config: libpng
#cgo CFLAGS: -I../internal/pngquant/lib
#cgo CFLAGS: -I../internal/pngquant
#cgo CFLAGS: -I../internal
#include "internal/pngquant/pngquant.h"
#include <stdio.h>
#include <stdlib.h>
*/
import "C"
import (
	"bytes"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"unsafe"

	"github.com/chrisfelesoid/cgostdio"
)

const pngquant = "pngquant"

type PngquantParams struct {
	NumColors  int
	Speed      int
	QualityMin int
	QualityMax int
}

func PngquantStruct(p PngquantParams, src []byte) ([]byte, error) {
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

func PngquantOneLine(str string, src []byte) ([]byte, error) {
	return Pngquant(strings.Split(strings.TrimSpace(str), " "), src)
}

func Pngquant(cmds []string, src []byte) ([]byte, error) {

	if len(cmds) == 0 {
		return src, errors.New("Nothing Parameters")
	}

	if strings.ToLower(cmds[0]) != pngquant {
		cmds = insertSlice(cmds, 0, pngquant)
	}

	for i, v := range cmds {
		if v == "-" {
			cmds = cmds[:i]
			break
		}
	}
	cmds = appendSlice(cmds, "-")

	argc := C.int(len(cmds))
	argv := make([]*C.char, int(argc))
	for i, arg := range cmds {
		argv[i] = C.CString(arg)
	}
	defer func() {
		for _, v := range argv {
			C.free(unsafe.Pointer(v))
		}
	}()

	f := cgostdio.NewStdinFromBuffer(src)
	defer f.Close()

	b, err := WriteCaptureWithCGo(func() error {
		res := C.pngquant(argc, (**C.char)(&argv[0]))
		if res != 0 {
			return errorType(res).sError()
		}
		return nil
	})
	if err != nil {
		return src, err
	}
	return b, nil
}

func WriteCaptureWithCGo(call func() error) ([]byte, error) {
	originalStdOut := os.Stdout
	originalCStdOut := C.stdout
	defer func() {
		os.Stdout, C.stdout = originalStdOut, originalCStdOut
	}()

	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	f := cgostdio.NewStdout(w)

	out := make(chan []byte)
	errch := make(chan error)
	go func() {
		var b bytes.Buffer

		_, err := io.Copy(&b, r)
		if err != nil {
			errch <- err
			return
		}
		out <- b.Bytes()
	}()

	if err = call(); err != nil {
		return nil, err
	}

	if err := f.Flush(); err != nil {
		return nil, err
	}
	if err := f.Close(); err != nil {
		return nil, err
	}

	select {
	case v := <-out:
		return v, nil
	case err := <-errch:
		return nil, err
	}
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
