package main

import (
	"regexp"
	"sort"
	"strconv"
	"strings"

	generators "github.com/terraform_runner/Generators"
)

var (
	ssh_defaultPort     int = 5119
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
	upf_sxuPort         int = 2152
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
	Core_Name           string
	Var_path            string
	Var_path_Comps      string
	Tls_path            string
	Inventory_hostname  string
	Diam_Realm          string
	FlagErr             bool
	Diam_groupNames     string
	Non_diam_groupNames string
	UserHost            string = "mos"
	Var_path_diameter   string
)

type HOST struct {
	ANSIBLE_HOST   string
	ManagementPort int
	User           string
	Password       string
	BecomePass     string
	KeepalivedRole string
	KeepalivedPrio int
}
type SGWCCluster struct {
	Master HOST
	Backup HOST

	S11Addr    string
	S11Port    int
	S11Subnet  int
	S11Gateway string

	S5cAddr    string
	S5cPort    int
	S5cSubnet  int
	S5cGateway string

	SxaAddr    string
	SxaPort    int
	SxaSubnet  int
	SxaGateway string

	Components []string
}
type SGWUCluster struct {
	Master HOST
	Backup HOST

	S1uAddr    string
	S1uPort    int
	S1uSubnet  int
	S1uGateway string

	S5uAddr    string
	S5uPort    int
	S5uSubnet  int
	S5uGateway string

	SxaAddr    string
	SxaPort    int
	SxaSubnet  int
	SxaGateway string

	Components []string
}

type MMECluster struct {
	Master HOST
	Backup HOST

	RouterID    int
	S1apAddr    string
	S1apPort    int
	S1apSubnet  int
	S1apGateway string

	S11Addr    string
	S11Port    int
	S11Subnet  int
	S11Gateway string

	S5cAddr    string
	S5cPort    int
	S5cSubnet  int
	S5cGateway string

	S6aAddr    string
	S6aPort    int
	S6aSubnet  int
	S6aGateway string

	Components []string
}
type SMFCluster struct {
	Master HOST
	Backup HOST

	S5cAddr    string
	S5cPort    int
	S5cSubnet  int
	S5cGateway string

	SxbAddr    string
	SxbPort    int
	SxbSubnet  int
	SxbGateway string

	SxuAddr    string
	SxuPort    int
	SxuSubnet  int
	SxuGateway string

	GxAddr    string
	GxPort    int
	GxSubnet  int
	GxGateway string

	Components []string
}
type UPFCluster struct {
	Master HOST
	Backup HOST

	S5uAddr    string
	S5uPort    int
	S5uSubnet  int
	S5uGateway string

	SxbAddr    string
	SxbPort    int
	SxbSubnet  int
	SxbGateway string

	SxuAddr    string
	SxuPort    int
	SxuSubnet  int
	SxuGateway string

	SgiAddr    string
	SgiPort    int
	SgiSubnet  int
	SgiGateway string

	Components []string
}
type HSSCluster struct {
	Master HOST
	Backup HOST

	S6aAddr    string
	S6aPort    int
	S6aSubnet  int
	S6aGateway string

	Components []string
}
type PCRFCluster struct {
	Master HOST
	Backup HOST

	GxAddr    string
	GxPort    int
	GxSubnet  int
	GxGateway string

	Components []string
}
type TemplateData struct {
	User                string
	Var_path            string
	Core_Name           string
	Tls_path            string
	Inventory_hostname  string
	Diam_Realm          string
	FlagErr             bool
	Diam_groupNames     string
	Non_diam_groupNames string
	UserHost            string
	Var_path_diameter   string

	MMEsCluster  []MMECluster
	HSSsCluster  []HSSCluster
	SGWCsCluster []SGWCCluster
	SGWUsCluster []SGWUCluster
	SMFsCluster  []SMFCluster
	UPFsCluster  []UPFCluster
	PCRFsCluster []PCRFCluster
}

// func whichSGWC(name string) string {
// 	numStr := strings.TrimPrefix(name, "MME")
// 	num, err := strconv.Atoi(numStr)
// 	if err != nil {
// 		return ""
// 	}

