package compmapper

import (
	"regexp"
	"sort"
	"strconv"
	"strings"

	ct "github.com/terraform_runner/internal/compmapper/compclustypes"
	generators "github.com/terraform_runner/internal/generators"
)

var (
	SSHDefaultPort     int    = 5119
	sgwcManagementPort int    = SSHDefaultPort
	sgwcS11Port        int    = 2123
	sgwcSxaPort        int    = 8805
	sgwcS5cPort        int    = 2124
	sgwuSxaPort        int    = 2152
	sgwuS5uPort        int    = 8805
	sgwuS1uPort        int    = 3333
	upfSxbPort         int    = 8805
	upfSxuPort         int    = 2152
	upfS5uPort         int    = 2153
	smfGxPort          int    = 2123
	gxSecPort          int    = 5868
	smfS5cPort         int    = 8805
	smfSxuPort         int    = 8806
	mmeS11Port         int    = 2123
	mmeS1apPort        int    = 36412
	mmeS6aPort         int    = 2221
	s6aSecPort         int    = 5868
	mmeManagementPort  int    = SSHDefaultPort
	hssS6aPort         int    = 2223
	USER               string = "{{ USER }}"
	PASSWORD           string = "{{ PASSWORD }}"
	CoreName           string
	VarPath            string
	VarPathComps       string
	TLSPath            string
	ConfigPath         string
	BinPath            string
	CoreNameV          string
	V                  string = ""
	CoreSourcePath     string
	InventoryHostname  string
	DiamRealm          string
	FlagErr            bool
	DiamGroupNames     string
	NonDiamGroupNames  string
	VarPathDiameter    string
	HardcodedDiamRealm string
)

func extractNumber(name string) int {
	re := regexp.MustCompile(`\d+`)
	numStr := re.FindString(name)

	if numStr == "" {
		return 0
	}

	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0
	}

	return num
}

func UpfCluster(upfsByName map[string]generators.ComponentData) []ct.UPFCluster {
	var (
		upfNames    []string
		upfClusters []ct.UPFCluster
	)

	for name := range upfsByName {
		if !strings.Contains(name, "-backup") && strings.Contains(name, "UPF") {
			upfNames = append(upfNames, name)
		}
	}
	sort.Slice(upfNames, func(i, j int) bool {
		return extractNumber(upfNames[i]) < extractNumber(upfNames[j])
	})

	for _, name := range upfNames {
		if !strings.Contains(name, "-backup") && strings.Contains(name, "UPF") {

			vm := upfsByName[name]

			clusterName := name

			master := vm
			backupName := clusterName + "-backup"
			backup, ok := upfsByName[backupName]

			if !ok {
				continue
			}

			cluster := ct.UPFCluster{
				Master: ct.HOST{
					ANSIBLEHost:    master.Networks["OAM"].IP,
					ManagementPort: mmeManagementPort,
					User:           USER,
					Password:       PASSWORD,
					BecomePass:     PASSWORD,
					KeepalivedRole: "BACKUP",
					KeepalivedPrio: 150,
				},

				Backup: ct.HOST{
					ANSIBLEHost:    backup.Networks["OAM"].IP,
					ManagementPort: mmeManagementPort,
					User:           USER,
					Password:       PASSWORD,
					BecomePass:     PASSWORD,
					KeepalivedRole: "BACKUP",
					KeepalivedPrio: 100,
				},

				S5uAddr:    master.Networks["s5u-"+strings.ToLower(clusterName)].IP,
				S5uPort:    upfS5uPort,
				S5uSubnet:  master.Networks["s5u-"+strings.ToLower(clusterName)].Subnet,
				S5uGateway: master.Networks["s5u-"+strings.ToLower(clusterName)].Gateway,

				SxbAddr:    master.Networks["sxb-"+strings.ToLower(clusterName)].IP,
				SxbPort:    upfSxbPort,
				SxbSubnet:  master.Networks["sxb-"+strings.ToLower(clusterName)].Subnet,
				SxbGateway: master.Networks["sxb-"+strings.ToLower(clusterName)].Gateway,

				SxuAddr:    master.Networks["sxu-"+strings.ToLower(clusterName)].IP,
				SxuPort:    upfSxuPort,
				SxuSubnet:  master.Networks["sxu-"+strings.ToLower(clusterName)].Subnet,
				SxuGateway: master.Networks["sxu-"+strings.ToLower(clusterName)].Gateway,

				// SgiAddr:    master.Networks["sgi-"+strings.ToLower(clusterName)].IP,
				// SgiPort:    upf_sgiPort,
				// SgiSubnet:  master.Networks["sgi-"+strings.ToLower(clusterName)].Subnet,
				// SgiGateway: master.Networks["sgi-"+strings.ToLower(clusterName)].Gateway,
			}

			upfClusters = append(upfClusters, cluster)
		}
	}
	return upfClusters
}

