package main

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/TwiN/go-color"
	"github.com/joho/godotenv"
	"github.com/terraform_runner/helper"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
)

var (
	additionalNetwork_name    string
	additionalNetwork_ip      string
	additionalNetwork_gateway string
	additionalNetwork_netmask string
	ManagementNetworkName     string
	ManagementNetworkIP       string
	ManagementNetworkNetmask  string
	vms                       []VM
	VmName                    string
	VmTargetHost              string
	DataStore                 string
	ClusterName               string
	numCPUstr                 string
	memoryGBstr               string
	VmGateway                 string
	VmDns                     string
	dnsStr                    string
	componentName             string
	additionalNetChoice       string
	numNetworkStr             string
	componentsToConnect       []string
	backupComps               string
	vipForBackupComps         string
	BackupTargetHost          string
	BackupDataStore           string
	vCenterURL                string
	vCenterUserName           string
	vCenterPass               string
)

var (
	MMEConnectedComps  []string = []string{"HSS1", "HSS2", "HSS3", "HSS4", "HSS5", "SGWC1", "SGWC2", "SMF1", "SMF2"}
	HSSConnectedComps  []string = []string{"MME1", "MME2", "MME3", "MME4", "MME5", "MME6", "MME7", "MME8", "MME9", "MME10", "MME11", "MME12", "MME13"}
	SGWCConnectedComps []string = []string{"MME1", "MME2", "MME3", "MME4", "MME5", "MME6", "MME7", "MME8", "MME9", "MME10", "MME11", "MME12", "MME13", "SMF1", "SMF2", "SGWU1", "SGWU2"}
	SGWUConnectedComps []string = []string{"SGWC1", "SGWC2", "UPF1", "UPF2"}
	SMFConnectedComps  []string = []string{"SGWC1", "SGWC2", "PCRF1", "UPF1", "UPF2"}
	UPFConnectedComps  []string = []string{"SMF1", "SMF2", "SGWU1", "SGWU2"}
	PCRFConnectedComps []string = []string{"SMF1", "SMF2"}
)

type Network struct {
	ID      string `json:"uuid"`
	Name    string `json:"name"`
	IP      string `json:"ip"`
	Gateway string `json:"gateway"`
	Netmask int    `json:"netmask"`
}

type VM struct {
	ID         string `json:"uuid"`
	Name       string `json:"name"`
	TargetHost string `json:"target_host"`
	DataStore  string `json:"datastore"`
	// ClusterName         string    `json:"cluster_name"`
	NumCPU              int       `json:"num_cpus"`
	MemoryGB            int       `json:"memory_gb"`
	Gateway             string    `json:"gateway"`
	DNSservers          []string  `json:"dns_servers"`
	Component           string    `json:"component"`
	ComponentsToConnect []string  `json:"componentsToConnect"`
	Networks            []Network `json:"network_adaptors"`
}

type TFvars struct {
	VMs []VM `json:"vms"`
}

