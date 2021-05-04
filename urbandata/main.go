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
    if len(argv) != 12 {
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

	
	pdensity := 50
	cdensity := 200
	host.Init()


	urbandata_file, err := os.Create(argv[10])
    if err != nil {
		fmt.Println("Error : urbandata : Failed to create file.")
		os.Exit(2)
    }

	areadata_file, err := os.Create(argv[11])
    if err != nil {
		fmt.Println("Error : urbandata : Failed to create file.")
		os.Exit(2)
    }

	

	var urbandata [][]int
	urbandata = make([][]int, host.Image_pixel_h)
	for i := 0; i < host.Image_pixel_h; i++{
		urbandata[i] = make([]int, host.Image_pixel_w)
		for j := 0; j < host.Image_pixel_w; j++{
			urbandata[i][j] = -1
		}
	}

	type Root struct{
		Root int
		Population int
	}
	
	var urbanroot []Root
	urbanroot = make([]Root, len(host.Citydata))
	for i := 0; i < len(host.Citydata); i++{
		urbanroot[i].Root = i
		urbanroot[i].Population = host.Citydata[i].Population
	}

	height_weight := 0.02
	dist_weight := 1000.0
	urban_weight:= 0.0
	pitv := 0.0022
	sea_weight := 1000000.0
	
	for n := 0; n<2; n++{

		for i := 0; i < len(host.Citydata); i++ {
			var density int
			density = pdensity
			if n == 1 { density = cdensity }



			path, _ := host.Make_aster_path(i, i, pitv, height_weight, dist_weight, urban_weight, sea_weight, host.Citydata[i].Population/density, false)
			for _, ptar := range path {
				yad := int(library.GetYFromLatitude(ptar.Latitude, host.Latitude_s, host.Latitude_e, host.Image_pixel_h))
				xad := int(library.GetXFromLongitude(ptar.Longitude, host.Longitude_s, host.Longitude_e, host.Image_pixel_w))
				
				if yad < 0 || yad >= host.Image_pixel_h || xad < 0 || xad >= host.Image_pixel_w{
					continue
				}
				if n == 0{
					cmp := urbandata[yad][xad]
					if cmp >= 0{
						if host.Citydata[cmp].Population > host.Citydata[i].Population{
							urbanroot[i].Root = cmp
							urbanroot[cmp].Population += urbanroot[i].Population
							urbanroot[i].Population = 0
						} else {
							urbanroot[cmp].Root = i
							urbanroot[i].Population += urbanroot[cmp].Population
							urbanroot[cmp].Population = 0
						}
					}
					urbandata[yad][xad] = i
				} else {
					urbandata[yad][xad] = i + len(host.Citydata)
				}
			}
		}

		if n == 0 {
			host.Writer = bufio.NewWriter(areadata_file)
			for i := 0; i < len(host.Citydata); i++{
				if urbanroot[i].Root == i {
					//fmt.Println(host.Citydata[i].Name, urbanroot[i].Population)
					line := host.Citydata[i].Name + "," + library.Itoa(urbanroot[i].Population)
					host.Write_line(line+"\n")
				}
			}
			host.Writer.Flush()
		} else {
			host.Writer = bufio.NewWriter(urbandata_file)
			for i := 0; i < host.Image_pixel_h; i++{
				line := ""
				for j := 0; j < host.Image_pixel_w; j++{
					if urbandata[i][j] >= len(host.Citydata){
						line += "2"
					} else if urbandata[i][j] >= 0 {
						line += "1"
					}
					line += ","
					
				}
				line = line[0:len(line)-1]
				host.Write_line(line+"\n")
			}
			host.Writer.Flush()
		}
	}


}