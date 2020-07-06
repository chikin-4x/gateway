package gateway

import (
	"os/exec"
)

func DiscoverGateway() (ip string, err error) {
	routeCmd := exec.Command("netstat", "-rn")
	output, err := routeCmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return parseSolarisNetstat(output)
}