func Nozaros_configure(wdir string) {
	reader := bufio.NewReader(os.Stdin)

	// getting vCenter server info
	if err := godotenv.Load(); err != nil {
		fmt.Printf("%s Can't read or load the vCenter Info , try again %s", color.Red, color.Reset)
	}

	vCenterURL = os.Getenv("vCenterURL")
	vCenterUserName = os.Getenv("vCenterUserName")
	vCenterPass = os.Getenv("vCenterPassword")

	fmt.Println()
	fmt.Println(color.Yellow + "\n================" + color.Reset)
	fmt.Println(color.Yellow + "\nOptions : \n" + color.Reset)
	var choice []string
	choice = []string{"create new VMs", "Modify existing VMs", "Delete VMs", "Generating Inventory.yml file", "Test enviroment", "Main menu", "Exit"}
	optionStr := helper.AskSelect(choice)
	// fmt.Println("1. create new VMs \n2. Modify existing VMs\n3. Delete VMs\n4. Generating Inventory.yml file\n5. Main menu\n6. Exit")
	// fmt.Print("\nSelect an option (1-5) : ")

	switch optionStr {
	case "create new VMs":
		createNewVMs(reader, wdir)
	case "Modify existing VMs":
		ModifyVMs(reader, wdir)
	case "Delete VMs":
		DeleteVMs(reader, wdir)
	case "Generating Inventory.yml file":
		Yml(wdir, vms)
		time.Sleep(1 * time.Second)
		Nozaros_configure(wdir)
	case "Test enviroment":
		err := godotenv.Load()
		if err != nil {
			panic(err)
		}

		vCneterURL := os.Getenv("vCenterURL")
		vCenterUser := os.Getenv("vCenterUserName")
		vCenterPass := os.Getenv("vCenterPassword")

		fmt.Println(vCneterURL, vCenterUser, vCenterPass)
	case "Main menu":
		fmt.Println(color.Yellow + "\nReturning to main menu..." + color.Reset)
		time.Sleep(1 * time.Second)
		main()
	case "Exit":
		fmt.Println("Exiting...")
		time.Sleep(1 * time.Second)
		os.Exit(0)
	default:
		fmt.Println(color.Yellow + "\nWarning : choose one of the above options ..." + color.Reset)
		fmt.Println(color.Yellow + "Returning to menu ..." + color.Reset)
		time.Sleep(1 * time.Second)
		Nozaros_configure(wdir)
	}
	time.Sleep(1 * time.Second)
	fmt.Printf("%s%s Updated Successfully. %s\n", color.Green, "terraform.tfvars.json", color.Reset)
	Nozaros_configure(wdir)
}

// =========================================================================== Creating New VMs ==========================================================================
func createNewVMs(reader *bufio.Reader, wdir string) {
	for {
		fmt.Println(color.Yellow + "\nCreating new VMs..." + color.Reset)
		time.Sleep(1 * time.Second)
		fmt.Print("\n\u2731 How many VMs you want to create ? ")

		numVMcount, err := helper.ReadInt(reader)
		if err != nil {
			fmt.Println(err)
			fmt.Println(color.Red + "Enter Numbers Please." + color.Reset)
			fmt.Println("\nReturning to menu ....")
			time.Sleep(1 * time.Second)
			Nozaros_configure(wdir)
		}

		vms, errVM := loadExistingVMs(wdir)
		if errVM != nil {
			panic(errVM)
		}
		for i := 0; i < numVMcount; i++ {
			fmt.Printf("%s%s Connecting to vCenter server ...%s%s" , color.Bold , color.Yellow , color.Reset , color.Reset)

			ds, dErr := GetAllDatastoreName(vCenterURL, vCenterUserName, vCenterPass)
			if dErr != nil {
				fmt.Printf("Error: %v\n", dErr)
				os.Exit(1)
			}
			fmt.Printf(color.Yellow+"\n--- VM %d ---\n"+color.Reset, i+1)

			collected := collectVM(reader , ds)

			for _, vm := range collected {
				vms = append(vms, vm)
			}
		}

		retrunedPreview := preview(vms, reader, wdir)
		if retrunedPreview == 0 {
			return
		}

	}
}

func HashGenerator(name string) string {
	if name == "" {
		return "unkown"
	}

	hash := sha256.Sum256([]byte(name))
	return hex.EncodeToString(hash[:])[:12]
}

