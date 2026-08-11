package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/TwiN/go-color"
	"github.com/terraform_runner/helper"
)

func Oranos_configure(wdir, hostname string) {
	var builder strings.Builder

	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	filename := currentDir + wdir + "/terraform.tfvars"
	content, _ := os.ReadFile(filename)
	text := string(content)

	vlans := getVlans(text)
	checked_vlans := vlanList(vlans, text, filename)

	text = removeVlansBlock(text)

	builder.WriteString(strings.TrimSpace(text))
	builder.WriteString("\n\nvlans = {\n")

	for name, id := range checked_vlans {
		builder.WriteString(fmt.Sprintf("  %s = %s\n", name, id))
	}

	builder.WriteString("}\n")

	os.WriteFile(filename, []byte(builder.String()), 0644)

	time.Sleep(2 * time.Second)
	fmt.Printf("%s%s Updated Successfully %s\n", color.Green, filename, color.Reset)
	MainStage(wdir, hostname, 1)
	// main()
}

func printSortedByID(vlans map[string]string) {
	keys := make([]string, 0, len(vlans))
	for k := range vlans {
		keys = append(keys, k)
	}

	// Sort keys based on numeric value of the string ID
	sort.Slice(keys, func(i, j int) bool {
		idI, _ := strconv.Atoi(vlans[keys[i]])
		idJ, _ := strconv.Atoi(vlans[keys[j]])
		return idI < idJ
	})
	// Print
	for _, name := range keys {
		id := vlans[name]
		fmt.Printf("- VLAN name :"+color.Yellow+" %s "+color.Reset+"==>"+" VLAN ID : "+color.Cyan+"%s\n"+color.Reset, name, id)
	}
	fmt.Printf("\n%sCount of Vlans : %d%s\n", color.Yellow, len(keys), color.Reset)
}

func ModifyVlans(vlans map[string]string, text string, filename string) {
	var vlanID string

	fmt.Println(color.Purple + "\nList of Existing VLANs : " + color.Reset)
	printSortedByID(vlans)

	vlanID = helper.Ask("Which VLAN ID do you want to Modify?", vlanID)

	var foundName string
	var found bool

	// 1. Search for the VLAN by ID
	for name, id := range vlans {
		if id == vlanID {
			foundName = name
			found = true
			break
		}
	}

	// 2. Handle the case where the ID doesn't exist
	if !found {
		fmt.Println(color.Red + "Error: VLAN ID not found." + color.Reset)
		time.Sleep(1 * time.Second)
		vlanList(vlans, text, filename)
		return
	}

	// 3. Prompt for new details
	fmt.Printf(color.Yellow+"\nModifying VLAN: %s (Current ID: %s)\n"+color.Reset, foundName, vlanID)

	var newName, newID string
	newName = helper.Ask("Enter new VLAN Name (or press Enter to keep current): ", newName)
	newID = helper.Ask("Enter new VLAN ID (or press Enter to keep current): ", newID)

	// 4. Keep existing values if the user leaves the input blank
	if strings.TrimSpace(newName) == "" {
		newName = foundName
	}
	if strings.TrimSpace(newID) == "" {
		newID = vlanID
	}

	// 5. Update the map
	// If the name (map key) changed, we must delete the old entry
	if newName != foundName {
		delete(vlans, foundName)
	}
	vlans[newName] = newID

	fmt.Println(color.Green + "VLAN Modified Successfully" + color.Reset)
	time.Sleep(1 * time.Second)

	// 6. Save changes using your existing refactor function
	refactorVlans(text, vlans, filename)
}
func vlanList(vlans map[string]string, text string, filename string) map[string]string {
	// var builder strings.Builder
	// filename := "terraform.tfvars"
	// reader := bufio.NewReader(os.Stdin)
	if len(vlans) == 0 {
		fmt.Println(color.Yellow + "\nNo VLANs configured yet." + color.Reset)
	}

	// fmt.Println(vlans)
	fmt.Println(color.Purple + "\nList of Existing VLANs : " + color.Reset)
	printSortedByID(vlans)
	// for name, id := range vlans {
	// 	fmt.Printf("- VLAN name :"+color.Yellow+" %s "+color.Reset+"==>"+" VLAN ID : "+color.Cyan+"%s\n"+color.Reset, name, id)
	// }

	fmt.Println()

	choice := helper.AskSelect([]string{"Add VLAN", "Modify VLANs", "Remove VLAN", "Removing Vlans on vCenter", "Main Menu"})

	switch choice {
	case "Add VLAN":
		var vlanName, IdString string
		// fmt.Printf("\nVLAN name : ")
		name := helper.Ask("VLAN name :", vlanName)
		idStr := helper.Ask("VLAN ID : ", IdString)

		vlans[name] = idStr
		fmt.Println(color.Green + "VLANs added Successfully" + color.Reset)
		time.Sleep(1 * time.Second)
		refactorVlans(text, vlans, filename)
	case "Remove VLAN":
		var rmID string
		removeID := helper.Ask("Enter VLAN ID to remove: ", rmID)

		found := false

		for name, id := range vlans {
			if removeID == id {
				delete(vlans, name)
				found = true
				break
			}
		}

		if found {
			fmt.Println(color.Green + "VLAN removed Successfully" + color.Reset)
			time.Sleep(1 * time.Second)
			refactorVlans(text, vlans, filename)
		} else {
			fmt.Println(color.Red + "Error : VLAN ID not found." + color.Reset)
			time.Sleep(1 * time.Second)
			vlanList(vlans, text, filename)
		}
	case "Modify VLANs":

		ModifyVlans(vlans, text, filename)

	case "Removing Vlans on vCenter":
		// RemoveVlansOnVCenter()

	case "Main Menu":
		fmt.Println(color.Yellow + "loading main menu ..." + color.Reset)
		time.Sleep(1 * time.Second)
		main()
	default:
		fmt.Println(color.Yellow + "Warning : Invalid choice, returning to options." + color.Reset)
		time.Sleep(1 * time.Second)
		vlanList(vlans, text, filename)
	}
	return vlans
}

