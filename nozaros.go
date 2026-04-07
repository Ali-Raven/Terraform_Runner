package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/TwiN/go-color"
)

var (
	additionalNetwork_name    string
	additionalNetwork_ip      string
	additionalNetwork_netmask string
	ManagementNetworkName     string
	ManagementNetworkIP       string
	ManagementNetworkNetmask  string
	vms                       []VM
)

type Network struct {
	Name    string `json:"name"`
	IP      string `json:"ip"`
	Netmask int    `json:"netmask"`
}

type VM struct {
	Name       string    `json:"name"`
	NumCPU     int       `json:"num_cpus"`
	MemoryGB   int       `json:"memory_gb"`
	Gateway    string    `json:"gateway"`
	DNSservers []string  `json:"dns_servers"`
	Networks   []Network `json:"networks"`
}

type TFvars struct {
	VMs []VM `json:"vms"`
}

func Nozaros_configure(wdir string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println()
	fmt.Println(color.Yellow + "\n================" + color.Reset)
	fmt.Println(color.Yellow + "\nOptions : " + color.Reset)
	fmt.Println("1. create new VMs \n2. Modify existing VMs\n3. Delete VMs\n4. Generating Inventory.yml file\n5. Main menu\n6. Exit")
	fmt.Print("\nSelect an option (1-5) : ")
	optionStr, _ := reader.ReadString('\n')
	optionStr = strings.TrimSpace(optionStr)

	switch optionStr {
	case "1":
		createNewVMs(reader, wdir)
	case "2":
		ModifyVMs(reader, wdir)
	case "3":
		DeleteVMs(reader, wdir)
	case "4":
		Yml(wdir, vms)
		time.Sleep(1 * time.Second)
		Nozaros_configure(wdir)
	case "5":
		fmt.Println(color.Yellow + "\nReturning to main menu..." + color.Reset)
		time.Sleep(1 * time.Second)
		main()
	case "6":
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

		// numVMstr, _ := reader.ReadString('\n')
		// numVMstr = strings.TrimSpace(numVMstr)
		// numVMcount := atoi(numVMstr)

		numVMcount, err := readInt(reader)
		if err != nil {
			fmt.Println(err)
			fmt.Println(color.Red + "Enter Numbers Please." + color.Reset)
			fmt.Println("\nReturning to menu ....")
			time.Sleep(1 * time.Second)
			Nozaros_configure(wdir)
		}

		// var vms []VM
		vms, errVM := loadExistingVMs(wdir)
		if errVM != nil {
			panic(errVM)
		}
		for i := 0; i < numVMcount; i++ {
			fmt.Printf(color.Yellow+"\n--- VM %d ---\n"+color.Reset, i+1)
			vm := collectVM(reader)
			vms = append(vms, vm)
		}

		retrunedPreview := preview(vms, reader, wdir)
		if retrunedPreview == 0 {
			return
		}

	}
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
		vm := collectVM(reader)
		vms = append(vms[:len(vms)-1], vm)

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

func collectVM(reader *bufio.Reader) VM {
	name := readRequired(reader, "Enter VM Name: ")
	numCPUstr := readRequired(reader, "Enter Number of CPUs: ")
	memoryGBstr := readRequired(reader, "Enter Memory in GB: ")
	gateway := readRequired(reader, "Enter Gateway: ")

	fmt.Print("Enter DNS servers : (Default : 1.1.1.1 , 1.0.0.1) ==> ")
	dnsStr, _ := reader.ReadString('\n')
	dnsStr = strings.TrimSpace(dnsStr)
	fmt.Println("--------")

	// Set default DNS servers if user input is empty
	dns := []string{"1.1.1.1", "1.0.0.1"}

	if dnsStr != "" {
		dns = strings.Split(dnsStr, ",")
		for i := range dns {
			dns[i] = strings.TrimSpace(dns[i])
		}
	}
	vm := VM{
		Name:       strings.TrimSpace(name),
		NumCPU:     atoi(numCPUstr),
		MemoryGB:   atoi(memoryGBstr),
		Gateway:    strings.TrimSpace(gateway),
		DNSservers: dns,
	}

	// Collecting Management Network ===========================================================================================================
	fmt.Println(color.Yellow + "Setting up the Management Network (VM Network Portgroup) ..." + color.Reset)

	time.Sleep(1 * time.Second)

	fmt.Print("\nEnter Management Network Name (default ==> VM Network) :  ")
	ManagementNetworkName, _ = reader.ReadString('\n')
	ManagementNetworkName = strings.TrimSpace(ManagementNetworkName)
	fmt.Println("--------")
	if ManagementNetworkName == "" {
		ManagementNetworkName = "VM Network"

		ManagementNetworkIP = readRequired(reader, "Enter Management Network IP : ")
		ManagementNetworkNetmask = readRequired(reader, "Enter Management Network Netmask : ")

		vm.Networks = append(vm.Networks, Network{
			Name:    strings.TrimSpace(ManagementNetworkName),
			IP:      strings.TrimSpace(ManagementNetworkIP),
			Netmask: atoi(ManagementNetworkNetmask),
		})
	} else {
		ManagementNetworkIP = readRequired(reader, "Enter Management Network IP : ")
		ManagementNetworkNetmask = readRequired(reader, "Enter Management Network Netmask : ")

		vm.Networks = append(vm.Networks, Network{
			Name:    strings.TrimSpace(ManagementNetworkName),
			IP:      strings.TrimSpace(ManagementNetworkIP),
			Netmask: atoi(ManagementNetworkNetmask),
		})
	}

	// end of collecting Management Network =====================================================================================================

	fmt.Println(color.Yellow + "Do you want to add additional Networks (VLANs or Portgroups) ? " + color.Reset)
	fmt.Print("Enter 'yes' or 'Enter' to add or 'no' or 'n' or press any key to skip: ")

	additionalNetChoice, _ := reader.ReadString('\n')
	additionalNetChoice = strings.TrimSpace(strings.ToLower(additionalNetChoice))

	if additionalNetChoice == "yes" || additionalNetChoice == "y" || additionalNetChoice == " " {
		vm.Networks = append(vm.Networks, readAdditionalNetworks(reader)...)
	} else {
		fmt.Println(color.Yellow + "\nSkipping additional Networks..." + color.Reset)
	}

	return vm
}

func readAdditionalNetworks(reader *bufio.Reader) []Network {
	var vmNetworks []Network
	fmt.Print(color.Yellow + "How many Networks do you want for your VMs ? " + color.Reset)
	numNetworkStr, _ := reader.ReadString('\n')
	netCount := atoi(numNetworkStr)

	for j := 0; j < netCount; j++ {
		fmt.Printf(color.Yellow+"\n--- Network %d ---\n"+color.Reset, j+1)
		additionalNetwork_name = readRequired(reader, "Network name: ")
		additionalNetwork_ip = readRequired(reader, "Network IP: ")
		additionalNetwork_netmask = readRequired(reader, "Network Netmask: ")

		// Yml(additionalNetwork_ip)
		vmNetworks = append(vmNetworks, Network{
			Name:    strings.TrimSpace(additionalNetwork_name),
			IP:      strings.TrimSpace(additionalNetwork_ip),
			Netmask: atoi(additionalNetwork_netmask),
		})
	}
	return vmNetworks
}

// =========================================================================== Creating New VMs (END) ==========================================================================

// =========================================================================== Helper Functions ==========================================================================
// for simple string to int conversion
func atoi(s string) int {
	i, _ := strconv.Atoi(strings.TrimSpace(s))
	return i
}

func readLine(reader *bufio.Reader) string {
	s, _ := reader.ReadString('\n')
	return strings.TrimSpace(s)
}

func readInt(reader *bufio.Reader) (int, error) {
	var input []byte
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return 0, err
		}

		if b == '\n' || b == ' ' || b == '\r' {
			break
		}
		input = append(input, b)
	}

	// converting bytes to strings => then to integer
	numStr := string(input)
	num, err := strconv.Atoi(strings.TrimSpace(numStr))
	if err != nil {
		return 0, fmt.Errorf("%sinvalid Integer (You enter => %s) %s", color.Red, numStr, color.Reset)
	}

	return num, nil
}

