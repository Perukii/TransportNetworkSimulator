package main
import (
	"fmt"
	"flag"
	"os"
	"strconv"
	"strings"
)

func atoi(it string) int{
	value, err := strconv.Atoi(it)
	if err != nil {
		fmt.Println("Error : heightdata_viewer : ", err)
		os.Exit(2)
	}
	return value
}

func main(){
	flag.Parse()
	argv := flag.Args()
    if len(argv) != 3 {
		fmt.Println("Error : heightdata_viewer : Invalid arguments.")
		os.Exit(2)
    }

    heightdata_file, err := os.Open(argv[0])
    if err != nil {
		fmt.Println("Error : heightdata_viewer : Failed to open file.")
		os.Exit(2)
    }
    defer heightdata_file.Close()

	image_pixel_w := atoi(argv[1])

	image_pixel_h := atoi(argv[2])

	var heightdata [][]int
	heightdata = make([][]int, image_pixel_h)
	for i := 0; i<image_pixel_h; i++ {
		heightdata[i] = make([]int, image_pixel_w)
	}

    buf := make([]byte, image_pixel_w*5)
	row := 0
    for{
		if row >= image_pixel_h { break }
        n, err := heightdata_file.Read(buf)
        if n == 0 {
            break
        }
        if err != nil {
            panic(err)
        }
        //fmt.Println(string(buf))
		
		column := 0
		slice := strings.Split(strings.Replace(string(buf), "\n", "", -1), ",")
		

		for _, it := range slice{
			if column >= image_pixel_w { break }
			if it == "" { break }
			value := atoi(it)
			heightdata[row][column] = value
			column++
		}
		row++
    }
	fmt.Println("ok", image_pixel_h)

	
}