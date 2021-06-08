package version

import (
	"fmt"
	"runtime"
)

// GetDefaultVersion returns the current version
func GetDefaultVersion() string {
	return defaultConsumerVersion
}

func String() string {
	return fmt.Sprintf("%s version %s Arch(%s) runtime(%s)", "NsqConsumer", GetDefaultVersion(), runtime.GOARCH, runtime.Version())
}