// =========================================================================== Helper Functions (END) ==========================================================================

// =========================================================================== Modify VMs ==========================================================================
func ModifyVMs(reader *bufio.Reader, wdir string) {
	fmt.Println(color.Yellow + "\nfetching list of existings vms ..." + color.Reset)
	time.Sleep(1 * time.Second)

	tfvars, err := loadTFvars(wdir)
	if err != nil {
		fmt.Println(color.Red + "Failed to load terraform.tfvars.json file" + color.Reset)
		return
	}

	// fetching the list of existing vms
	GettingVMsLists(tfvars)

	fmt.Printf("\nEnter the ID of the VM you want to modify :  %s(enter 0 to Return to menu )%s => ", color.Yellow, color.Reset)
	vmID, _ := reader.ReadString('\n')
	vmID = strings.TrimSpace(vmID)
	vmIDindex := atoi(vmID) - 1

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
	fmt.Printf("Modifying VM: %s (Functionality not yet implemented)\n", vmID)
	time.Sleep(2 * time.Second)

	tfvars.VMs[vmIDindex] = editVMs(reader, tfvars.VMs[vmIDindex])

	saveNewTFvars(tfvars, wdir)
}

func editVMs(reader *bufio.Reader, vm VM) VM {
	fmt.Println(color.Yellow + "\nPress ENTER to keep current value" + color.Reset)

	vm.Name = readOptionalValue(reader, "VM Name : ", vm.Name)
	vm.NumCPU = readOptionalINT(reader, "Number of CPU : ", vm.NumCPU)
	vm.MemoryGB = readOptionalINT(reader, "Memory (GB): ", vm.MemoryGB)
	vm.Gateway = readOptionalValue(reader, "Gateway : ", vm.Gateway)
	vm.DNSservers = readDNSserversValue(reader, "DNS servers : ", vm.DNSservers)
	vm.Networks = readNetworks(reader, vm.Networks)

	return vm
}
func readNetworks(reader *bufio.Reader, network []Network) []Network {
	// condition for checking the length of network array
	fmt.Println("\nNetwork Options : \n1. Modify Existing Networks value\n2. Add Network to the List\n3. Delete Network")

	fmt.Print("\nEnter your choice : (1/2/3) ")
	usrInput, _ := reader.ReadString('\n')
	usrInput = strings.TrimSpace(usrInput)

	switch usrInput {
	case "1":
		break
	case "2":
		network = append(network, readAdditionalNetworks(reader)...)
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
		id := atoi(readLine(reader)) - 1

		network = append(network[:id], network[id+1:]...)
		fmt.Println(color.Green + "Network is removed Successfully" + color.Reset)
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
	id := atoi(readLine(reader)) - 1

	if id < 0 || id >= len(network) {
		fmt.Println("Invalid ID")
		return network
	}
	// } else if (id + 1) == 0 {
	// 	network = append(network, readAdditionalNetworks(reader)...)
	// 	fmt.Println(color.Green + "Network Added Successfully ." + color.Reset)
	// 	return network
	// } else if id == 100 {
	// 	network = append(network[:id], network[id+1:]...)
	// 	fmt.Println(color.Green + "Network is removed Successfully" + color.Reset)
	// 	return network
	// }

	network[id].Name = readOptionalValue(reader, "Network Name : ", network[id].Name)
	network[id].IP = readOptionalValue(reader, "Network IP : ", network[id].IP)
	network[id].Netmask = readOptionalINT(reader, "Network Netmask : ", network[id].Netmask)

	return network
}
func deleteNetwork() {

}
func readOptionalValue(reader *bufio.Reader, label, current string) string {
	fmt.Printf("%s [current value => %s]: ", label, current)
	userInput, _ := reader.ReadString('\n')
	userInput = strings.TrimSpace(userInput)

	if userInput == "" {
		return current
	}
	return userInput
}
func readOptionalINT(reader *bufio.Reader, label string, current int) int {
	fmt.Printf("%s [current value => %d]: ", label, current)
	userInput, _ := reader.ReadString('\n')
	userInput = strings.TrimSpace(userInput)

	if userInput == "" {
		return current
	}
	return atoi(userInput)
}
func readDNSserversValue(reader *bufio.Reader, label string, current []string) []string {
	fmt.Printf("%s [current value => %s]: ", label, current)
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

func printVMBox(vm VM, index int) {
	fmt.Println("┌──────────────────────────────────────┐")
	fmt.Printf("│ %sVM ID          : %-20d %s│\n", color.Yellow, index+1, color.Reset)
	fmt.Printf("│ VM Name        : %-20s │\n", vm.Name)
	fmt.Printf("│ CPUs           : %-20d │\n", vm.NumCPU)
	fmt.Printf("│ Memory (GB)    : %-20d │\n", vm.MemoryGB)
	fmt.Printf("│ Gateway        : %-20s │\n", vm.Gateway)
	fmt.Printf("│ DNS Servers    : %-20s │\n", strings.Join(vm.DNSservers, ","))
	fmt.Println("├──────────────────────────────────────┤")
	fmt.Println("│ Networks                             │")
	fmt.Println("├──────────────┬──────────────┬────────┤")
	fmt.Println("│ Name         │ IP           │ Mask   │")
	fmt.Println("├──────────────┼──────────────┼────────┤")

	for _, n := range vm.Networks {
		fmt.Printf("│ %-12s │ %-12s │ %-6d │\n",
			n.Name, n.IP, n.Netmask)
	}

	fmt.Println("└──────────────┴──────────────┴────────┘")
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

	GettingVMsLists(tfvars)

	fmt.Printf("\nEnter VM ID that you want to Delete : %s(enter 0 to return to menu)%s => ", color.Yellow, color.Reset)
	// vmID := atoi(readLine(bufio.NewReader(os.Stdin))) - 1
	vmID, _ := reader.ReadString('\n')
	vmID = strings.TrimSpace(vmID)
	vmIDndex := atoi(vmID) - 1

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
