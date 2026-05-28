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
		fmt.Printf("\nUsage : \n\tgo run <file> command \n\t./terraform command \n\nthe commands are: \n\t%s--Helm%s      Setting up the Esxi product \n\t%s--Nozaros%s   creating multiple VMs with diffrent ips \n\t%s--Oranos%s    creating multiple VLANs \n\t%s--Cyborg%s    Configuring and installing packages with Ansible \n\t%s--Webui%s\t\t    UI for all Configuration\n", color.Yellow, color.Reset, color.Yellow, color.Reset, color.Yellow, color.Reset, color.Yellow, color.Reset , color.Yellow , color.Reset)
		return
	}

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	// var listNames = []string{"~/home/" , hostname , }
	switch os.Args[1] {
	case "--Helm" , "--helm":
		Helm(hostname, "/esxi_installer")
		return
	case "--Nozaros" , "--nozaros":
		Nozaros(hostname, "/final_terraform")
		return
	case "--Oranos" , "--oranos":
		Oranos(hostname, "/vlan_terraform")
		return
	case "--Cyborg" , "--cyborg":
		Cyborg(hostname, "/ansible-core-deploy")
	case "--Webui":
		Webui(hostname)
	case "--help" , "-- help" , "-h" , "- h" , "--h" , "-- h":
		showHelp()
	default:
		figure.NewColorFigure("BBDH DevOps", "", "cyan", true).Print()
		fmt.Println(color.Yellow + "\n choose one of commands ..." + color.Reset)
		// fmt.Println("\nUsage : \n\tgo run <file> command \n\t./terraform command \n\nthe commands are: \n\t--Nozaros creating multiple VMs with diffrent ips \n\t--Oranos  creating multiple VLANs \n\t--Moon    creating normal VMs")
		fmt.Printf("\nUsage : \n\tgo run <file> command \n\t./terraform command \n\nthe commands are: \n\t%s--Helm%s      Setting up the Esxi product \n\t%s--Nozaros%s   creating multiple VMs with diffrent ips \n\t%s--Oranos%s    creating multiple VLANs \n\t%s--Cyborg%s    Configuring and installing packages with Ansible \n\t%s--Webui%s\t\t    UI for all Configuration \n", color.Yellow, color.Reset, color.Yellow, color.Reset, color.Yellow, color.Reset, color.Yellow, color.Reset , color.Yellow , color.Reset)
	}
}


func showHelp() {
	fmt.Printf(`Terraform Runner is Automation Project for creating Dynamic and usable VMs , VLAN , and other amazing operation on vCenter (VMware Product).
	
Usage: 
	go run . [flags]
	(e.g.) ./terraform_runner [flags]

the flags are : 

	--Helm		Setting up the Esxi product		
	--Nozaros	creating multiple VMs with diffrent ips
	--Oranos	creating multiple VLANs
	--Cyborg	Configuring and installing packages with Ansible
	--Webui          	UI for all Configuration
 `)
}
