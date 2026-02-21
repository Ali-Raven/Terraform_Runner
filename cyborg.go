package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"github.com/TwiN/go-color"
	"github.com/common-nighthawk/go-figure"
)

func Cyborg(hostname, wdir string) {
	figure.NewColorFigure("CYBORG", "", "cyan", true).Print()
	reader := bufio.NewReader(os.Stdin)
	
	fmt.Println(color.Blue + "\nUsing Cyborg_starter ..." + color.Reset)
	fmt.Println(color.Yellow + "Ansible mode ..." + color.Reset)

	time.Sleep(1 * time.Second)

	fmt.Printf("\nOptions : \n\t\n\t1.Enter Configuration =>    %sEnter Ansible Configuration%s \n\t------------\t\n\t2.Launch Ansible =>         %sExecuting Ansbile on the VMs%s \n\t------------\n\t5.Exit" , color.Yellow , color.Reset , color.Yellow , color.Reset)
	fmt.Print()
	usrChoice := readRequired(reader , "\n\nchoose : ")

	switch usrChoice {
	case "1":
		Ansible_config(hostname , wdir , reader)
		return
	case "2":
		Running_Ansible(hostname , wdir , reader)
		return
	case "3":
		fmt.Println(color.Yellow + "Exiting ...." + color.Reset)
		time.Sleep(1 * time.Second)
		os.Exit(0)
	default:
		fmt.Println(color.Yellow + "Choose from the Options ..." + color.Reset)
		fmt.Println(color.Yellow + "Returning ..." + color.Reset)
		time.Sleep(1 * time.Second)
		Cyborg(hostname , wdir)
	}

}


func Ansible_config(hostname , wdir string , reader *bufio.Reader) {

}

func Running_Ansible(hostname , wdir string , reader *bufio.Reader) {

}

// =============================================================================================== Helper functions ==============================================================================================
func CurrentDir() string {
	currentDir , _ := os.Getwd()
	return currentDir
}
// =============================================================================================== Helper functions (END) ==============================================================================================