func SgwcCluster(sgwcsByName map[string]generators.ComponentData) []ct.SGWCCluster {
	var (
		sgwcNames    []string
		sgwcClusters []ct.SGWCCluster
	)

	for name := range sgwcsByName {
		if !strings.Contains(name, "-backup") && strings.Contains(name, "SGWC") {
			sgwcNames = append(sgwcNames, name)
		}
	}
	sort.Slice(sgwcNames, func(i, j int) bool {
		return extractNumber(sgwcNames[i]) < extractNumber(sgwcNames[j])
	})

	for _, name := range sgwcNames {
		// only master nodes (not backup)
		if !strings.Contains(name, "-backup") && strings.HasPrefix(name, "SGWC") {

			vm := sgwcsByName[name]
			clusterName := name // SGWC1, SGWC2, ...

			master := vm
			backupName := clusterName + "-backup"
			backup, ok := sgwcsByName[backupName]

			if !ok {
				continue
			}

			for i, name := range vm.ComponentsToConnect {
				vm.ComponentsToConnect[i] = strings.ToLower(name)
			}

			cluster := ct.SGWCCluster{
				Master: ct.HOST{
					ANSIBLEHost:    master.Networks["OAM"].IP,
					ManagementPort: sgwcManagementPort,
					User:           USER,
					Password:       PASSWORD,
					BecomePass:     PASSWORD,
					KeepalivedRole: "BACKUP",
					KeepalivedPrio: 150,
				},

				Backup: ct.HOST{
					ANSIBLEHost:    backup.Networks["OAM"].IP,
					ManagementPort: sgwcManagementPort,
					User:           USER,
					Password:       PASSWORD,
					BecomePass:     PASSWORD,
					KeepalivedRole: "BACKUP",
					KeepalivedPrio: 100,
				},

				S11Addr:    master.Networks["s11-"+strings.ToLower(clusterName)].IP,
				S11Port:    sgwcS11Port,
				S11Subnet:  master.Networks["s11-"+strings.ToLower(clusterName)].Subnet,
				S11Gateway: master.Networks["s11-"+strings.ToLower(clusterName)].Gateway,

				S5cAddr:    master.Networks["s5c-"+strings.ToLower(clusterName)].IP,
				S5cPort:    sgwcS5cPort,
				S5cSubnet:  master.Networks["s5c-"+strings.ToLower(clusterName)].Subnet,
				S5cGateway: master.Networks["s5c-"+strings.ToLower(clusterName)].Gateway,

				SxaAddr:    master.Networks["sxa-"+strings.ToLower(clusterName)].IP,
				SxaPort:    sgwcSxaPort,
				SxaSubnet:  master.Networks["sxa-"+strings.ToLower(clusterName)].Subnet,
				SxaGateway: master.Networks["sxa-"+strings.ToLower(clusterName)].Gateway,
				Components: master.ComponentsToConnect,
			}

			sgwcClusters = append(sgwcClusters, cluster)
		}
	}

	return sgwcClusters
}

