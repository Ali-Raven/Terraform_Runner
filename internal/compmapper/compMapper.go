// Package compmapper is for iterating over this Map that create in generators and spliting the network properties of that
package compmapper

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/TwiN/go-color"
	generators "github.com/terraform_runner/internal/generators"
	t "github.com/terraform_runner/internal/template"
	typesstructs "github.com/terraform_runner/internal/typesStructs"
	"github.com/terraform_runner/internal/vmstore"
)

type DataComp map[string]generators.ComponentData

type ComponentMapper struct {
	CompsMapsMME DataComp
	CompsMapsHSS DataComp
	CompsMapSGWC DataComp
	CompsMapSGWU DataComp
	CompsMapSMF  DataComp
	CompsMapUPF  DataComp
	CompsMapPCRF DataComp
}

// NetworkGeneratorComps function generates the network mapping for a given component (e.g., MME, HSS, etc.) based on the provided VM list.
func NetworkGeneratorComps(CompsName string, VMList map[string]typesstructs.VM) DataComp {
	CompsNetworksMaps := make(map[string]map[string]generators.NetworksStructure)
	result := make(DataComp)

	for vmName, vm := range VMList {

		// only take VMs that belong to this component (MME, HSS, etc.)
		if !strings.Contains(vmName, CompsName) {
			continue
		}

		nets := make(map[string]generators.NetworksStructure)

		for _, netw := range vm.Networks {
			nets[netw.Name] = generators.NetworksStructure{
				ID:      netw.ID,
				Name:    netw.Name,
				IP:      netw.IP,
				Gateway: netw.Gateway,
				Subnet:  netw.Netmask,
			}
		}

		CompsNetworksMaps[vmName] = nets

		result[vmName] = generators.ComponentData{
			Name:                vmName,
			Networks:            nets,
			ComponentsToConnect: vm.ComponentsToConnect,
		}
	}

	return result
}

func YmlCompMapper(wdir string, vms []typesstructs.VM) {
	// project_path := currentDir + "/ansible-core-deploy"
	// fmt.Println(vms[0].Networks[0].IP)

	fmt.Println(color.Yellow + "loading Existing VMs info ..." + color.Reset)
	vms, err := vmstore.LoadExistingVMs(wdir)
	if err != nil {
		fmt.Println(color.Red + "")
		panic(err)
	}

	fmt.Println(color.Green + "VMs loaded Successfully." + color.Reset)
	for i := range vms {
		if len(vms[i].Networks) < 2 {
			fmt.Printf("%sError : not enough networks interface for %s%s\n", color.Red, vms[i].Name, color.Reset)
			time.Sleep(300 * time.Millisecond)
			FlagErr = true
		}
	}
	if FlagErr {
		return
	}

	// mapping key for getting VMs list with Name of the vms not by ID
	VMList := make(map[string]typesstructs.VM)
	for _, name := range vms {
		VMList[name.Name] = name
	}

	MapperComp := ComponentMapper{
		CompsMapsMME: NetworkGeneratorComps("MME", VMList),
		CompsMapsHSS: NetworkGeneratorComps("HSS", VMList),
		CompsMapSGWC: NetworkGeneratorComps("SGWC", VMList),
		CompsMapSGWU: NetworkGeneratorComps("SGWU", VMList),
		CompsMapSMF:  NetworkGeneratorComps("SMF", VMList),
		CompsMapUPF:  NetworkGeneratorComps("UPF", VMList),
		CompsMapPCRF: NetworkGeneratorComps("PCRF", VMList),
	}

	mmesByName := generators.BuildAllMMEs(MapperComp.CompsMapsMME)
	hsssByName := generators.BuildAllHSSs(MapperComp.CompsMapsHSS)
	sgwcsByName := generators.BuildAllSGWCs(MapperComp.CompsMapSGWC)
	sgwusByName := generators.BuildAllSGWUs(MapperComp.CompsMapSGWU)
	smfsByName := generators.BuildAllSMFs(MapperComp.CompsMapSMF)
	upfsByName := generators.BuildAllUPFs(MapperComp.CompsMapUPF)
	pcrfsByName := generators.BuildAllPCRFs(MapperComp.CompsMapPCRF)

	givenDataTemplate := Data(mmesByName, hsssByName, sgwcsByName, sgwusByName, smfsByName, upfsByName, pcrfsByName)

	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"sub_teid": func(a int) int {
			return (((a+1)-1)*30000 + 1)
		},
		"add_teid": func(a int) int {
			return ((a + 1) * 30000)
		},
	}
	templateTest := template.Must(template.New("yaml").Funcs(funcMap).Parse(t.TemplateYML()))

	var buf bytes.Buffer

	if err := templateTest.Execute(&buf, givenDataTemplate); err != nil {
		panic(err)
	}

	inventoryPath, fileName := "ansible/ansible-core-deploy/inventory/", "maintest.yml"
	// fileName := "demo.yml"
	errOSwritefile := os.WriteFile(inventoryPath+fileName, buf.Bytes(), 0o644)
	if errOSwritefile != nil {
		log.Fatal(err)
	}
	// os.WriteFile(fileName, buf.Bytes(), 0644)
	fmt.Printf("\n%sGenerating %s  file ...%s", color.Yellow, fileName, color.Reset)
	time.Sleep(1 * time.Second)
	fmt.Printf("\n%s%s generated in the current path%s\n\n", color.Green, fileName, color.Reset)
}
