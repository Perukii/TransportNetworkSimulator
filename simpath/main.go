package main

import (
	"fmt"
	"flag"
	"os"
	"../library"

)


func main(){

	fmt.Println("simpath : processing...")

	flag.Parse()
	argv := flag.Args()
    if len(argv) != 10 {
		fmt.Println("Error : simpath : Invalid arguments.")
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

	host.Init()
	host.Init_writer(argv[9])
	
	path := host.Make_aster_path(host.Cityindex["鰺ヶ沢町"], host.Cityindex["むつ市"], 0.03, 2000, -1, true)
	host.Register_new_path(5.0, 1.0, 0.5, 0.3)
	for _, ptar := range path {
		host.Write_path_point(ptar.Longitude, ptar.Latitude)
	}
	host.Writer.Flush()
}