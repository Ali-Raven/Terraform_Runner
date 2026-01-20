package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func Nozaros(hostname string) {
	var userinput int8
	fmt.Println("\nUsing Nozaros_starter")
	fmt.Println("multi_VM creating mode ...")
	time.Sleep(1 * time.Second)
	fmt.Println("\nOptions : \n\t1.Plan =>    Show changes required by the current configuration \n\t2.apply =>   Create or update infrastructure \n\t3.destroy => Destroy previously-created infrastructure \n\t4.Exit")
	fmt.Print("choice : ")
	fmt.Scan(&userinput)
	mode := "plan"
	switch userinput {
	case 1:
		terraform_plan_nozaros(&mode)
		return
	case 2:
		terraform_apply_nozaros(&mode)
		return
	case 3:
		terraform_destroy_nozaros(&mode)
		return
	case 4:
		os.Exit(0)
	}
}

func terraform_plan_nozaros(mode *string) {
	fmt.Printf("command ==> terraform %v ==> executing ...\n\n", *mode)
	time.Sleep(2 * time.Second)
	cmd := exec.Command("terraform", "plan")

	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cmd.Dir = currentDir + "/final_terraform"

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error appeared during executing [terraform %v]\n", *mode)
	} else {
		fmt.Println("\nSuccessfully Executed.")
	}
	main()
}
func terraform_apply_nozaros(mode *string) {
	*mode = "apply"
	baseCommand_nozaros("terraform", "apply", "--auto-approve", mode)
	main()
}
func terraform_destroy_nozaros(mode *string) {
	*mode = "destroy"
	baseCommand_nozaros("terraform", "destroy", "--auto-approve", mode)
	main()
}

func baseCommand_nozaros(com1, com2, com3 string, mode *string) {
	fmt.Printf("command ==> terraform %v ==> executing ...\n\n", *mode)
	time.Sleep(2 * time.Second)
	cmd := exec.Command(com1, com2, com3)

	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cmd.Dir = currentDir + "/final_terraform"

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error appeared during executing [terraform %v --auto-approve]\n", *mode)
	} else {
		fmt.Println("\nSuccessfully Executed.")
	}
}
