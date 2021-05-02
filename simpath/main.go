package main

import (
	"fmt"
	"flag"
	"os"
	"bufio"
	"../library"
)

type SpHost struct{
	pathdata_file *os.File
	writer *bufio.Writer
}



func (host *SpHost) write_line(line string){
	if _, err := host.writer.Write([]byte(line)); err != nil {
		fmt.Println("Error : simpath : Failed to write file.")
		os.Exit(2)
	}
}
func (host *SpHost) register_new_path(width, r, g, b float64){
	host.write_line("#,"+library.Ftoa(width)+","+library.Ftoa(r)+","+library.Ftoa(g)+","+library.Ftoa(b)+"\n")
}

func (host *SpHost) write_path_point(lg, lt float64){
	host.write_line(library.Ftoa(lg)+","+library.Ftoa(lt)+"\n")
}

func (host *SpHost) init_writer(file string){
	var err error
	host.pathdata_file, err = os.Create(file)
    if err != nil {
		fmt.Println("Error : simpath : Failed to create file.")
		os.Exit(2)
    }

	host.writer = bufio.NewWriter(host.pathdata_file)
}

func main(){

	fmt.Println("simpath : processing...")

	flag.Parse()
	argv := flag.Args()
    if len(argv) != 10 {
		fmt.Println("Error : simpath : Invalid arguments.")
		os.Exit(2)
    }

	var host SpHost

	image_pixel_w := library.Atoi(argv[1])
	image_pixel_h := library.Atoi(argv[2])
	data_digit := library.Atoi(argv[3])
	longitude_s := library.Atof(argv[5])
	longitude_e := library.Atof(argv[6])
	latitude_s := library.Atof(argv[7])
	latitude_e := library.Atof(argv[8])
	

	heightdata := library.RequestHeightData(argv[0], image_pixel_w, image_pixel_h, data_digit)
	citydata := library.RequestCityData(argv[4], image_pixel_w, image_pixel_h, data_digit)

	_,_,_,_,_,_,_,_,_ = image_pixel_w, image_pixel_h,data_digit,longitude_s,longitude_e,latitude_s,latitude_e,heightdata,citydata
	
	host.init_writer(argv[9])
	host.register_new_path(1.0,1.0,0.5,0.3)
	host.write_path_point(40.0,50.0)
	host.writer.Flush()
}