func loadExistingVMs(wdir string) ([]VM, error) {
	tfvars, err := loadTFvars(wdir)
	if err != nil {
		fmt.Println(color.Yellow + "No existing VMs found. Starting fresh..." + color.Reset)
		return []VM{}, err
	}
	return tfvars.VMs, nil
}
func preview(vms []VM, reader *bufio.Reader, wdir string) int {
	data := TFvars{VMs: vms}
	jsonBytes, _ := json.MarshalIndent(data, "", "  ")

	fmt.Println(color.Green + "\n============ PREVIEW ============" + color.Reset)
	fmt.Println(string(jsonBytes))
	fmt.Println(color.Green + "==================================" + color.Reset)

	choice := confirmMenu(reader)

	switch choice {
	case "1":
		// Save the updated VMs to the file in the working directory
		currentDir, _ := os.Getwd()
		filename := currentDir + wdir + "/terraform.tfvars.json"
		os.WriteFile(filename, jsonBytes, 0644)
		// Yml(currentDir)
		return 0
	case "2":
		fmt.Println(color.Yellow + "\nRe-enter VM data...\n" + color.Reset)

		collected := collectVM(reader , )

		for _, vm := range collected {
			vms = append(vms, vm)
		}

		// fmt.Println(vms)
		preview(vms, reader, wdir)
	case "3":
		fmt.Println(color.Red + "Canceled ❌" + color.Reset)
		fmt.Println("returning to the main menu ...")
		time.Sleep(2 * time.Second)
		main()
	}
	return 0
}
func confirmMenu(reader *bufio.Reader) string {
	for {
		fmt.Println("\nOptions:")
		fmt.Println("1) Approve & write")
		fmt.Println("2) Edit")
		fmt.Println("3) Cancel")
		fmt.Print("Choose option: ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		if choice == "1" || choice == "2" || choice == "3" {
			return choice
		}

		fmt.Println("Invalid choice.")
	}
}

func readRequired(reader *bufio.Reader, label string) string {
	for {
		fmt.Print(label)
		value, _ := reader.ReadString('\n')
		value = strings.TrimSpace(value)
		fmt.Println("--------")

		if value != "" {
			return value
		}

		fmt.Println(color.Red + "This field is required." + color.Reset)
	}
}

func GetAllDatastoreName(vcURL, username, pass string) ([]string, error) {
	ctx := context.Background()

	// 1. Format the vCenter URL
	u, err := url.Parse(fmt.Sprintf("https://%s:%s@%s/sdk", username, pass, vcURL))

	if err != nil {
		return nil, fmt.Errorf("invalid vCenter URL: %w", err)
	}

	// 2. Connect to vCenter
	client, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to vCenter: %w", err)
	}
	defer client.Logout(ctx)

	// 3. Initialize the Finder
	finder := find.NewFinder(client.Client, true)

	// 4. Find the default datacenter to scope our search
	dc, err := finder.DefaultDatacenter(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find default datacenter: %w", err)
	}
	finder.SetDatacenter(dc)

	// 5. Find ALL datastores using the "*" wildcard
	datastores, err := finder.DatastoreList(ctx, "*")
	if err != nil {
		return nil, fmt.Errorf("failed to list datastores: %w", err)
	}

	// 6. Extract the names into a simple string slice
	var dsNames []string
	for _, ds := range datastores {
		dsNames = append(dsNames, ds.Name())
	}

	return dsNames, nil
}

func isVMNameIsCoreComponent(CompsConnectedItems []string) {
	componentsToConnect = helper.MultiSelect("which Components you want to Connect ?", CompsConnectedItems)
	backupComps = helper.Ask("Do you Want Backup Option for this Component ?", backupComps)
	if backupComps == "y" {
		vipForBackupComps = helper.Ask("Enter Your Management IP address for Backup Components :", vipForBackupComps)
		BackupTargetHost = helper.Ask("Enter Target Host for Backup VM :", BackupTargetHost)
		BackupDataStore = helper.Ask("Enter Datastore for Backup VM :", BackupDataStore)

		// checking target host and datastore of backup vm  equallity to main VM
		if BackupTargetHost == VmTargetHost || BackupDataStore == DataStore {
			fmt.Printf("%sError : Backup Target host or Datastore can't be the same with main Target host or Datastore Name.%s\n", color.Red, color.Reset)
			fmt.Println(color.Yellow+"try this section again with True value ", color.Reset)
			fmt.Println(color.Yellow+"Returning ...", color.Reset)
			time.Sleep(1 * time.Second)
			isVMNameIsCoreComponent(CompsConnectedItems)
		}
	} else {
		fmt.Println(color.Yellow + "Skipping ..." + color.Reset)
		time.Sleep(700 * time.Millisecond)
	}
}

