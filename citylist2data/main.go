package main
import (
	"fmt"
	"flag"
	"os"
	"strings"
	"bufio"
	"../library"
)

type city struct{
	name string
	longitude float64
	latitude float64
}

func Ftoa(it float64) string{

	return fmt.Sprintf("%f", it)
}

func tonum(code string) float64{
	var res float64
	var d0, dm1, dm2 int = -1, -1, -1
	i := 0
	for _, it := range code{
		if it=='°' { d0 = i }
		if it=='′' && dm1 == -1 {
			dm1 = i
		} else if it=='′' && dm1 != -1 && dm2 == -1 {
			dm2 = i
		}
		i++
	}
	res = library.Atof(code[0:d0])+library.Atof(code[d0+2:dm1+1])/60.0+library.Atof(code[dm1+4:dm2+3])/3600.0
	return res
}

func main(){
	fmt.Println("citylist2data : processing...")

	flag.Parse()
	argv := flag.Args()
    if len(argv) != 2 {
		fmt.Println("Error : citylist2data : Invalid arguments.")
		os.Exit(2)
    }
    
    citylist_file, err := os.Open(argv[0])
    if err != nil {
		fmt.Println("Error : citylist2data : Failed to open file.")
		os.Exit(2)
    }

    citydata_file, err := os.Create(argv[1])
    if err != nil {
		fmt.Println("Error : citylist2data : Failed to create file.")
		os.Exit(2)
    }

	buf := make([]byte, 255)
	writer := bufio.NewWriter(citydata_file)
	
    for {
        n, err := citylist_file.Read(buf)
        if n == 0 {
            break
        }
        if err != nil {
			fmt.Println("Error : citylist2data : Failed to read file.")
			os.Exit(2)
        }
		
		slice := strings.Split(string(buf), "\n")
		for _, it := range slice{
			itp := strings.Split(it, ",")
			var cp city
			cp.name = itp[0]
			cp.longitude = tonum(itp[1])
			cp.latitude = tonum(itp[2])

			line := cp.name +
					"," + library.Ftoa(cp.longitude) + 
					"," + library.Ftoa(cp.latitude) + 
					"\n"

			if _, err := writer.Write([]byte(line)); err != nil {
				fmt.Println("Error : citylist2data : Failed to write file.")
				os.Exit(2)
			}

		}
		
    }

	writer.Flush()

}