func SgwuCluster(sgwusByName map[string]generators.ComponentData) []ct.SGWUCluster {
	var (
		sgwuNames    []string
		sgwuClusters []ct.SGWUCluster
	)

	for name := range sgwusByName {
		if !strings.Contains(name, "-backup") && strings.Contains(name, "SGWU") {
			sgwuNames = append(sgwuNames, name)
		}
	}
	sort.Slice(sgwuNames, func(i, j int) bool {
		return extractNumber(sgwuNames[i]) < extractNumber(sgwuNames[j])
	})

	for _, name := range sgwuNames {
		// only master nodes (not backup)
		if !strings.Contains(name, "-backup") && strings.HasPrefix(name, "SGWU") {

			vm := sgwusByName[name]
			clusterName := name // SGWC1, SGWC2, ...

			master := vm
			backupName := clusterName + "-backup"
			backup, ok := sgwusByName[backupName]

			if !ok {
				continue
			}

			for i, name := range vm.ComponentsToConnect {
				vm.ComponentsToConnect[i] = strings.ToLower(name)
			}

			cluster := ct.SGWUCluster{
				Master: ct.HOST{
					ANSIBLEHost:    master.Networks["OAM"].IP,
					ManagementPort: sgwcManagementPort,
					User:           USER,
					Password:       PASSWORD,
					BecomePass:     PASSWORD,
					KeepalivedRole: "MASTER",
					KeepalivedPrio: 101,
				},

				Backup: ct.HOST{
					ANSIBLEHost:    backup.Networks["OAM"].IP,
					ManagementPort: sgwcManagementPort,
					User:           USER,
					Password:       PASSWORD,
					BecomePass:     PASSWORD,
					KeepalivedRole: "BACKUP",
					KeepalivedPrio: 100,
				},

				S1uAddr:    master.Networks["s1u-"+strings.ToLower(clusterName)].IP,
				S1uPort:    sgwuS1uPort,
				S1uSubnet:  master.Networks["s1u-"+strings.ToLower(clusterName)].Subnet,
				S1uGateway: master.Networks["s1u-"+strings.ToLower(clusterName)].Gateway,

				S5uAddr:    master.Networks["s5u-"+strings.ToLower(clusterName)].IP,
				S5uPort:    sgwuS5uPort,
				S5uSubnet:  master.Networks["s5u-"+strings.ToLower(clusterName)].Subnet,
				S5uGateway: master.Networks["s5u-"+strings.ToLower(clusterName)].Gateway,

				SxaAddr:    master.Networks["sxa-"+strings.ToLower(clusterName)].IP,
				SxaPort:    sgwuSxaPort,
				SxaSubnet:  master.Networks["sxa-"+strings.ToLower(clusterName)].Subnet,
				SxaGateway: master.Networks["sxa-"+strings.ToLower(clusterName)].Gateway,
				Components: master.ComponentsToConnect,
			}

			sgwuClusters = append(sgwuClusters, cluster)
		}
	}

	return sgwuClusters
}

