package main

import (
	"fmt"
	"flag"
	"os"
	"../library"
	"sort"
	"math"
)


func main(){

	fmt.Println("simpath : processing...")

	flag.Parse()
	argv := flag.Args()

	var host library.SpHost
	
    if host.ApplyCommonArgument(argv) < 0 {
		fmt.Println("Error : simpath : Invalid arguments.")
		os.Exit(2)
    }

	host.Init()
	host.Init_writer(argv[9])

	

	type Edge struct{
		a int
		b int
		score float64
		dist float64
	}
	
	var area_num int
	var group []int
	var edge_list []Edge
	var edge_board [][]int


	for i := 0; i<len(host.UrbanAreadata); i++ {
		if host.UrbanAreadata[i].Population >= host.Kruskal_path_min_population {
			group = append(group, i)
			area_num++
		} else {
			break
		}
	}
	
	edge_board = make([][]int, area_num)

	var get_dist = func(i, j int) float64{
		iad := host.Cityindex[host.UrbanAreadata[i].Name]
		jad := host.Cityindex[host.UrbanAreadata[j].Name]
		lgd := (host.Citydata[iad].Longitude - host.Citydata[jad].Longitude)/host.LgLt_ratio
		ltd := host.Citydata[iad].Latitude - host.Citydata[jad].Latitude
		return math.Sqrt(lgd*lgd+ltd*ltd)
	}

	get_score := func(edge Edge) float64{
		ca := host.Cityindex[host.UrbanAreadata[edge.a].Name]
		cb := host.Cityindex[host.UrbanAreadata[edge.b].Name]
		_, path_sc := host.Make_aster_path(ca, cb,
										   host.Path_draft_interval, -1, false)
		population := float64(host.Citydata[ca].Population+host.Citydata[cb].Population)
		
		if path_sc < 0 || path_sc > get_dist(edge.a, edge.b)*host.Max_path_distance_per_city_distance{
			return -1
		}

		return path_sc - population*host.Population_score
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
 


	var get_angle = func(a,b int) float64{
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
		
		listed := 0
		sort.Slice(edge_cmp, func(i, j int) bool { return edge_cmp[i].dist < edge_cmp[j].dist })
		for j := 0; j<host.Kruskal_path_max_cross; j++{
			if len(edge_cmp) <= j { break } 
			edge_cmp[j].score = get_score(edge_cmp[j])
			if (listed == 0 || edge_cmp[j].dist < host.Max_city_distance) && edge_cmp[j].score >= 0 {
				edge_list = append(edge_list, edge_cmp[j])
				listed++
			}
			
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

			if min_angle < host.Kruskal_path_max_angle_difference {
				continue
			}
		} else{
			group[gb] = edge_list[i].a
		}
	
		edge_board[edge_list[i].a] = append(edge_board[edge_list[i].a], edge_list[i].b)
		edge_board[edge_list[i].b] = append(edge_board[edge_list[i].b], edge_list[i].a)
		
		path, _ := host.Make_aster_path(
			host.Cityindex[host.UrbanAreadata[edge_list[i].a].Name],
			host.Cityindex[host.UrbanAreadata[edge_list[i].b].Name],
			host.Path_release_interval, -1, false)
		host.Register_new_path(3.0, 0.8, 0.4, 0.2)
		for _, ptar := range path {
			host.Write_path_point(ptar.Longitude, ptar.Latitude)
		}
		/*
		if host.UrbanAreadata[edge_list[i].a].Name == "福岡市" {
			fmt.Println(host.UrbanAreadata[edge_list[i].b].Name)
		}
		if host.UrbanAreadata[edge_list[i].b].Name == "福岡市" {
			fmt.Println(host.UrbanAreadata[edge_list[i].a].Name)
		}
		*/
	}

	host.Writer.Flush()
	
	
}