// 	switch {
// 	case num >= 1 && num <= 7:
// 		return "SMF2"
// 	case num >= 8 && num <= 13:
// 		return "SMF1"
// 	default:
// 		return ""
// 	}
// }

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
func UpfCluster(upfsByName map[string]generators.ComponentData) []UPFCluster {
	var (
		upfNames    []string
		upfClusters []UPFCluster
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

			cluster := UPFCluster{
				Master: HOST{
					ANSIBLE_HOST:   master.Networks["Core-DevOps-3"].IP,
					ManagementPort: mme_managementPort,
					User:           UserHost,
					Password:       "q",
					BecomePass:     "q",
					KeepalivedRole: "BACKUP",
					KeepalivedPrio: 101,
				},

				Backup: HOST{
					ANSIBLE_HOST:   backup.Networks["Core-DevOps-3"].IP,
					ManagementPort: mme_managementPort,
					User:           UserHost,
					Password:       "q",
					BecomePass:     "q",
					KeepalivedRole: "BACKUP",
					KeepalivedPrio: 100,
				},

				S5uAddr:    master.Networks["s5u-"+strings.ToLower(clusterName)].IP,
				S5uPort:    upf_s5uPort,
				S5uSubnet:  master.Networks["s5u-"+strings.ToLower(clusterName)].Subnet,
				S5uGateway: master.Networks["s5u-"+strings.ToLower(clusterName)].Gateway,

				SxbAddr:    master.Networks["sxb-"+strings.ToLower(clusterName)].IP,
				SxbPort:    upf_sxbPort,
				SxbSubnet:  master.Networks["sxb-"+strings.ToLower(clusterName)].Subnet,
				SxbGateway: master.Networks["sxb-"+strings.ToLower(clusterName)].Gateway,

				SxuAddr:    master.Networks["sxu-"+strings.ToLower(clusterName)].IP,
				SxuPort:    upf_sxuPort,
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
func SgwcCluster(sgwcsByName map[string]generators.ComponentData) []SGWCCluster {
	var (
		sgwcNames    []string
		sgwcClusters []SGWCCluster
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

			cluster := SGWCCluster{
				Master: HOST{
					ANSIBLE_HOST:   master.Networks["Core-DevOps-3"].IP,
					ManagementPort: sgwc_managementPort,
					User:           UserHost,
					Password:       "q",
					BecomePass:     "q",
					KeepalivedRole: "MASTER",
					KeepalivedPrio: 101,
				},

				Backup: HOST{
					ANSIBLE_HOST:   backup.Networks["Core-DevOps-3"].IP,
					ManagementPort: sgwc_managementPort,
					User:           UserHost,
					Password:       "q",
					BecomePass:     "q",
					KeepalivedRole: "BACKUP",
					KeepalivedPrio: 100,
				},

				S11Addr:    master.Networks["s11-"+strings.ToLower(clusterName)].IP,
				S11Port:    sgwc_s11Port,
				S11Subnet:  master.Networks["s11-"+strings.ToLower(clusterName)].Subnet,
				S11Gateway: master.Networks["s11-"+strings.ToLower(clusterName)].Gateway,

				S5cAddr:    master.Networks["s5c-"+strings.ToLower(clusterName)].IP,
				S5cPort:    sgwc_s5cPort,
				S5cSubnet:  master.Networks["s5c-"+strings.ToLower(clusterName)].Subnet,
				S5cGateway: master.Networks["s5c-"+strings.ToLower(clusterName)].Gateway,

				SxaAddr:    master.Networks["sxa-"+strings.ToLower(clusterName)].IP,
				SxaPort:    sgwc_sxaPort,
				SxaSubnet:  master.Networks["sxa-"+strings.ToLower(clusterName)].Subnet,
				SxaGateway: master.Networks["sxa-"+strings.ToLower(clusterName)].Gateway,
				Components: master.ComponentsToConnect,
			}

			sgwcClusters = append(sgwcClusters, cluster)
			// fmt.Println(sgwcClusters)
			// time.Sleep(100000 * time.Second)
		}
	}

	return sgwcClusters
}
func SgwuCluster(sgwusByName map[string]generators.ComponentData) []SGWUCluster {
	var (
		sgwuNames    []string
		sgwuClusters []SGWUCluster
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

			cluster := SGWUCluster{
				Master: HOST{
					ANSIBLE_HOST:   master.Networks["Core-DevOps-3"].IP,
					ManagementPort: sgwc_managementPort,
					User:           UserHost,
					Password:       "q",
					BecomePass:     "q",
					KeepalivedRole: "MASTER",
					KeepalivedPrio: 101,
				},

				Backup: HOST{
					ANSIBLE_HOST:   backup.Networks["Core-DevOps-3"].IP,
					ManagementPort: sgwc_managementPort,
					User:           UserHost,
					Password:       "q",
					BecomePass:     "q",
					KeepalivedRole: "BACKUP",
					KeepalivedPrio: 100,
				},

				S1uAddr:    master.Networks["s1u-"+strings.ToLower(clusterName)].IP,
				S1uPort:    sgwu_s1uPort,
				S1uSubnet:  master.Networks["s1u-"+strings.ToLower(clusterName)].Subnet,
				S1uGateway: master.Networks["s1u-"+strings.ToLower(clusterName)].Gateway,

				S5uAddr:    master.Networks["s5u-"+strings.ToLower(clusterName)].IP,
				S5uPort:    sgwu_s5uPort,
				S5uSubnet:  master.Networks["s5u-"+strings.ToLower(clusterName)].Subnet,
				S5uGateway: master.Networks["s5u-"+strings.ToLower(clusterName)].Gateway,

				SxaAddr:    master.Networks["sxa-"+strings.ToLower(clusterName)].IP,
				SxaPort:    sgwu_sxaPort,
				SxaSubnet:  master.Networks["sxa-"+strings.ToLower(clusterName)].Subnet,
				SxaGateway: master.Networks["sxa-"+strings.ToLower(clusterName)].Gateway,
				Components: master.ComponentsToConnect,
			}

			sgwuClusters = append(sgwuClusters, cluster)
			// fmt.Println(sgwuClusters)
			// time.Sleep(100000 * time.Second)
		}
	}

	return sgwuClusters
}
func MMeCluster(mmesByName, smfsByName map[string]generators.ComponentData) []MMECluster {

	var (
		mmeNames    []string
		mmeClusters []MMECluster
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
			cluster := MMECluster{
				Master: HOST{
					ANSIBLE_HOST:   master.Networks["Core-DevOps-3"].IP,
					ManagementPort: mme_managementPort,
					User:           UserHost,
					Password:       "q",
					BecomePass:     "q",
					KeepalivedRole: "BACKUP",
					KeepalivedPrio: 101,
				},

				Backup: HOST{
					ANSIBLE_HOST:   backup.Networks["Core-DevOps-3"].IP,
					ManagementPort: mme_managementPort,
					User:           UserHost,
					Password:       "q",
					BecomePass:     "q",
					KeepalivedRole: "BACKUP",
					KeepalivedPrio: 100,
				},

				S1apAddr:    master.Networks["s1ap-"+strings.ToLower(clusterName)].IP,
				S1apPort:    mme_s1apPort,
				S1apSubnet:  master.Networks["s1ap-"+strings.ToLower(clusterName)].Subnet,
				S1apGateway: master.Networks["s1ap-"+strings.ToLower(clusterName)].Gateway,

				S11Addr:    master.Networks["s11-"+strings.ToLower(clusterName)].IP,
				S11Port:    mme_s11Port,
				S11Subnet:  master.Networks["s11-"+strings.ToLower(clusterName)].Subnet,
				S11Gateway: master.Networks["s11-"+strings.ToLower(clusterName)].Gateway,

				S5cAddr:    master.Networks["s5c-"+strings.ToLower(clusterName)].IP,
				S5cPort:    smf_s5cPort,
				S5cSubnet:  master.Networks["s5c-"+strings.ToLower(clusterName)].Subnet,
				S5cGateway: master.Networks["s5c-"+strings.ToLower(clusterName)].Gateway,

				S6aAddr:    master.Networks["s6a-"+strings.ToLower(clusterName)].IP,
				S6aPort:    mme_s6aPort,
				S6aSubnet:  master.Networks["s6a-"+strings.ToLower(clusterName)].Subnet,
				S6aGateway: master.Networks["s6a-"+strings.ToLower(clusterName)].Gateway,
				Components: master.ComponentsToConnect,
			}

			mmeClusters = append(mmeClusters, cluster)

		}
	}
	return mmeClusters
}
func SMfCluster(smfsByName map[string]generators.ComponentData) []SMFCluster {
	var (
		smfNames    []string
		smfClusters []SMFCluster
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
			cluster := SMFCluster{
				Master: HOST{
					ANSIBLE_HOST:   master.Networks["Core-DevOps-3"].IP,
					ManagementPort: mme_managementPort,
					User:           UserHost,
					Password:       "q",
					BecomePass:     "q",
					KeepalivedRole: "BACKUP",
					KeepalivedPrio: 101,
				},

				Backup: HOST{
					ANSIBLE_HOST:   backup.Networks["Core-DevOps-3"].IP,
					ManagementPort: mme_managementPort,
					User:           UserHost,
					Password:       "q",
					BecomePass:     "q",
					KeepalivedRole: "BACKUP",
					KeepalivedPrio: 100,
				},

				SxbAddr:    master.Networks["sxb-"+strings.ToLower(clusterName)].IP,
				SxbPort:    smf_s5cPort,
				SxbSubnet:  master.Networks["sxb-"+strings.ToLower(clusterName)].Subnet,
				SxbGateway: master.Networks["sxb-"+strings.ToLower(clusterName)].Gateway,

				SxuAddr:    master.Networks["sxu-"+strings.ToLower(clusterName)].IP,
				SxuPort:    smf_sxuPort,
				SxuSubnet:  master.Networks["sxu-"+strings.ToLower(clusterName)].Subnet,
				SxuGateway: master.Networks["sxu-"+strings.ToLower(clusterName)].Gateway,

				S5cAddr:    master.Networks["s5c-"+strings.ToLower(clusterName)].IP,
				S5cPort:    smf_s5cPort,
				S5cSubnet:  master.Networks["s5c-"+strings.ToLower(clusterName)].Subnet,
				S5cGateway: master.Networks["s5c-"+strings.ToLower(clusterName)].Gateway,

				GxAddr:     master.Networks["gx-"+strings.ToLower(clusterName)].IP,
				GxPort:     smf_gxPort,
				GxSubnet:   master.Networks["gx-"+strings.ToLower(clusterName)].Subnet,
				GxGateway:  master.Networks["gx-"+strings.ToLower(clusterName)].Gateway,
				Components: master.ComponentsToConnect,
			}

			smfClusters = append(smfClusters, cluster)
			// fmt.Println(sgwcClusters)
			// time.Sleep(100000 * time.Second)
		}
	}
	// fmt.Println(smfClusters)
	// time.Sleep(10000 * time.Second)
	return smfClusters
}
func HSsCluster(hsssByName map[string]generators.ComponentData) []HSSCluster {
	var (
		hssNames    []string
		hssClusters []HSSCluster
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
			cluster := HSSCluster{
				Master: HOST{
					ANSIBLE_HOST:   master.Networks["Core-DevOps-3"].IP,
					ManagementPort: mme_managementPort,
					User:           UserHost,
					Password:       "q",
					BecomePass:     "q",
					KeepalivedRole: "BACKUP",
					KeepalivedPrio: 101,
				},

				Backup: HOST{
					ANSIBLE_HOST:   backup.Networks["Core-DevOps-3"].IP,
					ManagementPort: mme_managementPort,
					User:           UserHost,
					Password:       "q",
					BecomePass:     "q",
					KeepalivedRole: "BACKUP",
					KeepalivedPrio: 100,
				},

				S6aAddr:    master.Networks["s6a-"+strings.ToLower(name)].IP,
				S6aPort:    hss_s6aPort,
				S6aSubnet:  master.Networks["s6a-"+strings.ToLower(name)].Subnet,
				S6aGateway: master.Networks["s6a-"+strings.ToLower(name)].Gateway,
				Components: master.ComponentsToConnect,
			}

			hssClusters = append(hssClusters, cluster)
		}
	}
	return hssClusters
}
func PCRFsCluster(pcrfsByName map[string]generators.ComponentData) []PCRFCluster {
	var (
		pcrfNames    []string
		pcrfClusters []PCRFCluster
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
			cluster := PCRFCluster{
				Master: HOST{
					ANSIBLE_HOST:   master.Networks["Core-DevOps-3"].IP,
					ManagementPort: mme_managementPort,
					User:           UserHost,
					Password:       "q",
					BecomePass:     "q",
					KeepalivedRole: "BACKUP",
					KeepalivedPrio: 101,
				},

				Backup: HOST{
					ANSIBLE_HOST:   backup.Networks["Core-DevOps-3"].IP,
					ManagementPort: mme_managementPort,
					User:           UserHost,
					Password:       "q",
					BecomePass:     "q",
					KeepalivedRole: "BACKUP",
					KeepalivedPrio: 100,
				},

				GxAddr:     master.Networks["gx-"+strings.ToLower(name)].IP,
				GxPort:     hss_s6aPort,
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

	DataTempStruct := TemplateData{
		User:                UserHost,
		Core_Name:           Core_Name,
		Var_path:            Var_path,
		Tls_path:            Tls_path,
		Inventory_hostname:  Inventory_hostname,
		Diam_Realm:          Diam_Realm,
		Diam_groupNames:     Diam_groupNames,
		Non_diam_groupNames: Non_diam_groupNames,
		Var_path_diameter:   Var_path_diameter,
		SGWCsCluster:        sgwcCluster,
		MMEsCluster:         mmeCluster,
		SMFsCluster:         smfCluster,
		HSSsCluster:         hssCluster,
		PCRFsCluster:        pcrfCluster,
		UPFsCluster:         upfCluster,
		SGWUsCluster:        sgwuCluster,
	}

	return DataTempStruct
}
