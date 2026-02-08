package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)


type NormalVM struct {
	Name string `json:"name"`
	RAM  int    `json:"ram"`
	CPU  int    `json:"cpu_cores"`
}

type NormalNetwork struct {

}
func Moon_configure(wdir string) {
	var (
		RAM int8
		CPU int8
	)
	fmt.Print("RAM values : ")
	fmt.Scan(&RAM)
	fmt.Print("CPU core value : ")
	fmt.Scan(&CPU)

	valus := map[string]int{
		"cpu_cores": int(CPU),
		"ram":       int(RAM),
	}
	data, err := json.MarshalIndent(valus, "", " ")
	if err != nil {
		panic(data)
	}
	// getting working directory
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(currentDir + wdir + "/terraform.auto.tfvars.json", data, 0644)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Successfully written to terraform.auto.tfvars.json.")
		fmt.Println("going back to the menu ...")
		time.Sleep(1 * time.Second)
	}
	main()
}
