package main


import (
	"fmt"
	"os/exec"
	"time"
	"os"
)

func Moon() {
	var userinput int8
	fmt.Println("Options : \n\t1.Plan =>    showing all things and configs that you want to create \n\t2.apply =>   applying the changes and configs that you create to make VMs \n\t3.destroy => destroy all chagnes that you made like VMs or .. \n\t4.Exit =>    Exit from application")
	fmt.Print("choise : ")
	fmt.Scan(&userinput)

	switch userinput {
	case 1:
		terraform_plan_moon()
		return
	case 2:
		terraform_apply_moon()
		return
	case 3:
		terraform_destroy_moon()
		return
	case 4:
		os.Exit(0)
	}
}

func terraform_plan_moon() {
	fmt.Print("command => terraform plan => executing ...\n\n")
	time.Sleep(2 * time.Second)
	cmd := exec.Command("terraform", "plan")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Error appeared during executing [terraform plan]")
	} else {
		fmt.Println("\nSuccessfully Executed.")
	}
	main()
}
func terraform_apply_moon() {
	fmt.Print("command => terraform apply => executing ...\n\n")
	time.Sleep(2 * time.Second)
	cmd := exec.Command("terraform", "apply" , "--auto-approve")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Error appeared during executing [terraform apply --auto-approve]")
	} else {
		fmt.Println("\nSuccessfully Executed.")
	}
	main()
}
func terraform_destroy_moon() {
	fmt.Print("command => terraform destroy => executing ...\n\n")
	time.Sleep(2 * time.Second)
	cmd := exec.Command("terraform", "destroy" , "--auto-approve")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Error appeared during executing [terraform destroy --auto-approve]")
	} else {
		fmt.Println("\nSuccessfully Executed.")
	}
	main()
}