// func RemoveVlansOnVCenter() {

// 	vCenterURL := "https://administrator%40vsphere.local:Aa@123321@192.168.0.251/sdk"
// 	HostIP := "10.56.20.11"

// 	portGroupsToDelete := []string{
// 		"s11-sgwc1",
// 		"s6a-mme9",
// 		"s1ap-mme10",
// 		"s5c-sgwc2",
// 		"sxb-smf2",
// 		"sxu-smf1",
// 		"s6a-mme2",
// 		"sxa-sgwu2",
// 		"s11-mme8",
// 		"s6a-mme12",
// 		"s1ap-mme1",
// 		"s5c-mme8",
// 		"s6a-mme11",
// 		"s6a-mme3",
// 		"s6a-hss1",
// 		"sxa-sgwu1",
// 		"s11-mme4",
// 		"s5c-mme11",
// 		"s11-mme12",
// 		"s1ap-mme4",
// 		"s6a-hss2",
// 		"s6a-mme13",
// 		"s5u-sgwu2",
// 		"s11-sgwc2",
// 		"sxu-smf2",
// 		"s5c-mme6",
// 		"s1ap-mme13",
// 		"sxb-upf1",
// 		"s5c-sgwc1",
// 		"s1ap-mme6",
// 		"sxb-upf2",
// 		"s11-mme6",
// 		"gx-smf1",
// 		"s6a-hss3",
// 		"s1ap-mme8",
// 		"s5c-smf2",
// 		"gx-pcrf1",
// 		"s11-mme10",
// 		"s1u-sgwu1",
// 		"s5c-mme2",
// 		"s5c-mme9",
// 		"s5c-mme10",
// 		"s11-mme7",
// 		"s5c-mme3",
// 		"s1ap-mme11",
// 		"s6a-mme5",
// 		"gx-smf2",
// 		"s5c-mme7",
// 		"s5c-mme12",
// 		"s11-mme1",
// 		"sxa-sgwc2",
// 		"s11-mme5",
// 		"s5c-mme13",
// 		"s6a-mme6",
// 		"sxa-sgwc1",
// 		"sxu-upf2",
// 		"sxb-smf1",
// 		"s5c-mme5",
// 		"s11-mme2",
// 		"s6a-mme7",
// 		"s6a-mme4",
// 		"s1ap-mme12",
// 		"s5u-upf2",
// 		"s5c-mme4",
// 		"sxu-upf1",
// 		"s1ap-mme7",
// 		"s11-mme13",
// 		"s1u-sgwu2",
// 		"s6a-mme1",
// 		"s1ap-mme2",
// 		"s1ap-mme5",
// 		"s5c-smf1",
// 		"s11-mme3",
// 		"s5u-upf1",
// 		"s1ap-mme9",
// 		"s5u-sgwu1",
// 		"s5c-mme1",
// 		"s1ap-mme3",
// 		"s11-mme11",
// 		"s6a-mme10",
// 		"s11-mme9",
// 		"s6a-mme8",
// 		"gx-pcrf",
// 		"gx-smf1",
// 	}

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	u, err := soap.ParseURL(vCenterURL)
// 	if err != nil {
// 		log.Fatalf("Error parsing URL: %v", err)
// 	}

