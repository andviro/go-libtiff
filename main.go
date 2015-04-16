package main

import (
	"fmt"
	"gitlab.com/Rodionoff/go-libtiff/libtiff"
	"image/png"
	"os"
)

func main() {
	tiff, err := libtiff.Open("test.tiff")
	if err != nil {
		panic(err)
	}
	defer tiff.Close()

	i := 0
	tiff.Iter(func() {
		img, err := tiff.GetRGBA()
		if err != nil {
			panic(err)
		}

		w, _ := os.Create(fmt.Sprintf("ttt%d.png", i))
		png.Encode(w, &img)
		i++
	})

}
