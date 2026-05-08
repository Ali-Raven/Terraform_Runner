package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/TwiN/go-color"
	"github.com/common-nighthawk/go-figure"
	"github.com/terraform_runner/helper"
	"github.com/terraform_runner/servers"
)

// initializing global variables
var (
	usrHostname                string
	usrPassword                string
	usrIP                      string
	usrNetmask                 string
	usrGateway                 string
	usrVlan                    string
	deployment_options_vCenter string
	usrvCenterIp               string
	prefix_netmask             string
	vCenter_gateway            string
	system_name_vCenter        string
	vCenter_management_pass    string
	vCenter_login_pass         string
	usrApiUrl                  string
	usrApiTokenId              string
	usrApiTokenSecret          string
)

// this function is for
func Helm(hostname, wdir string) {
	reader := bufio.NewReader(os.Stdin)
	figure.NewColorFigure("HELM", "", "cyan", true).Print()

	fmt.Println(color.Blue + "\nUsing Helm_starter" + color.Reset)
	fmt.Println(color.Yellow + "ESXI and vCenter mode ..." + color.Reset)

	time.Sleep(1 * time.Second)
	fmt.Printf("\nOptions : \n\t\n\t1. Esxi Product =>\t%sInstallation & Config of ESXI product of vMVare%s \n\t------------\t\n\t2. vCenter Product =>\t%sInstallation of vCenter Product of vMware%s \n\t------------\n\t3.Exit", color.Yellow, color.Reset, color.Yellow, color.Reset)
	// fmt.Println("\n1\n2. vCenter Product (Installation)\n3. Exit")

	fmt.Print("\n\nchoice: (1/2/3): ")
	userinput, _ := reader.ReadString('\n')
	userinput = strings.TrimSpace(userinput)

	switch userinput {
	case "1":
		ESXI_setup(hostname, wdir, reader)
		return
	case "2":
		MinimalInputEsxi(reader, wdir)
		return
	case "3":
		fmt.Println(color.Yellow + "Exiting..." + color.Reset)
		os.Exit(0)
	default:
		fmt.Println(color.Yellow + "\nWarning : choose one of the above options ..." + color.Reset)
		fmt.Println(color.Yellow + "Returning to menu ..." + color.Reset)
		time.Sleep(1 * time.Second)
		Helm(hostname, wdir)
		// fmt.Println("Options : ")
		// fmt.Println("1. ESXI install & config\n2. vCenter installation")
	}
	// MainStage(wdir , 4)
}

func ESXI_setup(hostname, wdir string, reader *bufio.Reader) {
	figure.NewColorFigure("ESXI Setup", "", "green", true).Print()

	fmt.Println(color.Blue + "\nUsing ESXI_setup" + color.Reset)
	fmt.Println(color.Yellow + "Setting up ESXI Product ..." + color.Reset)
	time.Sleep(1 * time.Second)
	fmt.Println(color.Yellow + "Fetching list of Servers for installing ESXI products ..." + color.Reset)
	servers := servers.FetchServers()
	fmt.Println(color.Green + "\nAvailable Servers:" + color.Reset)
	for i, serverName := range servers {
		fmt.Printf("%d. %s\n", i+1, serverName)
	}
	fmt.Print(color.Yellow + "\nEnter the number corresponding to the server you want to use: " + color.Reset)
	userServer, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading user input:", err)
		return
	}
	fmt.Println(color.Green + "\nSelected Server: " + userServer + color.Reset)

	// getting iso file path
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	iso_path := currentDir + wdir + "/iso_files"

	// getting iso file path
	switch strings.TrimSpace(userServer) {
	case "1":
		fmt.Println(color.Yellow + "HP G9 Server Selected ..." + color.Reset)
		Custom_iso_maker(wdir, reader, "VMware-ESXi-7.0.3-20036589-HPE-703.0.0.10.9.1.5-Jul2022.iso", iso_path, currentDir, hostname)
	case "2":
		fmt.Println(color.Yellow + "Huawei V3 Server Selected ..." + color.Reset)
	case "3":
		fmt.Println(color.Yellow + "A5 Server Selected ..." + color.Reset)
	case "4":
		fmt.Println(color.Yellow + "HP G10 Server Selected ..." + color.Reset)
	default:
		fmt.Println(color.Yellow + "Invalid selection. Please choose a valid server number." + color.Reset)
	}
}

