package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/TwiN/go-color"
	"github.com/common-nighthawk/go-figure"
)


func Webui(hostname string) {
	figure.NewColorFigure("Webui" , "" , "purple" , true).Print()

	fmt.Println(color.Cyan + "opening WebUi ..." + color.Reset)
	time.Sleep(2 * time.Second)

	currentPath , _ := CurrentDir()

	execPy := exec.Command("./run_app")
	execPy.Dir = currentPath + "/webui/"
	execPy.Stderr = os.Stderr
	execPy.Stdout = os.Stdout

	if err := execPy.Run() ; err != nil {
		panic(err)	
	}
}