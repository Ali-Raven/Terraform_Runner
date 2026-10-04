// Package template presenting the yaml template for ansible (structure)
package template

var yamlData string

// TemplateYML function returns the YAML template as a string. It is used to generate the Ansible inventory file based on the provided data.
func TemplateYML() string {
	yamlData = `all:
  vars:
    core_name: "bbdh"
    version: "{{ .V }}"
    db_uri: mongodb://localhost/{{ .CoreName }}
    configs_path: {{ .ConfigPath }}
    core_name_version: "{{ .CoreNameV }}"
    core_source_path: {{ .CoreSourcePath }}
    binaries_path: {{ .BinPath }}
    var_path: {{ .VarPath }}
    var_path_diameter: {{ .VarPathDiameter }}
    tls_path: {{ .TLSPath }}
    diam_lib_dir: /usr/lib
    # user: "{{ .User }}"
    core_file_url: "http://192.168.0.37/mirror/open5gs/github/bbdh-2.6.6.zip"
    core_libtins_url: "http://192.168.0.37/mirror/open5gs/github/libtins.zip"
    core_freeDiameter_url: "http://192.168.0.37/mirror/open5gs/github/freeDiameter.zip"
    core_prometheus_client_c_url: "http://192.168.0.37/mirror/open5gs/github/prometheus-client-c.zip"

    # PLMN that use for most of the components
    plmn:
      mcc: 432
      mnc: 080
    diam_realm: {{ .DiamRealm }}

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
              ansible_host: {{ $c.Master.ANSIBLEHost }}
              ansible_port: {{ $c.Master.ManagementPort }}
              ansible_user: "{{ $c.Master.User }}"
              ansible_password: "{{ $c.Master.Password }}"
              ansible_become_pass: "{{ $c.Master.BecomePass }}"
              keepalived_role: {{ $c.Master.KeepalivedRole }}
              keepalived_priority: {{ $c.Master.KeepalivedPrio }}

            sgwc{{ add $i 1 }}_backup:
              ansible_host: {{ $c.Backup.ANSIBLEHost }}
              ansible_port: {{ $c.Backup.ManagementPort }}
              ansible_user: "{{ $c.Backup.User }}"
              ansible_password: "{{ $c.Backup.Password }}"
              ansible_become_pass: "{{ $c.Backup.BecomePass }}"
              keepalived_role: {{ $c.Backup.KeepalivedRole }}
              keepalived_priority: {{ $c.Backup.KeepalivedPrio }}

          vars:
            logger: "{{ $.VarPath }}/{{ $.NonDiamGroupNames }}.log"
            sgwc_id: {{ add $i 1 }}
            router_id: {{ add 90 $i }}
            max_ue: 30000
            max_peer: 30000
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
              ansible_host: {{ $c.Master.ANSIBLEHost }}
              ansible_port: {{ $c.Master.ManagementPort }}
              ansible_user: "{{ $c.Master.User }}"
              ansible_password: "{{ $c.Master.Password }}"
              ansible_become_pass: "{{ $c.Master.BecomePass }}"
              keepalived_role: {{ $c.Master.KeepalivedRole }}
              keepalived_priority: {{ $c.Master.KeepalivedPrio }}

            sgwu{{ add $i 1 }}_backup:
              ansible_host: {{ $c.Backup.ANSIBLEHost }}
              ansible_port: {{ $c.Backup.ManagementPort }}
              ansible_user: "{{ $c.Backup.User }}"
              ansible_password: "{{ $c.Backup.Password }}"
              ansible_become_pass: "{{ $c.Backup.BecomePass }}"
              keepalived_role: {{ $c.Backup.KeepalivedRole }}
              keepalived_priority: {{ $c.Backup.KeepalivedPrio }}

          vars:
            logger: "{{ $.VarPath }}/{{ $.NonDiamGroupNames }}.log"
            router_id: {{ add 100 $i }}
            max_ue: 30000
            max_peer: 30000
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
                  ansible_host: {{ $c.Master.ANSIBLEHost }}
                  ansible_port: {{ $c.Master.ManagementPort }}
                  ansible_user: "{{ $c.Master.User }}"
                  ansible_password: "{{ $c.Master.Password }}"
                  ansible_become_pass: "{{ $c.Master.BecomePass }}"
                  keepalived_role: {{ $c.Master.KeepalivedRole }}
                  keepalived_priority: {{ $c.Master.KeepalivedPrio }}

                upf{{ add $i 1 }}_backup:
                  ansible_host: {{ $c.Backup.ANSIBLEHost }}
                  ansible_port: {{ $c.Backup.ManagementPort }}
                  ansible_user: "{{ $c.Backup.User }}"
                  ansible_password: "{{ $c.Backup.Password }}"
                  ansible_become_pass: "{{ $c.Backup.BecomePass }}"
                  keepalived_role: {{ $c.Backup.KeepalivedRole }}
                  keepalived_priority: {{ $c.Backup.KeepalivedPrio }}

              vars:
                logger: "{{ $.VarPath }}/{{ $.DiamGroupNames }}.log"
                freeDiameter: "{{ $.VarPathDiameter }}{{ $.DiamGroupNames }}.conf"
                router_id: {{ add 120 $i }}
                max_ue: 30000
                max_peer: 30000
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
                  ansible_host: {{ $c.Master.ANSIBLEHost }}
                  ansible_port: {{ $c.Master.ManagementPort }}
                  ansible_user: "{{ $c.Master.User }}"
                  ansible_password: "{{ $c.Master.Password }}"
                  ansible_become_pass: "{{ $c.Master.BecomePass }}"
                  keepalived_role: {{ $c.Master.KeepalivedRole }}
                  keepalived_priority: {{ $c.Master.KeepalivedPrio }}

                mme{{ add $i 1 }}_backup:
                  ansible_host: {{ $c.Backup.ANSIBLEHost }}
                  ansible_port: {{ $c.Backup.ManagementPort }}
                  ansible_user: "{{ $c.Backup.User }}"
                  ansible_password: "{{ $c.Backup.Password }}"
                  ansible_become_pass: "{{ $c.Backup.BecomePass }}"
                  keepalived_role: {{ $c.Backup.KeepalivedRole }}
                  keepalived_priority: {{ $c.Backup.KeepalivedPrio }}

              vars:
                diam_host: "{{ $.DiamGroupNames }}.{{ $.HardcodedDiamRealm }}"
                logger: "{{ $.VarPath }}/{{ $.DiamGroupNames }}.log"
                freeDiameter: "{{ $.VarPathDiameter }}{{ $.DiamGroupNames }}.conf"
                router_id: {{ add $i 1 }}
                max_ue: 30000
                max_peer: 30000
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
                s6a_secport: {{ $c.S6aSecPort }}
                s6a_gateway: {{ $c.S6aGateway }}
                s6a_subnet: {{ $c.S6aSubnet }}


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
                teid_range:
                  status: true
                  min: {{ sub_teid $i }}
                  max: {{ add_teid $i }}
             
                components:
                {{- range $c.Components }}
                  - {{ . }}
                {{- end }}

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
                  ansible_host: {{ $c.Master.ANSIBLEHost }}
                  ansible_port: {{ $c.Master.ManagementPort }}
                  ansible_user: "{{ $c.Master.User }}"
                  ansible_password: "{{ $c.Master.Password }}"
                  ansible_become_pass: "{{ $c.Master.BecomePass }}"
                  keepalived_role: {{ $c.Master.KeepalivedRole }}
                  keepalived_priority: {{ $c.Master.KeepalivedPrio }}

                hss{{ add $i 1 }}_backup:
                  ansible_host: {{ $c.Backup.ANSIBLEHost }}
                  ansible_port: {{ $c.Backup.ManagementPort }}
                  ansible_user: "{{ $c.Backup.User }}"
                  ansible_password: "{{ $c.Backup.Password }}"
                  ansible_become_pass: "{{ $c.Backup.BecomePass }}"
                  keepalived_role: {{ $c.Backup.KeepalivedRole }}
                  keepalived_priority: {{ $c.Backup.KeepalivedPrio }}

              vars:
                diam_host: "{{ $.DiamGroupNames }}.{{ $.HardcodedDiamRealm }}"
                logger: "{{ $.VarPath }}/{{ $.DiamGroupNames }}.log"
                freeDiameter: "{{ $.VarPathDiameter }}{{ $.DiamGroupNames }}.conf"
                router_id: {{ add 50 $i }}
                hss_id: {{ add $i 1 }}
                max_ue: 30000
                max_peer: 30000
                s6a_addr: {{ $c.S6aAddr }}/{{ $c.S6aSubnet }}
                s6a_port: {{ $c.S6aPort }}
                s6a_secport: {{ $c.S6aSecPort }}
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
                  ansible_host: {{ $c.Master.ANSIBLEHost }}
                  ansible_port: {{ $c.Master.ManagementPort }}
                  ansible_user: "{{ $c.Master.User }}"
                  ansible_password: "{{ $c.Master.Password }}"
                  ansible_become_pass: "{{ $c.Master.BecomePass }}"
                  keepalived_role: {{ $c.Master.KeepalivedRole }}
                  keepalived_priority: {{ $c.Master.KeepalivedPrio }}

                smf{{ add $i 1 }}_backup:
                  ansible_host: {{ $c.Backup.ANSIBLEHost }}
                  ansible_port: {{ $c.Backup.ManagementPort }}
                  ansible_user: "{{ $c.Backup.User }}"
                  ansible_password: "{{ $c.Backup.Password }}"
                  ansible_become_pass: "{{ $c.Backup.BecomePass }}"
                  keepalived_role: {{ $c.Backup.KeepalivedRole }}
                  keepalived_priority: {{ $c.Backup.KeepalivedPrio }}

              vars:
                diam_host: "{{ $.DiamGroupNames }}.{{ $.HardcodedDiamRealm }}"
                logger: "{{ $.VarPath }}/{{ $.DiamGroupNames }}.log"
                freeDiameter: "{{ $.VarPathDiameter }}{{ $.DiamGroupNames }}.conf"
                smf_id: {{ add $i 1 }}
                router_id: {{ add 70 $i }}
                max_ue: 30000
                max_peer: 30000
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
                gx_secport: {{ $c.GxSecPort }}
                gx_gateway: {{ $c.GxGateway }}
                gx_subnet: {{ $c.GxSubnet }}

                subnets:
                  - subnet: 10.45.0.0/16
                    gateway: 10.45.0.1
                    apn: internet
                  - subnet: 2001:db8:cafe::/48
                    gateway: 2001:db8:cafe::1
                    apn: internetv6

                dnss:
                  - 8.8.8.8
                  - 8.8.4.4
                  - 2001:4860:4860::8888
                  - 2001:4860:4860::8844
                
                p_cscf:
                  - 10.60.0.20

                ctf:
                  enabled: no

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
                  ansible_host: {{ $c.Master.ANSIBLEHost }}
                  ansible_port: {{ $c.Master.ManagementPort }}
                  ansible_user: "{{ $c.Master.User }}"
                  ansible_password: "{{ $c.Master.Password }}"
                  ansible_become_pass: "{{ $c.Master.BecomePass }}"
                  keepalived_role: {{ $c.Master.KeepalivedRole }}
                  keepalived_priority: {{ $c.Master.KeepalivedPrio }}

                pcrf{{ add $i 1 }}_backup:
                  ansible_host: {{ $c.Backup.ANSIBLEHost }}
                  ansible_port: {{ $c.Backup.ManagementPort }}
                  ansible_user: "{{ $c.Backup.User }}"
                  ansible_password: "{{ $c.Backup.Password }}"
                  ansible_become_pass: "{{ $c.Backup.BecomePass }}"
                  keepalived_role: {{ $c.Backup.KeepalivedRole }}
                  keepalived_priority: {{ $c.Backup.KeepalivedPrio }}

              vars:
                diam_host: "{{ $.DiamGroupNames }}.{{ $.HardcodedDiamRealm }}"
                logger: "{{ $.VarPath }}/{{ $.DiamGroupNames }}.log"
                freeDiameter: "{{ $.VarPathDiameter }}{{ $.DiamGroupNames }}.conf"
                router_id: {{ add 80 $i }}
                pcrf_id: {{ add $i 1 }}
                max_ue: 30000
                max_peer: 30000
                gx_addr: {{ $c.GxAddr }}/{{ $c.GxSubnet }}
                gx_port: {{ $c.GxPort }}
                gx_secport: {{ $c.GxSecPort }}
                gx_gateway: {{ $c.GxGateway }}
                gx_subnet: {{ $c.GxSubnet }}
                components:
                {{- range $c.Components }}
                  - {{ . }}
                {{- end }}

            {{- end }}
               `
	return yamlData
}
