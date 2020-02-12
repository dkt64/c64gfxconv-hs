package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ErrCheck - obsługa błedów
// ================================================================================================
func ErrCheck(errNr error) bool {
	if errNr != nil {
		fmt.Println(errNr)
		return false
	}
	return true
}

// // colorDistance - calculate color distance
// // ================================================================================================
// func colorDistance(c color.Color, pal []int) int {
// 	return c.RGBA()
// }

// main - program entry point
// ================================================================================================
func main() {

	fmt.Println("===================================================================================")
	fmt.Println("= c64gfxconv-hs - C64 PNG -> Hires+Sprites gfx converter / Samar Productions 2020 =")
	fmt.Println("===================================================================================")

	inputFilename := flag.String("i", "", "input filename")

	flag.Parse()

	if *inputFilename != "" {

		fmt.Println("Loading file: ", *inputFilename)

		// Read image from file that already exists
		inputImageFile, err := os.Open(*inputFilename)
		ErrCheck(err)
		defer inputImageFile.Close()

		// Calling the generic image.Decode() will tell give us the data
		// and type of image it is as a string. We expect "png"
		inputImage, inputFormat, err := image.Decode(inputImageFile)
		ErrCheck(err)

		if inputFormat != "png" {
			log.Fatalln("Only png images are supported")
		}

		size := inputImage.Bounds().Size()

		fmt.Println("Image size :  ", size.X, "x", size.Y)

		if size.X != 320 && size.Y != 200 {
			log.Fatalln("Only 320 x 200 pixel image size is supported")
		}

		// paletteC64 := []int{0x000000, 0xffffff, 0x68372B, 0x70A4B2, 0x6F3D86, 0x588D43, 0x352879, 0xB8C76F, 0x6F4F25, 0x433900, 0x9A6759, 0x444444, 0x6C6C6C, 0x9AD284, 0x6C5EB5, 0x959595}

		rect := image.Rect(0, 0, size.X, size.Y)
		wImg := image.NewRGBA(rect)

		// loop though all the x
		for x := 0; x < size.X; x++ {
			// and now loop thorough all of this x's y
			for y := 0; y < size.Y; y++ {
				pixel := inputImage.At(x, y)
				originalColor := color.RGBAModel.Convert(pixel).(color.RGBA)

				// Offset colors a little, adjust it to your taste
				r := float64(originalColor.R) * 0.92126
				g := float64(originalColor.G) * 0.97152
				b := float64(originalColor.B) * 0.90722
				// average
				grey := uint8((r + g + b) / 3)
				c := color.RGBA{
					R: grey, G: grey, B: grey, A: originalColor.A,
				}
				wImg.Set(x, y, c)
			}
		}

		ext := filepath.Ext(*inputFilename)
		name := strings.TrimSuffix(filepath.Base(*inputFilename), ext)
		newImagePath := fmt.Sprintf("%s_gray%s", name, ext)
		fg, err := os.Create(newImagePath)
		defer fg.Close()
		ErrCheck(err)
		err = png.Encode(fg, wImg)
		ErrCheck(err)

		// colorDistance(inputImage.At(0, 0), paletteC64)
		// fmt.Println()

	} else {
		fmt.Println("Use -h flag to read usage")
	}
}
