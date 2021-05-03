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

func (host *SpHost) make_aster_path(index_a, index_b int, pitv, width, r, g, b float64){
	
	
	var point_list []PathPoint
	var point_index map[LgLt]int // contains addresses of point_list
	var open_list *rbt.Tree // contains addresses of point_list

	open_list = rbt.NewWithIntComparator()

	point_index = make(map[LgLt]int)
	city_a := &host.citydata[index_a]
	city_b := &host.citydata[index_b]
	ptar := toLgLt(city_a.Longitude, city_a.Latitude)

	get_height := func(tar LgLt) float64{
		height := 
			host.heightdata[
				int(library.GetYFromLatitude(tar.latitude, host.latitude_s, host.latitude_e, host.image_pixel_h))][
				int(library.GetXFromLongitude(tar.longitude, host.longitude_s, host.longitude_e, host.image_pixel_w))]
		fheight := float64(height)
		return fheight
	}

	get_score := func(tar LgLt, parent LgLt) float64{
		lgd := tar.longitude - city_b.Longitude
		ltd := tar.latitude - city_b.Latitude
		distance := math.Sqrt(lgd*lgd+ltd*ltd)

		height := get_height(tar)
		hdist := math.Abs(height-get_height(parent))
		return distance + hdist/2000
	}

	open_path_point := func(tar LgLt, parent LgLt){
		if _, ok := point_index[tar]; ok {
			//exists
			/*
			cmp := point_list[point_index[point_list[point_index[tar]].parent]]
			if cmp.score > point_list[point_index[parent]].score {
				point_list[point_index[tar]].parent = parent
			}
			*/
			
			return
		}

		var point PathPoint
		point.flag = 1
		point.parent = parent
		point.score = get_score(tar, parent)
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
	
	open_path_point(ptar, ptar)

	for {
		
		ad := open_list.Left().Value.(int)
		//point_list[ad].parent = ptar
		ptar = point_list[ad].lglt
		
		up := toLgLt(ptar.longitude     , ptar.latitude+pitv)
		dw := toLgLt(ptar.longitude     , ptar.latitude-pitv)
		lf := toLgLt(ptar.longitude-pitv, ptar.latitude     )
		rg := toLgLt(ptar.longitude+pitv, ptar.latitude     )
		
		open_path_point(up, ptar)
		open_path_point(dw, ptar)
		open_path_point(lf, ptar)
		open_path_point(rg, ptar)

		close_path_point(ptar)

		if math.Abs(ptar.longitude-city_b.Longitude) <= pitv && math.Abs(ptar.latitude-city_b.Latitude) <= pitv{
			break
		}

	}

	host.register_new_path(width, r, g, b)
	btar := ptar
	for {
		host.write_path_point((ptar.longitude+btar.longitude)/2, (ptar.latitude+btar.latitude)/2)
		
		if ptar.longitude == city_a.Longitude && ptar.latitude == city_a.Latitude{
			break
		}
		btar = ptar
		ptar = point_list[point_index[ptar]].parent
	}

	
	for _, iad := range open_list.Values() {
		ad := iad.(int)
		host.register_new_path(5.0, r*0.5, g*0.5, b*0.5)
		host.write_path_point(point_list[ad].lglt.longitude, point_list[ad].lglt.latitude)
		host.write_path_point(point_list[ad].lglt.longitude, point_list[ad].lglt.latitude+0.01)
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
	host.make_aster_path(host.cityindex["Akita"], host.cityindex["Sendai"], 0.03, 5.0, 1.0, 0.5, 0.3)
	host.writer.Flush()
}