package library

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type MapInfo struct{
	image_pixel_w int
	image_pixel_h int
	data_digit int

	longitude_s float64
	longitude_e float64
	latitude_s float64
	latitude_e  float64
}

type City struct{
	Name string
	Longitude float64
	Latitude float64
}

func Atoi(it string) int{
	value, err := strconv.Atoi(it)
	if err != nil {
		fmt.Println("Error : library : ", err)
		os.Exit(2)
	}
	return value
}

func Atof(it string) float64{
	value, err := strconv.ParseFloat(it, 64)
	if err != nil {
		fmt.Println("Error : library : ", err)
		os.Exit(2)
	}
	return value
}

func RequestHeightData(file string, image_pixel_w int, image_pixel_h int, data_digit int) [][]int{

    heightdata_file, err := os.Open(file)
    if err != nil {
		fmt.Println("Error : library : Failed to open file.")
		os.Exit(2)
    }
    defer heightdata_file.Close()

	heightdata := make([][]int, image_pixel_h)

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
			fmt.Println("Error : library : Failed to read file.")
			os.Exit(2)
        }
		
		column := 0
		slice := strings.Split(strings.Replace(string(buf), "\n", "", -1), ",")
		
		for _, it := range slice{
			if column >= image_pixel_w { break }
			if it == "" { break }
			value := Atoi(it)
			heightdata[row][column] = value
			column++
		}
    }
	return heightdata
}

func RequestCityData(file string, image_pixel_w int, image_pixel_h int, data_digit int) []City{
	citydata_file, err := os.Open(file)
    if err != nil {
		fmt.Println("Error : library : Failed to open file.")
		os.Exit(2)
    }
    defer citydata_file.Close()

	var citydata []City

	cuf := make([]byte, 255)

    for {
        n, err := citydata_file.Read(cuf)
        if n == 0 {
            break
        }
        if err != nil {
			fmt.Println("Error : library : Failed to read file.")
			os.Exit(2)
        }
		
		slice := strings.Split(string(cuf), "\n")
		
		for _, it := range slice{
			itp := strings.Split(it, ",")
			if len(itp) < 3 { continue }
			var cp City	
			cp.Name = itp[0]
			cp.Longitude = Atof(itp[1])
			cp.Latitude = Atof(itp[2])
			citydata = append(citydata, cp)
		}
	}

	return citydata
}