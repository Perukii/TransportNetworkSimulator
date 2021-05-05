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

	var get_angle = func(a, b int) float64{
		city_a := host.Citydata[a]
		city_b := host.Citydata[b]
		ltd := city_b.Latitude-city_a.Latitude
		lgd := city_b.Longitude-city_a.Longitude
		if lgd == 0 { lgd = 0.0001 }

		var abase float64
		if ltd >= 0 && lgd >= 0 { abase = 0 }
		if ltd >= 0 && lgd < 0  { abase = 90 }
		if ltd < 0 && lgd < 0   { abase = 180 }
		if ltd < 0 && lgd >= 0  { abase = 270 }

		angle := math.Atan(ltd/lgd)*180/3.1415
		for ; angle<0; {
			angle += 90
		}
		return angle+abase
	}

	var get_angle_dist = func(aa, ab float64) float64{

		da := aa-ab
		for ; da<0 ; { da += 360 }
		for ; da>=360 ; { da -= 360 }
		return math.Min(da, 360-da)
	}

	
	height_weight := 0.1
	height_diff_weight := 1.4
	dist_weight := 1300.0
	urban_weight:= 0.0
	cmp_pitv := 0.02
	res_pitv := 0.005
	sea_weight := 2000.0
	population_weight := -0.000001
	// クラスカル路を適用する都市の最小の人口
	min_area_pop := 20000
	// 非クラスカル路を採用するのに満たす必要のある、一つの都市に対する対象のパスと他都市間のパスとの角度の差の最小値の下限
	angle_limit := 80.0
	// パスのスコア計算にて、一つの都市が比較を行う都市の数
	search_limit := 8
	
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
	

	type Edge struct{
		a int
		b int
		score float64
		dist float64
	}

	

	get_score := func(edge Edge) float64{
		ca := host.Cityindex[host.UrbanAreadata[edge.a].Name]
		cb := host.Cityindex[host.UrbanAreadata[edge.b].Name]
		_, path_sc := host.Make_aster_path(ca, cb,
										   cmp_pitv, height_weight, height_diff_weight, dist_weight, urban_weight, sea_weight, -1, false)
		population := float64(host.Citydata[ca].Population+host.Citydata[cb].Population)
		return path_sc + population*population_weight
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

	var edge_list []Edge
	var edge_board [][]int
	edge_board = make([][]int, area_num)
	_, _, _ = get_angle, get_angle_dist, angle_limit

	for i := 0; i<area_num; i++ {
		var edge_cmp []Edge
		for j := 0; j<area_num; j++ {
			if i == j { continue }
			var edge Edge
			edge.a = i
			edge.b = j
			edge.dist = get_dist(i,j)
			edge_cmp = append(edge_cmp, edge)
		}
		sort.Slice(edge_cmp, func(i, j int) bool { return edge_cmp[i].dist < edge_cmp[j].dist })
		for j := 0; j<search_limit; j++{
			if len(edge_cmp) <= j { break } 
			edge_cmp[j].score = get_score(edge_cmp[j])
			edge_list = append(edge_list, edge_cmp[j])
		}
	}

	sort.Slice(edge_list, func(i, j int) bool { return edge_list[i].score < edge_list[j].score })

	for i := 0; i<len(edge_list); i++ {
		ga := get_group(edge_list[i].a)
		gb := get_group(edge_list[i].b)

		min_angle := 360.0
		if ga == gb {

			aindex := host.Cityindex[host.UrbanAreadata[edge_list[i].a].Name]
			bindex := host.Cityindex[host.UrbanAreadata[edge_list[i].b].Name]
			acomp := get_angle(aindex, bindex)
			bcomp := get_angle(bindex, aindex)
			
			for _, ac := range edge_board[edge_list[i].a]{
				if ac == edge_list[i].b{
					min_angle = 0
					break
				}
				index := host.Cityindex[host.UrbanAreadata[ac].Name]
				min_angle = math.Min(min_angle, get_angle_dist(acomp, get_angle(aindex, index)))
			}
			for _, ac := range edge_board[edge_list[i].b]{
				if ac == edge_list[i].a || min_angle == 0{
					break
				}
				index := host.Cityindex[host.UrbanAreadata[ac].Name]
				min_angle = math.Min(min_angle, get_angle_dist(bcomp, get_angle(bindex, index)))
			}

			if min_angle < angle_limit {
				continue
			}
		} else{
			group[gb] = edge_list[i].a
		}
		
		fmt.Println(min_angle,
			host.UrbanAreadata[edge_list[i].a].Name,
			host.UrbanAreadata[edge_list[i].b].Name,
		)
	
		edge_board[edge_list[i].a] = append(edge_board[edge_list[i].a], edge_list[i].b)
		edge_board[edge_list[i].b] = append(edge_board[edge_list[i].b], edge_list[i].a)
		
		path, _ := host.Make_aster_path(
			host.Cityindex[host.UrbanAreadata[edge_list[i].a].Name],
			host.Cityindex[host.UrbanAreadata[edge_list[i].b].Name],
			res_pitv, height_weight, height_diff_weight, dist_weight, urban_weight, sea_weight, -1, false)
		host.Register_new_path(3.0, 0.8, 0.4, 0.2)
		for _, ptar := range path {
			host.Write_path_point(ptar.Longitude, ptar.Latitude)
		}
	}

	host.Writer.Flush()
	
	
}