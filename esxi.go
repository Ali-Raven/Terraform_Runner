package main

import "time"

func fetchServers() []string {
	// Simulating fetching servers (replace with actual logic to fetch servers)
	time.Sleep(2 * time.Second)
	return []string{"HP G9 Server", "Huawei V3 Server", "A5 Server", "HP G10 Server"}
}
