package main

import "time"

// func ESXI_setup(hostname, wdir string) {
// 	reader := bufio.NewReader(os.Stdin)
// 	fmt.Println(color.Blue + "\nUsing ESXI_setup" + color.Reset)
// 	fmt.Println(color.Yellow + "Setting up ESXI Product ..." + color.Reset)
// 	time.Sleep(1 * time.Second)
// 	fmt.Println(color.Yellow + "Fetching list of Servers for installing ESXI products ..." + color.Reset)
// 	servers := fetchServers()
// 	fmt.Println(color.Green + "\nAvailable Servers:" + color.Reset)
// 	for i, serverName := range servers {
// 		fmt.Printf("%d. %s\n", i+1, serverName)
// 	}
// 	fmt.Print(color.Yellow + "\nEnter the number corresponding to the server you want to use: " + color.Reset)
// 	userServer, err := reader.ReadString('\n')
// 	if err != nil {
// 		fmt.Println("Error reading user input:", err)
// 		return
// 	}
// 	fmt.Println(color.Green + "\nSelected Server: " + userServer + color.Reset)

// }

func fetchServers() []string {
	// Simulating fetching servers (replace with actual logic to fetch servers)
	time.Sleep(2 * time.Second)
	return []string{"HP G9 Server", "Huawei V3 Server", "A5 Server", "HP G10 Server"}
}
