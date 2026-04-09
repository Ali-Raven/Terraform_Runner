package main

import (
	"bytes"
	"fmt"
	"github.com/TwiN/go-color"
	"os"
	"text/template"
	"time"
)

var (
	ssh_defaultPort     int = 22
	sgwc_managementIP   string
	sgwc_managementPort int = ssh_defaultPort
	sgwc_s11            string
	sgwc_s11Port        int = 2123
	sgwc_sxa            string
	sgwc_sxaPort        int = 8805
	sgwc_s5c            string
	sgwc_s5cPort        int = 2124
	sgwu_managementIP   string
	sgwu_managementPort int = ssh_defaultPort
	sgwu_sxa            string
	sgwu_sxaPort        int = 2152
	sgwu_s5u            string
	sgwu_s5uPort        int = 8805
	sgwu_s1u            string
	sgwu_s1uPort        int = 3333
	upf_managemetIP     string
	upf_managementPort  int = ssh_defaultPort
	upf_sxb             string
	upf_sxbPort         int = 8805
	upf_sxu             string
	upf_sxuPort         int = 8806
	upf_s5u             string
	upf_s5uPort         int = 2153
	upf_SGI             string
	upf_sgiPort         int = 2152
	smf_managementIP    string
	smf_managementPort  int = ssh_defaultPort
	smf_gx              string
	smf_gxPort          int = 2123
	gx_secPort          int = 5868
	smf_s5c             string
	smf_s5cPort         int = 8805
	smf_sxb             string
	smf_sxbPort         int = 2152
	smf_sxu             string
	smf_sxuPort         int = 8806
	mme_s11             string
	mme_s11Port         int = 2123
	mme_s1ap            string
	mme_s1apPort        int = 36412
	mme_s6a             string
	mme_s6aPort         int = 2221
	s6a_secPort         int = 5868
	mme_managementIP    string
	mme_managementPort  int = ssh_defaultPort
	hss_managementIP    string
	hss_managementPort  int = ssh_defaultPort
	hss_s6a             string
	hss_s6aPort         int = 2223
	pcrf_managementIP   string
	pcrf_managementPort int = ssh_defaultPort
	pcrf_gx             string
	pcrf_gxPort         int = 4434
	core_name           string
	var_path            string
	tls_path            string
	inventory_hostname  string
	diam_realm          string
	flagErr             bool
)

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

	time.Sleep(1 * time.Second)

	// Assigning MME vars
	mme_managementIP = vms[0].Networks[0].IP
	mme_s1ap = vms[0].Networks[1].IP
	mme_s6a = vms[0].Networks[2].IP
	mme_s11 = vms[0].Networks[3].IP

	// Assigning SGW-C
	sgwc_managementIP = vms[3].Networks[0].IP
	sgwc_s11 = vms[3].Networks[1].IP
	sgwc_sxa = vms[3].Networks[2].IP
	sgwc_s5c = vms[3].Networks[3].IP

	// Assigning SGW-U
	sgwu_managementIP = vms[4].Networks[0].IP
	sgwu_sxa = vms[4].Networks[1].IP
	sgwu_s5u = vms[4].Networks[2].IP
	sgwu_s1u = vms[4].Networks[3].IP

	// Assigning SMF
	smf_managementIP = vms[5].Networks[0].IP
	smf_gx = vms[5].Networks[1].IP
	smf_s5c = vms[5].Networks[2].IP
	smf_sxb = vms[5].Networks[3].IP
	smf_sxu = vms[5].Networks[4].IP

	// Assigning UPF
	upf_managemetIP = vms[6].Networks[0].IP
	upf_sxb = vms[6].Networks[1].IP
	upf_sxu = vms[6].Networks[2].IP
	upf_s5u = vms[6].Networks[3].IP
	upf_SGI = vms[6].Networks[4].IP

	// Assigning HSS
	hss_managementIP = vms[1].Networks[0].IP
	hss_s6a = vms[1].Networks[1].IP

	// Assigning PCRF
	pcrf_managementIP = vms[2].Networks[0].IP
	pcrf_gx = vms[2].Networks[1].IP

	core_name = "{{ core_name }}"
	var_path = "/var/log/" + core_name + "/"
	tls_path = "/etc/" + core_name + "/tls/"
	inventory_hostname = "{{ inventory_hostname }}"
	var_path_diameter := "/etc/" + core_name + "/freeDiameter/"
	diam_realm = "epc.mnc0{{ plmn.mnc }}.mcc{{ plmn.mcc }}.3gppnetwork.org"

	data := struct {
		SGWC_managementIP   string
		SGWC_managementPort int
		SGWC_s11            string
		SGWC_s11Port        int
		SGWC_sxa            string
		SGWC_sxaPort        int
		SGWC_s5c            string
		SGWC_s5cPort        int
		MME_s11             string
		MME_s11Port         int
		MME_managementIP    string
		MME_managementPort  int
		MME_s1ap            string
		MME_s1apPort        int
		MME_s6a             string
		MME_s6aPort         int
		SGWU_sxa            string
		SGWU_sxaPort        int
		SGWU_s5u            string
		SGWU_s5uPort        int
		SGWU_s1u            string
		SGWU_s1uPort        int
		SGWU_managementIP   string
		SGWU_managementPort int
		SMF_managementIP    string
		SMF_managementPort  int
		SMF_gx              string
		SMF_gxPort          int
		SMF_s5c             string
		SMF_s5cPort         int
		SMF_sxb             string
		SMF_sxbPort         int
		SMF_sxu             string
		SMF_sxuPort         int
		UPF_managementIP    string
		UPF_managementPort  int
		UPF_sxb             string
		UPF_sxbPort         int
		UPF_sxu             string
		UPF_sxuPort         int
		UPF_s5u             string
		UPF_s5uPort         int
		UPF_sgi             string
		UPF_sgiPort         int
		HSS_managementIP    string
		HSS_managementPort  int
		HSS_s6a             string
		HSS_s6aPort         int
		PCRF_managementIP   string
		PCRF_managementPort int
		PCRF_gx             string
		PCRF_gxPort         int
		Core_name           string
		Var_path            string
		Tls_path            string
		Inventory_hostname  string
		Diameter_path       string
		Diam_Realm          string
		Gx_secPort          int
		S6a_secPort         int
	}{sgwc_managementIP,
		sgwc_managementPort,
		sgwc_s11,
		sgwc_s11Port,
		sgwc_sxa,
		sgwc_sxaPort,
		sgwc_s5c,
		sgwc_s5cPort,
		mme_s11,
		mme_s11Port,
		mme_managementIP,
		mme_managementPort,
		mme_s1ap,
		mme_s1apPort,
		mme_s6a,
		mme_s6aPort,
		sgwu_sxa,
		sgwu_sxaPort,
		sgwu_s5u,
		sgwu_s5uPort,
		sgwu_s1u,
		sgwu_s1uPort,
		sgwu_managementIP,
		sgwu_managementPort,
		smf_managementIP,
		smf_managementPort,
		smf_gx,
		smf_gxPort,
		smf_s5c,
		smf_s5cPort,
		smf_sxb,
		smf_sxbPort,
		smf_sxu,
		smf_sxuPort,
		upf_managemetIP,
		upf_managementPort,
		upf_sxb,
		upf_sxbPort,
		upf_sxu,
		upf_sxuPort,
		upf_s5u,
		upf_s5uPort,
		upf_SGI,
		upf_sgiPort,
		hss_managementIP,
		hss_managementPort,
		hss_s6a,
		hss_s6aPort,
		pcrf_managementIP,
		pcrf_managementPort,
		pcrf_gx,
		pcrf_gxPort,
		core_name,
		var_path,
		tls_path,
		inventory_hostname,
		var_path_diameter,
		diam_realm,
		gx_secPort,
		s6a_secPort,
	}

	yamlData := `all:
  vars:
    core_name: bbdh
    db_uri: mongodb://localhost/{{ .Core_name }}
    configs_path: /etc/{{ .Core_name }}
    var_path: /var/log/{{ .Core_name }}/
    var_path_diameter: /etc/{{ .Core_name }}/freeDiameter/
    tls_path: /etc/{{ .Core_name }}/tls/ 
    diam_lib_dir: /usr/lib

	
    # PLMN that use for most of the components
    plmn:
      mcc: 432
      mnc: 085

  children:
    sgwc:
      hosts:
        sgwc1:
          ansible_host: {{ .SGWC_managementIP }}
          managementPort: {{ .SGWC_managementPort }}
          ansible_user: mos
          ansible_password: q 
          ansible_become_pass: q
          logger: "{{ .Var_path }}{{ .Inventory_hostname }}.log"
          s11_addr: {{ .SGWC_s11 }}
          s11_port: {{ .SGWC_s11Port }}
          s5c_addr: {{ .SGWC_s5c }}
          s5c_port: {{ .SGWC_s5cPort }}
          sxa_addr: {{ .SGWU_sxa }}
          sxa_port: {{ .SGWU_sxaPort }}

    sgwu:
      hosts:
        sgwu1:
          ansible_host: {{ .SGWU_managementIP }}
          managementPort: {{ .SGWU_managementPort }}
          ansible_user: mos
          ansible_password: q 
          ansible_become_pass: q
          logger: "{{ .Var_path }}{{ .Inventory_hostname }}.log"
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
          logger: "{{ .Var_path }}{{ .Inventory_hostname }}.log"
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
              logger: "{{ .Var_path }}{{ .Inventory_hostname }}.log"
              freeDiameter: "{{ .Diameter_path }}{{ .Inventory_hostname }}.conf"
              tac: 3 
              s11_addr: {{ .MME_s11 }}
              s11_port: {{ .MME_s11Port }}
              s1ap: {{ .MME_s1ap }}
              s1apPort: {{ .MME_s1apPort }}
              s6a_addr: {{ .MME_s6a }}
              s6a_port: {{ .MME_s6aPort }}
              s6a_secport: {{ .S6a_secPort }}

              # freeDiameter variables

              diam_realm: {{ .Diam_Realm }}
              diam_Id_host: "{{ .Inventory_hostname }}.{{ .Diam_Realm }}"

        hss:
          hosts:
            hss1:
              ansible_host: {{ .HSS_managementIP }}
              managementPort: {{ .HSS_managementPort }}
              ansible_user: mos
              ansible_password: q 
              ansible_become_pass: q
              logger: "{{ .Var_path }}{{ .Inventory_hostname }}.log"
              freeDiameter: "{{ .Diameter_path }}{{ .Inventory_hostname }}.conf"
              db_uri: mongodb://localhost/{{ .Core_name }}


              s6a_addr: {{ .HSS_s6a }}
              s6a_port: {{ .HSS_s6aPort }}
              s6a_secport: {{ .S6a_secPort }}

              # freeDiameter variables
              diam_realm: {{ .Diam_Realm }}
              diam_Id_host: "{{ .Inventory_hostname }}.{{ .Diam_Realm }}"

        smf:
          hosts:
            smf1:
              ansible_host: {{ .SMF_managementIP }}
              managementPort: {{ .SMF_managementPort }}
              ansible_user: mos
              ansible_password: q 
              ansible_become_pass: q
              logger: "{{ .Var_path }}{{ .Inventory_hostname }}.log"
              freeDiameter: "{{ .Diameter_path }}{{ .Inventory_hostname }}.conf"
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
              upf_sxb: {{ .UPF_sxb }}
              upf_sxbPort: {{ .UPF_sxbPort }}
              upf_sxu: {{ .UPF_sxu }}
              upf_sxuPort: {{ .UPF_sxuPort }}

              # freeDiameter variables
              diam_realm: {{ .Diam_Realm }}
              diam_Id_host: "{{ .Inventory_hostname }}.{{ .Diam_Realm }}"

        pcrf:
          hosts:
            pcrf1:
              ansible_host: {{ .PCRF_managementIP }}
              managementPort: {{ .PCRF_managementPort }}
              ansible_user: mos
              ansible_password: q 
              ansible_become_pass: q
              logger: "{{ .Var_path }}{{ .Inventory_hostname }}.log"
              freeDiameter: "{{ .Diameter_path }}{{ .Inventory_hostname }}.conf"
              db_uri: mongodb://localhost/bbdh

              gx_addr: {{ .PCRF_gx }}
              gx_port: {{ .PCRF_gxPort }}
              gx_secport: {{ .Gx_secPort }}
              # freeDiameter variables
              diam_realm: {{ .Diam_Realm }}
              diam_Id_host: "{{ .Inventory_hostname }}.{{ .Diam_Realm }}"
              diam_Port: 38888
              diam_SePort: 58162`

	templateTest := template.Must(template.New("yaml").Parse(yamlData))

	var buf bytes.Buffer

	if err := templateTest.Execute(&buf, data); err != nil {
		panic(err)
	}

	inventoryPath := "ansible-core-deploy/inventory/"
	os.WriteFile(inventoryPath + "Inventory.yml", buf.Bytes(), 0644)
	fmt.Println(color.Yellow + "\nGenerating Inventory.yml file ..." + color.Reset)
	time.Sleep(1 * time.Second)
	fmt.Print(color.Green + "Inventory.yml generated in the current path\n\n" + color.Reset)
}
