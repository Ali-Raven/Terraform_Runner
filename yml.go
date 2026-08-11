package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/TwiN/go-color"
	generators "github.com/terraform_runner/Generators"
	t "github.com/terraform_runner/Template"
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

// this function generates the network mapping for a given component (e.g., MME, HSS, etc.) based on the provided VM list.
func NetworkGeneratorComps(CompsName string, VmList map[string]VM) map[string]generators.ComponentData {
	CompsNetworksMaps := make(map[string]map[string]generators.NetworksStructure)
	result := make(map[string]generators.ComponentData)

	for vmName, vm := range VmList {

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

func Yml(wdir string, vms []VM) {
	// project_path := currentDir + "/ansible-core-deploy"
	// fmt.Println(vms[0].Networks[0].IP)

	fmt.Println(color.Yellow + "loading Existing VMs info ..." + color.Reset)
	vms, err := loadExistingVMs(wdir)
	if err != nil {
		fmt.Println(color.Red + "")
		panic(err)
	}

	fmt.Println(color.Green + "VMs loaded Successfully." + color.Reset)
	for i := 0; i < len(vms); i++ {

		if len(vms[i].Networks) < 2 {
			fmt.Printf("%sError : not enough networks interface for %s%s\n", color.Red, vms[i].Name, color.Reset)
			time.Sleep(300 * time.Millisecond)
			FlagErr = true
		}
	}
	if FlagErr == true {
		return
	}

	// mapping key for getting VMs list with Name of the vms not by ID
	VmList := make(map[string]VM)
	for _, name := range vms {
		VmList[name.Name] = name

	}

	MapperComp := ComponentMapper{
		CompsMapsMME: NetworkGeneratorComps("MME", VmList),
		CompsMapsHSS: NetworkGeneratorComps("HSS", VmList),
		CompsMapSGWC: NetworkGeneratorComps("SGWC", VmList),
		CompsMapSGWU: NetworkGeneratorComps("SGWU", VmList),
		CompsMapSMF:  NetworkGeneratorComps("SMF", VmList),
		CompsMapUPF:  NetworkGeneratorComps("UPF", VmList),
		CompsMapPCRF: NetworkGeneratorComps("PCRF", VmList),
	}

	// CompsMapsMME := NetworkGeneratorComps("MME", VmList)
	// CompsMapsHSS := NetworkGeneratorComps("HSS", VmList)
	// CompsMapSGWC := NetworkGeneratorComps("SGWC", VmList)
	// CompsMapSGWU := NetworkGeneratorComps("SGWU", VmList)
	// CompsMapSMF := NetworkGeneratorComps("SMF", VmList)
	// CompsMapUPF := NetworkGeneratorComps("UPF", VmList)
	// CompsMapPCRF := NetworkGeneratorComps("PCRF", VmList)

	mmesByName := generators.BuildAllMMEs(MapperComp.CompsMapsMME)
	hsssByName := generators.BuildAllHSSs(MapperComp.CompsMapsHSS)
	sgwcsByName := generators.BuildAllSGWCs(MapperComp.CompsMapSGWC)
	sgwusByName := generators.BuildAllSGWUs(MapperComp.CompsMapSGWU)
	smfsByName := generators.BuildAllSMFs(MapperComp.CompsMapSMF)
	upfsByName := generators.BuildAllUPFs(MapperComp.CompsMapUPF)
	pcrfsByName := generators.BuildAllPCRFs(MapperComp.CompsMapPCRF)

	// checking version of the core
	Core_Name = "{{ core_name }}"
	if V == "" {
		Core_name_v = Core_Name
	} else {
		Core_name_v = Core_Name + "-" + V
	}
	Var_path = "/var/log/{{ core_name }}"
	Tls_path = "/etc/{{ core_name }}/tls/"
	Config_path = "/etc/{{ core_name }}"
	Bin_path = "/usr/bin"
	Core_source_path = "/opt/{{ core_name_version }}"
	Inventory_hostname = "{{ inventory_hostname }}"
	Diam_groupNames = "{{ group_names[1] }}"
	Non_diam_groupNames = "{{ group_names[0] }}"
	Var_path_diameter = "/etc/{{ core_name }}/freeDiameter/"
	Diam_Realm = "epc.mnc{{ plmn.mnc }}.mcc{{ plmn.mcc }}.3gppnetwork.org"
	Var_path_Comps = "/var/log/{{ core_name }}/{{ group_names[0] }}.log"
	Hardcoded_Diam_Realm = "{{ diam_realm }}"

	givenDataTemplate := Data(mmesByName, hsssByName, sgwcsByName, sgwusByName, smfsByName, upfsByName, pcrfsByName)

	// getting the template from the Template package
	RecivedYamlData := t.TemplateYML()

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
	templateTest := template.Must(template.New("yaml").Funcs(funcMap).Parse(RecivedYamlData))

	var buf bytes.Buffer

	if err := templateTest.Execute(&buf, givenDataTemplate); err != nil {
		panic(err)
	}

	inventoryPath, fileName := "ansible/ansible-core-deploy/inventory/", "main.yml"
	// fileName := "demo.yml"
	os.WriteFile(inventoryPath+fileName, buf.Bytes(), 0644)
	// os.WriteFile(fileName, buf.Bytes(), 0644)
	fmt.Printf("\n%sGenerating %s  file ...%s", color.Yellow, fileName, color.Reset)
	time.Sleep(1 * time.Second)
	fmt.Printf("\n%s%s generated in the current path%s\n\n", color.Green, fileName, color.Reset)
}
