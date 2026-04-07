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
	sgwc_managementIP   string
	sgwc_managementPort string
	sgwc_gtpc           string
	sgwc_gtpcPort       int = 2123
	sgwc_pfcp           string
	sgwc_pfcpPort       int = 8805
	sgwc_gtpu           string
	sgwc_gtpuPort       int = 2152
	sgwu_managementIP   string
	sgwu_managementPort string
	sgwu_gtpu           string
	sgwu_gtpuPort       int = 2152
	sgwu_pfcp           string
	sgwu_pfcpPort       int = 8805
	upf_managemetIP     string
	upf_managementPort  string
	upf_pfcp            string
	upf_pfcpPort        int = 8805
	upf_gtpu            string
	upf_gtpuPort        int = 2152
	smf_managementIP    string
	smf_managementPort  string
	smf_gtpc            string
	smf_gtpcPort        int = 2123
	smf_pfcp            string
	smf_pfcpPort        int = 8805
	smf_gtpu            string
	smf_gtpuPort        int = 2152
	mme_gtpc            string
	mme_gtpcPort        int = 2123
	mme_s1ap            string
	mme_s1apPort        int = 36412
	mme_managementIP    string
	mme_managementPort  string
	hss_managementIP    string
	hss_managementPort  string
	pcrf_managementIP   string
	pcrf_managementPort string
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

		if len(vms[i].Networks) < 3 {
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
	mme_gtpc = vms[0].Networks[2].IP

	// Assigning SGW-C
	sgwc_managementIP = vms[3].Networks[0].IP
	sgwc_pfcp = vms[3].Networks[1].IP
	sgwc_gtpc = vms[3].Networks[2].IP
	sgwc_gtpu = vms[3].Networks[3].IP

	// Assigning SGW-U
	sgwu_managementIP = vms[4].Networks[0].IP
	sgwu_gtpu = vms[4].Networks[1].IP
	sgwu_pfcp = vms[4].Networks[2].IP

	// Assigning SMF
	smf_managementIP = vms[5].Networks[0].IP
	smf_gtpc = vms[5].Networks[1].IP
	smf_pfcp = vms[5].Networks[2].IP
	smf_gtpu = vms[5].Networks[3].IP

	// Assigning UPF
	upf_managemetIP = vms[6].Networks[0].IP
	upf_pfcp = vms[6].Networks[1].IP
	upf_gtpu = vms[6].Networks[2].IP

	// Assigning HSS
	hss_managementIP = vms[1].Networks[0].IP

	// Assigning PCRF
	pcrf_managementIP = vms[2].Networks[0].IP

	core_name = "{{ core_name }}"
	var_path = "/var/log/" + core_name + "/"
	tls_path = "/etc/" + core_name + "/tls/"
	inventory_hostname = "{{ inventory_hostname }}"
	var_path_diameter := "/etc/" + core_name + "/freeDiameter/"
	diam_realm = "epc.mnc0{{ plmn.mnc }}.mcc{{ plmn.mcc }}.3gppnetwork.org"

	data := struct {
		SGWC_managementIP  string
		SGWC_gtpc          string
		SGWC_gtpcPort      int
		SGWC_pfcp          string
		SGWC_pfcpPort      int
		MME_gtpc           string
		MME_gtpcPort       int
		MME_managementIP   string
		MME_s1ap           string
		MME_s1apPort       int
		SGWU_gtpu          string
		SGWU_gtpuPort      int
		SGWU_pfcp          string
		SGWU_pfcpPort      int
		SGWU_managementIP  string
		SMF_managementIP   string
		SMF_gtpc           string
		SMF_gtpcPort       int
		SMF_gtpu           string
		SMF_gtpuPort       int
		SMF_pfcp           string
		SMF_pfcpPort       int
		UPF_managementIP   string
		UPF_pfcp           string
		UPF_pfcpPort       int
		UPF_gtpu           string
		UPF_gtpuPort       int
		HSS_managementIP   string
		PCRF_managementIP  string
		Core_name          string
		Var_path           string
		Tls_path           string
		Inventory_hostname string
		Diameter_path      string
		Diam_Realm         string
	}{sgwc_managementIP,
		sgwc_gtpc,
		sgwc_gtpcPort,
		sgwc_pfcp,
		sgwc_pfcpPort,
		mme_gtpc,
		mme_gtpcPort,
		mme_managementIP,
		mme_s1ap,
		mme_s1apPort,
		sgwu_gtpu,
		sgwu_gtpuPort,
		sgwu_pfcp,
		sgwu_pfcpPort,
		sgwu_managementIP,
		smf_managementIP,
		smf_gtpc,
		smf_gtpcPort,
		smf_gtpu,
		smf_gtpuPort,
		smf_pfcp,
		smf_pfcpPort,
		upf_managemetIP,
		upf_pfcp,
		upf_pfcpPort,
		upf_gtpu,
		upf_gtpuPort,
		hss_managementIP,
		pcrf_managementIP,
		core_name,
		var_path,
		tls_path,
		inventory_hostname,
		var_path_diameter,
		diam_realm,
	}

	yamlData := `all:
  vars:
    core_name: bbdh
    db_uri: mongodb://localhost/{{ .Core_name }}
    configs_path: /etc/{{ .Core_name }}
    var_path: /var/log/{{ .Core_name }}/
    var_path_diameter: /etc/{{ .Core_name }}/freeDiameter/
    tls_path: /etc/{{ .Core_name }}/tls/ 

    # PLMN that use for most of the components
    plmn:
      mcc: 432
      mnc: 085

  children:
    sgwc:
      hosts:
        sgwc1:
          ansible_host: {{ .SGWC_managementIP }}
          ansible_user: mos
          ansible_password: q 
          ansible_become_pass: q
          logger: "{{ .Var_path }}{{ .Inventory_hostname }}.log"
          gtpc_addr: {{ .SGWC_gtpc }}
          pfcp_addr: {{ .SGWC_gtpc }}
          sgwu_pfcp: {{ .SGWC_pfcp }}

    sgwu:
      hosts:
        sgwu1:
          ansible_host: {{ .SGWU_managementIP }}
          ansible_user: mos
          ansible_password: q 
          ansible_become_pass: q
          logger: "{{ .Var_path }}{{ .Inventory_hostname }}.log"
          gtpu_addr: {{ .SGWU_gtpu }}
          pfcp_addr: {{ .SGWU_pfcp }}

    upf:
      hosts:
        upf1:
          ansible_host: {{ .UPF_managementIP }}
          ansible_user: mos
          ansible_password: q 
          ansible_become_pass: q
          logger: "{{ .Var_path }}{{ .Inventory_hostname }}.log"
          pfcp_addr: {{ .UPF_pfcp }}
          gtpu_addr: {{ .UPF_gtpu }}
          subnet:
              addr: 10.45.0.1/16
              dnn: internet
          smf_addr: {{ .SMF_managementIP }}

    # all diameter peers metagroup
    diam_peers:
      children:
        mme:
          hosts:
            mme1:
              ansible_host: {{ .MME_managementIP }}
              ansible_user: mos
              ansible_password: q 
              ansible_become_pass: q
              logger: "{{ .Var_path }}{{ .Inventory_hostname }}.log"
              freeDiameter: "{{ .Diameter_path }}{{ .Inventory_hostname }}.conf"
              tac: 3 
              gtpc_addr: {{ .MME_gtpc }}
              s1ap: {{ .MME_s1ap }}

              # freeDiameter variables
              diam_realm: {{ .Diam_Realm }}
              diam_Id_host: "{{ .Inventory_hostname }}.{{ .Diam_Realm }}"
              diam_tcp_port: 1111
              diam_tcpSec_port: 58162
              diam_listen_on: "{{ .SGWC_gtpc }}" # temporary

        hss:
          hosts:
            hss1:
              ansible_host: {{ .HSS_managementIP }}
              ansible_user: mos
              ansible_password: q 
              ansible_become_pass: q
              logger: "{{ .Var_path }}{{ .Inventory_hostname }}.log"
              freeDiameter: "{{ .Diameter_path }}{{ .Inventory_hostname }}.conf"
              db_uri: mongodb://localhost/bbdh

              # freeDiameter variables
              diam_realm: {{ .Diam_Realm }}
              diam_Id_host: "{{ .Inventory_hostname }}.{{ .Diam_Realm }}"
              diam_tcp_port: 38888
              diam_tcpSec_port: 58162
              diam_listen_on: "{{ .SGWC_gtpc }}" # temporary

        smf:
          hosts:
            smf1:
              ansible_host: {{ .SMF_managementIP }}
              ansible_user: mos
              ansible_password: q 
              ansible_become_pass: q
              logger: "{{ .Var_path }}{{ .Inventory_hostname }}.log"
              freeDiameter: "{{ .Diameter_path }}{{ .Inventory_hostname }}.conf"
              sbi_addr: 9877
              pfcp_addr: {{ .SMF_pfcp }}
              gtpc_addr: {{ .SMF_gtpc }}
              gtpu_addr: {{ .SMF_gtpu }}
              subnet:
                  addr: 10.45.0.1/16
                  dnn: internet
              dns:
                  primary: 8.8.8.8
                  secondary: 8.8.4.4
              upf_pfcp: {{ .UPF_managementIP }}

              # freeDiameter variables
              diam_realm: {{ .Diam_Realm }}
              diam_Id_host: "{{ .Inventory_hostname }}.{{ .Diam_Realm }}"
              diam_tcp_port: 38888
              diam_tcpSec_port: 58162
              diam_listen_on: "{{ .SGWC_gtpc }}" # temporary

        pcrf:
          hosts:
            pcrf1:
              ansible_host: {{ .PCRF_managementIP }}
              ansible_user: mos
              ansible_password: q 
              ansible_become_pass: q
              logger: "{{ .Var_path }}{{ .Inventory_hostname }}.log"
              freeDiameter: "{{ .Diameter_path }}{{ .Inventory_hostname }}.conf"
              db_uri: mongodb://localhost/bbdh

              # freeDiameter variables
              diam_realm: {{ .Diam_Realm }}
              diam_Id_host: "{{ .Inventory_hostname }}.{{ .Diam_Realm }}"
              diam_tcp_port: 38888
              diam_tcpSec_port: 58162
              diam_listen_on: "{{ .SGWC_gtpc }}" # temporary`

	templateTest := template.Must(template.New("yaml").Parse(yamlData))

	var buf bytes.Buffer

	if err := templateTest.Execute(&buf, data); err != nil {
		panic(err)
	}

	os.WriteFile("Inventory.yml", buf.Bytes(), 0644)
	fmt.Println(color.Yellow + "\nGenerating Inventory.yml file ..." + color.Reset)
	time.Sleep(1 * time.Second)
	fmt.Print(color.Green + "Inventory.yml generated in the current path\n\n" + color.Reset)
}
