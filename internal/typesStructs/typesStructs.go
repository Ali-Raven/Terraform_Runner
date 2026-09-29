package typesstructs

type Network struct {
	ID      string `json:"uuid"`
	Name    string `json:"name"`
	IP      string `json:"ip"`
	Gateway string `json:"gateway"`
	Netmask int    `json:"netmask"`
}

type VM struct {
	ID         string `json:"uuid"`
	Name       string `json:"name"`
	TargetHost string `json:"target_host"`
	DataStore  string `json:"datastore"`
	// ClusterName         string    `json:"cluster_name"`
	NumCPU              int       `json:"num_cpus"`
	MemoryGB            int       `json:"memory_gb"`
	Gateway             string    `json:"gateway"`
	DNSservers          []string  `json:"dns_servers"`
	Component           string    `json:"component"`
	ComponentsToConnect []string  `json:"componentsToConnect"`
	Networks            []Network `json:"network_adaptors"`
}

type TFvars struct {
	VMs []VM `json:"vms"`
}