func MMeCluster(mmesByName, smfsByName map[string]generators.ComponentData) []ct.MMECluster {
	var (
		mmeNames    []string
		mmeClusters []ct.MMECluster
	)
	for name := range mmesByName {
		if !strings.Contains(name, "-backup") && strings.HasPrefix(name, "MME") {
			mmeNames = append(mmeNames, name)
		}
	}

	sort.Slice(mmeNames, func(i, j int) bool {
		return extractNumber(mmeNames[i]) < extractNumber(mmeNames[j])
	})

	for _, name := range mmeNames {
		if !strings.Contains(name, "-backup") && strings.HasPrefix(name, "MME") {
			// whichSMF := whichSGWC(name)
			// fmt.Println(whichSMF , name)
			vm := mmesByName[name]
			clusterName := name

			master := vm
			backupName := clusterName + "-backup"
			backup, ok := mmesByName[backupName]

			if !ok {
				continue
			}

			for i, name := range vm.ComponentsToConnect {
				vm.ComponentsToConnect[i] = strings.ToLower(name)
			}
			cluster := ct.MMECluster{
				Master: ct.HOST{
					ANSIBLEHost:    master.Networks["OAM"].IP,
					ManagementPort: mmeManagementPort,
					User:           USER,
					Password:       PASSWORD,
					BecomePass:     PASSWORD,
					KeepalivedRole: "BACKUP",
					KeepalivedPrio: 150,
				},

				Backup: ct.HOST{
					ANSIBLEHost:    backup.Networks["OAM"].IP,
					ManagementPort: mmeManagementPort,
					User:           USER,
					Password:       PASSWORD,
					BecomePass:     PASSWORD,
					KeepalivedRole: "BACKUP",
					KeepalivedPrio: 100,
				},

				S1apAddr:    master.Networks["s1ap-"+strings.ToLower(clusterName)].IP,
				S1apPort:    mmeS1apPort,
				S1apSubnet:  master.Networks["s1ap-"+strings.ToLower(clusterName)].Subnet,
				S1apGateway: master.Networks["s1ap-"+strings.ToLower(clusterName)].Gateway,

				S11Addr:    master.Networks["s11-"+strings.ToLower(clusterName)].IP,
				S11Port:    mmeS11Port,
				S11Subnet:  master.Networks["s11-"+strings.ToLower(clusterName)].Subnet,
				S11Gateway: master.Networks["s11-"+strings.ToLower(clusterName)].Gateway,

				S5cAddr:    master.Networks["s5c-"+strings.ToLower(clusterName)].IP,
				S5cPort:    smfS5cPort,
				S5cSubnet:  master.Networks["s5c-"+strings.ToLower(clusterName)].Subnet,
				S5cGateway: master.Networks["s5c-"+strings.ToLower(clusterName)].Gateway,

				S6aAddr:    master.Networks["s6a-"+strings.ToLower(clusterName)].IP,
				S6aPort:    mmeS6aPort,
				S6aSecPort: s6aSecPort,
				S6aSubnet:  master.Networks["s6a-"+strings.ToLower(clusterName)].Subnet,
				S6aGateway: master.Networks["s6a-"+strings.ToLower(clusterName)].Gateway,
				Components: master.ComponentsToConnect,
			}

			mmeClusters = append(mmeClusters, cluster)

		}
	}
	return mmeClusters
}

func SMfCluster(smfsByName map[string]generators.ComponentData) []ct.SMFCluster {
	var (
		smfNames    []string
		smfClusters []ct.SMFCluster
	)
	for name := range smfsByName {
		if !strings.Contains(name, "-backup") && strings.HasPrefix(name, "SMF") {
			smfNames = append(smfNames, name)
		}
	}
	sort.Slice(smfNames, func(i, j int) bool {
		return extractNumber(smfNames[i]) < extractNumber(smfNames[j])
	})

	for _, name := range smfNames {
		if !strings.Contains(name, "-backup") && strings.HasPrefix(name, "SMF") {

			vm := smfsByName[name]
			clusterName := name

			master := vm
			backupName := clusterName + "-backup"
			backup, ok := smfsByName[backupName]

			if !ok {
				continue
			}

			for i, name := range vm.ComponentsToConnect {
				vm.ComponentsToConnect[i] = strings.ToLower(name)
			}
			cluster := ct.SMFCluster{
				Master: ct.HOST{
					ANSIBLEHost:    master.Networks["OAM"].IP,
					ManagementPort: mmeManagementPort,
					User:           USER,
					Password:       PASSWORD,
					BecomePass:     PASSWORD,
					KeepalivedRole: "BACKUP",
					KeepalivedPrio: 150,
				},

				Backup: ct.HOST{
					ANSIBLEHost:    backup.Networks["OAM"].IP,
					ManagementPort: mmeManagementPort,
					User:           USER,
					Password:       PASSWORD,
					BecomePass:     PASSWORD,
					KeepalivedRole: "BACKUP",
					KeepalivedPrio: 100,
				},

				SxbAddr:    master.Networks["sxb-"+strings.ToLower(clusterName)].IP,
				SxbPort:    smfS5cPort,
				SxbSubnet:  master.Networks["sxb-"+strings.ToLower(clusterName)].Subnet,
				SxbGateway: master.Networks["sxb-"+strings.ToLower(clusterName)].Gateway,

				SxuAddr:    master.Networks["sxu-"+strings.ToLower(clusterName)].IP,
				SxuPort:    smfSxuPort,
				SxuSubnet:  master.Networks["sxu-"+strings.ToLower(clusterName)].Subnet,
				SxuGateway: master.Networks["sxu-"+strings.ToLower(clusterName)].Gateway,

				S5cAddr:    master.Networks["s5c-"+strings.ToLower(clusterName)].IP,
				S5cPort:    smfS5cPort,
				S5cSubnet:  master.Networks["s5c-"+strings.ToLower(clusterName)].Subnet,
				S5cGateway: master.Networks["s5c-"+strings.ToLower(clusterName)].Gateway,

				GxAddr:     master.Networks["gx-"+strings.ToLower(clusterName)].IP,
				GxPort:     smfGxPort,
				GxSecPort:  gxSecPort,
				GxSubnet:   master.Networks["gx-"+strings.ToLower(clusterName)].Subnet,
				GxGateway:  master.Networks["gx-"+strings.ToLower(clusterName)].Gateway,
				Components: master.ComponentsToConnect,
			}

			smfClusters = append(smfClusters, cluster)
		}
	}
	return smfClusters
}

