package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/c-bata/go-prompt"
)

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "nmap -v1|2|3 ip_or_hostname", Description: "Scan the top 1000 ports of a remote host with various [v]erbosity levels:"},
		{Text: "nmap -T5 -sn 192.168.0.0/24|ip_or_hostname1,ip_or_hostname2,...", Description: "Run a ping sweep over an entire subnet or individual hosts very aggressively:"},
		{Text: "sudo nmap -A -iL path/to/file.txt", Description: "Enable OS detection, version detection, script scanning, and traceroute of hosts from a file:"},
		{Text: "nmap -p port1,port2,... ip_or_host1,ip_or_host2,...", Description: "Scan a specific list of ports (use `-p-` for all ports from 1 to 65535):"},
		{Text: "nmap -sC -sV -oA top-1000-ports ip_or_host1,ip_or_host2,...", Description: "Perform service and version detection of the top 1000 ports using default NSE scripts, writing results (`-oA`) to output files:"},
		{Text: "nmap --script 'default and safe' ip_or_host1,ip_or_host2,...", Description: "Scan target(s) carefully using `default and safe` NSE scripts:"},
		{Text: "nmap --script 'http-*' ip_or_host1,ip_or_host2,... -p 80,443", Description: "Scan for web servers running on standard ports 80 and 443 using all available `http-*` NSE scripts:"},
	}
	return CustomFilter(s, d.GetWordBeforeCursor())
}

func main() {
	fmt.Println("Please type `nmap` and press tab for options. 'exit' or 'quit' to exit the program.")
	for {
		t := prompt.Input("> ", completer)

		if strings.ToLower(t) == "exit" || strings.ToLower(t) == "quit" || strings.ToLower(t) == "q" {
			fmt.Println("Exiting program...")
			break
		}

		parts := strings.Fields(t)
		var cmd *exec.Cmd
		if len(parts) == 0 {
			fmt.Println("No command entered")
			continue
		} else if len(parts) == 1 && parts[0] == "nmap" {
			cmd = exec.Command("nmap", "-h")
		} else {
			cmd = exec.Command(parts[0], parts[1:]...)
		}

		output, err := cmd.Output()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		fmt.Println(string(output))
	}
}