func collectVM(reader *bufio.Reader , ds []string) []VM {

	VmName = helper.Ask("Enter VM Name:", VmName)
	VmTargetHost = helper.Ask("Enter Target Host of VM :", VmTargetHost)
	DataStore = helper.AskSelect(ds)
	// ClusterName = helper.AskSelect([]string{"Cluster1"})
	numCPUstr = helper.Ask("Enter Number of CPUs:", numCPUstr)
	memoryGBstr = helper.Ask("Enter Memory in GB:", memoryGBstr)
	VmGateway = helper.Ask("Enter Gateway:", VmGateway)
	componentName = helper.Ask("Enter Component Name:", componentName)

	switch {
	case strings.Contains(strings.ToUpper(componentName), "MME"):
		isVMNameIsCoreComponent(MMEConnectedComps)
	case strings.Contains(strings.ToUpper(componentName), "HSS"):
		isVMNameIsCoreComponent(HSSConnectedComps)
	case strings.Contains(strings.ToUpper(componentName), "SGWC"):
		isVMNameIsCoreComponent(SGWCConnectedComps)
	case strings.Contains(strings.ToUpper(componentName), "SGWU"):
		isVMNameIsCoreComponent(SGWUConnectedComps)
	case strings.Contains(strings.ToUpper(componentName), "SMF"):
		isVMNameIsCoreComponent(SMFConnectedComps)
	case strings.Contains(strings.ToUpper(componentName), "UPF"):
		isVMNameIsCoreComponent(UPFConnectedComps)
	case strings.Contains(strings.ToUpper(componentName), "PCRF"):
		isVMNameIsCoreComponent(PCRFConnectedComps)
	default:

	}

	// generate hash
	genHash := HashGenerator(VmName)
	// Set default DNS servers if user input is empty

	dnsStr = helper.AskThHasDefaultVal("Enter DNS servers:", dnsStr, "1.1.1.1 , 1.0.0.1")

	splitDnsStr := strings.Split(dnsStr, ",")
	for i := range splitDnsStr {
		splitDnsStr[i] = strings.TrimSpace(splitDnsStr[i])
	}

	vm := VM{
		ID:         genHash,
		Name:       strings.TrimSpace(VmName),
		TargetHost: strings.TrimSpace(VmTargetHost),
		DataStore:  strings.TrimSpace(DataStore),
		// ClusterName:         strings.TrimSpace(ClusterName),
		NumCPU:              helper.Atoi(numCPUstr),
		MemoryGB:            helper.Atoi(memoryGBstr),
		Gateway:             strings.TrimSpace(VmGateway),
		DNSservers:          splitDnsStr,
		Component:           componentName,
		ComponentsToConnect: componentsToConnect,
	}

	// Collecting Management Network ===========================================================================================================
	fmt.Println(color.Yellow + "\nSetting up the Management Network (VM Network Portgroup) ..." + color.Reset)

	time.Sleep(1 * time.Second)

	fmt.Print("\nEnter Management Network Name (default ==> VM Network) :  ")

	ManagementNetworkName = helper.AskThHasDefaultVal("Enter Management Network Name:", ManagementNetworkName, "VM Network")
	GenHashNetwork := HashGenerator(ManagementNetworkName)
	ManagementNetworkIP = helper.Ask("Enter Management Network IP:", ManagementNetworkIP)
	ManagementNetworkNetmask = helper.AskThHasDefaultVal("Enter Management Network Netmask:", ManagementNetworkNetmask, "24")
	// ManagementNetworkName, _ = reader.ReadString('\n')
	// ManagementNetworkName = strings.TrimSpace(ManagementNetworkName)
	// fmt.Println("--------")
	vm.Networks = append(vm.Networks, Network{
		ID:      GenHashNetwork,
		Name:    strings.TrimSpace(ManagementNetworkName),
		IP:      strings.TrimSpace(ManagementNetworkIP),
		Netmask: helper.Atoi(ManagementNetworkNetmask),
	})

	// end of collecting Management Network =====================================================================================================

	fmt.Println(color.Yellow + "Do you want to add additional Networks (VLANs or Portgroups) ? " + color.Reset)

	// create validator
	Validator := func(ans interface{}) error {
		value := strings.ToLower(strings.TrimSpace(ans.(string)))

		switch value {
		case "yes", "y", "Y", "":
			vm.Networks = append(vm.Networks, readAdditionalNetworks()...)
		case "no", "n", "NO", "N":
			return nil
		default:
			return fmt.Errorf("Undefined value")
		}
		return nil
	}

	prompt := survey.Input{
		Message: "Enter 'yes' or 'Enter' to add or 'no' or 'n' or press any key to skip:",
	}

	err := survey.AskOne(
		&prompt,
		&additionalNetChoice,
		survey.WithValidator(Validator),
	)

	if err != nil {
		panic(err)
	}

	vmList := []VM{vm}

	if backupComps == "y" && vipForBackupComps != "" {
		backupVM := createBackupVM(vm, vipForBackupComps, BackupTargetHost, BackupDataStore)
		vmList = append(vmList, backupVM)
		fmt.Println(color.Green + "✓ Backup VM created successfully!" + color.Reset)
	}

	return vmList
}

