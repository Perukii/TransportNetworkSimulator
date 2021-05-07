package main

import (
	"fmt"
	"flag"
	"os"
	"io"
	"strings"
	"bufio"
	"../library"
)

func main(){

	fmt.Println("citydababase2data : processing...")
	
	flag.Parse()
	argv := flag.Args()
    if len(argv) != 7 {
		fmt.Println("Error : citydababase2data : Invalid arguments.")
		os.Exit(2)
    }
	
    population_file, err := os.Open(argv[0])
    if err != nil {
		fmt.Println("Error : citydababase2data : Failed to open file.")
		os.Exit(2)
    }
	
    position_file, err := os.Open(argv[1])
    if err != nil {
		fmt.Println("Error : citydababase2data : Failed to open file.")
		os.Exit(2)
    }

    citydata_file, err := os.Create(argv[2])
    if err != nil {
		fmt.Println("Error : citydababase2data : Failed to create file.")
		os.Exit(2)
    }

	longitude_s := library.Atof(argv[3])
	longitude_e := library.Atof(argv[4])
	latitude_s := library.Atof(argv[5])
	latitude_e := library.Atof(argv[6])

	

	_,_=position_file,citydata_file

	var data_index map[int]int
	data_index = make(map[int]int)
	var data_list []library.City

	pop_header := true

	pop_reader := bufio.NewReaderSize(population_file, 4096)

    for {
		buf, _, err := pop_reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error : citydababase2data : Failed to read file.")
			os.Exit(2)
        }

		if pop_header == true {
			pop_header = false
			continue
		}
		
		slice := strings.Split(string(buf), "\n")
		//fmt.Println(slice)
		for _, it := range slice{
			itp := strings.Split(it, ",")
			if len(itp) < 7 { continue }

			ad := len(data_list)
			data_index[library.Atoi(itp[2])] = ad

			var cp library.City
			cp.Population = library.Atoi(itp[5])
			data_list = append(data_list, cp)
		}
	}

	pos_reader := bufio.NewReaderSize(position_file, 4096)
	pos_header := true

    for {
		buf, _, err := pos_reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error : citydababase2data : Failed to read file.")
			os.Exit(2)
        }

		if pos_header == true {
			pos_header = false
			continue
		}
		
		slice := strings.Split(string(buf), "\n")
		//fmt.Println(slice)
		for _, it := range slice{
			itp := strings.Split(it, ",")
			if len(itp) < 11 { continue }
			if strings.Contains(itp[1], "区") {
				continue
				//data_list[data_index[library.Atoi(itp[0])/100*100]].Population = 0
			}
			ad := data_index[library.Atoi(itp[0])]
			data_list[ad].Name = itp[1]

			data_list[ad].Longitude = library.Atof(itp[9])
			data_list[ad].Latitude = library.Atof(itp[8])
			lg := data_list[ad].Longitude
			lt := data_list[ad].Latitude 

			if lg < longitude_s || lg > longitude_e || lt < latitude_s || lt > latitude_e {
				data_list[ad].Population = 0
			}

		}
	}

	writer := bufio.NewWriter(citydata_file)

	for _, cp := range data_list {

		if cp.Name == "" || cp.Population == 0 { continue }

		line := cp.Name +
			"," + library.Ftoa(cp.Longitude) + 
			"," + library.Ftoa(cp.Latitude) + 
			"," + library.Itoa(cp.Population) + 
			"\n"

		if _, err := writer.Write([]byte(line)); err != nil {
			fmt.Println("Error : citydababase2data : Failed to write file.")
			os.Exit(2)
		}
	}

	writer.Flush()

}