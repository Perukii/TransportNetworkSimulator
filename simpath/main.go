package main

import (
	"fmt"
	"flag"
	"os"
	"../library"
	"sort"
	"math"
)

type Edge struct{
	a int
	b int
	score float64
	dist float64
}

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

	for i := 0; i<area_num; i++ {
		var edge_cmp []Edge
		for j := 0; j<area_num; j++ {
			if i == j { continue }
			var edge Edge
			edge.a = i
			edge.b = j
			edge.dist = Get_dist(&host, i,j)
			edge_cmp = append(edge_cmp, edge)
		}
		
		listed := 0
		sort.Slice(edge_cmp, func(i, j int) bool { return edge_cmp[i].dist < edge_cmp[j].dist })
		for j := 0; j<host.Kruskal_path_max_cross; j++{
			if len(edge_cmp) <= j { break } 
			edge_cmp[j].score = Get_score(&host, edge_cmp[j])
			if (listed == 0 || edge_cmp[j].dist < host.Max_city_distance) && edge_cmp[j].score >= 0 {
				edge_list = append(edge_list, edge_cmp[j])
				listed++
			}
			
		}
	}

	sort.Slice(edge_list, func(i, j int) bool { return edge_list[i].score < edge_list[j].score })

	for i := 0; i<len(edge_list); i++ {
		ga := Get_group(&host, &group, edge_list[i].a)
		gb := Get_group(&host, &group, edge_list[i].b)

		min_angle := 360.0
		
		if ga == gb {

			aindex := host.Cityindex[host.UrbanAreadata[edge_list[i].a].Name]
			bindex := host.Cityindex[host.UrbanAreadata[edge_list[i].b].Name]
			acomp := Get_angle(&host, aindex, bindex)
			bcomp := Get_angle(&host, bindex, aindex)
			
			for _, ac := range edge_board[edge_list[i].a]{
				if ac == edge_list[i].b{
					min_angle = 0
					break
				}
				index := host.Cityindex[host.UrbanAreadata[ac].Name]
				min_angle = math.Min(min_angle, Get_angle_dist(&host, acomp, Get_angle(&host, aindex, index)))
			}
			for _, ac := range edge_board[edge_list[i].b]{
				if ac == edge_list[i].a || min_angle == 0{
					break
				}
				index := host.Cityindex[host.UrbanAreadata[ac].Name]
				min_angle = math.Min(min_angle, Get_angle_dist(&host, bcomp, Get_angle(&host, bindex, index)))
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
		host.Register_new_path(host.Path_width, host.Path_r, host.Path_g, host.Path_b)
		for _, ptar := range path {
			host.Write_path_point(ptar.Longitude, ptar.Latitude)
		}

		//fmt.Println("simpath : processing...("+library.Itoa(i)+"/"+library.Itoa(len(edge_list))+")")

	}

	host.Writer.Flush()
	
	
}