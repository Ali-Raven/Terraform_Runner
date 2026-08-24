package apis

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/TwiN/go-color"
	"github.com/joho/godotenv"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/ovf"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
)

type ProgressReader struct {
	Reader io.Reader
	Total  int64
	Reads  int64
}

func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	pr.Reads += int64(n)
	if pr.Total > 0 {
		percent := float64(pr.Reads) / float64(pr.Total) * 100
		fmt.Printf("\rUploading... %s%.2f%%%s", color.Yellow, percent, color.Reset)
	}
	return n, err
}

// this function is for validating the selected OVF file and its associated files in the same directory
func ValidateOvfDirectory(ovfPath string) error {
	fileInfo, err := os.Stat(ovfPath)
	if err != nil {
		return fmt.Errorf("%sfile does not exist: %w%s", color.Red, err, color.Reset)
	}
	if fileInfo.Size() == 0 {
		return fmt.Errorf("%sselected OVF file is empty%s", color.Red, color.Reset)
	}
	return nil
}

func GetAllDatastoreName(vcURL, username, pass string) ([]string, error) {
	log.Default().Printf("%s Getting list of datasotres from vCenter ...%s\n", color.Yellow, color.Reset)
	time.Sleep(1 * time.Second)
	ctx := context.Background()

	// 1. Format the vCenter URL
	u, err := url.Parse(fmt.Sprintf("https://%s:%s@%s/sdk", username, pass, vcURL))

	if err != nil {
		return nil, fmt.Errorf("invalid vCenter URL: %w", err)
	} else {
		log.Default().Printf("%s URL valication is done. %s", color.Green, color.Reset)
	}

	// 2. Connect to vCenter
	client, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		return nil, fmt.Errorf("%sfailed to connect to vCenter: %w%s", color.Red, err, color.Reset)
	} else {
		log.Default().Printf("%s Connected to vCenter Successfully. %s", color.Green, color.Reset)
	}
	defer client.Logout(ctx)

	// 3. Initialize the Finder
	finder := find.NewFinder(client.Client, true)

	// 4. Find the default datacenter to scope our search
	dc, err := finder.DefaultDatacenter(ctx)
	if err != nil {
		return nil, fmt.Errorf("%sfailed to find default datacenter: %w%s", color.Red, err, color.Reset)
	} else {
		log.Default().Printf("%s default datacenter is identified : %s%s", color.Green, dc.Name(), color.Reset)
	}
	finder.SetDatacenter(dc)

	// 5. Find ALL datastores using the "*" wildcard
	datastores, err := finder.DatastoreList(ctx, "*")
	if err != nil {
		return nil, fmt.Errorf("%sfailed to list datastores: %w%s", color.Red, err, color.Reset)
	} else {
		log.Default().Printf("%s all datasotres is identified Successfully.%s", color.Green, color.Reset)
	}

	// 6. Extract the names into a simple string slice
	var dsNames []string
	for _, ds := range datastores {
		dsNames = append(dsNames, ds.Name())
	}

	return dsNames, nil
}

