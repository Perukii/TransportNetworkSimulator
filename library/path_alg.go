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

	Height_score float64
	Height_difference_score float64
	Distance_score float64
	Urban_area_score float64
	Sea_area_score float64
	Population_score float64
	Kruskal_path_min_population int
	Kruskal_path_max_angle_difference float64
	Kruskal_path_max_cross int
	Path_release_interval float64
	Path_draft_interval float64
	Urban_area_interval float64
	Urban_wide_area_density int
	Urban_central_area_density int
	Urban_area_height_difference_score float64
	Max_path_distance_per_city_distance float64
	Max_bridge_distance float64
	Max_city_distance float64

	Path_r float64
	Path_g float64
	Path_b float64
	Path_width float64
	Mark_r float64
	Mark_g float64
	Mark_b float64
	Mark_width float64

	Project_name string
	Latitude_per_pixel float64

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

func (host *SpHost) ApplyCommonArgument(argv []string) int{
	
	if len(argv) != 38 { return -1 }
	host.Image_pixel_w = Atoi(argv[1])
	host.Image_pixel_h = Atoi(argv[2])
	host.Data_digit = Atoi(argv[3])
	host.Longitude_s = Atof(argv[5])
	host.Longitude_e = Atof(argv[6])
	host.Latitude_s = Atof(argv[7])
	host.Latitude_e = Atof(argv[8])
	host.LgLt_ratio = ((host.Longitude_s-host.Longitude_e)/float64(host.Image_pixel_w))/((host.Latitude_s-host.Latitude_e)/float64(host.Image_pixel_h))
	host.Heightdata = RequestHeightData(argv[0], host.Image_pixel_w, host.Image_pixel_h, host.Data_digit)
	host.Citydata = RequestCityData(argv[4], host.Image_pixel_w, host.Image_pixel_h, host.Data_digit)
	host.Cityindex = make(map[string]int)
	host.Urbandata = RequestUrbanData(argv[10], host.Image_pixel_w, host.Image_pixel_h, host.Data_digit)
	host.UrbanAreadata = RequestUrbanAreaData(argv[11], host.Image_pixel_w, host.Image_pixel_h, host.Data_digit)
	
	host.Latitude_per_pixel = (host.Latitude_s-host.Latitude_e)/float64(host.Image_pixel_h)
	host.Height_score = Atof(argv[12])
	host.Height_difference_score = Atof(argv[13])
	host.Distance_score = Atof(argv[14])
	host.Urban_area_score = Atof(argv[15])
	host.Sea_area_score = Atof(argv[16])
	host.Population_score = Atof(argv[17])
	host.Kruskal_path_min_population = Atoi(argv[18])
	host.Kruskal_path_max_angle_difference = Atof(argv[19])
	host.Kruskal_path_max_cross = Atoi(argv[20])
	host.Path_release_interval = Atof(argv[21])
	host.Path_draft_interval = Atof(argv[22])
	host.Urban_area_interval = host.Latitude_per_pixel
	host.Urban_area_height_difference_score = Atof(argv[23])
	host.Urban_wide_area_density = Atoi(argv[24])
	host.Urban_central_area_density = Atoi(argv[25])
	host.Max_path_distance_per_city_distance = Atof(argv[26])
	host.Max_bridge_distance = Atof(argv[27])
	host.Max_city_distance = Atof(argv[28])
	host.Project_name = argv[29]

	host.Path_r = Atof(argv[30])
	host.Path_g = Atof(argv[31])
	host.Path_b = Atof(argv[32])
	host.Path_width = math.Abs(Atof(argv[33])/host.Latitude_per_pixel)
	host.Mark_r = Atof(argv[34])
	host.Mark_g = Atof(argv[35])
	host.Mark_b = Atof(argv[36])
	host.Mark_width = math.Abs(Atof(argv[37])/host.Latitude_per_pixel)
	return 1
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
	host.Pathdata_file, err = os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
		fmt.Println("Error : simpath : Failed to create file.")
		os.Exit(2)
    }

	host.Writer = bufio.NewWriter(host.Pathdata_file)
}

func (host *SpHost) Make_aster_path(index_a, index_b int, interval float64, loop int, debug bool) ([]LgLt, float64){
	
	
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
		if host.Urban_area_score == 0{
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
		seascore := 0.0
		if height == 0 {
			seascore = host.Sea_area_score
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
		
		return height*host.Height_score +
			   distance*host.Distance_score +
			   hdist*host.Height_difference_score +
			   uscore*host.Urban_area_score +
			   seascore
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
		
		up := ToLgLt(ptar.Longitude     				, ptar.Latitude+interval)
		dw := ToLgLt(ptar.Longitude     				, ptar.Latitude-interval)
		lf := ToLgLt(ptar.Longitude-interval*host.LgLt_ratio, ptar.Latitude     )
		rg := ToLgLt(ptar.Longitude+interval*host.LgLt_ratio, ptar.Latitude     )

		
		open_path_point(up, ptar)
		open_path_point(dw, ptar)
		open_path_point(lf, ptar)
		open_path_point(rg, ptar)

		close_path_point(ptar)

		if loop < 0 &&
		   math.Abs(ptar.Longitude-city_b.Longitude) <= interval &&
		   math.Abs(ptar.Latitude-city_b.Latitude) <= interval{
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

		bridge := 0.0
		for i := 1; i < len(res); i++ {
			lgd := (res[i].Longitude - res[i-1].Longitude)/host.LgLt_ratio
			ltd := res[i].Latitude - res[i-1].Latitude
			distance := math.Sqrt(lgd*lgd+ltd*ltd)
			length += distance

			if get_height(res[i]) == 0{
				bridge += distance
				if bridge > host.Max_bridge_distance {
					length = -1
					break
				}
			} else {
				bridge = 0
			}
		}
		
	} else {
		for _, tar := range point_list{
			if tar.Flag != 2 { continue }
			res = append(res, tar.LgLt)
		} 
	}
	return res, length

}