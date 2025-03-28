package utils

import (
	"github.com/go-logr/logr"

	"github.com/TencentBlueking/bkpaas/cnb-builder-shim/pkg/logging"
)

var logger = logging.Default()

// GetLogger returns the default logger
func GetLogger() logr.Logger {
	return logger
}
