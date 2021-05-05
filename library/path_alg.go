package library
import (
	"fmt"
	"math"
	"os"
	"bufio"
	rbt "github.com/emirpasic/gods/trees/redblacktree"
)
type SpHost struct{
	Pathdata_file *os.File
	Writer *bufio.Writer
	Image_pixel_w int
	Image_pixel_h int
	Data_digit int
	Longitude_s float64
	Longitude_e float64
	Latitude_s float64
	Latitude_e float64
	LgLt_ratio float64

	Heightdata [][]int
	Citydata []City
	Cityindex map[string]int
	Urbandata [][]int
	UrbanAreadata []UrbanArea
}

type LgLt struct{
	Longitude float64
	Latitude  float64
}

type LgLtFix struct{
	LongitudeF int
	LatitudeF  int
}

type PathPoint struct{
	Flag int // 1:opened 2:closed
	Score float64
	LgLt LgLt
	Parent LgLt
}

func ToLgLt(longitude, latitude float64) LgLt{
	var res LgLt
	res.Longitude = longitude
	res.Latitude = latitude
	return res
}

func ToLgLtFix(lglt LgLt) LgLtFix{
	var res LgLtFix
	res.LongitudeF = int(math.Floor(lglt.Longitude*1000))
	res.LatitudeF = int(math.Floor(lglt.Latitude*1000))
	return res
}

func (host *SpHost) Init(){
	for i := 0; i < len(host.Citydata); i++{
		host.Cityindex[host.Citydata[i].Name] = i
	}
}

func (host *SpHost) Write_line(line string){
	if _, err := host.Writer.Write([]byte(line)); err != nil {
		fmt.Println("Error : simpath : Failed to write file.")
		os.Exit(2)
	}
}

func (host *SpHost) Register_new_path(width, r, g, b float64){
	host.Write_line("#,"+Ftoa(width)+","+Ftoa(r)+","+Ftoa(g)+","+Ftoa(b)+"\n")
}

func (host *SpHost) Write_path_point(lg, lt float64){
	host.Write_line(Ftoa(lg)+","+Ftoa(lt)+"\n")
}

func (host *SpHost) Init_writer(file string){
	var err error
	host.Pathdata_file, err = os.Create(file)
    if err != nil {
		fmt.Println("Error : simpath : Failed to create file.")
		os.Exit(2)
    }

	host.Writer = bufio.NewWriter(host.Pathdata_file)
}

