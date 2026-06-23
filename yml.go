package main

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
	"time"

	"github.com/TwiN/go-color"
)

var (
	ssh_defaultPort     int = 22
	sgwc_managementPort int = ssh_defaultPort
	sgwc_s11Port        int = 2123
	sgwc_sxaPort        int = 8805
	sgwc_s5cPort        int = 2124
	sgwu_managementPort int = ssh_defaultPort
	sgwu_sxaPort        int = 2152
	sgwu_s5uPort        int = 8805
	sgwu_s1uPort        int = 3333
	upf_managementPort  int = ssh_defaultPort
	upf_sxbPort         int = 8805
	upf_sxuPort         int = 8806
	upf_s5uPort         int = 2153
	upf_SGI             string
	upf_sgiPort         int = 2152
	smf_managementPort  int = ssh_defaultPort
	smf_gxPort          int = 2123
	gx_secPort          int = 5868
	smf_s5cPort         int = 8805
	smf_sxbPort         int = 2153
	smf_sxuPort         int = 8806
	mme_s11Port         int = 2123
	mme_s1apPort        int = 36412
	mme_s6aPort         int = 2221
	s6a_secPort         int = 5868
	mme_managementPort  int = ssh_defaultPort
	hss_managementPort  int = ssh_defaultPort
	hss_s6aPort         int = 2223
	pcrf_managementPort int = ssh_defaultPort
	pcrf_gxPort         int = 4434
	core_name           string
	var_path            string
	tls_path            string
	inventory_hostname  string
	diam_realm          string
	flagErr             bool
	diam_groupNames     string
	non_diam_groupNames string
)

type ComponentData struct {
	Name                string
	Networks            map[string]string
	ComponentsToConnect []string
}

func BuildAllMMEs(componentData map[string]ComponentData) map[string]ComponentData {
	// var mmes []ComponentData
	mmesByName := make(map[string]ComponentData)

	// fmt.Println(CompConnMMEs)
	// time.Sleep(10000 * time.Second)
	for name, nets := range componentData {
		mme := ComponentData{
			Name:                name,
			Networks:            nets.Networks,
			ComponentsToConnect: nets.ComponentsToConnect,
		}
		// mmes = append(mmes, mme)
		mmesByName[name] = mme
	}
	return mmesByName
}

func BuildAllHSSs(componentData map[string]ComponentData) map[string]ComponentData {
	hsssByName := make(map[string]ComponentData)

	for name, nets := range componentData {
		hss := ComponentData{
			Name:                name,
			Networks:            nets.Networks,
			ComponentsToConnect: nets.ComponentsToConnect,
		}
		// hsss = append(hsss, hss)
		hsssByName[name] = hss
	}
	return hsssByName
}
func BuildAllSGWCs(componentData map[string]ComponentData) map[string]ComponentData {
	sgwcsByName := make(map[string]ComponentData)

	for name, nets := range componentData {
		sgwc := ComponentData{
			Name:                name,
			Networks:            nets.Networks,
			ComponentsToConnect: nets.ComponentsToConnect,
		}
		// sgwcs = append(sgwcs, sgwc)
		sgwcsByName[name] = sgwc
	}
	return sgwcsByName
}
func BuildAllSGWUs(componentData map[string]ComponentData) map[string]ComponentData {
	sgwusByName := make(map[string]ComponentData)

	for name, nets := range componentData {
		sgwu := ComponentData{
			Name:                name,
			Networks:            nets.Networks,
			ComponentsToConnect: nets.ComponentsToConnect,
		}
		// sgwus = append(sgwus, sgwu)
		sgwusByName[name] = sgwu
	}
	return sgwusByName
}
func BuildAllSMFs(componentData map[string]ComponentData) map[string]ComponentData {
	smfsByName := make(map[string]ComponentData)

	for name, nets := range componentData {
		smf := ComponentData{
			Name:                name,
			Networks:            nets.Networks,
			ComponentsToConnect: nets.ComponentsToConnect,
		}
		// smfs = append(smfs, smf)
		smfsByName[name] = smf
	}
	return smfsByName
}
func BuildAllUPFs(componentData map[string]ComponentData) map[string]ComponentData {
	upfsByName := make(map[string]ComponentData)

	for name, nets := range componentData {
		upf := ComponentData{
			Name:                name,
			Networks:            nets.Networks,
			ComponentsToConnect: nets.ComponentsToConnect,
		}
		// upfs = append(upfs, upf)
		upfsByName[name] = upf
	}
	return upfsByName
}
func BuildAllPCRFs(componentData map[string]ComponentData) map[string]ComponentData {
	pcrfsByName := make(map[string]ComponentData)

	for name, nets := range componentData {
		pcrf := ComponentData{
			Name:                name,
			Networks:            nets.Networks,
			ComponentsToConnect: nets.ComponentsToConnect,
		}
		// pcrfs = append(pcrfs, pcrf)
		pcrfsByName[name] = pcrf
	}
	return pcrfsByName
}