func HSsCluster(hsssByName map[string]generators.ComponentData) []ct.HSSCluster {
	var (
		hssNames    []string
		hssClusters []ct.HSSCluster
	)

	for name := range hsssByName {
		if !strings.Contains(name, "-backup") && strings.Contains(name, "HSS") {
			hssNames = append(hssNames, name)
		}
	}
	sort.Slice(hssNames, func(i, j int) bool {
		return extractNumber(hssNames[i]) < extractNumber(hssNames[j])
	})

	for _, name := range hssNames {
		if !strings.Contains(name, "-backup") && strings.Contains(name, "HSS") {

			vm := hsssByName[name]
			clusterName := name

			master := vm
			backupName := clusterName + "-backup"
			backup, ok := hsssByName[backupName]

			if !ok {
				continue
			}
			for i, name := range vm.ComponentsToConnect {
				vm.ComponentsToConnect[i] = strings.ToLower(name)
			}
			cluster := ct.HSSCluster{
				Master: ct.HOST{
					ANSIBLEHost:    master.Networks["OAM"].IP,
					ManagementPort: mmeManagementPort,
					User:           USER,
					Password:       PASSWORD,
					BecomePass:     PASSWORD,
					KeepalivedRole: "BACKUP",
					KeepalivedPrio: 150,
				},

				Backup: ct.HOST{
					ANSIBLEHost:    backup.Networks["OAM"].IP,
					ManagementPort: mmeManagementPort,
					User:           USER,
					Password:       PASSWORD,
					BecomePass:     PASSWORD,
					KeepalivedRole: "BACKUP",
					KeepalivedPrio: 100,
				},

				S6aAddr:    master.Networks["s6a-"+strings.ToLower(name)].IP,
				S6aPort:    hssS6aPort,
				S6aSecPort: s6aSecPort,
				S6aSubnet:  master.Networks["s6a-"+strings.ToLower(name)].Subnet,
				S6aGateway: master.Networks["s6a-"+strings.ToLower(name)].Gateway,
				Components: master.ComponentsToConnect,
			}

			hssClusters = append(hssClusters, cluster)
		}
	}
	return hssClusters
}

