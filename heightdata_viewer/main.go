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

type city struct{
	name string
	longitude float64
	latitude float64
}

func atoi(it string) int{
	value, err := strconv.Atoi(it)
	if err != nil {
		fmt.Println("Error : heightdata_viewer : ", err)
		os.Exit(2)
	}
	return value
}


func atof(it string) float64{
	value, err := strconv.ParseFloat(it, 64)
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
    if len(argv) != 9 {
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

    citydata_file, err := os.Open(argv[4])
    if err != nil {
		fmt.Println("Error : heightdata_viewer : Failed to open file.")
		os.Exit(2)
    }
    defer citydata_file.Close()
	
	longitude_s := atof(argv[5])
	longitude_e := atof(argv[6])
	latitude_s := atof(argv[7])
	latitude_e := atof(argv[8])
	

	var heightdata [][]int
	var citydata []city	
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
			fmt.Println("Error : heightdata_viewer: Failed to read file.")
			os.Exit(2)
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

	cuf := make([]byte, 255)

    for {
        n, err := citydata_file.Read(cuf)
        if n == 0 {
            break
        }
        if err != nil {
			fmt.Println("Error : heightdata_viewer : Failed to read file.")
			os.Exit(2)
        }
		
		slice := strings.Split(string(cuf), "\n")
		
		for _, it := range slice{
			itp := strings.Split(it, ",")
			if len(itp) < 3 { continue }
			var cp city
			cp.name = itp[0]
			cp.longitude = atof(itp[1])
			cp.latitude = atof(itp[2])
			citydata = append(citydata, cp)
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
	for _, cp := range citydata {
		point_lg := (cp.longitude-longitude_s)/(longitude_e-longitude_s)
		point_lt := 1.0-(cp.latitude-latitude_s)/(latitude_e-latitude_s)
		surface.SetSourceRGB(0.9, 0.2, 0.2)
		surface.Rectangle(float64(image_pixel_w)*point_lg, float64(image_pixel_h)*point_lt, 20, 20)
		surface.Fill()
	}
	surface.WriteToPNG("../view.png")
	surface.Finish()
}