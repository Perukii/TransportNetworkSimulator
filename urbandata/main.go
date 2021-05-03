package main

import (
	"fmt"
	"flag"
	"os"
	"../library"
	"bufio"
)

func main(){
	fmt.Println("urbandata : processing...")

	flag.Parse()
	argv := flag.Args()
    if len(argv) != 11 {
		fmt.Println("Error : urbandata : Invalid arguments.")
		os.Exit(2)
    }

	var host library.SpHost

	host.Image_pixel_w = library.Atoi(argv[1])
	host.Image_pixel_h = library.Atoi(argv[2])
	host.Data_digit = library.Atoi(argv[3])
	host.Longitude_s = library.Atof(argv[5])
	host.Longitude_e = library.Atof(argv[6])
	host.Latitude_s = library.Atof(argv[7])
	host.Latitude_e = library.Atof(argv[8])
	host.LgLt_ratio = ((host.Longitude_s-host.Longitude_e)/float64(host.Image_pixel_w))/((host.Latitude_s-host.Latitude_e)/float64(host.Image_pixel_h))
	host.Heightdata = library.RequestHeightData(argv[0], host.Image_pixel_w, host.Image_pixel_h, host.Data_digit)
	host.Citydata = library.RequestCityData(argv[4], host.Image_pixel_w, host.Image_pixel_h, host.Data_digit)
	host.Cityindex = make(map[string]int)

	
	pdensity := 1000
	host.Init()


	urbandata_file, err := os.Create(argv[10])
    if err != nil {
		fmt.Println("Error : urbandata : Failed to create file.")
		os.Exit(2)
    }

	host.Writer = bufio.NewWriter(urbandata_file)

	var urbandata [][]int
	urbandata = make([][]int, host.Image_pixel_h)
	for i := 0; i < host.Image_pixel_h; i++{
		urbandata[i] = make([]int, host.Image_pixel_w)
	}
	
	_ = pdensity
	for i := 1; i < len(host.Citydata); i++ {
		//if host.Citydata[i].Population == 0 { continue }
		path := host.Make_aster_path(i, i, 0.002, 300, host.Citydata[i].Population/pdensity, false)
		for _, ptar := range path {
			//fmt.Println(ptar.Longitude, ptar.Latitude)
			yad := int(library.GetYFromLatitude(ptar.Latitude, host.Latitude_s, host.Latitude_e, host.Image_pixel_h))
			xad := int(library.GetXFromLongitude(ptar.Longitude, host.Longitude_s, host.Longitude_e, host.Image_pixel_w))
			
			if yad < 0 || yad >= host.Image_pixel_h || xad < 0 || xad >= host.Image_pixel_w{
				continue
			}
			
			urbandata[yad][xad] = 1
		}
	}

	for i := 0; i < host.Image_pixel_h; i++{
		line := ""
		for j := 0; j < host.Image_pixel_w; j++{
			line += library.Itoa(urbandata[i][j])+","
		}
		line = line[0:len(line)-1]
		host.Write_line(line+"\n")
	}
	host.Writer.Flush()



}