func createBackupVM(originalVM VM, vIP, bTarget, bDataStore string) VM {

	backup := VM{
		ID:         HashGenerator(originalVM.Name + "-backup"),
		Name:       originalVM.Name + "-backup",
		TargetHost: bTarget,
		DataStore:  bDataStore,
		// ClusterName:         originalVM.ClusterName,
		NumCPU:              originalVM.NumCPU,
		MemoryGB:            originalVM.MemoryGB,
		Gateway:             originalVM.Gateway,
		DNSservers:          originalVM.DNSservers,
		Component:           originalVM.Component,
		ComponentsToConnect: originalVM.ComponentsToConnect,
	}

	// Update only Management IP
	for _, net := range originalVM.Networks {
		newNet := net
		if strings.Contains(strings.ToUpper(net.Name), ManagementNetworkName) ||
			net.Name == ManagementNetworkName {
			newNet.IP = strings.TrimSpace(vIP)
			newNet.ID = HashGenerator(net.Name + "-backup")
		}
		backup.Networks = append(backup.Networks, newNet)
	}

	backup.Component = originalVM.Component + "(Backup)"
	return backup
}

func readAdditionalNetworks() []Network {
	var vmNetworks []Network
	fmt.Println(color.Bold + "\n\nadding additional networks ...\n" + color.Reset)
	time.Sleep(1 * time.Second)
	// fmt.Print(color.Yellow + "How many Networks do you want for your VMs ? " + color.Reset)

	numNetworkStr = helper.Ask("How many Networks do you want for your VMs ?", numNetworkStr)
	// numNetworkStr, _ := reader.ReadString('\n')
	netCount := helper.Atoi(numNetworkStr)

	for j := 0; j < netCount; j++ {

		fmt.Printf(color.Yellow+"\n--- Network %d ---\n"+color.Reset, j+1)
		additionalNetwork_name = helper.Ask("Network name:", additionalNetwork_name)
		additionalNetwork_ip = helper.Ask("Network IP:", additionalNetwork_ip)
		additionalNetwork_gateway = helper.Ask("Network Gateway :", additionalNetwork_gateway)
		additionalNetwork_netmask = helper.Ask("Network Netmask:", additionalNetwork_netmask)

		GenHashAdditionNetwork := HashGenerator(additionalNetwork_name)
		// Yml(additionalNetwork_ip)
		vmNetworks = append(vmNetworks, Network{
			ID:      GenHashAdditionNetwork,
			Name:    strings.TrimSpace(additionalNetwork_name),
			IP:      strings.TrimSpace(additionalNetwork_ip),
			Gateway: strings.TrimSpace(additionalNetwork_gateway),
			Netmask: helper.Atoi(additionalNetwork_netmask),
		})
	}
	return vmNetworks
}

// =========================================================================== Creating New VMs (END) ==========================================================================

