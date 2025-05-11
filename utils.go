package main

import (
	"fmt"
	"strings"

	"github.com/c-bata/go-prompt"
)

func GetHelp(category string) {
	switch category {
	case "ssh":
		fmt.Println("SSH Enumeration tips")
		fmt.Println("1. Figure out what auh methods are used - `--script ssh-auth-methods`")
		fmt.Println("2. If password auth is enabled, try brute force - `--script ssh-brute --script-args userdb=<usernames.txt>,passdb=<passwords.txt>`")
	case "snmp":
		fmt.Println("SNMP Enumeration tips")
		fmt.Println("SNMP runs on UDP port 161!")
		fmt.Println("1. Try brute force the community string - `--script snmp-brute --script-args snmp-brute.communitiesdb=<wordlist.txt>`")
		fmt.Println("2. Query network interfaces - `--script snmp-interfaces --script-args creds.snmp=<secret>`")
		fmt.Println("3. Query active network connections - `--script snmp-netstat --script-args creds.snmp=<secret>`")
		fmt.Println("4. Query running processes - `--script snmp-processes --script-args creds.snmp=<secret>`")
	case "smb":
		fmt.Println("SMB Enumeration tips")
		fmt.Println("1. Investigate security configuration - `--script smb-security-mode`")
		fmt.Println("2. List users - `--script smb-enum-users`")
		fmt.Println("3. Brute force users - `--script smb-brute --script-args smbtype=v2`")
		fmt.Println("4. List shares - `--script smb-enum-shares`")
		fmt.Println("5. Enumerate and list shares - `--script smb-enum-shares,smb-ls`")
	case "nfs":
		fmt.Println("NFS Enumeration tips")
		fmt.Println("1. Retrieve all configured exports - `--script nfs-showmount` ")
		fmt.Println("2. List shares - `--script nfs-ls` ")
	case "target specification":
		fmt.Println("TARGET SPECIFICATION:")
		fmt.Println("Can pass hostnames, IP addresses, networks, etc.")
		fmt.Println("Ex: scanme.nmap.org, microsoft.com/24, 192.168.0.1; 10.0.0-255.1-254")
		fmt.Println("-iL <inputfilename>: Input from list of hosts/networks")
		fmt.Println("-iR <num hosts>: Choose random targets")
		fmt.Println("--exclude <host1[,host2][,host3],...>: Exclude hosts/networks")
		fmt.Println("--excludefile <exclude_file>: Exclude list from file")
	case "host discovery":
		fmt.Println("HOST DISCOVERY:")
		fmt.Println("-sL: List Scan - simply list targets to scan")
		fmt.Println("-sn: Ping Scan - disable port scan")
		fmt.Println("-Pn: Treat all hosts as online -- skip host discovery")
		fmt.Println("-PS/PA/PU/PY[portlist]: TCP SYN, TCP ACK, UDP or SCTP discovery to given ports")
		fmt.Println("-PE/PP/PM: ICMP echo, timestamp, and netmask request discovery probes")
		fmt.Println("-PO[protocol list]: IP Protocol Ping")
		fmt.Println("-n/-R: Never do DNS resolution/Always resolve [default: sometimes]")
		fmt.Println("--dns-servers <serv1[,serv2],...>: Specify custom DNS servers")
		fmt.Println("--system-dns: Use OS's DNS resolver")
		fmt.Println("--traceroute: Trace hop path to each host")
	case "scan techniques":
		fmt.Println("SCAN TECHNIQUES:")
		fmt.Println("-sS/sT/sA/sW/sM: TCP SYN/Connect()/ACK/Window/Maimon scans")
		fmt.Println("-sU: UDP Scan")
		fmt.Println("-sN/sF/sX: TCP Null, FIN, and Xmas scans")
		fmt.Println("--scanflags <flags>: Customize TCP scan flags")
		fmt.Println("-sI <zombie host[:probeport]>: Idle scan")
		fmt.Println("-sY/sZ: SCTP INIT/COOKIE-ECHO scans")
		fmt.Println("-sO: IP protocol scan")
		fmt.Println("-b <FTP relay host>: FTP bounce scan")
	case "port specification":
		fmt.Println("PORT SPECIFICATION AND SCAN ORDER:")
		fmt.Println("-p <port ranges>: Only scan specified ports")
		fmt.Println("Ex: -p22; -p1-65535; -p U:53,111,137,T:21-25,80,139,8080,S:9")
		fmt.Println("--exclude-ports <port ranges>: Exclude the specified ports from scanning")
		fmt.Println("-F: Fast mode - Scan fewer ports than the default scan")
		fmt.Println("-r: Scan ports sequentially - don't randomize")
		fmt.Println("--top-ports <number>: Scan <number> most common ports")
		fmt.Println("--port-ratio <ratio>: Scan ports more common than <ratio>")
	case "version detection":
		fmt.Println("SERVICE/VERSION DETECTION:")
		fmt.Println("-sV: Probe open ports to determine service/version info")
		fmt.Println("--version-intensity <level>: Set from 0 (light) to 9 (try all probes)")
		fmt.Println("--version-light: Limit to most likely probes (intensity 2)")
		fmt.Println("--version-all: Try every single probe (intensity 9)")
		fmt.Println("--version-trace: Show detailed version scan activity (for debugging)")
	case "scripts":
		fmt.Println("SCRIPT SCAN:")
		fmt.Println("-sC: equivalent to --script=default")
		fmt.Println("--script=<Lua scripts>: <Lua scripts> is a comma separated list of")
		fmt.Println("directories, script-files or script-categories")
		fmt.Println("--script-args=<n1=v1,[n2=v2,...]>: provide arguments to scripts")
		fmt.Println("--script-args-file=filename: provide NSE script args in a file")
		fmt.Println("--script-trace: Show all data sent and received")
		fmt.Println("--script-updatedb: Update the script database.")
		fmt.Println("--script-help=<Lua scripts>: Show help about scripts.")
		fmt.Println("<Lua scripts> is a comma-separated list of script-files or")
		fmt.Println("script-categories.")
	case "OS detection":
		fmt.Println("OS DETECTION:")
		fmt.Println("-O: Enable OS detection")
		fmt.Println("--osscan-limit: Limit OS detection to promising targets")
		fmt.Println("--osscan-guess: Guess OS more aggressively")
	case "timing and performance":
		fmt.Println("TIMING AND PERFORMANCE:")
		fmt.Println("Options which take <time> are in seconds, or append 'ms' (milliseconds),")
		fmt.Println("'s' (seconds), 'm' (minutes), or 'h' (hours) to the value (e.g. 30m).")
		fmt.Println("-T<0-5>: Set timing template (higher is faster)")
		fmt.Println("--min-hostgroup/max-hostgroup <size>: Parallel host scan group sizes")
		fmt.Println("--min-parallelism/max-parallelism <numprobes>: Probe parallelization")
		fmt.Println("--min-rtt-timeout/max-rtt-timeout/initial-rtt-timeout <time>: Specifies")
		fmt.Println("probe round trip time.")
		fmt.Println("--max-retries <tries>: Caps number of port scan probe retransmissions.")
		fmt.Println("--host-timeout <time>: Give up on target after this long")
		fmt.Println("--scan-delay/--max-scan-delay <time>: Adjust delay between probes")
		fmt.Println("--min-rate <number>: Send packets no slower than <number> per second")
		fmt.Println("--max-rate <number>: Send packets no faster than <number> per second")
	case "firewall evasion and spoofing":
		fmt.Println("FIREWALL/IDS EVASION AND SPOOFING:")
		fmt.Println("-f; --mtu <val>: fragment packets (optionally w/given MTU)")
		fmt.Println("-D <decoy1,decoy2[,ME],...>: Cloak a scan with decoys")
		fmt.Println("-S <IP_Address>: Spoof source address")
		fmt.Println("-e <iface>: Use specified interface")
		fmt.Println("-g/--source-port <portnum>: Use given port number")
		fmt.Println("--proxies <url1,[url2],...>: Relay connections through HTTP/SOCKS4 proxies")
		fmt.Println("--data <hex string>: Append a custom payload to sent packets")
		fmt.Println("--data-string <string>: Append a custom ASCII string to sent packets")
		fmt.Println("--data-length <num>: Append random data to sent packets")
		fmt.Println("--ip-options <options>: Send packets with specified ip options")
		fmt.Println("--ttl <val>: Set IP time-to-live field")
		fmt.Println("--spoof-mac <mac address/prefix/vendor name>: Spoof your MAC address")
		fmt.Println("--badsum: Send packets with a bogus TCP/UDP/SCTP checksum")
	case "output":
		fmt.Println("OUTPUT:")
		fmt.Println("-oN/-oX/-oS/-oG <file>: Output scan in normal, XML, s|<rIpt kIddi3,")
		fmt.Println("and Grepable format, respectively, to the given filename.")
		fmt.Println("-oA <basename>: Output in the three major formats at once")
		fmt.Println("-v: Increase verbosity level (use -vv or more for greater effect)")
		fmt.Println("-d: Increase debugging level (use -dd or more for greater effect)")
		fmt.Println("--reason: Display the reason a port is in a particular state")
		fmt.Println("--open: Only show open (or possibly open) ports")
		fmt.Println("--packet-trace: Show all packets sent and received")
		fmt.Println("--iflist: Print host interfaces and routes (for debugging)")
		fmt.Println("--append-output: Append to rather than clobber specified output files")
		fmt.Println("--resume <filename>: Resume an aborted scan")
		fmt.Println("--noninteractive: Disable runtime interactions via keyboard")
		fmt.Println("--stylesheet <path/URL>: XSL stylesheet to transform XML output to HTML")
		fmt.Println("--webxml: Reference stylesheet from Nmap.Org for more portable XML")
		fmt.Println("--no-stylesheet: Prevent associating of XSL stylesheet w/XML output")
	case "misc":
		fmt.Println("MISC:")
		fmt.Println("-6: Enable IPv6 scanning")
		fmt.Println("-A: Enable OS detection, version detection, script scanning, and traceroute")
		fmt.Println("--datadir <dirname>: Specify custom Nmap data file location")
		fmt.Println("--send-eth/--send-ip: Send using raw ethernet frames or IP packets")
		fmt.Println("--privileged: Assume that the user is fully privileged")
		fmt.Println("--unprivileged: Assume the user lacks raw socket privileges")
		fmt.Println("-V: Print version number")
		fmt.Println("-h: Print this help summary page.")
	default:
		fmt.Println("Refer to the manpage for more info")
	}
}

func CustomFilter(suggestions []prompt.Suggest, sub string) []prompt.Suggest {
	return filterSuggestions(suggestions, sub, true, strings.Contains)
}

func filterSuggestions(suggestions []prompt.Suggest, sub string, ignoreCase bool, function func(string, string) bool) []prompt.Suggest {
	if sub == "" {
		return suggestions
	}
	if ignoreCase {
		sub = strings.ToUpper(sub)
	}

	ret := make([]prompt.Suggest, 0, len(suggestions))
	for i := range suggestions {
		c := suggestions[i].Text + suggestions[i].Description
		if ignoreCase {
			c = strings.ToUpper(c)
		}
		if function(c, sub) {
			ret = append(ret, suggestions[i])
		}
	}
	return ret
}
