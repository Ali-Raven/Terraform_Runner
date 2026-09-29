package vmstore

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/TwiN/go-color"
	typesstructs "github.com/terraform_runner/internal/typesStructs"
)

func LoadExistingVMs(wdir string) ([]typesstructs.VM, error) {
	tfvars, err := LoadTFvars(wdir)
	if err != nil {
		fmt.Println(color.Yellow + "No existing VMs found. Starting fresh..." + color.Reset)
		return []typesstructs.VM{}, err
	}
	return tfvars.VMs, nil
}

func LoadTFvars(wdir string) (typesstructs.TFvars, error) {
	currentDir, _ := os.Getwd()
	filename := currentDir + wdir + "/terraform.tfvars.json"
	data, err := os.ReadFile(filename)
	if err != nil {
		return typesstructs.TFvars{}, err
	}

	var tfvars typesstructs.TFvars
	err = json.Unmarshal(data, &tfvars)
	if err != nil {
		fmt.Println(color.Red+"Error unmarshaling JSON:"+color.Reset, err)
		return tfvars, err
	}
	return tfvars, err
}
