package main

import (
	"fmt"
	"flag"
	"os"
	"bufio"
	"../library"
	"math"
	rbt "github.com/emirpasic/gods/trees/redblacktree"
)

type SpHost struct{
	pathdata_file *os.File
	writer *bufio.Writer
	image_pixel_w int
	image_pixel_h int
	data_digit int
	longitude_s float64
	longitude_e float64
	latitude_s float64
	latitude_e float64

	heightdata [][]int
	citydata []library.City
	cityindex map[string]int
}

type LgLt struct{
	longitude float64
	latitude  float64
}

type PathPoint struct{
	flag int // 1:opened 2:closed
	score float64
	lglt LgLt
	parent LgLt
}

func toLgLt(longitude, latitude float64) LgLt{
	var res LgLt
	res.longitude = longitude
	res.latitude = latitude
	return res
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

func (host *SpHost) make_aster_path(index_a, index_b int, width, r, g, b float64){
	
	
	var point_list []PathPoint
	var point_index map[LgLt]int // contains addresses of point_list
	var open_list *rbt.Tree // contains addresses of point_list

	open_list = rbt.NewWithIntComparator()

	point_index = make(map[LgLt]int)
	city_a := &host.citydata[index_a]
	city_b := &host.citydata[index_b]
	ptar := toLgLt(city_a.Longitude, city_a.Latitude)
	pitv := 0.01

	get_score := func(tar LgLt) float64{
		lgd := tar.longitude - city_b.Longitude
		ltd := tar.latitude - city_b.Latitude
		return math.Sqrt(lgd*lgd+ltd*ltd)
	}

	open_path_point := func(tar LgLt){
		if _, ok := point_index[tar]; ok {
			//exists
			return
		}

		var point PathPoint
		point.flag = 1
		point.score = get_score(tar)
		point.lglt = tar

		ad := len(point_list)
		point_index[tar] = ad
		idscore := int(math.Floor(point.score*10000))
		for{
			_, found := open_list.Get(idscore)
			if found == true {
				idscore += 1
				continue
			} else {
				open_list.Put(idscore, ad)
				break
			}
		}
		
		point_list = append(point_list, point)

	}

	close_path_point := func(tar LgLt){
		ad := point_index[tar]
		if point_list[ad].flag != 1 { return }
		point_list[ad].flag = 2

		idscore := int(math.Floor(point_list[ad].score))
		for{
			cad, found := open_list.Get(idscore)
			if found == false || cad != ad {
				idscore += 1
				continue
			} else {
				open_list.Remove(idscore)
				break
			}
		}
		
	}

	count := 0

	
	open_path_point(ptar)

	for {
		
		ad := open_list.Left().Value.(int)
		point_list[ad].parent = ptar
		ptar = point_list[ad].lglt
		
		up := toLgLt(ptar.longitude     , ptar.latitude+pitv)
		dw := toLgLt(ptar.longitude     , ptar.latitude-pitv)
		lf := toLgLt(ptar.longitude-pitv, ptar.latitude     )
		rg := toLgLt(ptar.longitude+pitv, ptar.latitude     )
		
		open_path_point(up)
		open_path_point(dw)
		open_path_point(lf)
		open_path_point(rg)

		close_path_point(ptar)

		if math.Abs(ptar.longitude-city_b.Longitude) <= pitv && math.Abs(ptar.latitude-city_b.Latitude) <= pitv{
			break
		}

		count++
		
	}

	fmt.Println(count)

	host.register_new_path(width, r, g, b)
	for {

		host.write_path_point(ptar.longitude, ptar.latitude)
		
		if ptar.longitude == city_a.Longitude && ptar.latitude == city_a.Latitude{
			break
		}

		ptar = point_list[point_index[ptar]].parent
		
	}
	

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

	host.image_pixel_w = library.Atoi(argv[1])
	host.image_pixel_h = library.Atoi(argv[2])
	host.data_digit = library.Atoi(argv[3])
	host.longitude_s = library.Atof(argv[5])
	host.longitude_e = library.Atof(argv[6])
	host.latitude_s = library.Atof(argv[7])
	host.latitude_e = library.Atof(argv[8])
	
	host.heightdata = library.RequestHeightData(argv[0], host.image_pixel_w, host.image_pixel_h, host.data_digit)
	host.citydata = library.RequestCityData(argv[4], host.image_pixel_w, host.image_pixel_h, host.data_digit)
	host.cityindex = make(map[string]int)

	for i := 0; i < len(host.citydata); i++{
		host.cityindex[host.citydata[i].Name] = i
	}

	host.init_writer(argv[9])
	host.make_aster_path(host.cityindex["Akita"], host.cityindex["Sendai"], 5.0, 1.0, 0.5, 0.3)
	host.writer.Flush()
}