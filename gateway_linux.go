package gateway

import (
	"fmt"
	"io/ioutil"
	"os"
)

const (
	// See http://man7.org/linux/man-pages/man8/route.8.html
	file = "/proc/net/route"
)

func DiscoverGateway() (ip string, err error) {
	f, err := os.Open(file)
	if err != nil {
		return "", fmt.Errorf("Can't access %s", file)
	}
	defer f.Close()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("Can't read %s", file)
	}
	return parseLinuxProcNetRoute(bytes)
}