// ==================================================================================== Custom ISO maker ==========================================================================================
func Custom_iso_maker(wdir string, reader *bufio.Reader, isoName, isoPath, currenDir, hostname string) {
	fmt.Printf(color.Yellow+"Using Custom Image maker for ISO : %s\n"+color.Reset, isoName)
	time.Sleep(1 * time.Second)

	fmt.Printf("\nEnter ESXI Hostname: ")
	usrHostname, _ = reader.ReadString('\n') // readRequired(reader, "\nEnter ESXI Hostname: ")
	usrHostname = strings.TrimSpace(usrHostname)

	fmt.Printf("\nEnter ESXI Password: ")
	usrPassword, _ = reader.ReadString('\n') // readRequired(reader, "Enter ESXI Password: ")
	usrPassword = strings.TrimSpace(usrPassword)

	fmt.Printf("\nEnter ESXI IP Address: ")
	usrIP, _ = reader.ReadString('\n') // readRequired(reader, "Enter ESXI IP Address: ")
	usrIP = strings.TrimSpace(usrIP)

	fmt.Printf("\nEnter Netmask: ")
	usrNetmask, _ = reader.ReadString('\n') // readRequired(reader, "Enter Netmask: ")
	usrNetmask = strings.TrimSpace(usrNetmask)

	fmt.Printf("\nEnter ESXI Gateway: ")
	usrGateway, _ = reader.ReadString('\n') // readRequired(reader, "Enter ESXI Gateway: ")
	usrGateway = strings.TrimSpace(usrGateway)

	fmt.Printf("\nEnter ESXI Management Network Vlan: ")
	usrVlan, _ = reader.ReadString('\n') // readRequired(reader, "Enter ESXI Management Network Vlan: ")
	usrVlan = strings.TrimSpace(usrVlan)

	if usrHostname == "" || usrPassword == "" || usrIP == "" || usrNetmask == "" || usrGateway == "" || usrVlan == "" {
		usrHostname = "esxi2.example.org"
		usrPassword = "Aa@123321"
		usrIP = "192.168.0.250"
		usrNetmask = "255.255.255.0"
		usrGateway = "192.168.0.254"
		usrVlan = "0"
	}

	vars := map[string]string{
		"esxi_hostname":   usrHostname,
		"esxi_password":   usrPassword,
		"esxi_ip_address": usrIP,
		"esxi_gateway":    usrGateway,
		"esxi_netmask":    usrNetmask,
		"esxi_vlan":       usrVlan,
	}

	fmt.Println(color.Yellow + "Generating env file ..." + color.Reset)
	time.Sleep(1 * time.Second)

	cmd := exec.Command("bash", isoPath+"/maker.sh", isoPath+"/"+isoName)
	cmd.Env = os.Environ()
	for k, v := range vars {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(color.Green + "Custom image setup is Successfully completed ..." + color.Reset)

	// Ansible_injection_esxi(reader, wdir, currenDir, hostname)
	// Terraform_proxmox(reader, wdir, currenDir, hostname)
}

// ==================================================================================== Custom ISO maker (END) ==========================================================================================

// ==================================================================================== Terraform Proxmox (ESXI)  ==========================================================================================
// func Terraform_proxmox(reader *bufio.Reader, wdir, crrentDir, hostname string) {
// 	fmt.Println(color.Yellow + "\nTerraform Section is running ..." + color.Reset)
// 	time.Sleep(1 * time.Second)
// 	fmt.Println(color.Yellow + "Enter Authentication Parameters :" + color.Reset)

// 	usrApiUrl = readRequired(reader, "\nEnter Proxmox API URL : ")
// 	usrApiTokenId = readRequired(reader, "Enter Api TOKEN ID : ")
// 	usrApiTokenSecret = readRequired(reader, "Enter Api TOKEN Secret : ")

// 	fmt.Println(color.Green + "\nParameters set Successfully ." + color.Reset)
// 	time.Sleep(1 * time.Second)

// 	// ============= running terraform ==============================

// 	// ============= running terraform (END) ==============================

// 	// to continue and Setup
// 	fmt.Print(color.Yellow + "\nDo you want to install vCenter Product on your Esxi ? (y/n) " + color.Reset)
// 	usrChoise, _ := reader.ReadString('\n')
// 	usrChoise = strings.TrimSpace(usrChoise)

// 	switch usrChoise {
// 	case "y", "yes", "", "Yes", "YES":
// 		fmt.Println(color.Yellow + "vCenter installation Processing ..." + color.Reset)
// 		time.Sleep(1 * time.Second)
// 		Vcenter_setup(wdir, reader)
// 	case "n", "no", "NO", "No":
// 		MinimalInputEsxi(reader, wdir)
// 	default:
// 		fmt.Println(color.Red + "Ivalid choice." + color.Reset)
// 		Helm(hostname, wdir)
// 	}

// 	fmt.Println(color.Yellow + "\nReturning to Main menu ..." + color.Reset)
// 	time.Sleep(1 * time.Second)
// 	figure.NewColorFigure("HELM", "", "blue", true).Print()
// 	Helm(wdir, hostname)
// }

// ==================================================================================== Terraform Proxmox (ESXI)(END)  ==========================================================================================

// ==================================================================================== Ansible injection (ESXI)  ==========================================================================================

// func Ansible_injection_esxi(reader *bufio.Reader, wdir, currentDir, hostname string) {
// 	fmt.Println(color.Yellow + "Running Ansible .." + color.Reset)
// 	time.Sleep(1 * time.Second)

// 	playbooks_path := currentDir + "/ansible-vmware-config/playbooks/"
// 	parentAnsible := currentDir + "/ansible-vmware-config"
// 	cmd := exec.Command("ansible-playbook", playbooks_path+"upload-iso-to-datastore.yml")
// 	cmd.Dir = parentAnsible
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr

// 	if err := cmd.Run(); err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println(color.Green + "Upload Successfully ..." + color.Reset)
// 	time.Sleep(1 * time.Second)

// 	fmt.Println(color.Green + "Creating ESXI ...." + color.Reset)
// 	time.Sleep(1 * time.Second)
// 	cmd2 := exec.Command("ansible-playbook", playbooks_path+"create-vm-on-vcenter.yml")
// 	cmd2.Dir = parentAnsible
// 	cmd2.Stdout = os.Stdout
// 	cmd2.Stderr = os.Stderr

// 	if err2 := cmd2.Run(); err2 != nil {
// 		log.Fatal(err2)
// 	}

// 	fmt.Println(color.Green + "ESXI Deployed Successfully ..." + color.Reset)
// 	time.Sleep(1 + time.Second)

// asking user to install vCenter or not

// func Exec_ansible() {
// 	// implement soon
// }

func MinimalInputEsxi(reader *bufio.Reader, wdir string) {
	figure.NewColorFigure("vCenter Setup", "", "green", true).Print()
	fmt.Println(color.Blue + "\nUsing vCenter_setup" + color.Reset)
	fmt.Println(color.Yellow + "\n\nYou have Existing ESXI host " + color.Reset)
	fmt.Println(color.Yellow + "Enter you vCenter configs : " + color.Reset)
	usrHostname = helper.ReadRequired(reader, "\nEnter ESXI Hostname: ")
	usrIP = helper.ReadRequired(reader, "Enter ESXI IP Address: ")
	usrNetmask = helper.ReadRequired(reader, "Enter Netmask: ")
	usrGateway = helper.ReadRequired(reader, "Enter ESXI Gateway: ")
	usrVlan = helper.ReadRequired(reader, "Enter ESXI Management Network Vlan: ")

	fmt.Println(color.Green + "\nGetting Esxi Parameters Successfully Completed " + color.Reset)
	time.Sleep(1 * time.Second)
	Vcenter_setup(wdir, reader)
}

// ==================================================================================== Ansible injection (ESXI) (END) ==========================================================================================

// ==================================================================================== vCenter setup ==========================================================================================
func Vcenter_setup(wdir string, reader *bufio.Reader) {

	fmt.Println(color.Yellow + "Setting up vCenter Product ..." + color.Reset)
	fmt.Println()
	time.Sleep(1 * time.Second)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)

	fmt.Fprintln(w, color.Purple+"ID\tSIZE\tVCPUs\tMEMORY(GB)\tSTORAGE(GB)\tHOSTS\tVMs"+color.Reset)
	fmt.Fprintln(w, color.Yellow+"1\tTiny\t2\t10\t300\t10\t100"+color.Reset)
	fmt.Fprintln(w, color.Yellow+"2\tSmall\t4\t16\t340\t100\t1000"+color.Reset)
	fmt.Fprintln(w, color.Yellow+"3\tMedium\t8\t24\t525\t400\t4000"+color.Reset)
	fmt.Fprintln(w, color.Yellow+"4\tLarge\t16\t32\t740\t1000\t10000"+color.Reset)
	fmt.Fprintln(w, color.Yellow+"5\tX-Large\t24\t48\t1180\t2000\t35000"+color.Reset)

	w.Flush()

	deployment_options_vCenter = helper.ReadRequired(reader, "\nEnter Deployment Options : ")
	// switch for checking the answer of user for deployment options
	switch deployment_options_vCenter {
	case "1":
		deployment_options_vCenter = "tiny"
	case "2":
		deployment_options_vCenter = "small"
	case "3":
		deployment_options_vCenter = "medium"
	case "4":
		deployment_options_vCenter = "large"
	case "5":
		deployment_options_vCenter = "x-large"
	default:
		fmt.Println(color.Red + "Choose from the optinos ...." + color.Reset)
		Vcenter_setup(wdir, reader)

	}
	// ---------------------end of switch---------------------------------

	usrvCenterIp = helper.ReadRequired(reader, "Enter vCenter IP : ")
	prefix_netmask = helper.ReadRequired(reader, "Enter Netmask (e.g. 24) : ")
	vCenter_gateway = helper.ReadRequired(reader, "Enter vCenter Gateway : ")
	system_name_vCenter = helper.ReadRequired(reader, "Enter System name : ")
	vCenter_management_pass = helper.ReadRequired(reader, "Enter vCenter Management Password : ")
	vCenter_login_pass = helper.ReadRequired(reader, "Enter vCenter login Password : ")

	JsonData := map[string]any{
		"__version":  "2.13.0",
		"__comments": "Production-ready embedded VCSA 8.0.3 deployment on ESXi host. Includes NTP, DNS, and security configs for enterprise reliability.",
		"new_vcsa": map[string]any{
			"esxi": map[string]any{
				"hostname":           "192.168.0.250",
				"username":           "root",
				"password":           "Aa@123321",
				"deployment_network": "VM Network",
				"datastore":          "datastore1",
			},
			"appliance": map[string]any{
				"thin_disk_mode":    true,
				"deployment_option": deployment_options_vCenter,
				"name":              "vcsa-prod",
			},
			"network": map[string]any{
				"ip_family":   "ipv4",
				"mode":        "static",
				"ip":          usrvCenterIp,
				"dns_servers": []string{"1.1.1.1", "8.8.8.8"},
				"prefix":      prefix_netmask,
				"gateway":     vCenter_gateway,
				"system_name": system_name_vCenter,
			},
			"os": map[string]any{
				"password":    vCenter_management_pass,
				"ntp_servers": []string{"0.pool.ntp.org", "1.pool.ntp.org"},
				"ssh_enable":  false,
			},
			"sso": map[string]any{
				"password":    vCenter_login_pass,
				"domain_name": "vsphere.local",
			},
		},
		"ceip": map[string]any{
			"settings": map[string]any{
				"ceip_enabled": false,
			},
		},
	}
	// fmt.Println(JsonData)
	// creating json file from jsonFile_content
	jsonBytes, err := json.MarshalIndent(JsonData, "", "  ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("production-vcsa.json", jsonBytes, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println(color.Green + "production-vcsa.json writed Successfully ..." + color.Reset)
	time.Sleep(1 * time.Second)

	fmt.Println(color.Blue + "=========================================================================================================" + color.Reset)
	fmt.Println(color.Blue + "=========================================================================================================" + color.Reset)
	fmt.Println(color.Blue + "=========================================================================================================" + color.Reset)
	fmt.Println(color.Yellow + "\nInstalling vCenter Product ..." + color.Reset)
	time.Sleep(1 * time.Second)
	Installing_vCenter(wdir)
}

func Installing_vCenter(wdir string) {

	// ./vcsa-deploy install --accept-eula --no-ssl-certificate-verification ~/myfiles/test/vcsa-auto-deploy/production-vcsa.json
	currentDir, _ := os.Getwd()
	jsonPath := currentDir + "/production-vcsa.json"

	fmt.Println(color.Yellow + "Executing ..." + color.Reset)
	cmd := exec.Command("./vcsa-deploy", "install", "--accept-eula", "--no-ssl-certificate-verification", jsonPath)
	cmd.Dir = currentDir + wdir + "/iso_files/source/vcsa-cli-installer/lin64"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(color.Green + "vCenter installed Successfully ...." + color.Reset)
	fmt.Println(color.Yellow + "\nExiting ..." + color.Reset)
	time.Sleep(2 * time.Second)
	os.Exit(0)
}

// ==================================================================================== vCenter setup (END) ==========================================================================================
