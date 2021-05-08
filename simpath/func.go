package main

import (
	"../library"
	"math"
)

func Get_dist(host *library.SpHost, i, j int) float64{
	iad := host.Cityindex[host.UrbanAreadata[i].Name]
	jad := host.Cityindex[host.UrbanAreadata[j].Name]
	lgd := (host.Citydata[iad].Longitude - host.Citydata[jad].Longitude)/host.LgLt_ratio
	ltd := host.Citydata[iad].Latitude - host.Citydata[jad].Latitude
	return math.Sqrt(lgd*lgd+ltd*ltd)
}

func Get_score(host *library.SpHost, edge Edge) float64{
	ca := host.Cityindex[host.UrbanAreadata[edge.a].Name]
	cb := host.Cityindex[host.UrbanAreadata[edge.b].Name]
	_, path_sc := host.Make_aster_path(ca, cb,
									   host.Path_draft_interval, -1, false)
	population := float64(host.Citydata[ca].Population+host.Citydata[cb].Population)
	
	if path_sc < 0 || path_sc > Get_dist(host, edge.a, edge.b)*host.Max_path_distance_per_city_distance{
		return -1
	}

	return path_sc - population*host.Population_score
}

func Get_group(host *library.SpHost, group *[]int, ad int) int {
	for {
		if (*group)[ad] == ad {
			break
		} else {
			ad = (*group)[ad]
		}
	}
	return ad
}

func Get_angle(host *library.SpHost, a,b int) float64{
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

func Get_angle_dist(host *library.SpHost, aa, ab float64) float64{

	da := aa-ab
	for ; da<0 ; { da += 360 }
	for ; da>=360 ; { da -= 360 }
	return math.Min(da, 360-da)
}