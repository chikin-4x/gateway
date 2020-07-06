package gateway

import (
	"bufio"
	"bytes"
	"errors"
	"runtime"
	"strings"
)

var errNoGateway = errors.New("no gateway found")

func parseWindowsRoutePrint(output []byte) (string, error) {
	// Windows route output format is always like this:
	// ===========================================================================
	// Interface List
	// 8 ...00 12 3f a7 17 ba ...... Intel(R) PRO/100 VE Network Connection
	// 1 ........................... Software Loopback Interface 1
	// ===========================================================================
	// IPv4 Route Table
	// ===========================================================================
	// Active Routes:
	// Network Destination        Netmask          Gateway       Interface  Metric
	//           0.0.0.0          0.0.0.0      192.168.1.1    192.168.1.100     20
	// ===========================================================================
	//
	// Windows commands are localized, so we can't just look for "Active Routes:" string
	// I'm trying to pick the active route,
	// then jump 2 lines and pick the third IP
	// Not using regex because output is quite standard from Windows XP to 8 (NEEDS TESTING)

	// lines := strings.Split(string(output), "\n")
	// sep := 0
	// for idx, line := range lines {
	// 	if sep == 3 {
	// 		// We just entered the 2nd section containing "Active Routes:"
	// 		if len(lines) <= idx+2 {
	// 			return "nil", errNoGateway
	// 		}

	// 		fields := strings.Fields(lines[idx+2])
	// 		if len(fields) < 3 {
	// 			return "nil", errNoGateway
	// 		}

	// 		ip := net.ParseIP(fields[2])
	// 		if ip != nil {
	// 			return ip, errNoGateway
	// 		}
	// 	}
	// 	if strings.HasPrefix(line, "=======") {
	// 		sep++
	// 		continue
	// 	}
	// }
	// return "nil", errNoGateway

	return "nil", errors.New("DiscoverGateway not implemented for OS " + runtime.GOOS)
}

func parseLinuxProcNetRoute(f []byte) (string, error) {
	/* /proc/net/route file:
	   Iface   Destination Gateway     Flags   RefCnt  Use Metric  Mask
	   eno1    00000000    C900A8C0    0003    0   0   100 00000000    0   00
	   eno1    0000A8C0    00000000    0001    0   0   100 00FFFFFF    0   00
	*/
	const (
		sep   = "\t" // field separator
		field = 0    // field containing hex gateway address
	)
	scanner := bufio.NewScanner(bytes.NewReader(f))
	for scanner.Scan() {
		// Skip header line
		if !scanner.Scan() {
			return "nil", errors.New("Invalid linux route file")
		}

		// get field containing gateway address. if len == 1, then there are no routes
		tokens := strings.Split(scanner.Text(), sep)
		if len(tokens) <= field || len(tokens) == 1 {
			return "nil", errors.New("Invalid linux route file")
		}

		return tokens[field], nil
	}
	return "nil", errors.New("Failed to parse linux route file")
}

func parseDarwinRouteGet(output []byte) (string, error) {
	// Darwin route out format is always like this:
	//    route to: default
	// destination: default
	//        mask: default
	//     gateway: 192.168.1.1
	//   interface: en0
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[0] == "interface:" {
			return fields[1], nil
		}
	}

	return "nil", errNoGateway
}

func parseBSDNetstat(output []byte) (string, error) {
	// netstat -rn produces the following on FreeBSD:
	// Routing tables
	//
	// Internet:
	// Destination        Gateway            Flags      Netif Expire
	// default            10.88.88.2         UGS         em0
	// 10.88.88.0/24      link#1             U           em0
	// 10.88.88.148       link#1             UHS         lo0
	// 127.0.0.1          link#2             UH          lo0
	//
	// Internet6:
	// Destination                       Gateway                       Flags      Netif Expire
	// ::/96                             ::1                           UGRS        lo0
	// ::1                               link#2                        UH          lo0
	// ::ffff:0.0.0.0/96                 ::1                           UGRS        lo0
	// fe80::/10                         ::1                           UGRS        lo0
	// ...
	outputLines := strings.Split(string(output), "\n")
	for _, line := range outputLines {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[0] == "default" {
			return fields[3], nil
		}
	}

	return "nil", errNoGateway
}

func parseSolarisNetstat(output []byte) (string, error) {
	// netstat -rn produces the following on Solaris:
	//   Routing Table: IPv4
	//   Destination           Gateway           Flags  Ref     Use     Interface
	//   -------------------- -------------------- ----- ----- ---------- ---------
	//   default              172.16.32.1          UG        2      76419 net0
	//   127.0.0.1            127.0.0.1            UH        2         36 lo0
	//   172.16.32.0          172.16.32.17         U         4       8100 net0

	//   Routing Table: IPv6
	// 	Destination/Mask            Gateway                   Flags Ref   Use    If
	//   --------------------------- --------------------------- ----- --- ------- -----
	//   ::1                         ::1                         UH      3   75382 lo0
	//   2001:470:deeb:32::/64       2001:470:deeb:32::17        U       3    2744 net0
	//   fe80::/10                   fe80::6082:52ff:fedc:7df0   U       3    8430 net0
	outputLines := strings.Split(string(output), "\n")
	for _, line := range outputLines {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[0] == "default" {
			return fields[5], nil
		}
	}

	return "nil", errNoGateway
}
