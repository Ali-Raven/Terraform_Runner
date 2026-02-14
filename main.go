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
		fmt.Printf("\nUsage : \n\tgo run <file> command \n\t./terraform command \n\nthe commands are: \n\t%s--ESXI_setup%s      Setting up the Esxi product \n\t%s--Vcenter_setup%s   Setting up vCenter Product \n\t%s--Nozaros_starter%s creating multiple VMs with diffrent ips \n\t%s--Oranos_starter%s  creating multiple VLANs \n",color.Yellow , color.Reset , color.Yellow , color.Reset , color.Yellow, color.Reset, color.Yellow, color.Reset)
		return
	}

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	// var listNames = []string{"~/home/" , hostname , }
	switch os.Args[1] {
	case "--ESXI_setup":
		ESXI_setup(hostname, "/esxi_setup")
		return
	case "--Vcenter_setup":
		Vcenter_setup(hostname, "/vcenter_setup")
		return
	case "--Nozaros_starter":
		Nozaros(hostname, "/final_terraform")
		return
	case "--Oranos_starter":
		Oranos(hostname, "/vlan_terraform")
		return
	default:
		fmt.Println(color.Yellow + "\n choose one of commands ..." + color.Reset)
		fmt.Println("\nUsage : \n\tgo run <file> command \n\t./terraform command \n\nthe commands are: \n\t--Nozaros_starter creating multiple VMs with diffrent ips \n\t--Oranos_starter  creating multiple VLANs \n\t--Moon_starter    creating normal VMs")
	}
}
