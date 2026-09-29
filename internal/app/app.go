package app

import (
	"fmt"
	"os"

	"github.com/TwiN/go-color"
	"github.com/common-nighthawk/go-figure"
	"github.com/terraform_runner/helper"
	"github.com/terraform_runner/internal/cyborg"
	helm "github.com/terraform_runner/internal/helm"
	s "github.com/terraform_runner/internal/stage"
	"github.com/terraform_runner/web"
)

var (
	items  []string
	choice string
)

func Run() {
	if len(os.Args) < 2 {
		figure.NewColorFigure("GoProvision", "", "cyan", true).Print()
		fmt.Printf("\n\n%s%sProjects :\n%s%s\n", color.Bold, color.Cyan, color.Reset, color.Reset)
		items = []string{"Nozaros", "Oranos", "Helm", "Cyborg", "Web UI", "Exit"}
		choice = helper.AskSelect(items)
	} else {
		choice = os.Args[1]
	}

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	switch choice {
	case "--Helm", "--helm", "Helm":
		helm.Helm(hostname, "/esxi_installer")
		return
	case "--Nozaros", "--nozaros", "Nozaros":
		s.Nozaros(hostname, "/terraform/final_terraform")
		return
	case "--Oranos", "--oranos", "Oranos":
		s.Oranos(hostname, "terraform/vlan_trraform")
		return
	case "--Cyborg", "--cyborg", "Cyborg":
		cyborg.Cyborg(hostname, "/ansible/ansible-core-deploy")
	case "--Webui", "Web UI":
		web.Webui(hostname)
	case "Exit":
		os.Exit(0)
	case "--help", "-- help", "-h", "- h", "--h", "-- h", "-H", "--H":
		showHelp()
	default:
		figure.NewColorFigure("GoProvision", "", "cyan", true).Print()
		fmt.Println(color.Yellow + "\n choose one of commands ..." + color.Reset)
		// fmt.Println("\nUsage : \n\tgo run <file> command \n\t./terraform command \n\nthe commands are: \n\t--Nozaros creating multiple VMs with diffrent ips \n\t--Oranos  creating multiple VLANs \n\t--Moon    creating normal VMs")
		fmt.Printf("\nUsage : \n\tgo run <file> command \n\t./terraform command \n\nthe commands are: \n\t%s--Helm%s      Setting up the Esxi product \n\t%s--Nozaros%s   creating multiple VMs with diffrent ips \n\t%s--Oranos%s    creating multiple VLANs \n\t%s--Cyborg%s    Configuring and installing packages with Ansible \n\t%s--Webui%s\t\t    UI for all Configuration \n", color.Yellow, color.Reset, color.Yellow, color.Reset, color.Yellow, color.Reset, color.Yellow, color.Reset, color.Yellow, color.Reset)
	}
}

func showHelp() {
	fmt.Printf(`%s%sTerraform Runner is Automation Project for creating Dynamic and usable VMs , VLAN , and other amazing operation on vCenter (VMware Product).%s%s
	
%sUsage:%s 
	go run . [flags]
	(e.g.) ./terraform_runner [flags]

%sthe flags are :%s 

	--Helm		Setting up the Esxi product		
	--Nozaros	creating multiple VMs with diffrent ips
	--Oranos	creating multiple VLANs
	--Cyborg	Configuring and installing packages with Ansible
	--Webui     	UI for all Configuration
 `, color.Bold, color.Yellow, color.Reset, color.Reset, color.Bold, color.Reset, color.Bold, color.Reset)
}
