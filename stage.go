package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
	"github.com/TwiN/go-color"
)

func Oranos(hostname, wdir string) {
	fmt.Println(color.Blue + "\nUsing Oranos_starter" + color.Reset)
	fmt.Println(color.Yellow + "VLAN creating mode ..." + color.Reset)
	MainStage(wdir , 1)
}

func Nozaros(hostname, wdir string) {
	fmt.Println(color.Blue + "\nUsing Nozaros_starter" + color.Reset)
	fmt.Println(color.Yellow + "multi_VM creating mode ..." + color.Reset)
	MainStage(wdir , 2)
}

func MainStage(wdir string , componentID int8) {
	var userinput int8
	time.Sleep(1 * time.Second)
	fmt.Println("\nOptions : \n\t1.Enter Configuration =>    user configuration for VMs \n\t------------\t\n\t2.Plan =>    Show changes required by the current configuration \n\t------------\t\n\t3.apply =>   Create or update infrastructure \n\t------------\t\n\t4.destroy => Destroy previously-created infrastructure \n\t------------\t\n\t5.Exit")
	fmt.Print("\nchoice: ")
	fmt.Scan(&userinput)
	mode := "plan"
	switch userinput {
	case 1:
		switch componentID  {
			case 1:
				Oranos_configure(wdir)
			case 2:
				Nozaros_configure(wdir)
		}
		return
	case 2:
		terraform_plan(&mode, wdir)
		return
	case 3:
		terraform_apply(&mode, wdir)
		return
	case 4:
		terraform_destroy(&mode, wdir)
		return
	case 5:
		os.Exit(0)
	}
}

func terraform_plan(mode *string, wdir string) {
	fmt.Printf(color.Yellow + "command ==> terraform %v ==> executing ...\n\n" + color.Reset, *mode)
	time.Sleep(2 * time.Second)
	cmd := exec.Command("terraform", "plan")
	// execute dst directory and getting current dir
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cmd.Dir = currentDir + wdir

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf(color.Red + "Error appeared during executing [terraform %v]\n" + color.Reset, *mode)
	} else {
		fmt.Println(color.Green + "\nSuccessfully Executed." + color.Reset)
	}
	main()
}
func terraform_apply(mode *string, wdir string) {
	*mode = "apply"
	baseCommand("terraform", "apply", "--auto-approve" , wdir , mode)
	main()
}
func terraform_destroy(mode *string, wdir string) {
	*mode = "destroy"
	baseCommand("terraform", "destroy", "--auto-approve" , wdir, mode)
	main()
}

func baseCommand(com1, com2, com3 , wdir string, mode *string) {
	fmt.Printf(color.Yellow + "command ==> terraform %v ==> executing ...\n\n" + color.Reset, *mode)
	time.Sleep(2 * time.Second)
	cmd := exec.Command(com1, com2, com3)

	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cmd.Dir = currentDir + wdir

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf(color.Red + "Error appeared during executing [terraform %v --auto-approve]\n" + color.Reset, *mode)
	} else {
		fmt.Println(color.Green + "\nSuccessfully Executed." + color.Reset)
	}
}
