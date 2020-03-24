package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	_ "image/png"
	"log"
	"os"
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

// main - program entry point
// ================================================================================================
func main() {

	fmt.Println("===================================================================================")
	fmt.Println("= c64gfxconv-hs - C64 PNG -> Hires+Sprites gfx converter / Samar Productions 2020 =")
	fmt.Println("===================================================================================")

	inputFilename1 := flag.String("u", "", "Unrestricted hires PNG filename")
	inputFilename2 := flag.String("h", "", "Hires PNG filename")
	outFilename := flag.String("o", "", "Difference PNG filename")
	backgroundColorString := flag.String("b", "", "Background color in RRGGBB hex format")

	flag.Parse()

	backColor, err := hex.DecodeString(*backgroundColorString)
	ErrCheck(err)

	fmt.Printf("Background color is R=%x G=%x B=%x\n", backColor[0], backColor[1], backColor[2])

	if *inputFilename1 != "" && *inputFilename2 != "" && *outFilename != "" {

		fmt.Println("Loading files:  ", *inputFilename1, *inputFilename2)

		inputImageFile1, err := os.Open(*inputFilename1)
		ErrCheck(err)
		defer inputImageFile1.Close()
		inputImageFile2, err := os.Open(*inputFilename2)
		ErrCheck(err)
		defer inputImageFile2.Close()

		imgUnres, inputFormat1, err := image.Decode(inputImageFile1)
		ErrCheck(err)
		imgHires, inputFormat2, err := image.Decode(inputImageFile2)
		ErrCheck(err)

		if inputFormat1 != "png" || inputFormat2 != "png" {
			log.Fatalln("Only png images are supported")
		}

		size1 := imgUnres.Bounds().Size()
		size2 := imgHires.Bounds().Size()

		fmt.Println("Image 1 size :  ", size1.X, "x", size1.Y)
		fmt.Println("Image 2 size :  ", size2.X, "x", size2.Y)

		if size1.X != 320 && size1.Y != 200 || size2.X != 320 && size2.Y != 200 {
			log.Fatalln("Only 320 x 200 pixel image size is supported")
		}

		// paletteC64 := []int{0x000000, 0xffffff, 0x68372B, 0x70A4B2, 0x6F3D86, 0x588D43, 0x352879, 0xB8C76F, 0x6F4F25, 0x433900, 0x9A6759, 0x444444, 0x6C6C6C, 0x9AD284, 0x6C5EB5, 0x959595}

		rect := image.Rect(0, 0, 320, 200)
		wImg := image.NewRGBA(rect)

		// Kolor nie występujacy w palecie dla pokazania miejsc gdzie musi być uzupełnione tło
		//
		indepColor := color.RGBA{0x20, 0x20, 0, 0xff}

		// loop though all the x
		for x := 0; x < 320; x++ {
			// and now loop thorough all of this x's y
			for y := 0; y < 200; y++ {
				pixelUnres := imgUnres.At(x, y)
				pixelUnresColor := color.RGBAModel.Convert(pixelUnres).(color.RGBA)
				pixelHires := imgHires.At(x, y)
				pixelHiresColor := color.RGBAModel.Convert(pixelHires).(color.RGBA)

				if pixelUnresColor != pixelHiresColor {
					wImg.Set(x, y, pixelUnresColor)
				} else {
					wImg.Set(x, y, indepColor)
				}
			}
		}

		fg, err := os.Create(*outFilename)
		defer fg.Close()
		ErrCheck(err)
		err = png.Encode(fg, wImg)
		ErrCheck(err)

	} else if *inputFilename1 != "" && *inputFilename2 == "" && *outFilename == "" {

		// ===================================================================================
		// Razem z konwersją
		// ===================================================================================
		//
		fmt.Println("Loading file:  ", *inputFilename1)

		inputImageFile1, err := os.Open(*inputFilename1)
		ErrCheck(err)
		defer inputImageFile1.Close()

		img, imgFormat, err := image.Decode(inputImageFile1)
		ErrCheck(err)

		if imgFormat != "png" {
			log.Fatalln("Only png images are supported")
		}

		size1 := img.Bounds().Size()

		fmt.Println("Image size :  ", size1.X, "x", size1.Y)

		if size1.X != 320 && size1.Y != 200 {
			log.Fatalln("Only 320 x 200 pixel image size is supported")
		}

		// paletteC64 := []int{0x000000, 0xffffff, 0x68372B, 0x70A4B2, 0x6F3D86, 0x588D43, 0x352879, 0xB8C76F, 0x6F4F25, 0x433900, 0x9A6759, 0x444444, 0x6C6C6C, 0x9AD284, 0x6C5EB5, 0x959595}

		// rect := image.Rect(0, 0, 320, 200)
		// wImg := image.NewRGBA(rect)

		// Kolor nie występujacy w palecie dla pokazania miejsc gdzie musi być uzupełnione tło
		//
		// indepColor := color.RGBA{0x20, 0x20, 0, 0xff}

		var colors []color.RGBA

		// loop though all the x
		for x := 0; x < 320; x++ {
			// and now loop thorough all of this x's y
			for y := 0; y < 200; y++ {
				pixel := img.At(x, y)
				pixelColor := color.RGBAModel.Convert(pixel).(color.RGBA)

				found := false
				for _, col := range colors {
					if col == pixelColor {
						found = true
					}
				}
				if !found && pixelColor.R != backColor[0] && pixelColor.G != backColor[1] && pixelColor.B != backColor[2] {
					colors = append(colors, pixelColor)
				}
			}
		}

		fmt.Println("Amount of colors = ", len(colors))
	} else {
		fmt.Println("Use -h flag to read usage")
	}
}
