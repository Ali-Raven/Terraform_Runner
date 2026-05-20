package main

import (
	"fmt"
	"github.com/TwiN/go-color"
	"github.com/common-nighthawk/go-figure"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		figure.NewColorFigure("BBDH DevOps", "", "cyan", true).Print()
		fmt.Printf("\nUsage : \n\tgo run <file> command \n\t./terraform command \n\nthe commands are: \n\t%s--Helm_starter%s      Setting up the Esxi product \n\t%s--Nozaros_starter%s   creating multiple VMs with diffrent ips \n\t%s--Oranos_starter%s    creating multiple VLANs \n\t%s--Cyborg_starter%s    Configuring and installing packages with Ansible \n", color.Yellow, color.Reset, color.Yellow, color.Reset, color.Yellow, color.Reset, color.Yellow, color.Reset)
		return
	}

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	// var listNames = []string{"~/home/" , hostname , }
	switch os.Args[1] {
	case "--Helm_starter":
		Helm(hostname, "/esxi_installer")
		return
	case "--Nozaros_starter":
		Nozaros(hostname, "/final_terraform")
		return
	case "--Oranos_starter":
		Oranos(hostname, "/vlan_terraform")
		return
	case "--Cyborg_starter":
		Cyborg(hostname, "/ansible-core-deploy")
	case "--Webui":
		Webui(hostname)
	default:
		figure.NewColorFigure("BBDH DevOps", "", "cyan", true).Print()
		fmt.Println(color.Yellow + "\n choose one of commands ..." + color.Reset)
		// fmt.Println("\nUsage : \n\tgo run <file> command \n\t./terraform command \n\nthe commands are: \n\t--Nozaros_starter creating multiple VMs with diffrent ips \n\t--Oranos_starter  creating multiple VLANs \n\t--Moon_starter    creating normal VMs")
		fmt.Printf("\nUsage : \n\tgo run <file> command \n\t./terraform command \n\nthe commands are: \n\t%s--Helm_starter%s      Setting up the Esxi product \n\t%s--Nozaros_starter%s   creating multiple VMs with diffrent ips \n\t%s--Oranos_starter%s    creating multiple VLANs \n\t%s--Cyborg_starter%s    Configuring and installing packages with Ansible \n", color.Yellow, color.Reset, color.Yellow, color.Reset, color.Yellow, color.Reset, color.Yellow, color.Reset)
	}
}
