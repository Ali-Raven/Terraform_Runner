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
)

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

	CompsMapsMME := NetworkGeneratorComps("MME", VmList)
	CompsMapsHSS := NetworkGeneratorComps("HSS", VmList)
	CompsMapSGWC := NetworkGeneratorComps("SGWC", VmList)
	CompsMapSGWU := NetworkGeneratorComps("SGWU", VmList)
	CompsMapSMF := NetworkGeneratorComps("SMF", VmList)
	CompsMapUPF := NetworkGeneratorComps("UPF", VmList)
	CompsMapPCRF := NetworkGeneratorComps("PCRF", VmList)

	mmesByName := generators.BuildAllMMEs(CompsMapsMME)
	hsssByName := generators.BuildAllHSSs(CompsMapsHSS)
	sgwcsByName := generators.BuildAllSGWCs(CompsMapSGWC)
	sgwusByName := generators.BuildAllSGWUs(CompsMapSGWU)
	smfsByName := generators.BuildAllSMFs(CompsMapSMF)
	upfsByName := generators.BuildAllUPFs(CompsMapUPF)
	pcrfsByName := generators.BuildAllPCRFs(CompsMapPCRF)

	// fmt.Println(pcrfsByName)
	// time.Sleep(10000 * time.Second)

	Core_Name = "bbdh"
	Var_path = "/var/log/" + Core_Name + "/"
	Tls_path = "/etc/" + Core_Name + "/tls/"
	Inventory_hostname = "{{ inventory_hostname }}"
	Diam_groupNames = "{{ group_names[1] }}"
	Non_diam_groupNames = "{{ group_names[0] }}"
	Var_path_diameter = "/etc/" + Core_Name + "/freeDiameter/"
	Diam_Realm = "epc.mnc{{ plmn.mnc }}.mcc{{ plmn.mcc }}.3gppnetwork.org"
	Var_path_Comps = "/var/log/{{ core_name }}/{{ group_names[0] }}.log"

	// calling Data function

	givenDataTemplate := Data(mmesByName, hsssByName, sgwcsByName, sgwusByName, smfsByName, upfsByName, pcrfsByName)

	yamlData := `all:
  vars:
    core_name: "{{ .Core_Name }}"
    db_uri: mongodb://localhost/{{ .Core_Name }}
    configs_path: /etc/{{ .Core_Name }}
    var_path: {{ .Var_path }}
    var_path_diameter: {{ .Var_path_diameter }}
    tls_path: {{ .Tls_path }}
    diam_lib_dir: /usr/lib
    user: {{ .User }}
    max_ue: 1024
    # PLMN that use for most of the components
    plmn:
      mcc: 432
      mnc: 080
    diam_realm: {{ .Diam_Realm }}

  children:
    # ============================================================
    # SGWC CLUSTER
    # ============================================================
    sgwc:
      children:
        {{- range $i, $c := .SGWCsCluster }}
        sgwc{{ add $i 1 }}_cluster:
          hosts:
            sgwc{{ add $i 1 }}:
              ansible_host: {{ $c.Master.ANSIBLE_HOST }}
              managementPort: {{ $c.Master.ManagementPort }}
              ansible_user: {{ $c.Master.User }}
              ansible_password: q
              ansible_become_pass: q
              keepalived_role: MASTER
              keepalived_priority: 101

            sgwc{{ add $i 1 }}_backup:
              ansible_host: {{ $c.Backup.ANSIBLE_HOST }}
              managementPort: 22
              ansible_user: {{ $c.Backup.User }}
              ansible_password: q
              ansible_become_pass: q
              keepalived_role: BACKUP
              keepalived_priority: 100

          vars:
            logger: "{{ $.Var_path }}{{ $.Non_diam_groupNames }}.log"
            routerId: {{ add 90 $i }}
            s11_addr: {{ $c.S11Addr }}/{{ $c.S11Subnet }}
            s11_port: {{ $c.S11Port }}
            s11_gateway: {{ $c.S11Gateway }}
            s11_subnet: {{ $c.S11Subnet}}

            s5c_addr: {{ $c.S5cAddr }}/{{ $c.S5cSubnet }}
            s5c_port: {{ $c.S5cPort }}
            s5c_gateway: {{ $c.S5cGateway }}
            s5c_subnet: {{ $c.S5cSubnet }} 

            sxa_addr: {{ $c.SxaAddr }}/{{ $c.SxaSubnet }}
            sxa_port: {{ $c.SxaPort }}
            sxa_gateway: {{ $c.SxaGateway }}
            sxa_subnet: {{ $c.SxaSubnet }}

            components:
            {{- range $c.Components }}
              - {{ . }}
            {{- end }}

        {{- end }}

    # ============================================================
    # SGWU CLUSTER
    # ============================================================
    sgwu:
      children:
        {{- range $i, $c := .SGWUsCluster }}
        sgwu{{ add $i 1 }}_cluster:
          hosts:
            sgwu{{ add $i 1 }}:
              ansible_host: {{ $c.Master.ANSIBLE_HOST }}
              managementPort: {{ $c.Master.ManagementPort }}
              ansible_user: {{ $c.Master.User }}
              ansible_password: q
              ansible_become_pass: q
              keepalived_role: MASTER
              keepalived_priority: 101

            sgwu{{ add $i 1 }}_backup:
              ansible_host: {{ $c.Backup.ANSIBLE_HOST }}
              managementPort: 22
              ansible_user: {{ $c.Backup.User }}
              ansible_password: q
              ansible_become_pass: q
              keepalived_role: BACKUP
              keepalived_priority: 100

          vars:
            logger: "{{ $.Var_path }}{{ $.Non_diam_groupNames }}.log"
            routerId: {{ add 100 $i }}
            s1u_addr: {{ $c.S1uAddr }}/{{ $c.S1uSubnet }}
            s1u_port: {{ $c.S1uPort }}
            s1u_gateway: {{ $c.S1uGateway }}
            s1u_subnet: {{ $c.S1uSubnet}}

            s5u_addr: {{ $c.S5uAddr }}/{{ $c.S5uSubnet }}
            s5u_port: {{ $c.S5uPort }}
            s5u_gateway: {{ $c.S5uGateway }}
            s5u_subnet: {{ $c.S5uSubnet }} 

            sxa_addr: {{ $c.SxaAddr }}/{{ $c.SxaSubnet }}
            sxa_port: {{ $c.SxaPort }}
            sxa_gateway: {{ $c.SxaGateway }}
            sxa_subnet: {{ $c.SxaSubnet }}

            components:
            {{- range $c.Components }}
              - {{ . }}
            {{- end }}

        {{- end }}

    # ============================================================
    # UPF CLUSTER
    # ============================================================
    upf:
      children:
        {{- range $i, $c := .UPFsCluster }}
            upf{{ add $i 1 }}_cluster:
              hosts:
                upf{{ add $i 1 }}:
                  ansible_host: {{ $c.Master.ANSIBLE_HOST }}
                  managementPort: {{ $c.Master.ManagementPort }}
                  ansible_user: {{ $c.Master.User }}
                  ansible_password: {{ $c.Master.Password }}
                  ansible_become_pass: {{ $c.Master.BecomePass }}
                  keepalived_role: {{ $c.Master.KeepalivedRole }}
                  keepalived_priority: {{ $c.Master.KeepalivedPrio }}

                upf{{ add $i 1 }}_backup:
                  ansible_host: {{ $c.Backup.ANSIBLE_HOST }}
                  managementPort: {{ $c.Backup.ManagementPort }}
                  ansible_user: {{ $c.Backup.User }}
                  ansible_password: {{ $c.Backup.Password }}
                  ansible_become_pass: {{ $c.Backup.BecomePass }}
                  keepalived_role: {{ $c.Backup.KeepalivedRole }}
                  keepalived_priority: {{ $c.Backup.KeepalivedPrio }}

              vars:
                logger: "{{ $.Var_path }}{{ $.Diam_groupNames }}.log"
                freeDiameter: "{{ $.Var_path_diameter }}{{ $.Diam_groupNames }}.conf"
                routerId: {{ add 120 $i }}
                sxb_addr: {{ $c.SxbAddr }}/{{ $c.SxbSubnet }}
                sxb_port: {{ $c.SxbPort }}
                sxb_gateway: {{ $c.SxbGateway }}
                sxb_subnet: {{ $c.SxbSubnet }}

                sxu_addr: {{ $c.SxuAddr }}/{{ $c.SxuSubnet}}
                sxu_port: {{ $c.SxuPort }}
                sxu_gateway: {{ $c.SxuGateway }}
                sxu_subnet: {{ $c.SxuSubnet}}

                s5c_addr: {{ $c.S5uAddr }}/{{ $c.S5uSubnet }}
                s5c_port: {{ $c.S5uPort }}
                s5c_gateway: {{ $c.S5uGateway }}
                s5c_subnet: {{ $c.S5uSubnet }} 

                components:
                {{- range $c.Components }}
                  - {{ . }}
                {{- end }}

            {{- end }}
         
    # ============================================================
    # DIAMETER PEERS METAGROUP
    # ============================================================
    diam_peers:
      children:
        # --------------------------------------------------------
        # MME CLUSTER
        # --------------------------------------------------------
        mme:
          children:
            {{- range $i, $c := .MMEsCluster }}
            mme{{ add $i 1 }}_cluster:
              hosts:
                mme{{ add $i 1 }}:
                  ansible_host: {{ $c.Master.ANSIBLE_HOST }}
                  managementPort: {{ $c.Master.ManagementPort }}
                  ansible_user: {{ $c.Master.User }}
                  ansible_password: {{ $c.Master.Password }}
                  ansible_become_pass: {{ $c.Master.BecomePass }}
                  keepalived_role: {{ $c.Master.KeepalivedRole }}
                  keepalived_priority: {{ $c.Master.KeepalivedPrio }}

                mme{{ add $i 1 }}_backup:
                  ansible_host: {{ $c.Backup.ANSIBLE_HOST }}
                  managementPort: {{ $c.Backup.ManagementPort }}
                  ansible_user: {{ $c.Backup.User }}
                  ansible_password: {{ $c.Backup.Password }}
                  ansible_become_pass: {{ $c.Backup.BecomePass }}
                  keepalived_role: {{ $c.Backup.KeepalivedRole }}
                  keepalived_priority: {{ $c.Backup.KeepalivedPrio }}

              vars:
                logger: "{{ $.Var_path }}{{ $.Diam_groupNames }}.log"
                freeDiameter: "{{ $.Var_path_diameter }}{{ $.Diam_groupNames }}.conf"
                routerId: {{ add $i 1 }}
                s1ap_addr: {{ $c.S1apAddr }}/{{ $c.S1apSubnet }}
                s1ap_port: {{ $c.S1apPort }}
                s1ap_gateway: {{ $c.S1apGateway }}
                s1ap_subnet: {{ $c.S1apSubnet }}

                s11_addr: {{ $c.S11Addr }}/{{ $c.S11Subnet}}
                s11_port: {{ $c.S11Port }}
                s11_gateway: {{ $c.S11Gateway }}
                s11_subnet: {{ $c.S11Subnet}}

                s5c_addr: {{ $c.S5cAddr }}/{{ $c.S5cSubnet }}
                s5c_port: {{ $c.S5cPort }}
                s5c_gateway: {{ $c.S5cGateway }}
                s5c_subnet: {{ $c.S5cSubnet }} 

                s6a_addr: {{ $c.S6aAddr }}/{{ $c.S6aSubnet }}
                s6a_port: {{ $c.S6aPort }}
                s6a_gateway: {{ $c.S6aGateway }}
                s6a_subnet: {{ $c.S6aSubnet }}

                components:
                {{- range $c.Components }}
                  - {{ . }}
                {{- end }}

                gummei:
                  - plmn_id:
                      mcc: 432
                      mnc: 80
                    mme_gid: 1111
                    mme_code: 111
                
                tai:
                  - plmn_id:
                      mcc: 432
                      mnc: 80
                    tac: [30511, 30512, 30513, 30514, 30516, 30517, 30519, 30581, 30582, 30583, 30585, 30590]
                  - plmn_id:
                      mcc: 432
                      mnc: 11
                    tac: [30509]

                non_restrict_plmn:
                  - plmn_id:
                      mcc: 432
                      mnc: 11
                      decision_digits: 29997
                  - plmn_id:
                      mcc: 432
                      mnc: 11
                      decision_digits: 00000

            {{- end }}

        # --------------------------------------------------------
        # HSS CLUSTER
        # --------------------------------------------------------
        hss:
          children:
            {{- range $i, $c := .HSSsCluster }}
            hss{{ add $i 1 }}_cluster:
              hosts:
                hss{{ add $i 1 }}:
                  ansible_host: {{ $c.Master.ANSIBLE_HOST }}
                  managementPort: {{ $c.Master.ManagementPort }}
                  ansible_user: {{ $c.Master.User }}
                  ansible_password: {{ $c.Master.Password }}
                  ansible_become_pass: {{ $c.Master.BecomePass }}
                  keepalived_role: {{ $c.Master.KeepalivedRole }}
                  keepalived_priority: {{ $c.Master.KeepalivedPrio }}

                hss{{ add $i 1 }}_backup:
                  ansible_host: {{ $c.Backup.ANSIBLE_HOST }}
                  managementPort: {{ $c.Backup.ManagementPort }}
                  ansible_user: {{ $c.Backup.User }}
                  ansible_password: {{ $c.Backup.Password }}
                  ansible_become_pass: {{ $c.Backup.BecomePass }}
                  keepalived_role: {{ $c.Backup.KeepalivedRole }}
                  keepalived_priority: {{ $c.Backup.KeepalivedPrio }}

              vars:
                logger: "{{ $.Var_path }}{{ $.Diam_groupNames }}.log"
                freeDiameter: "{{ $.Var_path_diameter }}{{ $.Diam_groupNames }}.conf"
                routerId: {{ add 50 $i }}
                hss_id: {{ add $i 1 }}
                s6a_addr: {{ $c.S6aAddr }}/{{ $c.S6aSubnet }}
                s6a_port: {{ $c.S6aPort }}
                s6a_gateway: {{ $c.S6aGateway }}
                s6a_subnet: {{ $c.S6aSubnet }}
                components:
                {{- range $c.Components }}
                  - {{ . }}
                {{- end }}

            {{- end }} 

        # --------------------------------------------------------
        # SMF CLUSTER
        # --------------------------------------------------------
        smf:
          children:
            {{- range $i, $c := .SMFsCluster }}
            smf{{ add $i 1 }}_cluster:
              hosts:
                smf{{ add $i 1 }}:
                  ansible_host: {{ $c.Master.ANSIBLE_HOST }}
                  managementPort: {{ $c.Master.ManagementPort }}
                  ansible_user: {{ $c.Master.User }}
                  ansible_password: {{ $c.Master.Password }}
                  ansible_become_pass: {{ $c.Master.BecomePass }}
                  keepalived_role: {{ $c.Master.KeepalivedRole }}
                  keepalived_priority: {{ $c.Master.KeepalivedPrio }}

                smf{{ add $i 1 }}_backup:
                  ansible_host: {{ $c.Backup.ANSIBLE_HOST }}
                  managementPort: {{ $c.Backup.ManagementPort }}
                  ansible_user: {{ $c.Backup.User }}
                  ansible_password: {{ $c.Backup.Password }}
                  ansible_become_pass: {{ $c.Backup.BecomePass }}
                  keepalived_role: {{ $c.Backup.KeepalivedRole }}
                  keepalived_priority: {{ $c.Backup.KeepalivedPrio }}

              vars:
                logger: "{{ $.Var_path }}{{ $.Diam_groupNames }}.log"
                freeDiameter: "{{ $.Var_path_diameter }}{{ $.Diam_groupNames }}.conf"
                routerId: {{ add 70 $i }}
                sxb_addr: {{ $c.SxbAddr }}/{{ $c.SxbSubnet }}
                sxb_port: {{ $c.SxbPort }}
                sxb_gateway: {{ $c.SxbGateway }}
                sxb_subnet: {{ $c.SxbSubnet }}

                sxu_addr: {{ $c.SxuAddr }}/{{ $c.SxuSubnet}}
                sxu_port: {{ $c.SxuPort }}
                sxu_gateway: {{ $c.SxuGateway }}
                sxu_subnet: {{ $c.SxuSubnet}}

                s5c_addr: {{ $c.S5cAddr }}/{{ $c.S5cSubnet }}
                s5c_port: {{ $c.S5cPort }}
                s5c_gateway: {{ $c.S5cGateway }}
                s5c_subnet: {{ $c.S5cSubnet }} 

                gx_addr: {{ $c.GxAddr }}/{{ $c.GxSubnet }}
                gx_port: {{ $c.GxPort }}
                gx_gateway: {{ $c.GxGateway }}
                gx_subnet: {{ $c.GxSubnet }}

                components:
                {{- range $c.Components }}
                  - {{ . }}
                {{- end }}

            {{- end }}

        # --------------------------------------------------------
        # PCRF CLUSTER
        # --------------------------------------------------------
        pcrf:
          children:
            {{- range $i, $c := .PCRFsCluster }}
            pcrf{{ add $i 1 }}_cluster:
              hosts:
                pcrf{{ add $i 1 }}:
                  ansible_host: {{ $c.Master.ANSIBLE_HOST }}
                  managementPort: {{ $c.Master.ManagementPort }}
                  ansible_user: {{ $c.Master.User }}
                  ansible_password: {{ $c.Master.Password }}
                  ansible_become_pass: {{ $c.Master.BecomePass }}
                  keepalived_role: {{ $c.Master.KeepalivedRole }}
                  keepalived_priority: {{ $c.Master.KeepalivedPrio }}

                pcrf{{ add $i 1 }}_backup:
                  ansible_host: {{ $c.Backup.ANSIBLE_HOST }}
                  managementPort: {{ $c.Backup.ManagementPort }}
                  ansible_user: {{ $c.Backup.User }}
                  ansible_password: {{ $c.Backup.Password }}
                  ansible_become_pass: {{ $c.Backup.BecomePass }}
                  keepalived_role: {{ $c.Backup.KeepalivedRole }}
                  keepalived_priority: {{ $c.Backup.KeepalivedPrio }}

              vars:
                logger: "{{ $.Var_path }}{{ $.Diam_groupNames }}.log"
                freeDiameter: "{{ $.Var_path_diameter }}{{ $.Diam_groupNames }}.conf"
                routerId: {{ add 80 $i }}
                gx_addr: {{ $c.GxAddr }}/{{ $c.GxSubnet }}
                gx_port: {{ $c.GxPort }}
                gx_gateway: {{ $c.GxGateway }}
                gx_subnet: {{ $c.GxSubnet }}
                components:
                {{- range $c.Components }}
                  - {{ . }}
                {{- end }}

            {{- end }}
               `

	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	}
	templateTest := template.Must(template.New("yaml").Funcs(funcMap).Parse(yamlData))

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
