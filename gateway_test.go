package gateway

import (
	"testing"
)

type testcase struct {
	output  []byte
	ok      bool
	gateway string
}

func TestParseWindowsRoutePrint(t *testing.T) {
	correctData := []byte(`
===========================================================================
Interface List
  8 ...00 12 3f a7 17 ba ...... Intel(R) PRO/100 VE Network Connection
  1 ........................... Software Loopback Interface 1
===========================================================================
IPv4 Route Table
===========================================================================
Active Routes:
Network Destination        Netmask          Gateway       Interface  Metric
          0.0.0.0          0.0.0.0       10.88.88.2     10.88.88.149     10
===========================================================================
Persistent Routes:
`)
	localizedData := []byte(
		`===========================================================================
Liste d'Interfaces
 17...00 28 f8 39 61 6b ......Microsoft Wi-Fi Direct Virtual Adapter
  1...........................Software Loopback Interface 1
===========================================================================
IPv4 Table de routage
===========================================================================
Itinéraires actifs :
Destination réseau    Masque réseau  Adr. passerelle   Adr. interface Métrique
          0.0.0.0          0.0.0.0      10.88.88.2     10.88.88.149     10
===========================================================================
Itinéraires persistants :
  Aucun
`)
	randomData := []byte(`
Lorem ipsum dolor sit amet, consectetur adipiscing elit,
sed do eiusmod tempor incididunt ut labore et dolore magna
aliqua. Ut enim ad minim veniam, quis nostrud exercitation
`)
	noRoute := []byte(`
===========================================================================
Interface List
  8 ...00 12 3f a7 17 ba ...... Intel(R) PRO/100 VE Network Connection
  1 ........................... Software Loopback Interface 1
===========================================================================
IPv4 Route Table
===========================================================================
Active Routes:
`)
	badRoute1 := []byte(`
===========================================================================
Interface List
  8 ...00 12 3f a7 17 ba ...... Intel(R) PRO/100 VE Network Connection
  1 ........................... Software Loopback Interface 1
===========================================================================
IPv4 Route Table
===========================================================================
Active Routes:
===========================================================================
Persistent Routes:
`)
	badRoute2 := []byte(`
===========================================================================
Interface List
  8 ...00 12 3f a7 17 ba ...... Intel(R) PRO/100 VE Network Connection
  1 ........................... Software Loopback Interface 1
===========================================================================
IPv4 Route Table
===========================================================================
Active Routes:
Network Destination        Netmask          Gateway       Interface  Metric
          0.0.0.0          0.0.0.0          foo           10.88.88.149     10
===========================================================================
Persistent Routes:
`)

	testcases := []testcase{
		{correctData, false, "10.88.88.2"},
		{localizedData, false, "10.88.88.2"},
		{randomData, false, ""},
		{noRoute, false, ""},
		{badRoute1, false, ""},
		{badRoute2, false, ""},
	}

	test(t, testcases, parseWindowsRoutePrint)
}

func TestParseLinuxProcNetRoute(t *testing.T) {
	correctData := []byte(`Iface	Destination	Gateway 	Flags	RefCnt	Use	Metric	Mask		MTU	Window	IRTT
wlp4s0	00000000	0108A8C0	0003	0	0	600	00000000	0	0	0
wlp4s0	0000FEA9	00000000	0001	0	0	1000	0000FFFF	0	0	0
docker0	000011AC	00000000	0001	0	0	0	0000FFFF	0	0	0
docker_gwbridge	000012AC	00000000	0001	0	0	0	0000FFFF	0	0	0
wlp4s0	0008A8C0	00000000	0001	0	0	600	00FFFFFF	0	0	0
`)
	noRoute := []byte(`
Iface   Destination     Gateway         Flags   RefCnt  Use     Metric  Mask            MTU     Window  IRTT
`)

	testcases := []testcase{
		{correctData, true, "wlp4s0"},
		{noRoute, false, ""},
	}

	test(t, testcases, parseLinuxProcNetRoute)
}

