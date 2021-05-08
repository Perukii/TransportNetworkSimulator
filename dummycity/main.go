package main
import (
	"fmt"
	"../library"
	"flag"
	"os"
	"bufio"
)

func main(){
	fmt.Println("dummycity : processing...")

	var host library.SpHost
	
	flag.Parse()
	argv := flag.Args()

    if host.ApplyCommonArgument(argv) < 0 {
		fmt.Println("Error : dummycity : Invalid arguments.", len(argv))
		os.Exit(2)
    }

    citydata_file, err := os.OpenFile(argv[4], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
		fmt.Println("Error : dummycity : Failed to create file.")
		os.Exit(2)
    }

	writer := bufio.NewWriter(citydata_file)

	//y_interval := host.Dummy_city_interval_latitude/host.Latitude_per_pixel

	var get_x = func(lg float64) int{
		return int(library.GetXFromLongitude(lg, host.Longitude_s, host.Longitude_e, host.Image_pixel_w))
	}

	var get_y = func(lt float64) int{
		return int(library.GetYFromLatitude(lt, host.Latitude_s, host.Latitude_e, host.Image_pixel_h))
	}
	dummy_longitude_interval := host.Dummy_city_interval_latitude*host.LgLt_ratio

	for lt := host.Latitude_s; lt<host.Latitude_e; lt += host.Dummy_city_interval_latitude{
		for lg := host.Longitude_s; lg<host.Longitude_e; lg += dummy_longitude_interval{

			yb := get_y(lt+host.Dummy_city_interval_latitude)
			yf := get_y(lt)
			
			xb := get_x(lg)
			xf := get_x(lg+dummy_longitude_interval)
			
			min_height := host.Dummy_city_max_height
			lgn := lg
			ltn := lt
			
			for y := yb; y<yf; y++{
				for x := xb; x<xf; x++{
					
					if x < 0 || y < 0 || y >= host.Image_pixel_h || x >= host.Image_pixel_w { continue }
					if host.Heightdata[y][x] == 0 { continue }

					if min_height > host.Heightdata[y][x] {
						min_height = host.Heightdata[y][x]
						lgn = float64(x-xb)/float64(xf-xb)*dummy_longitude_interval+lg
						ltn = (1.0-float64(y-yb)/float64(yf-yb))*host.Dummy_city_interval_latitude+lt
					}
				}
			}

			if min_height >= host.Dummy_city_max_height {
				continue
			}

			_,_=lgn,ltn

			lga := library.Ftoa(lgn)
			lta := library.Ftoa(ltn)
			
			line := "DUMMY_"+lga+"_"+lta+","+lga+","+lta+",0\n"

			if _, err := writer.Write([]byte(line)); err != nil {
				fmt.Println("Error : dummycity : Failed to write file.")
				os.Exit(2)
			}
			
		}
	}

	writer.Flush()

}