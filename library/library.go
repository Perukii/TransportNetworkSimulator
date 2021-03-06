package library

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"bufio"
	"io"
	"sort"
	"math"
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

type UrbanArea struct{
	Name string
	Population int
}

type City struct{
	Name string
	Longitude float64
	Latitude float64
	Population int
}

func Atoi(it string) int{
	if it == "" { return -1 }
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

func Itoa(it int) string{

	return fmt.Sprintf("%d", it)
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

	reader := bufio.NewReaderSize(citydata_file, 4096)

    for {
		buf, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error : library : Failed to read file.")
			os.Exit(2)
        }
		
		slice := strings.Split(string(buf), "\n")
		
		for _, it := range slice{
			itp := strings.Split(it, ",")
			if len(itp) < 4 { continue }
			var cp City	
			cp.Name = itp[0]
			cp.Longitude = Atof(itp[1])
			cp.Latitude = Atof(itp[2])
			cp.Population = Atoi(itp[3])
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

	reader := bufio.NewReaderSize(pathdata_file, 4096)

	psize := 0

    for {

        buf, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error : library : Failed to read file.")
			os.Exit(2)
        }
		
		slice := strings.Split(string(buf), "\n")
		
		
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
				if psize <= 0 { continue }
				pathdata[psize-1].Longitude = append(pathdata[psize-1].Longitude, Atof(itp[0]))
				pathdata[psize-1].Latitude =  append(pathdata[psize-1].Latitude,  Atof(itp[1]))

				
			}
		}
	}

	return pathdata
}

func RequestUrbanData(file string, image_pixel_w int, image_pixel_h int, data_digit int) [][]int{

    urbandata_file, err := os.Open(file)
    if err != nil {
		fmt.Println("Error : library : Failed to open file.")
		os.Exit(2)
    }
    defer urbandata_file.Close()

	urbandata := make([][]int, image_pixel_h)

	for i := 0; i<image_pixel_h; i++ {
		urbandata[i] = make([]int, image_pixel_w)
	}

	reader := bufio.NewReaderSize(urbandata_file, 3*image_pixel_w)

	row := 0

    for {

        buf, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error : library : Failed to read file.")
			os.Exit(2)
        }
		
		slice := strings.Split(string(buf), "\n")
		
		
		
		for _, it := range slice{
			if it == "" { continue }
			itp := strings.Split(it, ",")
			column := 0
			for _, pitp := range itp{
				urbandata[row][column] = Atoi(pitp)
				column++
			}
			row++
		}
		
	}

	return urbandata
}

func RequestUrbanAreaData(file string, image_pixel_w int, image_pixel_h int, data_digit int) []UrbanArea{

    areadata_file, err := os.Open(file)
    if err != nil {
		fmt.Println("Error : library : Failed to open file.")
		os.Exit(2)
    }
    defer areadata_file.Close()

	var areadata []UrbanArea

	reader := bufio.NewReaderSize(areadata_file, 4096)

    for {

        buf, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error : library : Failed to read file.")
			os.Exit(2)
        }
		
		slice := strings.Split(string(buf), "\n")

		var area UrbanArea
		
		for _, it := range slice{
			if it == "" { continue }
			itp := strings.Split(it, ",")
			area.Name = itp[0]
			area.Population = Atoi(itp[1])
		}
		areadata = append(areadata, area)
		
	}

	sort.Slice(areadata, func(i, j int) bool { return areadata[i].Population > areadata[j].Population })

	return areadata
}


func GetXFromLongitude(tar_longitude float64, longitude_s float64, longitude_e float64,
						image_pixel_w int)float64{

	return (tar_longitude-longitude_s)/(longitude_e-longitude_s)*float64(image_pixel_w)
}


func GetYFromLatitude(tar_latitude float64, latitude_s float64, latitude_e float64,
						image_pixel_h int)float64{
							
	var f = func(lt float64) float64{
		ltr := lt/180*3.1415
		return math.Log(math.Abs(math.Tan(3.1415/4+ltr/2)))
	}

	return (f(latitude_e)-f(tar_latitude))/(f(latitude_e)-f(latitude_s))*float64(image_pixel_h)

}