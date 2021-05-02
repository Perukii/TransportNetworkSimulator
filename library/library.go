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

type Path struct{
	Width float64
	R float64
	G float64
	B float64
	Longitude []float64
	Latitude  []float64
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

func Ftoa(it float64) string{
	return fmt.Sprintf("%f", it)
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

	buf := make([]byte, 255)

    for {
        n, err := citydata_file.Read(buf)
        if n == 0 {
            break
        }
        if err != nil {
			fmt.Println("Error : library : Failed to read file.")
			os.Exit(2)
        }
		
		slice := strings.Split(string(buf), "\n")
		
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

func RequestPathData(file string, image_pixel_w int, image_pixel_h int, data_digit int) []Path{
	pathdata_file, err := os.Open(file)
    if err != nil {
		fmt.Println("Error : library : Failed to open file.")
		os.Exit(2)
    }
    defer pathdata_file.Close()

	var pathdata []Path

	buf := make([]byte, 255)

    for {
        n, err := pathdata_file.Read(buf)
        if n == 0 {
            break
        }
        if err != nil {
			fmt.Println("Error : library : Failed to read file.")
			os.Exit(2)
        }
		
		slice := strings.Split(string(buf), "\n")
		psize := 0
		
		for _, it := range slice{
			if it == "" { continue }
			itp := strings.Split(it, ",")
			if itp[0][0] == '#'{
				if len(itp) != 5 { continue }
				
				var pp Path
				pp.Width = Atof(itp[1])
				pp.R = Atof(itp[2])
				pp.G = Atof(itp[3])
				pp.B = Atof(itp[4])
				pathdata = append(pathdata, pp)
				psize++

			} else {
				if len(itp) != 2 { continue }
				pathdata[psize-1].Longitude = append(pathdata[psize-1].Longitude, Atof(itp[0]))
				pathdata[psize-1].Latitude =  append(pathdata[psize-1].Latitude,  Atof(itp[1]))
			}
		}
	}

	return pathdata
}


func GetXFromLongitude(tar_longitude float64, longitude_s float64, longitude_e float64,
						image_pixel_w int)float64{
	return (tar_longitude-longitude_s)/(longitude_e-longitude_s)*float64(image_pixel_w)
}

func GetYFromLatitude(tar_latitude float64, latitude_s float64, latitude_e float64,
						image_pixel_h int)float64{
	return (1.0-(tar_latitude-latitude_s)/(latitude_e-latitude_s))*float64(image_pixel_h)
}