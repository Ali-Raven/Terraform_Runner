package main

import (
	"fmt"
	"os"

	"github.com/common-nighthawk/go-figure"
)

func main() {
	// go-figure (art banner)
	figure.NewColorFigure("DevOps", "", "cyan", true).Print()
	// go-figure (art banner)

	if len(os.Args) < 2 {
		fmt.Println("\nUsage : \n\tgo run <file> command \n\t./terraform command \n\nthe commands are: \n\t--Nozaros_starter creating multiple VMs with diffrent ips \n\t--Oranos_starter  creating multiple VLANs \n\t--Moon_starter    creating normal VMs")
		return
	}
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	// fmt.Println(os.Hostname())
	// path
	PATH := "~/home/" + hostname + "/Desktop/Workstation/terraform_project/final_terraform/"
	// var listNames = []string{"~/home/" , hostname , }
	switch os.Args[1] {
	case "--Nozaros_starter":
		Nozaros(PATH)
		return
	case "--Oranos_starter":
		Oranos()
		return
	case "--Moon_starter":
		Moon()
		return
	default:
		fmt.Println("\n choose one of commands ...")
		fmt.Println("\nUsage : \n\tgo run <file> command \n\t./terraform command \n\nthe commands are: \n\t--Nozaros_starter creating multiple VMs with diffrent ips \n\t--Oranos_starter  creating multiple VLANs \n\t--Moon_starter    creating normal VMs")
	}
}
