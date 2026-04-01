package main

import (
	"bufio"
	"fmt"
	"github.com/TwiN/go-color"
	"os"
	"strings"
	"time"
)

func Oranos_configure(wdir string) {
	var builder strings.Builder

	//#####################################
	//TEST
	// filename := "terraform.tfvars"
	//#####################################

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

func vlanList(vlans map[string]string, text string, filename string) map[string]string {
	// var builder strings.Builder
	// filename := "terraform.tfvars"
	reader := bufio.NewReader(os.Stdin)
	if len(vlans) == 0 {
		fmt.Println(color.Yellow + "\nNo VLANs configured yet." + color.Reset)
	}
	fmt.Println(color.Purple + "\nList of Existing VLANs : " + color.Reset)
	for name, id := range vlans {
		fmt.Printf("- VLAN name :"+color.Yellow+" %s "+color.Reset+"==>"+" VLAN ID : "+color.Cyan+"%s\n"+color.Reset, name, id)
	}

	fmt.Println("\n1) Add VLAN")
	fmt.Println("2) Remove VLAN")
	fmt.Println("3) Main menu")
	fmt.Print("\nChoose option: ")

	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		fmt.Printf("\nVLAN name : ")
		name, _ := reader.ReadString('\n')
		name = strings.TrimSpace(name)

		fmt.Printf("VLAN ID : ")
		idStr, _ := reader.ReadString('\n')
		idStr = strings.TrimSpace(idStr)

		vlans[name] = idStr
		fmt.Println(color.Green + "VLANs added Successfully" + color.Reset)
		time.Sleep(1 * time.Second)
		refactorVlans(text, vlans, filename)
	case "2":
		fmt.Print("Enter VLAN ID to remove: ")
		removeID, _ := reader.ReadString('\n')
		removeID = strings.TrimSpace(removeID)

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

	case "3":
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

	time.Sleep(2 * time.Second)
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
