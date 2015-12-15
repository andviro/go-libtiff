package main

import (
	"fmt"
	"github.com/andviro/go-libtiff/libtiff"
	"image/png"
	"os"
)

func main() {
	tiff, err := libtiff.Open("test.tiff")
	if err != nil {
		panic(err)
	}
	defer tiff.Close()

	n := tiff.Iter(func(n int) {
		img, err := tiff.GetRGBA()
		if err != nil {
			panic(err)
		}

		w, _ := os.Create(fmt.Sprintf("page%d.png", n+1))
		png.Encode(w, &img)
	})
	fmt.Printf("Total pages: %d\n", n)
}
