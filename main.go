package main

import (
	"fmt"
	"os"
	"github.com/TwiN/go-color"
	"github.com/common-nighthawk/go-figure"
)

func main() {
	// go-figure (art banner)
	figure.NewColorFigure("BBDH DevOps", "", "cyan", true).Print()
	// go-figure (art banner)

	if len(os.Args) < 2 {
		fmt.Printf("\nUsage : \n\tgo run <file> command \n\t./terraform command \n\nthe commands are: \n\t%s--Nozaros_starter%s creating multiple VMs with diffrent ips \n\t%s--Oranos_starter%s  creating multiple VLANs \n\t%s--Moon_starter%s    creating normal VMs\n", color.Yellow, color.Reset, color.Yellow, color.Reset, color.Yellow, color.Reset)
		return
	}
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	
	// var listNames = []string{"~/home/" , hostname , }
	switch os.Args[1] {
	case "--Nozaros_starter":
		Nozaros(hostname , "/final_terraform")
		return
	case "--Oranos_starter":
		Oranos(hostname , "/vlan_terraform")
		return
	case "--Moon_starter":
		Moon(hostname , "/normal")
		return
	default:
		fmt.Println(color.Yellow + "\n choose one of commands ..." + color.Reset)
		fmt.Println("\nUsage : \n\tgo run <file> command \n\t./terraform command \n\nthe commands are: \n\t--Nozaros_starter creating multiple VMs with diffrent ips \n\t--Oranos_starter  creating multiple VLANs \n\t--Moon_starter    creating normal VMs")
	}
}
