package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"github.com/TwiN/go-color"
)

func Oranos_configure(wdir string) {
	var builder strings.Builder
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("How many VLANs you want to adds ? ")
	countStr, _ := reader.ReadString('\n')
	countStr = strings.TrimSpace(countStr)
	count, _ := strconv.Atoi(countStr)

	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	filename := currentDir + wdir + "/terraform.tfvars"
	content, err := os.ReadFile(filename)
	text := ""
	if err == nil {
		text = string(content)
	}
	
	text = removeVlansBlock(text)
	
	//checking existing vlans block
	hasVlansBlock := strings.Contains(text, "vlans = {")
	if hasVlansBlock {
		text = strings.TrimRight(text, "\n \t")
		text = strings.TrimSuffix(text, "}")
		builder.WriteString(text + "\n")
	} else {
		builder.WriteString(strings.TrimSpace(text))
		builder.WriteString("\n\nvlans = {\n")
	}

	// builder.WriteString(strings.TrimSpace(text))
	// builder.WriteString("\n\nvlans = {\n")

	for i := 0; i < count; i++ {
		fmt.Printf("\nVLAN %d name : ", i+1)
		name, _ := reader.ReadString('\n')
		name = strings.TrimSpace(name)

		fmt.Printf("VLAN %d ID : ", i+1)
		idStr, _ := reader.ReadString('\n')
		idStr = strings.TrimSpace(idStr)

		builder.WriteString(fmt.Sprintf("  %s = %s\n", name, idStr))
	}

	builder.WriteString("}\n")

	os.WriteFile(filename, []byte(builder.String()), 0644)

	time.Sleep(2 * time.Second)
	fmt.Println(color.Green + "VLANS added Successfully" + color.Reset)
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
