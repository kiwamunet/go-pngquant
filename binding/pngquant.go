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
	"errors"
	"io"
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

	rb := C.CString("rb")
	defer C.free(unsafe.Pointer(rb))

	C.stdin = C.fmemopen(unsafe.Pointer(&src[0]), C.size_t(len(src)), (*C.char)(unsafe.Pointer(rb)))

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

func OutputFile(b []byte, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(b)
	return nil
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

	cw := C.CString("w")
	defer C.free(unsafe.Pointer(cw))

	f := C.fdopen((C.int)(w.Fd()), cw)
	os.Stdout, C.stdout = w, f

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

	C.fflush(f)

	err = w.Close()
	if err != nil {
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