// =========================================================================== Modify VMs ==========================================================================
func ModifyVMs(reader *bufio.Reader, wdir string) {
	fmt.Println(color.Bold + color.Yellow + "\nfetching list of existings vms ..." + color.Reset + color.Reset)
	time.Sleep(1 * time.Second)
	// loadExistingVMs(wdir)

	tfvars, err := loadTFvars(wdir)
	if err != nil {
		fmt.Println(color.Red + "Failed to load terraform.tfvars.json file" + color.Reset)
		return
	}

	// fmt.Println(tfvars.VMs == )
	// checking existing VMs
	if len(tfvars.VMs) == 0 {
		fmt.Println(color.Yellow + "No existing VMs for Modifying. Starting fresh..." + color.Reset)
		return
	}

	// fetching the list of existing vms
	GettingVMsLists(tfvars)

	fmt.Printf("\nEnter the ID of the VM you want to modify :  %s(enter 0 to Return to menu )%s => ", color.Yellow, color.Reset)
	vmID, _ := reader.ReadString('\n')
	vmID = strings.TrimSpace(vmID)
	vmIDindex := helper.Atoi(vmID) - 1

	// checking VM index ID is valid or not
	if vmIDindex >= len(tfvars.VMs) {
		fmt.Println(color.Red + "Invalid VM ID" + color.Reset)
		return
	} else if vmID == "0" {
		fmt.Println(color.Yellow + "\nReturning to menu ..." + color.Reset)
		time.Sleep(1 * time.Second)
		Nozaros_configure(wdir)
	}
	// Further implementation to modify the VM with the given name
	fmt.Printf("Modifying VM:%s%s %s %s%s\n", color.Bold, color.Yellow, vmID, color.Reset, color.Reset)
	time.Sleep(600 * time.Millisecond)

	tfvars.VMs[vmIDindex] = editVMs(reader, tfvars.VMs[vmIDindex], wdir)

	saveNewTFvars(tfvars, wdir)
	fmt.Printf("\n%s%sUpdating VMs list ...%s%s\n", color.Bold, color.Yellow, color.Reset, color.Reset)
	time.Sleep(1 * time.Second)
	ModifyVMs(reader, wdir)
}

func editVMs(reader *bufio.Reader, vm VM, wdir string) VM {
	fmt.Println(color.Yellow + "\nPress ENTER to keep current value" + color.Reset)

	vm.Name = readOptionalValue(reader, "VM Name : ", vm.Name)
	vm.TargetHost = readOptionalValue(reader, "VM Target Host : ", vm.TargetHost)
	vm.DataStore = readOptionalValue(reader, "VM Datastore Name : ", vm.DataStore)
	// vm.ClusterName = readOptionalValue(reader  , "VM Cluster Name : " , vm.ClusterName)
	vm.NumCPU = readOptionalINT(reader, "Number of CPU : ", vm.NumCPU)
	vm.MemoryGB = readOptionalINT(reader, "Memory (GB): ", vm.MemoryGB)
	vm.Gateway = readOptionalValue(reader, "Gateway : ", vm.Gateway)
	vm.DNSservers = readDNSserversValue(reader, "DNS servers : ", vm.DNSservers)
	vm.Component = readOptionalValue(reader, "Component Name : ", vm.Component)
	vm.ComponentsToConnect = readDNSserversValue(reader, "Which component to Connect ?", vm.ComponentsToConnect)
	vm.Networks = readNetworks(reader, vm.Networks, wdir)

	return vm
}
func readNetworks(reader *bufio.Reader, network []Network, wdir string) []Network {
	// condition for checking the length of network array
	fmt.Println("\nNetwork Options : \n1. Modify Existing Networks value\n2. Add Network to the List\n3. Delete Network\n4. Update and Exit")

	fmt.Print("\nEnter your choice : (1/2/3) ")
	usrInput, _ := reader.ReadString('\n')
	usrInput = strings.TrimSpace(usrInput)

	switch usrInput {
	case "1":
		break
	case "2":
		network = append(network, readAdditionalNetworks()...)
		fmt.Println(color.Green + "Network Added Successfully ." + color.Reset)
		return network
	case "3":
		if len(network) == 0 {
			fmt.Println(color.Yellow + "No Networks to edit" + color.Reset)
			return network
		}
		fmt.Println(color.Yellow + "\nlisting Networks ..." + color.Reset)
		time.Sleep(1 * time.Second)
		for i, n := range network {
			fmt.Printf("%d) Name : %s ,  IP ; (%s/%d)\n", i+1, n.Name, n.IP, n.Netmask)
		}

		fmt.Print("\nwhich one of Network you want to Delete ? (Enter ID) ")
		id := helper.Atoi(helper.ReadLine(reader)) - 1

		network = append(network[:id], network[id+1:]...)
		fmt.Println(color.Green + "Network is removed Successfully" + color.Reset)
		return network
	case "4":
		return network
	}
	if len(network) == 0 {
		fmt.Println(color.Yellow + "No Networks to edit" + color.Reset)
		return network
	}

	fmt.Println(color.Yellow + "\nlisting Networks ..." + color.Reset)
	time.Sleep(1 * time.Second)
	for i, n := range network {
		fmt.Printf("%d) Name : %s ,  IP ; (%s/%d)\n", i+1, n.Name, n.IP, n.Netmask)
	}
	// fmt.Println("0. Add Network")

	fmt.Print("\nEnter the network ID to edit (Or Add New Network) : ")
	id := helper.Atoi(helper.ReadLine(reader)) - 1

	if id < 0 || id >= len(network) {
		fmt.Println("Invalid ID")
		return network
	}

	network[id].Name = readOptionalValue(reader, "Network Name : ", network[id].Name)
	network[id].IP = readOptionalValue(reader, "Network IP : ", network[id].IP)
	network[id].Gateway = readOptionalValue(reader, "Network Gateway : ", network[id].Gateway)
	network[id].Netmask = readOptionalINT(reader, "Network Netmask : ", network[id].Netmask)

	return network
}

