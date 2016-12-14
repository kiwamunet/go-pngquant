# go-pngquant


Go go-pngquant is a Go bind to [pornel/pngquant](https://github.com/pornel/pngquant).


# Install

```
go get github.com/kiwamunet/go-pngquant
```

#### make static library

```
git clong git@github.com:kiwamunet/pngquant.git./configure --without-cocoamake ar rcs libpngquant.a *.o lib/*.ocp libpngquant.a vender/cp pngquant.h vender/
```

# Use

```
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
```

## Options

See `pngquant -h` for full list.

### `--quality min-max`

`min` and `max` are numbers in range 0 (worst) to 100 (perfect), similar to JPEG. pngquant will use the least amount of colors required to meet or exceed the `max` quality. If conversion results in quality below the `min` quality the image won't be saved (if outputting to stdin, 24-bit original will be output) and pngquant will exit with status code 99.

    pngquant --quality=65-80 image.png

### `--ext new.png`

Set custom extension (suffix) for output filename. By default `-or8.png` or `-fs8.png` is used. If you use `--ext=.png --force` options pngquant will overwrite input files in place (use with caution).

### `-o out.png` or `--output out.png`

Writes converted file to the given path. When this option is used only single input file is allowed.

### `--skip-if-larger`

Don't write converted files if the conversion isn't worth it.

### `--speed N`

Speed/quality trade-off from 1 (brute-force) to 11 (fastest). The default is 3. Speed 10 has 5% lower quality, but is 8 times faster than the default. Speed 11 disables dithering and lowers compression level.

### `--nofs`

Disables Floyd-Steinberg dithering.

### `--floyd=0.5`

Controls level of dithering (0 = none, 1 = full). Note that the `=` character is required.

### `--posterize bits`

Reduce precision of the palette by number of bits. Use when the image will be displayed on low-depth screens (e.g. 16-bit displays or compressed textures in ARGB444 format).

### `--strip`

Don't copy optional PNG chunks. Metadata is always removed on Mac (when using Cocoa reader).

### `--version`

Print version information to stdout.

### `-`

Read image from stdin and send result to stdout.

### `--`

Stops processing of arguments. This allows use of file names that start with `-`. If you're using pngquant in a script, it's advisable to put this before file names:

    pngquant $OPTIONS -- "$FILE"
    
## License

pngquant is dual-licensed:

* GPL v3 or later, and additional copyright notice must be kept for older parts of the code. See [COPYRIGHT](https://github.com/pornel/pngquant/blob/master/COPYRIGHT) for details.

* For use in non-GPL software (e.g. closed-source or App Store distribution) please ask kornel@pngquant.org for a commercial license.