func PCRFsCluster(pcrfsByName map[string]generators.ComponentData) []ct.PCRFCluster {
	var (
		pcrfNames    []string
		pcrfClusters []ct.PCRFCluster
	)

	for name := range pcrfsByName {
		if !strings.Contains(name, "-backup") && strings.Contains(name, "PCRF") {
			pcrfNames = append(pcrfNames, name)
		}
	}
	sort.Slice(pcrfNames, func(i, j int) bool {
		return extractNumber(pcrfNames[i]) < extractNumber(pcrfNames[j])
	})

	for _, name := range pcrfNames {
		if !strings.Contains(name, "-backup") && strings.Contains(name, "PCRF") {
			vm := pcrfsByName[name]
			clusterName := name

			master := vm
			backupName := clusterName + "-backup"
			backup, ok := pcrfsByName[backupName]

			if !ok {
				continue
			}
			for i, name := range vm.ComponentsToConnect {
				vm.ComponentsToConnect[i] = strings.ToLower(name)
			}
			cluster := ct.PCRFCluster{
				Master: ct.HOST{
					ANSIBLEHost:    master.Networks["OAM"].IP,
					ManagementPort: mmeManagementPort,
					User:           USER,
					Password:       PASSWORD,
					BecomePass:     PASSWORD,
					KeepalivedRole: "BACKUP",
					KeepalivedPrio: 150,
				},

				Backup: ct.HOST{
					ANSIBLEHost:    backup.Networks["OAM"].IP,
					ManagementPort: mmeManagementPort,
					User:           USER,
					Password:       PASSWORD,
					BecomePass:     PASSWORD,
					KeepalivedRole: "BACKUP",
					KeepalivedPrio: 100,
				},

				GxAddr:     master.Networks["gx-"+strings.ToLower(name)].IP,
				GxPort:     hssS6aPort,
				GxSecPort:  gxSecPort,
				GxSubnet:   master.Networks["gx-"+strings.ToLower(name)].Subnet,
				GxGateway:  master.Networks["gx-"+strings.ToLower(name)].Gateway,
				Components: master.ComponentsToConnect,
			}

			pcrfClusters = append(pcrfClusters, cluster)
		}
	}
	return pcrfClusters
}

func Data(mmesByName, hsssByName, sgwcsByName, sgwusByName, smfsByName, upfsByName, pcrfsByName map[string]generators.ComponentData) any {
	sgwcCluster := SgwcCluster(sgwcsByName)
	mmeCluster := MMeCluster(mmesByName, smfsByName)
	smfCluster := SMfCluster(smfsByName)
	hssCluster := HSsCluster(hsssByName)
	pcrfCluster := PCRFsCluster(pcrfsByName)
	upfCluster := UpfCluster(upfsByName)
	sgwuCluster := SgwuCluster(sgwusByName)

	CoreName = "{{ core_name }}"
	if V == "" {
		CoreNameV = CoreName
	} else {
		CoreNameV = CoreName + "-" + V
	}
	VarPath = "/var/log/{{ core_name }}"
	TLSPath = "/etc/{{ core_name }}/tls/"
	ConfigPath = "/etc/{{ core_name }}"
	BinPath = "/usr/bin"
	CoreSourcePath = "/opt/{{ core_name_version }}"
	InventoryHostname = "{{ inventory_hostname }}"
	DiamGroupNames = "{{ group_names[1] }}"
	NonDiamGroupNames = "{{ group_names[0] }}"
	VarPathDiameter = "/etc/{{ core_name }}/freeDiameter/"
	DiamRealm = "epc.mnc{{ plmn.mnc }}.mcc{{ plmn.mcc }}.3gppnetwork.org"
	VarPathComps = "/var/log/{{ core_name }}/{{ group_names[0] }}.log"
	HardcodedDiamRealm = "{{ diam_realm }}"

	DataTempStruct := ct.TemplateData{
		User:               USER,
		CoreName:           CoreName,
		VarPath:            VarPath,
		ConfigPath:         ConfigPath,
		BinPath:            BinPath,
		CoreNameV:          CoreNameV,
		CoreSourcePath:     CoreSourcePath,
		V:                  V,
		TLSPath:            TLSPath,
		InventoryHostname:  InventoryHostname,
		DiamGroupNames:     DiamGroupNames,
		NonDiamGroupNames:  NonDiamGroupNames,
		VarPathDiameter:    VarPathDiameter,
		HardcodedDiamRealm: HardcodedDiamRealm,
		DiamRealm:          DiamRealm,
		SGWCsCluster:       sgwcCluster,
		MMEsCluster:        mmeCluster,
		SMFsCluster:        smfCluster,
		HSSsCluster:        hssCluster,
		PCRFsCluster:       pcrfCluster,
		UPFsCluster:        upfCluster,
		SGWUsCluster:       sgwuCluster,
	}
	return DataTempStruct
}
