package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/TwiN/go-color"
	"github.com/common-nighthawk/go-figure"
)

func Helm(hostname, wdir string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(color.Blue + "Using Helm_starter" + color.Reset)
	fmt.Println(color.Yellow + "ESXI and vCenter mode ..." + color.Reset)

	time.Sleep(1 * time.Second)
	fmt.Println("\nOptions : ")
	fmt.Println("\n1. Esxi Product\n2. vCenter Product\n3. Exit")

	fmt.Print("\nchoice: (1/2/3): ")
	userinput, _ := reader.ReadString('\n')
	userinput = strings.TrimSpace(userinput)

	switch userinput {
	case "1":
		ESXI_setup(hostname, wdir , reader)
		return
	case "2":
		Vcenter_setup(hostname, wdir , reader)
		return
	case "3":
		fmt.Println(color.Yellow + "Exiting..." + color.Reset)
		os.Exit(0)
	default:
		fmt.Println(color.Yellow + "\n choose one of options ..." + color.Reset)
		fmt.Println("Options : ")
		fmt.Println("1. Esxi Product\n2. vCenter Product")

	}
	// MainStage(wdir , 4)
}

func ESXI_setup(hostname, wdir string , reader *bufio.Reader) {
	figure.NewColorFigure("ESXI Setup", "", "green" , true).Print()

	fmt.Println(color.Blue + "\nUsing ESXI_setup" + color.Reset)
	fmt.Println(color.Yellow + "Setting up ESXI Product ..." + color.Reset)
	time.Sleep(1 * time.Second)
	fmt.Println(color.Yellow + "Fetching list of Servers for installing ESXI products ..." + color.Reset)
	servers := fetchServers()
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

}
func Vcenter_setup(hostname, wdir string , reader *bufio.Reader) {
	figure.NewColorFigure("vCenter Setup", "", "green" , true).Print()
	fmt.Println(color.Blue + "\nUsing vCenter_setup" + color.Reset)
	fmt.Println(color.Yellow + "Setting up vCenter Product ..." + color.Reset)


}
