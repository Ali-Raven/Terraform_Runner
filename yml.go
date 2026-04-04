package main

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/TwiN/go-color"
)

func Yml(vms []VM) {
	// project_path := currentDir + "/ansible-core-deploy"
	// fmt.Println(vms[0].Networks[0].IP)

	additionVmManageIP := vms[8].Networks[0].IP
	additionVmIP1 := vms[8].Networks[1].IP
	additionVmIP2 := vms[8].Networks[2].IP
	core_name := "bbdh"
    var_path := "/var/log/" + core_name + "/"
    var_path_diameter := "/etc/" + core_name + "/freeDiameter/"
    tls_path := "/etc/" + core_name + "/tls/"

	data := struct {
		VmManageIP string
		VmIP1      string
		VmIP2      string
		Core_name  string
	}{additionVmManageIP, additionVmIP1, additionVmIP2 , core_name}

	yamlData := `all:
  vars:
    # core_name: bbdh
    # db_uri: mongodb://localhost/{{ .Core_name }}
    configs_path: /etc/{{ .Core_name }}
    var_path: /var/log/{{ .Core_name }}/
    var_path_diameter: /etc/{{ .Core_name }}/freeDiameter/
    tls_path: /etc/{{ .Core_name }}/tls/ 

    # PLMN that use for most of the components
    plmn:
      mcc: 432
      mnc: 85

  children:
    sgwc:
      hosts:
        sgwc1:
          ansible_host: {{ .VmManageIP }}
          ansible_user: mos
          ansible_password: q 
          ansible_become_pass: q
          logger: "{{ var_path }}{{ inventory_hostname }}.log"
          gtpc_addr: {{ .VmIP1 }}
          pfcp_addr: {{ .VmIP2 }}
          sgwu_pfcp: 192.168.0.145

    sgwu:
      hosts:
        sgwu1:
          ansible_host: 192.168.0.145
          ansible_user: mos
          ansible_password: q 
          ansible_become_pass: q
          logger: "{{ var_path }}{{ inventory_hostname }}.log"
          gtpu_addr: 192.168.0.1455
          pfcp_addr: 192.168.0.145

    upf:
      hosts:
        upf1:
          ansible_host: 192.168.0.147
          ansible_user: mos
          ansible_password: q 
          ansible_become_pass: q
          logger: "{{ var_path }}{{ inventory_hostname }}.log"
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
              logger: "{{ var_path }}{{ inventory_hostname }}.log"
              freeDiameter: "{{ var_path_diameter }}{{ inventory_hostname }}.conf"
              tac: 3 
              gtpc_addr: 192.168.0.144
              s1ap: 192.168.0.155

              # freeDiameter variables
              diam_realm: epc.mnc0{{ plmn.mnc }}.mcc{{ plmn.mcc }}.3gppnetwork.org
              diam_Id_host: "{{ inventory_hostname }}.{{ diam_realm }}"
              diam_tcp_port: 1111
              diam_tcpSec_port: 58162
              diam_listen_on: "{{ ansible_host }}" # temporary

        hss:
          hosts:
            hss1:
              ansible_host: 192.168.0.142
              ansible_user: mos
              ansible_password: q 
              ansible_become_pass: q
              logger: "{{ var_path }}{{ inventory_hostname }}.log"
              freeDiameter: "{{ var_path_diameter }}{{ inventory_hostname }}.conf"
              db_uri: mongodb://localhost/bbdh

              # freeDiameter variables
              diam_realm: epc.mnc0{{ plmn.mnc }}.mcc{{ plmn.mcc }}.3gppnetwork.org
              diam_Id_host: "{{ inventory_hostname }}.{{ diam_realm }}"
              diam_tcp_port: 38888
              diam_tcpSec_port: 58162
              diam_listen_on: "{{ ansible_host }}" # temporary

        smf:
          hosts:
            smf1:
              ansible_host: 192.168.0.146
              ansible_user: mos
              ansible_password: q 
              ansible_become_pass: q
              logger: "{{ var_path }}{{ inventory_hostname }}.log"
              freeDiameter: "{{ var_path_diameter }}{{ inventory_hostname }}.conf"
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
              diam_realm: epc.mnc0{{ plmn.mnc }}.mcc{{ plmn.mcc }}.3gppnetwork.org
              diam_Id_host: "{{ inventory_hostname }}.{{ diam_realm }}"
              diam_tcp_port: 38888
              diam_tcpSec_port: 58162
              diam_listen_on: "{{ ansible_host }}" # temporary

        pcrf:
          hosts:
            pcrf1:
              ansible_host: 192.168.0.143
              ansible_user: mos
              ansible_password: q 
              ansible_become_pass: q
              logger: "{{ var_path }}{{ inventory_hostname }}.log"
              freeDiameter: "{{ var_path_diameter }}{{ inventory_hostname }}.conf"
              db_uri: mongodb://localhost/bbdh

              # freeDiameter variables
              diam_realm: epc.mnc0{{ plmn.mnc }}.mcc{{ plmn.mcc }}.3gppnetwork.org
              diam_Id_host: "{{ inventory_hostname }}.{{ diam_realm }}"
              diam_tcp_port: 38888
              diam_tcpSec_port: 58162
              diam_listen_on: "{{ ansible_host }}" # temporary`

	templateTest := template.Must(template.New("yaml").Parse(yamlData))

	var buf bytes.Buffer

	if err := templateTest.Execute(&buf, data); err != nil {
		panic(err)
	}

	os.WriteFile("Inventory.yml", buf.Bytes(), 0644)
	fmt.Print(color.Green + "\nInventory.yml generated in the current path\n\n" + color.Reset)
}
