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

func Oranos_configure(wdir string) {
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
	main()
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

	choice := helper.AskSelect([]string{"Add VLAN", "Remove VLAN", "Main Menu"})

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
		// fmt.Print("Enter VLAN ID to remove: ")
		// removeID, _ := reader.ReadString('\n')
		// removeID = strings.TrimSpace(removeID)

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
