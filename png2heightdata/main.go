package main

import (
	"fmt"
	image "image/png"
	"os"
	"../library"
	"flag"
	"bufio"
	"math"
)

func main(){
	fmt.Println("png2heightdata : processing...")

	flag.Parse()
	argv := flag.Args()
    if len(argv) != 6 {
		fmt.Println("Error : png2heightdata : Invalid arguments.")
		os.Exit(2)
    }

	height_diff := library.Atoi(argv[2])
	data_digit := library.Atoi(argv[5])
	
    png_file, err := os.Open(argv[0])
    if err != nil {
		fmt.Println("Error : png2heightdata : Failed to open file.")
		os.Exit(2)
    }

    heightdata_file, err := os.Create(argv[1])
    if err != nil {
		fmt.Println("Error : png2heightdata : Failed to create file.")
		os.Exit(2)
    }

    img, err := image.Decode(png_file)
    if err != nil {
		fmt.Println("Error : png2heightdata : Failed to road file.")
		os.Exit(2)
    }

	bounds := img.Bounds()

	writer := bufio.NewWriter(heightdata_file)

	get_height := func(x, y int) float64 {
		r, g, b, _ := img.At(x, y).RGBA()

		if r == 65535 && g == 65535 && b == 65535 {
			return 0
		} else {
			r /= 256
			g /= 256
			b /= 256
			return float64(int((255-g)+r+b+1)*height_diff)
		}
	}

	bias := 30.0

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		line := ""
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			target := get_height(x, y)
			if x >= bounds.Min.X+1 && x < bounds.Max.X-1 {
				lf, rg := get_height(x-1, y), get_height(x+1, y)
				if target > lf && target > rg{
					if math.Max(math.Abs(target-lf), math.Abs(target-rg)) > bias{
						target = math.Max(lf, rg)
					}
				}
			}
			if y >= bounds.Min.Y+1 && y < bounds.Max.Y-1 {
				up, dw := get_height(x, y-1), get_height(x, y+1)
				if target > up && target > dw{
					if math.Max(math.Abs(target-up), math.Abs(target-dw)) > bias{
						target = math.Max(up, dw)
					}
				}
			}
			snum := library.Itoa(int(target))

			for {
				if len(snum) == data_digit{
					break
				} else if len(snum) > data_digit {
					fmt.Println("Error : png2heightdata : Too few digit space to output data.")
					os.Exit(2)
				} else {
					snum = "0" + snum
				}
			}

			line += snum
			if x != bounds.Max.X-1 { line += "," }
			
			
		}

		line += "\n"

		if _, err := writer.Write([]byte(line)); err != nil {
			fmt.Println("Error : png2heightdata : Failed to write file.")
			os.Exit(2)
		}

		if y%100 == 0 {
			fmt.Println("png2heightdata : creating heightdata... ("+library.Itoa(y)+"x"+library.Itoa(bounds.Max.X)+")")
		}
		
	}

	writer.Flush()

}