func DeployOvfTemplateOnVcenter(hostname, wdir, ovfFilePath, vmName , diskProvisionType string) {
	fmt.Printf("\n%s%sDeploying VMs section ...%s%s\n", color.Bold, color.Yellow, color.Reset, color.Reset)
	time.Sleep(1 * time.Second)

	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		fmt.Printf("%s Can't read or load the vCenter Info , try again %s\n", color.Red, color.Reset)
	}

	vCenterURL := os.Getenv("vCenterURL")
	vCenterUserName := os.Getenv("vCenterUserName")
	vCenterPass := os.Getenv("vCenterPassword")
	targetHost := os.Getenv("targetHost")

	u, err := url.Parse(fmt.Sprintf("https://%s:%s@%s/sdk", vCenterUserName, vCenterPass, vCenterURL))

	if err != nil {
		log.Panicf("%sinvalid vCenter URL: %v%s", color.Red, err, color.Reset)
	}

	// Connecting to vCenter
	client, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		log.Fatalf("Failed to connect to vCenter: %v", err)
	} else {
		log.Default().Printf("%sConnected to vCenter successfully!%s", color.Green, color.Reset)
	}
	defer client.Logout(ctx)

	finder := find.NewFinder(client.Client, true)

	// Resolve Target Infrastructure Objects
	dc, err := finder.DefaultDatacenter(ctx)
	if err != nil {
		log.Fatalf("Failed to find default datacenter: %v", err)
	} else {
		log.Default().Printf("%sFound default datacenter: %s%s", color.Green, dc.Name(), color.Reset)
	}
	finder.SetDatacenter(dc)

	// pinpointing the target host from .env file
	host, err := finder.HostSystem(ctx, targetHost)
	if err != nil {
		log.Fatalf("Failed to find target host: %v", err)
	} else {
		log.Default().Printf("%sFound target host: %s%s", color.Green, host.Name(), color.Reset)
	}

	// finding the resource pool of the target host that we find
	rp, err := host.ResourcePool(ctx)
	if err != nil {
		log.Fatalf("Failed to find default resource pool: %v", err)
	} else {
		log.Default().Printf("%sFound default resource pool ... %s%s", color.Green, rp.Name(), color.Reset)
	}

	var hostProps mo.HostSystem

	// Ask vCenter specifically for the datastore property of this host
	err = host.Properties(ctx, host.Reference(), []string{"datastore"}, &hostProps)
	if err != nil {
		log.Fatalf("Failed to get host properties: %v", err)
	} else {
		log.Default().Printf("%sRetrieved host properties successfully!%s", color.Green, color.Reset)
	}

	// Ensure the host actually has at least one datastore attached
	if len(hostProps.Datastore) == 0 {
		log.Fatalf("Host '%s' has no datastores attached!", targetHost)
	}

	// Grab the first datastore in the host's list
	dsRef := hostProps.Datastore[0]

	// Convert that raw reference into a usable Datastore object
	ds := object.NewDatastore(client.Client, dsRef)

	// fetching datastore name for showing in the log
	var dsProps mo.Datastore
	err = ds.Properties(ctx, ds.Reference(), []string{"name"}, &dsProps)
	if err == nil {
		log.Default().Printf("%sAutomatically selected Datastore: %s%s\n", color.Green, dsProps.Name, color.Reset)
	}

	// finding the folder for deploying vm properties
	folder, err := finder.DefaultFolder(ctx)
	if err != nil {
		log.Fatalf("Failed to find default folder: %v", err)
	} else {
		log.Default().Printf("%sFound default folder: %s%s", color.Green, folder.Name(), color.Reset)
	}

	// Read & Parse OVF Descriptor
	ovfContent, err := os.ReadFile(ovfFilePath)
	if err != nil {
		log.Fatalf("Failed to read OVF file: %v", err)
	} else {
		log.Default().Printf("%sSuccessfully read OVF file: %s%s", color.Green, ovfFilePath, color.Reset)
	}

	ovfManager := ovf.NewManager(client.Client)

	// Build Import Spec
	crd := types.OvfCreateImportSpecParams{
		EntityName: vmName,
		DiskProvisioning: diskProvisionType,
	}

	spec, err := ovfManager.CreateImportSpec(ctx, string(ovfContent), rp, ds, &crd)
	if err != nil {
		log.Fatalf("Failed to create import spec: %v", err)
	} else {
		log.Default().Printf("%sSuccessfully created import spec for VM: %s%s", color.Green, vmName, color.Reset)
	}
	if spec.Error != nil {
		log.Fatalf("Import spec error: %v", spec.Error[0].LocalizedMessage)
	} else {
		log.Default().Printf("%sImport spec is valid and ready for deployment.%s", color.Green, color.Reset)
	}

	// Initiate Import (NFC Lease creation)
	lease, err := rp.ImportVApp(ctx, spec.ImportSpec, folder, host)
	if err != nil {
		log.Fatalf("Failed to initiate vApp import: %v", err)
	} else {
		log.Default().Printf("%sSuccessfully initiated vApp import. NFC Lease created.%s", color.Green, color.Reset)
	}

	// Wait for NFC lease to be ready
	info, err := lease.Wait(ctx, spec.FileItem)
	if err != nil {
		log.Fatalf("Lease wait failed: %v", err)
	} else {
		log.Default().Printf("%sNFC Lease is ready for file upload.%s", color.Green, color.Reset)
	}

	// Proper standard upload loop
	updater := lease.StartUpdater(ctx, info)
	defer updater.Done()


	// fmt.Println(info.Items)
	// time.Sleep(10000 * time.Second)
	// Loop through the items (disks/files) the lease expects and upload them
	for _, item := range info.Items {
		// Assume the disk files (.vmdk) are in the same directory as the .ovf file
		diskFilePath := filepath.Join(filepath.Dir(ovfFilePath), item.Path)

		file, err := os.Open(diskFilePath)
		if err != nil {
			log.Fatalf("Failed to open disk file %s: %v", diskFilePath, err)
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			log.Fatalf("Failed to stat file: %v", err)
		}

		fmt.Println()
		log.Default().Printf("%sStarting upload for %s (Size: %d bytes)...%s", color.Green, item.Path, fileInfo.Size(), color.Reset)

		pr := &ProgressReader{
			Reader: file,
			Total:  fileInfo.Size(),
		}
		// Upload the specific file
		err = lease.Upload(ctx, item, pr, soap.DefaultUpload)
		if err != nil {
			log.Fatalf("Disk upload failed for %s: %v", item.Path, err)
		}
		file.Close()
	}

	
	// Complete the lease setup
	if err := lease.Complete(ctx); err != nil {
		log.Fatalf("Failed to complete lease: %v", err)
	}

	log.Default().Printf("%s\nSuccessfully deployed OVF template as VM '%s'!\n", color.Green, vmName)
}