// 	client, err := govmomi.NewClient(ctx, u, true)
// 	if err != nil {
// 		log.Fatalf("Error connecting to vCenter: %v", err)
// 	}

// 	defer client.Logout(ctx)
// 	fmt.Println(color.Green + "Successfully connected to vCenter." + color.Reset)
// 	finder := find.NewFinder(client.Client, true)

// 	dc, _ := finder.DatacenterOrDefault(ctx, "Datacenter")
// 	finder.SetDatacenter(dc)
// 	host, err := finder.HostSystem(ctx, HostIP)
// 	if err != nil {
// 		log.Fatalf("Error finding host %s: %v", HostIP, err)
// 	}
// 	fmt.Printf("Found host: %s\n", HostIP)

// 	// 3. Get the Host Network System
// 	netSystem, err := host.ConfigManager().NetworkSystem(ctx)
// 	if err != nil {
// 		log.Fatalf("Error getting network system for host: %v", err)
// 	}

// 	// 4. Iterate through the port groups and delete them
// 	for _, pg := range portGroupsToDelete {
// 		fmt.Printf("Attempting to delete port group '%s'... ", pg)

//			err = netSystem.RemovePortGroup(ctx, pg)
//			if err != nil {
//				// It might fail if the port group doesn't exist or is still in use by a VM
//				fmt.Printf("FAILED: %v\n", err)
//			} else {
//				fmt.Println("SUCCESS")
//			}
//		}
//	}
func refactorVlans(text string, vlans map[string]string, filename string) {
	var builder strings.Builder
	text = removeVlansBlock(text)

	builder.WriteString(strings.TrimSpace(text))
	builder.WriteString("\n\nvlans = {\n")

	for name, id := range vlans {
		builder.WriteString(fmt.Sprintf("  %s = %s\n", name, id))
	}

	builder.WriteString("}\n")

	os.WriteFile(filename, []byte(builder.String()), 0644)

	time.Sleep(1 * time.Second)
	fmt.Printf("%s%s Updated Successfully %s\n", color.Green, filename, color.Reset)
	vlanList(vlans, text, filename)
}

func getVlans(text string) map[string]string {
	vlans := make(map[string]string)

	start := strings.Index(text, "vlans = {")
	if start == -1 {
		return vlans
	}

	end := strings.Index(text[start:], "}")
	if end == -1 {
		return vlans
	}

	block := text[start : start+end]
	lines := strings.Split(block, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "vlans = {" || line == "}" || line == "" {
			continue
		}
		if strings.Contains(line, "=") {
			parts := strings.Split(line, "=")
			name := strings.TrimSpace(parts[0])
			id := strings.TrimSpace(parts[1])

			if name != "" && id != "" {
				vlans[name] = id
			}
		}
	}

	return vlans
}

func removeVlansBlock(text string) string {
	start := strings.Index(text, "vlans = {")
	if start == -1 {
		return text
	}

	end := strings.Index(text[start:], "}")
	if end == -1 {
		return text
	}

	end = start + end + 1
	return text[:start] + text[end:]
}
