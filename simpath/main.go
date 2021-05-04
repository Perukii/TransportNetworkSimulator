package main

import (
	"fmt"
	"flag"
	"os"
	"../library"
	"sort"
	"math"
	//rbt "github.com/emirpasic/gods/trees/redblacktree"
)


func main(){

	fmt.Println("simpath : processing...")

	flag.Parse()
	argv := flag.Args()
    if len(argv) != 12 {
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
	host.Urbandata = library.RequestUrbanData(argv[10], host.Image_pixel_w, host.Image_pixel_h, host.Data_digit)
	host.UrbanAreadata = library.RequestUrbanAreaData(argv[11], host.Image_pixel_w, host.Image_pixel_h, host.Data_digit)
	
	host.Init()
	host.Init_writer(argv[9])
	/*
	path, distance := host.Make_aster_path(host.Cityindex["佐渡市"], host.Cityindex["仙台市青葉区"], 0.008, 1500, 0, 10, -1, true)
	host.Register_new_path(3.0, 0.8, 0.4, 0.2)
	for _, ptar := range path {
		host.Write_path_point(ptar.Longitude, ptar.Latitude)
	}
	fmt.Println(distance)
	host.Writer.Flush()
	*/
	
	height_weight := 0.01
	dist_weight := 1000.0
	urban_weight:= -10.0
	cmp_pitv := 0.01
	res_pitv := 0.002
	sea_weight := 10.0

	
	

	min_area_pop := 100000
	var area_num int
	var group []int

	for i := 0; i<len(host.UrbanAreadata); i++ {
		if host.UrbanAreadata[i].Population >= min_area_pop {
			group = append(group, i)
			area_num++
		} else {
			break
		}
	}
	
	fmt.Println(area_num)

	type Edge struct{
		a int
		b int
		score float64
	}

	var edge_list []Edge

	get_score := func(edge Edge) float64{
		_, score := host.Make_aster_path(host.Cityindex[host.UrbanAreadata[edge.a].Name],
										 host.Cityindex[host.UrbanAreadata[edge.b].Name],
										 cmp_pitv, height_weight, dist_weight, urban_weight, sea_weight, -1, false)
		return score
	}

	get_group := func(ad int) int {
		for {
			if group[ad] == ad {
				break
			} else {
				ad = group[ad]
			}
		}
		return ad
	}
 
	var get_dist = func(i, j int) float64{
		iad := host.Cityindex[host.UrbanAreadata[i].Name]
		jad := host.Cityindex[host.UrbanAreadata[j].Name]
		lgd := host.Citydata[iad].Longitude - host.Citydata[jad].Longitude
		ltd := host.Citydata[iad].Latitude - host.Citydata[jad].Latitude
		return math.Sqrt(lgd*lgd+ltd*ltd)
	}

	max_nearestdist := 0.0

	for i := 0; i<area_num; i++ {
		nearestdist := 0.0
		for j := i+1; j<area_num; j++ {
			dist := get_dist(i,j)
			if nearestdist == 0.0 || dist < nearestdist {
				nearestdist = dist
			}
		}
		if max_nearestdist < nearestdist{
			max_nearestdist = nearestdist
		}
	}


	for i := 0; i<area_num; i++ {
		for j := i+1; j<area_num; j++ {
			var edge Edge
			edge.a = i
			edge.b = j
			if get_dist(i,j) >= max_nearestdist { continue }
			edge.score = get_score(edge)
			edge_list = append(edge_list, edge)
		}
	}

	sort.Slice(edge_list, func(i, j int) bool { return edge_list[i].score < edge_list[j].score })

	for i := 0; i<len(edge_list); i++ {
		ga := get_group(edge_list[i].a)
		gb := get_group(edge_list[i].b)
		if ga == gb { continue }
		group[gb] = edge_list[i].a
		
		path, _ := host.Make_aster_path(
			host.Cityindex[host.UrbanAreadata[edge_list[i].a].Name],
			host.Cityindex[host.UrbanAreadata[edge_list[i].b].Name],
			res_pitv, height_weight, dist_weight, urban_weight, sea_weight, -1, false)
		host.Register_new_path(3.0, 0.8, 0.4, 0.2)
		for _, ptar := range path {
			host.Write_path_point(ptar.Longitude, ptar.Latitude)
		}
	}

	
	
	
	
	host.Writer.Flush()
	
	
}