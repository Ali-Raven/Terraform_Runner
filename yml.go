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
	sgwc_managementIP string = "null"
	sgwc_gtpc         string = "null"
	sgwc_pfcp         string
	sgwc_gtpu         string
	sgwu_gtpu         string
	sgwu_pfcp         string
	upf_pfcp          string
	upf_gtpu          string
	smf_gtpc          string
	smf_pfcp          string
	smf_gtpu          string
	mme_sgwc          string
	mme_s1ap          string

	core_name          string
	var_path           string
	tls_path           string
	inventory_hostname string
	diam_realm         string
	flagErr            bool
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
		// for j := 0 ; i <  len(vms[i].Networks) ; j++ {

		// }

		if len(vms[i].Networks) < 3 {
			fmt.Printf("%sError : not enough networks interface for %s%s\n",color.Red ,  vms[i].Name , color.Reset)
      time.Sleep(300 * time.Millisecond)
			flagErr = true
		}
	}
  if flagErr == true {
    return 
  }

	time.Sleep(1 * time.Second)

	sgwc_managementIP = vms[1].Networks[0].IP
	sgwc_pfcp = vms[1].Networks[1].IP
	sgwc_gtpc = vms[1].Networks[2].IP

	core_name = "{{ core_name }}"
	var_path = "/var/log/" + core_name + "/"
	tls_path = "/etc/" + core_name + "/tls/"
	inventory_hostname = "{{ inventory_hostname }}"
	var_path_diameter := "/etc/" + core_name + "/freeDiameter/"
	diam_realm = "epc.mnc0{{ plmn.mnc }}.mcc{{ plmn.mcc }}.3gppnetwork.org"

	// checking value of vms list

	// if additionVmManageIP == "" || additionVmIP1 == "" || additionVmIP2 == "" {
	// 	additionVmIP1 = "not assigned"
	// 	additionVmIP2 = "not assigned"
	// 	additionVmManageIP = "not assigned"
	// }

	data := struct {
		SGWC_managementIP  string
		SGWC_gtpc          string
		SGWC_pfcp          string
		Core_name          string
		Var_path           string
		Tls_path           string
		Inventory_hostname string
		Diameter_path      string
		Diam_Realm         string
	}{sgwc_managementIP, sgwc_gtpc, sgwc_pfcp, core_name, var_path, tls_path, inventory_hostname, var_path_diameter, diam_realm}

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
          ansible_host: 192.168.0.145
          ansible_user: mos
          ansible_password: q 
          ansible_become_pass: q
          logger: "{{ .Var_path }}{{ .Inventory_hostname }}.log"
          gtpu_addr: 192.168.0.1455
          pfcp_addr: 192.168.0.145

    upf:
      hosts:
        upf1:
          ansible_host: 192.168.0.147
          ansible_user: mos
          ansible_password: q 
          ansible_become_pass: q
          logger: "{{ .Var_path }}{{ .Inventory_hostname }}.log"
          pfcp_addr: 192.168.0.147
          gtpu_addr: 192.168.0.147
          subnet:
              addr: 10.45.0.1/16
              dnn: internet
          smf_addr: 192.168.0.146

    # all diameter peers metagroup
    diam_peers:
      children:
        mme:
          hosts:
            mme1:
              ansible_host: 192.168.0.141
              ansible_user: mos
              ansible_password: q 
              ansible_become_pass: q
              logger: "{{ .Var_path }}{{ .Inventory_hostname }}.log"
              freeDiameter: "{{ .Diameter_path }}{{ .Inventory_hostname }}.conf"
              tac: 3 
              gtpc_addr: 192.168.0.144
              s1ap: 192.168.0.155

              # freeDiameter variables
              diam_realm: {{ .Diam_Realm }}
              diam_Id_host: "{{ .Inventory_hostname }}.{{ .Diam_Realm }}"
              diam_tcp_port: 1111
              diam_tcpSec_port: 58162
              diam_listen_on: "{{ .SGWC_gtpc }}" # temporary

        hss:
          hosts:
            hss1:
              ansible_host: 192.168.0.142
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
              ansible_host: 192.168.0.146
              ansible_user: mos
              ansible_password: q 
              ansible_become_pass: q
              logger: "{{ .Var_path }}{{ .Inventory_hostname }}.log"
              freeDiameter: "{{ .Diameter_path }}{{ .Inventory_hostname }}.conf"
              sbi_addr: 9877
              pfcp_addr: 192.168.0.145
              gtpc_addr: 192.168.0.146
              gtpu_addr: 192.168.0.146
              subnet:
                  addr: 10.45.0.1/16
                  dnn: internet
              dns:
                  primary: 8.8.8.8
                  secondary: 8.8.4.4
              upf_pfcp: 192.168.0.147

              # freeDiameter variables
              diam_realm: {{ .Diam_Realm }}
              diam_Id_host: "{{ .Inventory_hostname }}.{{ .Diam_Realm }}"
              diam_tcp_port: 38888
              diam_tcpSec_port: 58162
              diam_listen_on: "{{ .SGWC_gtpc }}" # temporary

        pcrf:
          hosts:
            pcrf1:
              ansible_host: 192.168.0.143
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
