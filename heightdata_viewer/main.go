package main
import (
	"fmt"
	"flag"
	"os"
	"strconv"
	"strings"
	//"math"
)

import "github.com/ungerik/go-cairo"

func atoi(it string) int{
	value, err := strconv.Atoi(it)
	if err != nil {
		fmt.Println("Error : heightdata_viewer : ", err)
		os.Exit(2)
	}
	return value
}

func main(){

	fmt.Println("heightdata_viewer : processing...")

	flag.Parse()
	argv := flag.Args()
    if len(argv) != 4 {
		fmt.Println("Error : heightdata_viewer : Invalid arguments.")
		os.Exit(2)
    }

    heightdata_file, err := os.Open(argv[0])
    if err != nil {
		fmt.Println("Error : heightdata_viewer : Failed to open file.")
		os.Exit(2)
    }
    defer heightdata_file.Close()

	image_pixel_w := atoi(argv[1])
	image_pixel_h := atoi(argv[2])
	data_digit := atoi(argv[3])

	var heightdata [][]int
	heightdata = make([][]int, image_pixel_h)
	for i := 0; i<image_pixel_h; i++ {
		heightdata[i] = make([]int, image_pixel_w)
	}

    buf := make([]byte, image_pixel_w*(data_digit+1))
	
    for row := 0; row < image_pixel_h; row++{
        n, err := heightdata_file.Read(buf)
        if n == 0 {
            break
        }
        if err != nil {
            panic(err)
        }
		
		column := 0
		slice := strings.Split(strings.Replace(string(buf), "\n", "", -1), ",")
		
		for _, it := range slice{
			if column >= image_pixel_w { break }
			if it == "" { break }
			value := atoi(it)
			heightdata[row][column] = value
			column++
		}
		
    }

	surface := cairo.NewSurface(cairo.FORMAT_ARGB32, image_pixel_w, image_pixel_h)
	for row := 0; row < image_pixel_h; row++{
		for column := 0; column < image_pixel_w; column++{
			drow := float64(row)
			dcolumn := float64(column)
			if heightdata[row][column] == 0 {
				continue
			} else {
				color := float64(heightdata[row][column])/2000

				surface.SetSourceRGB(0.6-color*0.2, 0.9-color*0.7, 0.4-color*0.4)
			}
			surface.Rectangle(dcolumn, drow, 2, 2)
			surface.Fill()
		}
	}
	surface.WriteToPNG("../view.png")
	surface.Finish()
}