package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
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

	fmt.Printf("\nOptions : \n\t\n\t1.Enter Configuration =>\t%sEnter Ansible Configuration%s \n\t------------\t\n\t2.Launch Ansible =>\t\t%sExecuting Ansbile on the VMs%s \n\t------------\n\t5.Exit" , color.Yellow , color.Reset , color.Yellow , color.Reset)
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
	fmt.Println(color.Yellow + "Executing Ansible ...." + color.Reset)
	fmt.Print(color.Yellow + "\nStarting ...\n" + color.Reset)
	time.Sleep(1 * time.Second)

	// getting current directory
	currenDir := CurrentDir()

	fmt.Println(color.Yellow + "Starting installing LTE core ..." + color.Reset)
	time.Sleep(1 * time.Second)
	cmd := exec.Command("ansible-playbook" , "-i" , "inventory/hosts" , "playbooks/install-core-new.yml")
	cmd.Dir = currenDir + wdir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run() ; err != nil {
		log.Fatal(err)
	}
	fmt.Println(color.Green + "Installing core is Successfully Completed." + color.Reset)

	fmt.Println(color.Yellow + "\nStarting Configuring LTE core ..." + color.Reset)
	time.Sleep(1 * time.Second)
	cmd2 := exec.Command("ansible-playbook" , "-i" , "inventory/hosts" , "playbooks/config_core.yaml")
	cmd2.Dir = currenDir + wdir
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr
	if err2 := cmd2.Run() ; err2 != nil {
		log.Fatal(err2)
	}
	fmt.Println(color.Green + "Configuring core is Successfully Completed." + color.Reset)
	time.Sleep(1 * time.Second)

	fmt.Println(color.Cyan + "Returning to main menu ..." + color.Cyan)
	Cyborg(hostname , wdir)

}

// =============================================================================================== Helper functions ==============================================================================================
func CurrentDir() string {
	currentDir , _ := os.Getwd()
	return currentDir
}
// =============================================================================================== Helper functions (END) ==============================================================================================
