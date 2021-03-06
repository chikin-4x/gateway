package gateway

import (
	"os/exec"
	"syscall"
)

func DiscoverGateway() (ip string, err error) {
	routeCmd := exec.Command("route", "print", "0.0.0.0")
	routeCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := routeCmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return parseWindowsRoutePrint(output)
}