func (host *SpHost) Make_aster_path(index_a, index_b int, pitv, height_weight, height_diff_weight, dist_weight, urban_weight, sea_weight float64, loop int, debug bool) ([]LgLt, float64){
	
	
	var point_list []PathPoint
	var point_index map[LgLtFix]int // contains an address of point_list
	var open_list *rbt.Tree // contains an address of point_list

	open_list = rbt.NewWithIntComparator()

	point_index = make(map[LgLtFix]int)
	city_a := &host.Citydata[index_a]
	city_b := &host.Citydata[index_b]
	ptar := ToLgLt(city_a.Longitude, city_a.Latitude)

	get_height_and_urban := func(tar LgLt) (float64, int){

		yad := int(GetYFromLatitude(tar.Latitude, host.Latitude_s, host.Latitude_e, host.Image_pixel_h))
		xad := int(GetXFromLongitude(tar.Longitude, host.Longitude_s, host.Longitude_e, host.Image_pixel_w))
		if yad < 0 || yad >= host.Image_pixel_h || xad < 0 || xad >= host.Image_pixel_w{
			return 0.0, 0
		}
			
		height := float64(host.Heightdata[yad][xad])
		var urban int
		if urban_weight == 0{
			urban = 0
		} else {
			urban = host.Urbandata[yad][xad]
		}
		
		return height, urban
	}

	get_height := func(tar LgLt) float64{
		height,_ := get_height_and_urban(tar)
		return height
	}

	get_score := func(tar LgLt, parent LgLt) float64{

		height, urban := get_height_and_urban(tar)
		hdist := math.Abs(height-get_height(parent))
		sea_point := 0.0
		if height == 0 {
			sea_point = sea_weight
		}

		var uscore float64
		if urban == -1 {
			uscore = 1.0
		} else {
			uscore = 0.0
		}

		lgd := (tar.Longitude - city_b.Longitude)/host.LgLt_ratio
		ltd := tar.Latitude - city_b.Latitude
		distance := math.Sqrt(lgd*lgd+ltd*ltd)
		
		return height*height_weight + distance*dist_weight + hdist*height_diff_weight + uscore*urban_weight + sea_point
	}

	open_path_point := func(tar LgLt, parent LgLt){
		if _, ok := point_index[ToLgLtFix(tar)]; ok {
			//exists
			return
		}

		var point PathPoint
		point.Flag = 1
		point.Parent = parent
		point.Score = get_score(tar, parent)
		point.LgLt = tar

		if point.Score < 0{
			point.Flag = 3
			return 
		}

		ad := len(point_list)
		point_index[ToLgLtFix(tar)] = ad
		idscore := int(math.Floor(point.Score))
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
		ad := point_index[ToLgLtFix(tar)]
		if point_list[ad].Flag != 1 { return }
		point_list[ad].Flag = 2

		idscore := int(math.Floor(point_list[ad].Score))
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
		if open_list.Size() == 0 { break }
		ad := open_list.Left().Value.(int)

		ptar = point_list[ad].LgLt
		
		up := ToLgLt(ptar.Longitude     				, ptar.Latitude+pitv)
		dw := ToLgLt(ptar.Longitude     				, ptar.Latitude-pitv)
		lf := ToLgLt(ptar.Longitude-pitv*host.LgLt_ratio, ptar.Latitude     )
		rg := ToLgLt(ptar.Longitude+pitv*host.LgLt_ratio, ptar.Latitude     )

		
		open_path_point(up, ptar)
		open_path_point(dw, ptar)
		open_path_point(lf, ptar)
		open_path_point(rg, ptar)

		close_path_point(ptar)

		if loop < 0 &&
		   math.Abs(ptar.Longitude-city_b.Longitude) <= pitv &&
		   math.Abs(ptar.Latitude-city_b.Latitude) <= pitv{
			break
		} else {
			if loop == 0 {
				break
			} else {
				loop--
			}
		}

	}

	if debug == true {
		for _, iad := range open_list.Values() {
			ad := iad.(int)
			host.Register_new_path(5.0, 0, 0, 0)
			host.Write_path_point(point_list[ad].LgLt.Longitude, point_list[ad].LgLt.Latitude)
			host.Write_path_point(point_list[ad].LgLt.Longitude, point_list[ad].LgLt.Latitude+0.01)
		}
	}
	var res []LgLt
	length := 0.0
	
	if loop < 0{
		
		res = append(res, ToLgLt(city_b.Longitude, city_b.Latitude))
		btar := ptar
		for {
			res = append(res, ToLgLt((ptar.Longitude+btar.Longitude)/2, (ptar.Latitude+btar.Latitude)/2))
			
			if ptar.Longitude == city_a.Longitude && ptar.Latitude == city_a.Latitude{
				break
			}
			btar = ptar
			ptar = point_list[point_index[ToLgLtFix(ptar)]].Parent
		}
		res = append(res, ToLgLt(city_a.Longitude, city_a.Latitude))

		for i := 1; i < len(res); i++ {
			lgd := (res[i].Longitude - res[i-1].Longitude)/host.LgLt_ratio
			ltd := res[i].Latitude - res[i-1].Latitude
			distance := math.Sqrt(lgd*lgd+ltd*ltd)
			length += distance
		}
		
	} else {
		for _, tar := range point_list{
			if tar.Flag != 2 { continue }
			res = append(res, tar.LgLt)
		} 
	}
	return res, length

}