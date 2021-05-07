package main
import (
	"fmt"
	"flag"
	"os"
	"../library"
	"math"
)

import "github.com/ungerik/go-cairo"


func main(){

	fmt.Println("data_output : processing...")

	flag.Parse()
	argv := flag.Args()
	
	var host library.SpHost
	
    if host.ApplyCommonArgument(argv) < 0 {
		fmt.Println("Error : data_output : Invalid arguments.")
		os.Exit(2)
    }

	fix_length := float64(host.Image_pixel_h)/1875.0

	//host.Heightdata := library.RequestHeightData(argv[0], host.Image_pixel_w, host.Image_pixel_h, host.Data_digit)
	//host.Citydata := library.RequestCityData(argv[4], host.Image_pixel_w, host.Image_pixel_h, host.Data_digit)
	pathdata := library.RequestPathData(argv[9], host.Image_pixel_w, host.Image_pixel_h, host.Data_digit)
	//host.Urbandata := library.RequestUrbanData(argv[10], host.Image_pixel_w, host.Image_pixel_h, host.Data_digit)

	surface := cairo.NewSurface(cairo.FORMAT_ARGB32, host.Image_pixel_w, host.Image_pixel_h)


	for row := 0; row < host.Image_pixel_h; row++{
		for column := 0; column < host.Image_pixel_w; column++{
			drow := float64(row)
			dcolumn := float64(column)
			if host.Heightdata[row][column] == 0 {
				continue
			} else {
				color := float64(host.Heightdata[row][column])/2000

				surface.SetSourceRGB(0.6-color*0.2, 0.9-color*0.7, 0.4-color*0.4)
			}
			surface.Rectangle(dcolumn, drow, 2, 2)
			surface.Fill()
	
		}
	}



	for row := 0; row < host.Image_pixel_h; row++{
		for column := 0; column < host.Image_pixel_w; column++{
			drow := float64(row)
			dcolumn := float64(column)
			if host.Urbandata[row][column] == -1 {
				continue
			}
			if host.Urbandata[row][column] == 1 {
				surface.SetSourceRGB(0.7,0.8,0.3)
			}
			if host.Urbandata[row][column] == 2 {
				surface.SetSourceRGB(1.0,0.8,0.2)
			}
			
			surface.Rectangle(dcolumn, drow, 1, 1)
			surface.Fill()
		}
	}



	for _, pp := range pathdata {
		surface.SetLineWidth(pp.Width*fix_length)
		surface.SetSourceRGB(pp.R, pp.G, pp.B)
		
		for i := 0; i<len(pp.Longitude); i++ {
			point_lg := library.GetXFromLongitude(pp.Longitude[i], host.Longitude_s, host.Longitude_e, host.Image_pixel_w)
			point_lt := library.GetYFromLatitude(pp.Latitude[i], host.Latitude_s, host.Latitude_e, host.Image_pixel_h)
			if i == 0{
				surface.MoveTo(point_lg, point_lt)
			}else{
				surface.LineTo(point_lg, point_lt)
			}
		}
		surface.Stroke()
	}

	for _, cp := range host.Citydata {
		point_lg := library.GetXFromLongitude(cp.Longitude, host.Longitude_s, host.Longitude_e, host.Image_pixel_w)
		point_lt := library.GetYFromLatitude(cp.Latitude, host.Latitude_s, host.Latitude_e, host.Image_pixel_h)
		surface.SetSourceRGB(host.Mark_r, host.Mark_g, host.Mark_b)
		mark_size := host.Mark_width*math.Sqrt(float64(cp.Population)/1000000.0)
		mark_size_f := mark_size*fix_length
		surface.Rectangle(point_lg-mark_size_f/2,
						  point_lt-mark_size_f/2,
						  mark_size_f, mark_size_f)
		surface.Fill()
	}
	
	surface.WriteToPNG("../"+host.Project_name+".png")
	surface.Finish()
}