func readOptionalValue(reader *bufio.Reader, label, current string) string {
	fmt.Printf("%s%s [current value => %s]: %s", color.Bold, label, current, color.Reset)
	userInput, _ := reader.ReadString('\n')
	userInput = strings.TrimSpace(userInput)

	if userInput == "" {
		return current
	}
	return userInput
}
func readOptionalINT(reader *bufio.Reader, label string, current int) int {
	fmt.Printf("%s%s [current value => %d]: %s", color.Bold, label, current, color.Reset)
	userInput, _ := reader.ReadString('\n')
	userInput = strings.TrimSpace(userInput)

	if userInput == "" {
		return current
	}
	return helper.Atoi(userInput)
}
func readDNSserversValue(reader *bufio.Reader, label string, current []string) []string {
	fmt.Printf("%s%s [current value => %s]: %s", color.Bold, label, current, color.Reset)
	userInput, _ := reader.ReadString('\n')
	userInput = strings.TrimSpace(userInput)

	if userInput == "" {
		return current
	}

	dns := strings.Split(userInput, ",")
	for i := range dns {
		dns[i] = strings.TrimSpace(dns[i])
	}

	return dns
}

func saveNewTFvars(tfvars TFvars, wdir string) {
	data, _ := json.MarshalIndent(tfvars, "", " ")

	currentDir, _ := os.Getwd()
	filename := currentDir + wdir + "/terraform.tfvars.json"
	os.WriteFile(filename, data, 0644)
}

func loadTFvars(wdir string) (TFvars, error) {
	currentDir, _ := os.Getwd()
	filename := currentDir + wdir + "/terraform.tfvars.json"
	data, err := os.ReadFile(filename)
	if err != nil {
		return TFvars{}, err
	}

	var tfvars TFvars
	err = json.Unmarshal(data, &tfvars)
	if err != nil {
		fmt.Println(color.Red+"Error unmarshaling JSON:"+color.Reset, err)
		return tfvars, err
	}
	return tfvars, err
}
func GettingVMsLists(tfvars TFvars) {
	for i, vm := range tfvars.VMs {
		printVMBox(vm, i)
		fmt.Println(color.Green + "================================================" + color.Reset)
	}
}

func truncate(s string, max int) string {
	if len(s) > max {
		return s[:max-3] + "..."
	}
	return fmt.Sprintf("%-*s", max, s)
}