func NetworkGeneratorComps(indexCount int, CompsName string, VmList map[string]VM) map[string]ComponentData {
	CompsNetworksMaps := make(map[string]map[string]string)
	result := make(map[string]ComponentData)
	var name string

	for i := 1; i <= indexCount; i++ {
		if CompsName == "PCRF" {
			name = fmt.Sprintf("%s", CompsName)
		} else {
			name = fmt.Sprintf("%s%d", CompsName, i)
		}

		vm, exists := VmList[name]

		if !exists {
			fmt.Printf("%s%d not found in the %s list !", CompsName, indexCount, CompsName)
			continue
		}

		nets := make(map[string]string)
		for _, netw := range vm.Networks {
			nets[netw.Name] = netw.IP
		}
		// CompsNetworksMaps = append(CompsNetworksMaps, CompsNetwork)
		CompsNetworksMaps[name] = nets

		result[name] = ComponentData{
			Name:                name,
			Networks:            nets,
			ComponentsToConnect: vm.ComponentsToConnect,
		}

		// fmt.Println(CompConn)
		// time.Sleep(1 * time.Second)
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
			flagErr = true
		}
	}
	if flagErr == true {
		return
	}

	// mapping key for getting VMs list with Name of the vms not by ID
	VmList := make(map[string]VM)
	for _, name := range vms {
		VmList[name.Name] = name
	}

	time.Sleep(1 * time.Second)

	CompsMapsMME := NetworkGeneratorComps(13, "MME", VmList)
	CompsMapsHSS := NetworkGeneratorComps(3, "HSS", VmList)
	CompsMapSGWC := NetworkGeneratorComps(2, "SGWC", VmList)
	CompsMapSGWU := NetworkGeneratorComps(2, "SGWU", VmList)
	CompsMapSMF := NetworkGeneratorComps(2, "SMF", VmList)
	CompsMapUPF := NetworkGeneratorComps(2, "UPF", VmList)
	CompsMapPCRF := NetworkGeneratorComps(1, "PCRF", VmList)

	mmesByName := BuildAllMMEs(CompsMapsMME)
	hsssByName := BuildAllHSSs(CompsMapsHSS)
	sgwcsByName := BuildAllSGWCs(CompsMapSGWC)
	sgwusByName := BuildAllSGWUs(CompsMapSGWU)
	smfsByName := BuildAllSMFs(CompsMapSMF)
	upfsByName := BuildAllUPFs(CompsMapUPF)
	pcrfsByName := BuildAllPCRFs(CompsMapPCRF)

	core_name = "{{ core_name }}"
	var_path = "/var/log/" + core_name + "/"
	tls_path = "/etc/" + core_name + "/tls/"
	inventory_hostname = "{{ inventory_hostname }}"
	diam_groupNames = "{{ group_names[1] }}"
	non_diam_groupNames = "{{ group_names[0] }}"
	var_path_diameter := "/etc/" + core_name + "/freeDiameter/"
	diam_realm = "epc.mnc{{ plmn.mnc }}.mcc{{ plmn.mcc }}.3gppnetwork.org"

	data := struct {
		SGWC_managementIP       string
		SGWC_managementPort     int
		SGWC_s11                string
		SGWC_s11Port            int
		SGWC_sxa                string
		SGWC_sxaPort            int
		SGWC_s5c                string
		SGWC_s5cPort            int
		SGWC_componentToConnect []string
		MME_s11                 string
		MME_s11Port             int
		MME_managementIP        string
		MME_managementPort      int
		MME_s1ap                string
		MME_s1apPort            int
		MME_s6a                 string
		MME_s6aPort             int
		MME1_componentToConnect []string
		SGWU_sxa                string
		SGWU_sxaPort            int
		SGWU_s5u                string
		SGWU_s5uPort            int
		SGWU_s1u                string
		SGWU_s1uPort            int
		SGWU_managementIP       string
		SGWU_managementPort     int
		SGWU_componentToConnect []string
		SMF_managementIP        string
		SMF_managementPort      int
		SMF_gx                  string
		SMF_gxPort              int
		SMF_s5c                 string
		SMF_s5cPort             int
		SMF_sxb                 string
		SMF_sxbPort             int
		SMF_sxu                 string
		SMF_sxuPort             int
		SMF_componentToConnect  []string
		UPF_managementIP        string
		UPF_managementPort      int
		UPF_sxb                 string
		UPF_sxbPort             int
		UPF_sxu                 string
		UPF_sxuPort             int
		UPF_s5u                 string
		UPF_s5uPort             int
		UPF_sgi                 string
		UPF_sgiPort             int
		UPF_componentToConnect  []string
		HSS_managementIP        string
		HSS_managementPort      int
		HSS_s6a                 string
		HSS_s6aPort             int
		HSS_componentToConnect  []string
		PCRF_managementIP       string
		PCRF_managementPort     int
		PCRF_gx                 string
		PCRF_gxPort             int
		PCRF_componentToConnect []string
		Core_name               string
		Var_path                string
		Tls_path                string
		Inventory_hostname      string
		Diam_groupNames         string
		Non_diam_groupNames     string
		Diameter_path           string
		Diam_Realm              string
		Gx_secPort              int
		S6a_secPort             int
	}{sgwcsByName["SGWC1"].Networks["VM Network"],
		sgwc_managementPort,
		sgwcsByName["SGWC1"].Networks["s11-sgwc1"],
		sgwc_s11Port,
		sgwcsByName["SGWC1"].Networks["sxa-sgwc1"],
		sgwc_sxaPort,
		sgwcsByName["SGWC1"].Networks["s5c-sgwc1"],
		sgwc_s5cPort,
		sgwcsByName["SGWC1"].ComponentsToConnect,
		mmesByName["MME1"].Networks["s11-mme1"],
		mme_s11Port,
		mmesByName["MME1"].Networks["VM Network"],
		mme_managementPort,
		mmesByName["MME1"].Networks["s1ap-mme1"],
		mme_s1apPort,
		mmesByName["MME1"].Networks["s6a-mme1"],
		mme_s6aPort,
		mmesByName["MME1"].ComponentsToConnect,
		sgwusByName["SGWU1"].Networks["sxa-sgwu1"],
		sgwu_sxaPort,
		sgwusByName["SGWU1"].Networks["s5u-sgwu1"],
		sgwu_s5uPort,
		sgwusByName["SGWU1"].Networks["s1u-sgwu1"],
		sgwu_s1uPort,
		sgwusByName["SGWU1"].Networks["VM Network"],
		sgwu_managementPort,
		sgwusByName["SGWU1"].ComponentsToConnect,
		smfsByName["SMF1"].Networks["VM Network"],
		smf_managementPort,
		smfsByName["SMF1"].Networks["gx-smf1"],
		smf_gxPort,
		smfsByName["SMF1"].Networks["s5c-smf1"],
		smf_s5cPort,
		smfsByName["SMF1"].Networks["sxb-smf1"],
		smf_sxbPort,
		smfsByName["SMF1"].Networks["sxu-smf1"],
		smf_sxuPort,
		smfsByName["SMF1"].ComponentsToConnect,
		upfsByName["UPF1"].Networks["VM Network"],
		upf_managementPort,
		upfsByName["UPF1"].Networks["sxb-upf1"],
		upf_sxbPort,
		upfsByName["UPF1"].Networks["sxu-upf1"],
		upf_sxuPort,
		upfsByName["UPF1"].Networks["s5u-upf1"],
		upf_s5uPort,
		upf_SGI,
		upf_sgiPort,
		upfsByName["UPF1"].ComponentsToConnect,
		hsssByName["HSS1"].Networks["VM Network"],
		hss_managementPort,
		hsssByName["HSS1"].Networks["s6a-hss1"],
		hss_s6aPort,
		hsssByName["HSS1"].ComponentsToConnect,
		pcrfsByName["PCRF"].Networks["VM Network"],
		pcrf_managementPort,
		pcrfsByName["PCRF"].Networks["gx-pcrf"],
		pcrf_gxPort,
		pcrfsByName["PCRF"].ComponentsToConnect,
		core_name,
		var_path,
		tls_path,
		inventory_hostname,
		diam_groupNames,
		non_diam_groupNames,
		var_path_diameter,
		diam_realm,
		gx_secPort,
		s6a_secPort,
	}

	// fmt.Printf("%#v\n", data.MME1_componentToConnect)
	// fmt.Printf("len=%d\n", len(data.MME1_componentToConnect))

	// time.Sleep(1000 * time.Second)
	yamlData := `all:
  vars:
    core_name: bbdh
    db_uri: mongodb://localhost/{{ .Core_name }}
    configs_path: /etc/{{ .Core_name }}
    var_path: /var/log/{{ .Core_name }}/
    var_path_diameter: /etc/{{ .Core_name }}/freeDiameter/
    tls_path: /etc/{{ .Core_name }}/tls/ 
    diam_lib_dir: /usr/lib
    max_ue: 1024
    # PLMN that use for most of the components
    plmn:
      mcc: 432
      mnc: 080
    diam_realm: {{ .Diam_Realm }}

  children:
    sgwc:
      hosts:
        sgwc1:
          ansible_host: {{ .SGWC_managementIP }}
          managementPort: {{ .SGWC_managementPort }}
          ansible_user: mos
          ansible_password: q 
          ansible_become_pass: q
          logger: "{{ .Var_path }}{{ .Non_diam_groupNames }}.log"
          s11_addr: {{ .SGWC_s11 }}
          s11_port: {{ .SGWC_s11Port }}
          s5c_addr: {{ .SGWC_s5c }}
          s5c_port: {{ .SGWC_s5cPort }}
          sxa_addr: {{ .SGWC_sxa }}
          sxa_port: {{ .SGWC_sxaPort }}
          components:
          {{- range .SGWC_componentToConnect }}
            - {{ . }}
          {{- end }}
    sgwu:
      hosts:
        sgwu1:
          ansible_host: {{ .SGWU_managementIP }}
          managementPort: {{ .SGWU_managementPort }}
          ansible_user: mos
          ansible_password: q 
          ansible_become_pass: q
          logger: "{{ .Var_path }}{{ .Non_diam_groupNames }}.log"
          s5u_addr: {{ .SGWU_s5u }}
          s5u_port: {{ .SGWU_s5uPort }}
          sxa_addr: {{ .SGWU_sxa }}
          sxa_port: {{ .SGWU_sxaPort }}
          s1u_addr: {{ .SGWU_s1u }}
          s1u_port: {{ .SGWU_s1uPort }}
    upf:
      hosts:
        upf1:
          ansible_host: {{ .UPF_managementIP }}
          managementPort: {{ .UPF_managementPort }}
          ansible_user: mos
          ansible_password: q 
          ansible_become_pass: q
          logger: "{{ .Var_path }}{{ .Non_diam_groupNames }}.log"
          sxb_addr: {{ .UPF_sxb }}
          sxb_port: {{ .UPF_sxbPort }}
          sxu_addr: {{ .UPF_sxu }}
          sxu_port: {{ .UPF_sxuPort }}
          s5u_addr: {{ .UPF_s5u }}
          s5u_port: {{ .UPF_s5uPort }}
          sgi_addr: {{ .UPF_sgi }}
          sgi_port: {{ .UPF_sgiPort }}
          subnet:
              addr: 10.45.0.1/16
              dev: ogstun
              apn: internet
          smf_addr: {{ .SMF_managementIP }}

    # all diameter peers metagroup
    diam_peers:
      children:
        mme:
          hosts:
            mme1:
              ansible_host: {{ .MME_managementIP }}
              managementPort: {{ .MME_managementPort }}
              ansible_user: mos
              ansible_password: q 
              ansible_become_pass: q
              logger: "{{ .Var_path }}{{ .Diam_groupNames }}.log"
              freeDiameter: "{{ .Diameter_path }}{{ .Diam_groupNames }}.conf"
              tac: 1 
              s11_addr: {{ .MME_s11 }}
              s11_port: {{ .MME_s11Port }}
              s1ap_addr: {{ .MME_s1ap }}
              s1ap_port: {{ .MME_s1apPort }}
              s6a_addr: {{ .MME_s6a }}
              s6a_port: {{ .MME_s6aPort }}
              s6a_secport: {{ .S6a_secPort }}
              components:
              {{- range .MME1_componentToConnect}}
                - {{ . }}
			  {{- end }}

              # freeDiameter variables
              diam_Id_host: "{{ .Inventory_hostname }}.{{ .Diam_Realm }}"

        hss:
          hosts:
            hss1:
              ansible_host: {{ .HSS_managementIP }}
              managementPort: {{ .HSS_managementPort }}
              ansible_user: mos
              ansible_password: q 
              ansible_become_pass: q
              logger: "{{ .Var_path }}{{ .Diam_groupNames }}.log"
              freeDiameter: "{{ .Diameter_path }}{{ .Diam_groupNames }}.conf"
              db_uri: mongodb://localhost/{{ .Core_name }}


              s6a_addr: {{ .HSS_s6a }}
              s6a_port: {{ .HSS_s6aPort }}
              s6a_secport: {{ .S6a_secPort }}
              components:
              {{- range .HSS_componentToConnect }}
                - {{ . }}
              {{- end }}

              # freeDiameter variables
              diam_Id_host: "{{ .Inventory_hostname }}.{{ .Diam_Realm }}"

        smf:
          hosts:
            smf1:
              ansible_host: {{ .SMF_managementIP }}
              managementPort: {{ .SMF_managementPort }}
              ansible_user: mos
              ansible_password: q 
              ansible_become_pass: q
              logger: "{{ .Var_path }}{{ .Diam_groupNames }}.log"
              freeDiameter: "{{ .Diameter_path }}{{ .Diam_groupNames }}.conf"
              sbi_addr: 9877
              gx_addr: {{ .SMF_gx }}
              gx_port: {{ .SMF_gxPort }}
              gx_secport: {{ .Gx_secPort }}
              s5c_addr: {{ .SMF_s5c }}
              s5c_port: {{ .SMF_s5cPort }}
              sxb_addr: {{ .SMF_sxb }}
              sxb_port: {{ .SMF_sxbPort }}
              sxu_addr: {{ .SMF_sxu }}
              sxu_port: {{ .SMF_sxuPort }}
              subnet:
                  addr: 10.45.0.1/16
                  dev: ogstun
                  apn: internet
              dns:
                  primary: 8.8.8.8
                  secondary: 8.8.4.4

              # freeDiameter variables
              diam_Id_host: "{{ .Inventory_hostname }}.{{ .Diam_Realm }}"
              components:
              {{- range .SMF_componentToConnect }}
                - {{ . }}
              {{- end }}

        pcrf:
          hosts:
            pcrf1:
              ansible_host: {{ .PCRF_managementIP }}
              managementPort: {{ .PCRF_managementPort }}
              ansible_user: mos
              ansible_password: q 
              ansible_become_pass: q
              logger: "{{ .Var_path }}{{ .Diam_groupNames }}.log"
              freeDiameter: "{{ .Diameter_path }}{{ .Diam_groupNames }}.conf"
              db_uri: mongodb://localhost/bbdh

              gx_addr: {{ .PCRF_gx }}
              gx_port: {{ .PCRF_gxPort }}
              gx_secport: {{ .Gx_secPort }}
              components:
              {{- range .PCRF_componentToConnect }}
                - {{ . }}
              {{- end }}

              # freeDiameter variables
              diam_Id_host: "{{ .Inventory_hostname }}.{{ .Diam_Realm }}"
               `

	templateTest := template.Must(template.New("yaml").Parse(yamlData))

	var buf bytes.Buffer

	if err := templateTest.Execute(&buf, data); err != nil {
		panic(err)
	}

	// inventoryPath, fileName := "ansible/ansible-core-deploy/inventory/", "main.yml"
	fileName := "demo.yml"
	// os.WriteFile(inventoryPath+fileName, buf.Bytes(), 0644)
	os.WriteFile(fileName, buf.Bytes(), 0644)
	fmt.Printf("\n%sGenerating %s  file ...%s", color.Yellow, fileName, color.Reset)
	time.Sleep(1 * time.Second)
	fmt.Printf("\n%s%s generated in the current path%s\n\n", color.Green, fileName, color.Reset)
}
