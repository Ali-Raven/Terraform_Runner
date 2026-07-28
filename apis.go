package main

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/TwiN/go-color"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
)

func GetAllDatastoreName(vcURL, username, pass string) ([]string, error) {
	fmt.Printf("%s\n Getting list of datasotres from vCenter ...%s\n", color.Yellow, color.Reset)
	time.Sleep(1 * time.Second)
	ctx := context.Background()

	// 1. Format the vCenter URL
	u, err := url.Parse(fmt.Sprintf("https://%s:%s@%s/sdk", username, pass, vcURL))

	if err != nil {
		return nil, fmt.Errorf("invalid vCenter URL: %w", err)
	}

	// 2. Connect to vCenter
	client, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		return nil, fmt.Errorf("%sfailed to connect to vCenter: %w%s", color.Red, err, color.Reset)
	}
	defer client.Logout(ctx)

	// 3. Initialize the Finder
	finder := find.NewFinder(client.Client, true)

	// 4. Find the default datacenter to scope our search
	dc, err := finder.DefaultDatacenter(ctx)
	if err != nil {
		return nil, fmt.Errorf("%sfailed to find default datacenter: %w%s", color.Red, err, color.Reset)
	}
	finder.SetDatacenter(dc)

	// 5. Find ALL datastores using the "*" wildcard
	datastores, err := finder.DatastoreList(ctx, "*")
	if err != nil {
		return nil, fmt.Errorf("%sfailed to list datastores: %w%s", color.Red, err, color.Reset)
	}

	// 6. Extract the names into a simple string slice
	var dsNames []string
	for _, ds := range datastores {
		dsNames = append(dsNames, ds.Name())
	}

	return dsNames, nil
}
