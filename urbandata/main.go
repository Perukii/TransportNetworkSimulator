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
	var host library.SpHost

    if host.ApplyCommonArgument(argv) < 0 {
		fmt.Println("Error : urbandata : Invalid arguments.")
		os.Exit(2)
    }

	host.Height_difference_score = host.Urban_area_height_difference_score

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

	var get_root = func(i int)int{
		res := i
		for {
			if res < 0 { return -1 }
			if urbanroot[res].Root == res{
				return res
			} else {
				res = urbanroot[res].Root
			}
		}
	}

	for n := 0; n<2; n++{

		for i := 0; i < len(host.Citydata); i++ {
			var density int
			density = host.Urban_wide_area_density
			if n == 1 { density = host.Urban_central_area_density }

			path, _ := host.Make_aster_path(i, i, host.Urban_area_interval, host.Citydata[i].Population/density, false)
			for _, ptar := range path {
				yad := int(library.GetYFromLatitude(ptar.Latitude, host.Latitude_s, host.Latitude_e, host.Image_pixel_h))
				xad := int(library.GetXFromLongitude(ptar.Longitude, host.Longitude_s, host.Longitude_e, host.Image_pixel_w))
				data := -1
				if yad < 0 || yad >= host.Image_pixel_h || xad < 0 || xad >= host.Image_pixel_w{
					continue
				}
				if n == 0{
					cmp := urbandata[yad][xad]
				
					if get_root(i) != get_root(cmp) && cmp >= 0{
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
					data = i

					
				} else {
					data = i + len(host.Citydata)
				}
				
				urbandata[yad][xad] = data
				/*
				if xad >= 1 { urbandata[yad][xad-1] = data }
				if yad >= 1 { urbandata[yad-1][xad] = data }
				if xad < host.Image_pixel_w-1 { urbandata[yad][xad+1] = data }
				if yad < host.Image_pixel_h-1 { urbandata[yad+1][xad] = data }
				*/

			}
		}

		
		for i := 0; i < host.Image_pixel_h; i++{
			for j := 0; j < host.Image_pixel_w; j++{
				var void int
				if n == 0 {
					void = 0
				} else {
					void = len(host.Citydata)
				}
				if urbandata[i][j] >= void {
					data := urbandata[i][j]
					if j >= 1 && urbandata[i][j-1] < void { urbandata[i][j-1] = data }
					if i >= 1 && urbandata[i-1][j] < void { urbandata[i-1][j] = data }
					//if j < host.Image_pixel_w-1 && urbandata[i][j+1] == -1{ urbandata[i][j+1] = data }
					//if i < host.Image_pixel_h-1 && urbandata[i+1][j] == -1{ urbandata[i+1][j] = data }
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