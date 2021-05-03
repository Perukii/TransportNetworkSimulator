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
    if len(argv) != 10 {
		fmt.Println("Error : data_output : Invalid arguments.")
		os.Exit(2)
    }


	image_pixel_w := library.Atoi(argv[1])
	image_pixel_h := library.Atoi(argv[2])
	data_digit := library.Atoi(argv[3])
	longitude_s := library.Atof(argv[5])
	longitude_e := library.Atof(argv[6])
	latitude_s := library.Atof(argv[7])
	latitude_e := library.Atof(argv[8])

	heightdata := library.RequestHeightData(argv[0], image_pixel_w, image_pixel_h, data_digit)
	citydata := library.RequestCityData(argv[4], image_pixel_w, image_pixel_h, data_digit)
	pathdata := library.RequestPathData(argv[9], image_pixel_w, image_pixel_h, data_digit)

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


	for _, pp := range pathdata {
		surface.SetLineWidth(pp.Width)
		surface.SetSourceRGB(pp.R, pp.G, pp.B)
		
		for i := 0; i<len(pp.Longitude); i++ {
			point_lg := library.GetXFromLongitude(pp.Longitude[i], longitude_s, longitude_e, image_pixel_w)
			point_lt := library.GetYFromLatitude(pp.Latitude[i], latitude_s, latitude_e, image_pixel_h)
			if i == 0{
				surface.MoveTo(point_lg, point_lt)
			}else{
				surface.LineTo(point_lg, point_lt)
			}
		}
		surface.Stroke()
	}
	
	for _, cp := range citydata {
		point_lg := library.GetXFromLongitude(cp.Longitude, longitude_s, longitude_e, image_pixel_w)
		point_lt := library.GetYFromLatitude(cp.Latitude, latitude_s, latitude_e, image_pixel_h)
		surface.SetSourceRGB(0.9, 0.2, 0.2)
		mark_size := 30.0*math.Sqrt(float64(cp.Population)/1000000.0)
		surface.Rectangle(point_lg-mark_size/2,
						  point_lt-mark_size/2,
						  mark_size, mark_size)
		surface.Fill()
	}
	surface.WriteToPNG("../view.png")
	surface.Finish()
}