func TestParseDarwinRouteGet(t *testing.T) {
	correctData := []byte(`
   route to: 0.0.0.0
destination: default
       mask: default
    gateway: 172.16.32.1
  interface: en0
      flags: <UP,GATEWAY,DONE,STATIC,PRCLONING>
 recvpipe  sendpipe  ssthresh  rtt,msec    rttvar  hopcount      mtu     expire
       0         0         0         0         0         0      1500         0
`)
	randomData := []byte(`
test
Lorem ipsum dolor sit amet, consectetur adipiscing elit,
sed do eiusmod tempor incididunt ut labore et dolore magna
aliqua. Ut enim ad minim veniam, quis nostrud exercitation
`)
	noRoute := []byte(`
   route to: 0.0.0.0
destination: default
       mask: default
`)

	testcases := []testcase{
		{correctData, true, "en0"},
		{randomData, false, ""},
		{noRoute, false, ""},
	}

	test(t, testcases, parseDarwinRouteGet)
}

func TestParseBSDNetstat(t *testing.T) {
	correctDataFreeBSD := []byte(`
Routing tables

Internet:
Destination        Gateway            Flags      Netif Expire
default            10.88.88.2         UGS         em0
10.88.88.0/24      link#1             U           em0
10.88.88.148       link#1             UHS         lo0
127.0.0.1          link#2             UH          lo0

Internet6:
Destination                       Gateway                       Flags      Netif Expire
::/96                             ::1                           UGRS        lo0
::1                               link#2                        UH          lo0
::ffff:0.0.0.0/96                 ::1                           UGRS        lo0
fe80::/10                         ::1                           UGRS        lo0
`)
	randomData := []byte(`
Lorem ipsum dolor sit amet, consectetur adipiscing elit,
sed do eiusmod tempor incididunt ut labore et dolore magna
aliqua. Ut enim ad minim veniam, quis nostrud exercitation
`)
	noRoute := []byte(`
Internet:
Destination        Gateway            Flags      Netif Expire
10.88.88.0/24      link#1             U           em0
10.88.88.148       link#1             UHS         lo0
127.0.0.1          link#2             UH          lo0
`)

	testcases := []testcase{
		{correctDataFreeBSD, true, "em0"},
		{randomData, false, ""},
		{noRoute, false, ""},
	}

	test(t, testcases, parseBSDNetstat)
}

func TestParseSolarisNetstat(t *testing.T) {
	correctDataSolaris := []byte(`
Routing Table: IPv4
  Destination           Gateway           Flags  Ref     Use     Interface
-------------------- -------------------- ----- ----- ---------- ---------
default              172.16.32.1          UG        2      76419 net0
127.0.0.1            127.0.0.1            UH        2         36 lo0
172.16.32.0          172.16.32.17         U         4       8100 net0

Routing Table: IPv6
  Destination/Mask            Gateway                   Flags Ref   Use    If
--------------------------- --------------------------- ----- --- ------- -----
::1                         ::1                         UH      3   75382 lo0
2001:470:deeb:32::/64       2001:470:deeb:32::17        U       3    2744 net0
fe80::/10                   fe80::6082:52ff:fedc:7df0   U       3    8430 net0
`)
	randomData := []byte(`
Lorem ipsum dolor sit amet, consectetur adipiscing elit,
sed do eiusmod tempor incididunt ut labore et dolore magna
aliqua. Ut enim ad minim veniam, quis nostrud exercitation
`)

	testcases := []testcase{
		{correctDataSolaris, true, "net0"},
		{randomData, false, ""},
	}

	test(t, testcases, parseSolarisNetstat)
}

func test(t *testing.T, testcases []testcase, fn func([]byte) (string, error)) {
	for i, tc := range testcases {
		net, err := fn(tc.output)
		if tc.ok {
			if err != nil {
				t.Errorf("Unexpected error in test #%d: %v", i, err)
			}
			if net != tc.gateway {
				t.Errorf("Unexpected gateway address %v != %s", net, tc.gateway)
			}
		} else if err == nil {
			t.Errorf("Unexpected error in test #%d", i)
		}
	}
}
