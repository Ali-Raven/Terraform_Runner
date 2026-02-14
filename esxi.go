package main

import (
	"fmt"
	"github.com/TwiN/go-color"
)

func ESXI_setup(hostname, wdir string) {
	fmt.Println(color.Blue + "\nUsing ESXI_setup" + color.Reset)
	fmt.Println(color.Yellow + "Setting up ESXI Product ..." + color.Reset)
	MainStage(wdir, 4)
}
