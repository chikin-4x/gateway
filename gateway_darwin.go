package gateway

import (
	"os/exec"
)

func DiscoverGateway() (string, error) {
	routeCmd := exec.Command("/sbin/route", "-n", "get", "0.0.0.0")
	output, err := routeCmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return parseDarwinRouteGet(output)
}
