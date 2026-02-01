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

type Network struct {
	Name    string `json:"name"`
	IP      string `json:"ip"`
	Netmask int `json:"netmask"`
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

	fmt.Println("\nOptions : ")
	fmt.Println("1. create new VMs \n2. Modify existing VMs (Not implemented yet)\n3. Delete VMs (Not implemented yet)\n4. Exit")
	fmt.Print("\nSelect an option (1-4) : ")
	optionStr, _ := reader.ReadString('\n')
	optionStr = strings.TrimSpace(optionStr)

	switch optionStr {
	case "1":
		createNewVMs(reader, wdir)
	case "2":

	}
	time.Sleep(1 * time.Second)
	fmt.Printf("%s%s Updated Successfully %s\n", color.Green, "terraform.tfvars.json", color.Reset)
}

func createNewVMs(reader *bufio.Reader, wdir string) {
	for {
		fmt.Println("Creating new VMs...")
		fmt.Print("\nHow many VMs you want to create ? ")

		numVMstr, _ := reader.ReadString('\n')
		numVMstr = strings.TrimSpace(numVMstr)
		numVMcount := atoi(numVMstr)

		var vms []VM

		for i := 0; i < numVMcount; i++ {
			fmt.Printf(color.Yellow+"\n--- VM %d ---\n"+color.Reset, i+1)
			vm := collectVM(reader)
			vms = append(vms, vm)
		}

		data := TFvars{VMs: vms}
		jsonBytes, _ := json.MarshalIndent(data, "", "  ")

		fmt.Println(color.Green + "\n============ PREVIEW ============" + color.Reset)
		fmt.Println(string(jsonBytes))
		fmt.Println(color.Green + "==================================" + color.Reset)

		choice := confirmMenu(reader)

		if choice == "1" {
			os.WriteFile("terraform.tfvars.json", jsonBytes, 0644)
			return
		}

		if choice == "3" {
			fmt.Println(color.Red + "Canceled ❌" + color.Reset)
			fmt.Println("returning to the main menu ...")
			time.Sleep(2 * time.Second)
			main()
		}

		// choice == "2" → loop again (edit)
		fmt.Println(color.Yellow + "\nRe-enter VM data...\n" + color.Reset)
	}
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

		fmt.Println("❌ This field is required.")
	}
}

func collectVM(reader *bufio.Reader) VM {
	name := readRequired(reader, "Enter VM Name: ")
	numCPUstr := readRequired(reader, "Enter Number of CPUs: ")
	memoryGBstr := readRequired(reader, "Enter Memory in GB: ")
	gateway := readRequired(reader, "Enter Gateway : ")

	fmt.Print("Enter DNS servers : (Default : 1.1.1.1 , 1.0.0.1) ==> ")
	dnsStr, _ := reader.ReadString('\n')
	dnsStr = strings.TrimSpace(dnsStr)
	fmt.Println("--------")

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

	fmt.Print(color.Yellow + "How many Networks do you want for your VMs ? " + color.Reset)
	numNetworkStr, _ := reader.ReadString('\n')
	netCount := atoi(numNetworkStr)

	for j := 0; j < netCount; j++ {
		n := readRequired(reader, "Network name: ")
		ip := readRequired(reader, "Network IP: ")
		netmask := readRequired(reader, "Network Netmask: ")

		vm.Networks = append(vm.Networks, Network{
			Name:    strings.TrimSpace(n),
			IP:      strings.TrimSpace(ip),
			Netmask: atoi(netmask),
		})
	}
	return vm
}
// for simple string to int conversion
func atoi(s string) int {
	i, _ := strconv.Atoi(strings.TrimSpace(s))
	return i
}
