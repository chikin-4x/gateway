// +build !darwin,!linux,!windows,!solaris,!freebsd

package gateway

import (
	"fmt"
	"runtime"
)

func DiscoverGateway() (string, err error) {
	err = fmt.Errorf("DiscoverGateway not implemented for OS %s", runtime.GOOS)
	return
}
