package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func Oranos(hostname, wdir string) {
	fmt.Println("\nUsing Oranos_starter")
	fmt.Println("VLAN creating1 mode ...")
	MainStage(wdir)
}

func Moon(hostname, wdir string) {
	fmt.Println("\nUsing Moon_starter")
	fmt.Println("enabling normal mode ...")
	MainStage(wdir)
}

func Nozaros(hostname, wdir string) {
	fmt.Println("\nUsing Nozaros_starter")
	fmt.Println("multi_VM creating mode ...")
	MainStage(wdir)
}

func MainStage(wdir string) {
	var userinput int8
	time.Sleep(1 * time.Second)
	fmt.Println("\nOptions : \n\t1.Enter Configuration =>    user configuration for VMs \n\t------------\t\n\t2.Plan =>    Show changes required by the current configuration \n\t------------\t\n\t3.apply =>   Create or update infrastructure \n\t------------\t\n\t4.destroy => Destroy previously-created infrastructure \n\t------------\t\n\t5.Exit")
	fmt.Print("\nchoice: ")
	fmt.Scan(&userinput)
	mode := "plan"
	switch userinput {
	case 1:
		Moon_configure()
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
	fmt.Printf("command ==> terraform %v ==> executing ...\n\n", *mode)
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
		fmt.Printf("Error appeared during executing [terraform %v]\n", *mode)
	} else {
		fmt.Println("\nSuccessfully Executed.")
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
	fmt.Printf("command ==> terraform %v ==> executing ...\n\n", *mode)
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
		fmt.Printf("Error appeared during executing [terraform %v --auto-approve]\n", *mode)
	} else {
		fmt.Println("\nSuccessfully Executed.")
	}
}
