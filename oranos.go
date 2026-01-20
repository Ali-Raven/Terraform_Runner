package main


import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func Oranos() {
	var userinput int8
	fmt.Println("\nUsing Oranos_starter ...")
	time.Sleep(1 * time.Second)
	fmt.Println("\nOptions : \n\t1.Plan =>    Show changes required by the current configuration \n\t2.apply =>   Create or update infrastructure \n\t3.destroy => Destroy previously-created infrastructure \n\t4.Exit")
	fmt.Print("\nyour option :")
	fmt.Scan(&userinput)

	switch userinput {
	case 1:
		terraform_plan_oranos()
		return
	case 2:
		terraform_apply_oranos()
		return
	case 3:
		terraform_destroy_oranos()
		return
	case 4:
		os.Exit(0)
	}
}


func terraform_plan_oranos() {
	fmt.Print("command ==> terraform plan ==> executing ...\n\n")
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
func terraform_apply_oranos() {
	fmt.Print("command ==> terraform apply ==> executing ...\n\n")
	time.Sleep(2 * time.Second)
	cmd := exec.Command("terraform", "apply" , "--auto-aprove")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Error appeared during executing [terraform apply]")
	} else {
		fmt.Println("\nSuccessfully Executed.")
	}
	main()
}
func terraform_destroy_oranos() {
	fmt.Print("command ==> terraform destroy ==> executing ...\n\n")
	time.Sleep(2 * time.Second)
	cmd := exec.Command("terraform", "destroy" , "--auto-approve")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Error appeared during executing [terraform destroy]")
	} else {
		fmt.Println("\nSuccessfully Executed.")
	}
	main()
}