package pngquant

/*
#include "internal/pngquant/rwpng.h"
*/
import "C"
import "fmt"

type errorType int

const (
	SUCCESS                 errorType = C.SUCCESS
	MISSING_ARGUMENT        errorType = C.MISSING_ARGUMENT
	READ_ERROR              errorType = C.READ_ERROR
	INVALID_ARGUMENT        errorType = C.INVALID_ARGUMENT
	NOT_OVERWRITING_ERROR   errorType = C.NOT_OVERWRITING_ERROR
	CANT_WRITE_ERROR        errorType = C.CANT_WRITE_ERROR
	OUT_OF_MEMORY_ERROR     errorType = C.OUT_OF_MEMORY_ERROR
	WRONG_ARCHITECTURE      errorType = C.WRONG_ARCHITECTURE
	PNG_OUT_OF_MEMORY_ERROR errorType = C.PNG_OUT_OF_MEMORY_ERROR
	LIBPNG_FATAL_ERROR      errorType = C.LIBPNG_FATAL_ERROR
	WRONG_INPUT_COLOR_TYPE  errorType = C.WRONG_INPUT_COLOR_TYPE
	LIBPNG_INIT_ERROR       errorType = C.LIBPNG_INIT_ERROR
	TOO_LARGE_FILE          errorType = C.TOO_LARGE_FILE
	TOO_LOW_QUALITY         errorType = C.TOO_LOW_QUALITY
)

var exceptionTypeStrings = map[errorType]string{
	SUCCESS:                 "SUCCESS",
	MISSING_ARGUMENT:        "MISSING_ARGUMENT",
	READ_ERROR:              "READ_ERROR",
	INVALID_ARGUMENT:        "INVALID_ARGUMENT",
	NOT_OVERWRITING_ERROR:   "NOT_OVERWRITING_ERROR",
	CANT_WRITE_ERROR:        "CANT_WRITE_ERROR",
	OUT_OF_MEMORY_ERROR:     "OUT_OF_MEMORY_ERROR",
	WRONG_ARCHITECTURE:      "WRONG_ARCHITECTURE",
	PNG_OUT_OF_MEMORY_ERROR: "PNG_OUT_OF_MEMORY_ERROR",
	LIBPNG_FATAL_ERROR:      "LIBPNG_FATAL_ERROR",
	WRONG_INPUT_COLOR_TYPE:  "WRONG_INPUT_COLOR_TYPE",
	LIBPNG_INIT_ERROR:       "LIBPNG_INIT_ERROR",
	TOO_LARGE_FILE:          "TOO_LARGE_FILE",
	TOO_LOW_QUALITY:         "TOO_LOW_QUALITY",
}

type Error struct {
	code        errorType
	description string
}

func (err *Error) Error() string {
	return fmt.Sprintf("Err: %s [code=%d]", err.description, err.code)
}

func (et errorType) sError() error {
	if v, ok := exceptionTypeStrings[errorType(et)]; ok {
		return &Error{description: v, code: et}
	}
	return &Error{description: fmt.Sprintf("UnknownError[%d]", et), code: et}

}