func printVMBox(vm VM, index int) {
	// Total width: 52 chars. Inner width: 50 chars.
	fmt.Println("┌──────────────────────────────────────────────────┐")

	// Label column is 20 chars. Value column is 25 chars.
	// We put color codes *outside* the padded variables so they don't break terminal spacing math.
	fmt.Printf("│ %-20s : %s%-25d%s │\n", "VM ID", color.Yellow, index+1, color.Reset)
	fmt.Printf("│ %-20s : %s │\n", "VM Name", truncate(vm.Name, 25))
	fmt.Printf("│ %-20s : %s │\n", "Target Host", truncate(vm.TargetHost, 25))
	fmt.Printf("│ %-20s : %s │\n", "Data Store", truncate(vm.DataStore, 25))
	fmt.Printf("│ %-20s : %-25d │\n", "CPUs", vm.NumCPU)
	fmt.Printf("│ %-20s : %-25d │\n", "Memory (GB)", vm.MemoryGB)
	fmt.Printf("│ %-20s : %s │\n", "Gateway", truncate(vm.Gateway, 25))

	dnsJoined := strings.Join(vm.DNSservers, ",")
	fmt.Printf("│ %-20s : %s │\n", "DNS Servers", truncate(dnsJoined, 25))

	fmt.Printf("│ %-20s : %s │\n", "Component Name", truncate(vm.Component, 25))

	compJoined := strings.Join(vm.ComponentsToConnect, ",")
	fmt.Printf("│ %-20s : %s │\n", "Components To Connect", truncate(compJoined, 25))

	fmt.Println("├──────────────────────────────────────────────────┤")
	fmt.Println("│ Networks                                         │")

	// Columns: 16 + 18 + 14 = 48. Plus 2 inner dividers = 50 chars exactly.
	fmt.Println("├────────────────┬──────────────────┬──────────────┤")
	fmt.Println("│ Name           │ IP               │ Mask         │")
	fmt.Println("├────────────────┼──────────────────┼──────────────┤")

	for _, n := range vm.Networks {
		fmt.Printf("│ %-14s │ %-16s │ %-12d │\n", truncate(n.Name, 14), truncate(n.IP, 16), n.Netmask)
	}

	fmt.Println("└────────────────┴──────────────────┴──────────────┘")
}

// =========================================================================== Modify VMs (END) ==========================================================================

// =========================================================================== Delete VMs ==========================================================================
func DeleteVMs(reader *bufio.Reader, wdir string) {
	fmt.Println(color.Yellow + "\nDelete VMs" + color.Reset)
	time.Sleep(2 * time.Second)
	tfvars, err := loadTFvars(wdir)
	if err != nil {
		fmt.Println(color.Red + "Failed to load terraform.tfvars.json file" + color.Reset)
		return
	}

	// checking existings of VMs
	if len(tfvars.VMs) == 0 {
		fmt.Println(color.Yellow + "No existing VMs for Deleting. Starting fresh..." + color.Reset)
		return
	}

	GettingVMsLists(tfvars)

	// fmt.Printf("\nEnter VM ID that you want to Delete : %s(enter 0 to return to menu)%s => ", color.Yellow, color.Reset)
	var vmID string
	vmID = helper.Ask("Enter VM ID that you want to Delete : (enter 0 to return to menu) =>", vmID)
	// vmID := helper.Atoi(readLine(bufio.NewReader(os.Stdin))) - 1
	// vmID, _ := reader.ReadString('\n')
	vmID = strings.TrimSpace(vmID)
	vmIDndex := helper.Atoi(vmID) - 1

	if vmIDndex >= len(tfvars.VMs) {
		fmt.Println(color.Red + "Invalid VM ID" + color.Reset)
		return
	} else if vmID == "0" {
		fmt.Println(color.Yellow + "\nReturning to menu ..." + color.Reset)
		time.Sleep(1 * time.Second)
		Nozaros_configure(wdir)
	}

	fmt.Printf("Deleting VM ==> %s%s%s\n", color.Cyan, vmID, color.Reset)
	time.Sleep(2 * time.Second)

	// deleting operation
	tfvars.VMs = append(tfvars.VMs[:vmIDndex], tfvars.VMs[vmIDndex+1:]...)
	saveNewTFvars(tfvars, wdir)

	fmt.Printf("\n%sVM with ID %s has been deleted Successfully%s\n\n", color.Green, vmID, color.Reset)
}

// =========================================================================== Delete VMs (END) ==========================================================================
