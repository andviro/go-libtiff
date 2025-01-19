package libtiff

// #cgo LDFLAGS: -ltiff
// #include <stdlib.h>
// #include <tiffio.h>
//
//int getWH(TIFF *tif, int *w, int *h) {
//if (!TIFFGetField(tif, TIFFTAG_IMAGEWIDTH, w)) {
//return 1;
//}
//if (!TIFFGetField(tif, TIFFTAG_IMAGELENGTH, h)) {
//return 2;
//}
//return 0;
//}
//
//int getRes(TIFF *tif, float *x, float *y) {
//if (!TIFFGetField(tif, TIFFTAG_XRESOLUTION, x)) {
//return 1;
//}
//if (!TIFFGetField(tif, TIFFTAG_YRESOLUTION, y)) {
//return 2;
//}
//return 0;
//}
import "C"

import (
	"errors"
	"fmt"
	"image"
)

type Tiff struct {
	data *C.struct_tiff
}

func (t Tiff) Close() {
	C.TIFFClose(t.data)
}

func (t Tiff) Iter(cb func(int)) int {
	n := 0
	for {
		cb(n)
		n++
		if C.TIFFReadDirectory(t.data) == 0 {
			break
		}
	}
	return n
}

func (t Tiff) SetDir(n int) error {
	if int(C.TIFFSetDirectory(t.data, C.tdir_t(n))) != 1 {
		return errors.New("Invalid directory")
	}
	return nil
}

func (t Tiff) GetDPI() (float32, float32, error) {
	var x, y C.float
	if errCode := int(C.getRes(t.data, &x, &y)); errCode != 0 {
		return 0, 0, errors.New(fmt.Sprintf("Error getting image resolution: %d", errCode))
	}

	return float32(x), float32(y), nil
}

func (t Tiff) GetRGBA() (image.RGBA, error) {
	var w, h C.int
	var res image.RGBA

	if errCode := int(C.getWH(t.data, &w, &h)); errCode != 0 {
		return res, errors.New(fmt.Sprintf("Error getting image width/height: %d", errCode))
	}
	res.Rect = image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{int(w), int(h)}}
	res.Stride = res.Rect.Max.X * 4
	nBytes := C.size_t(res.Rect.Max.X * res.Rect.Max.Y * 4)
	data := C.malloc(nBytes)
	defer C.free(data)
	if r := int(C.TIFFReadRGBAImageOriented(t.data, C.uint32(res.Rect.Max.X), C.uint32(res.Rect.Max.Y), (*C.uint32)(data), C.ORIENTATION_TOPLEFT, C.int(0))); r == 0 {
		return res, errors.New("Error reading image data")
	}
	res.Pix = C.GoBytes(data, C.int(nBytes))
	return res, nil
}

func Open(name string) (Tiff, error) {
	var res Tiff
	res.data = C.TIFFOpen(C.CString(name), C.CString("r"))
	if res.data == nil {
		return res, errors.New("Error opening file")
	}
	return res, nil
}
