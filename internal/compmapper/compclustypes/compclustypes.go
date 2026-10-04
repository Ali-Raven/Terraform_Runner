// Package compclustypes define all components cluster types for inventory.yaml file in ansibe
package compclustypes

type HOST struct {
	ANSIBLEHost    string
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
	S6aSecPort int
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
	GxSecPort int
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
	S6aSecPort int
	S6aSubnet  int
	S6aGateway string

	Components []string
}
type PCRFCluster struct {
	Master HOST
	Backup HOST

	GxAddr    string
	GxPort    int
	GxSecPort int
	GxSubnet  int
	GxGateway string

	Components []string
}
type TemplateData struct {
	User               string
	VarPath            string
	CoreName           string
	ConfigPath         string
	BinPath            string
	CoreNameV          string
	V                  string
	CoreSourcePath     string
	TLSPath            string
	InventoryHostname  string
	FlagErr            bool
	DiamGroupNames     string
	NonDiamGroupNames  string
	UserHost           string
	VarPathDiameter    string
	DiamRealm          string
	HardcodedDiamRealm string

	MMEsCluster  []MMECluster
	HSSsCluster  []HSSCluster
	SGWCsCluster []SGWCCluster
	SGWUsCluster []SGWUCluster
	SMFsCluster  []SMFCluster
	UPFsCluster  []UPFCluster
	PCRFsCluster []PCRFCluster
}
