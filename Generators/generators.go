package generators

type ComponentData struct {
	Name                string
	Networks            map[string]NetworksStructure
	ComponentsToConnect []string
}

type NetworksStructure struct {
	ID      string
	Name    string
	IP      string
	Gateway string
	Subnet  int
}

func BuildAllMMEs(componentData map[string]ComponentData) map[string]ComponentData {
	// var mmes []ComponentData
	mmesByName := make(map[string]ComponentData)

	for name, nets := range componentData {
		mme := ComponentData{
			Name:                name,
			Networks:            nets.Networks,
			ComponentsToConnect: nets.ComponentsToConnect,
		}
		// mmes = append(mmes, mme)
		mmesByName[name] = mme
	}
	return mmesByName
}

func BuildAllHSSs(componentData map[string]ComponentData) map[string]ComponentData {
	hsssByName := make(map[string]ComponentData)

	for name, nets := range componentData {
		hss := ComponentData{
			Name:                name,
			Networks:            nets.Networks,
			ComponentsToConnect: nets.ComponentsToConnect,
		}
		// hsss = append(hsss, hss)
		hsssByName[name] = hss
	}
	return hsssByName
}
func BuildAllSGWCs(componentData map[string]ComponentData) map[string]ComponentData {
	sgwcsByName := make(map[string]ComponentData)

	for name, nets := range componentData {
		sgwc := ComponentData{
			Name:                name,
			Networks:            nets.Networks,
			ComponentsToConnect: nets.ComponentsToConnect,
		}
		// sgwcs = append(sgwcs, sgwc)
		sgwcsByName[name] = sgwc
	}
	return sgwcsByName
}
func BuildAllSGWUs(componentData map[string]ComponentData) map[string]ComponentData {
	sgwusByName := make(map[string]ComponentData)

	for name, nets := range componentData {
		sgwu := ComponentData{
			Name:                name,
			Networks:            nets.Networks,
			ComponentsToConnect: nets.ComponentsToConnect,
		}
		// sgwus = append(sgwus, sgwu)
		sgwusByName[name] = sgwu
	}
	return sgwusByName
}
func BuildAllSMFs(componentData map[string]ComponentData) map[string]ComponentData {
	smfsByName := make(map[string]ComponentData)

	for name, nets := range componentData {
		smf := ComponentData{
			Name:                name,
			Networks:            nets.Networks,
			ComponentsToConnect: nets.ComponentsToConnect,
		}
		// smfs = append(smfs, smf)
		smfsByName[name] = smf
	}
	return smfsByName
}
func BuildAllUPFs(componentData map[string]ComponentData) map[string]ComponentData {
	upfsByName := make(map[string]ComponentData)

	for name, nets := range componentData {
		upf := ComponentData{
			Name:                name,
			Networks:            nets.Networks,
			ComponentsToConnect: nets.ComponentsToConnect,
		}
		// upfs = append(upfs, upf)
		upfsByName[name] = upf
	}
	return upfsByName
}
func BuildAllPCRFs(componentData map[string]ComponentData) map[string]ComponentData {
	pcrfsByName := make(map[string]ComponentData)

	for name, nets := range componentData {
		pcrf := ComponentData{
			Name:                name,
			Networks:            nets.Networks,
			ComponentsToConnect: nets.ComponentsToConnect,
		}
		// pcrfs = append(pcrfs, pcrf)
		pcrfsByName[name] = pcrf
	}
	return pcrfsByName
}
