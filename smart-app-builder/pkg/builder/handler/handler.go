package handler

import (
	"fmt"

	"github.com/TencentBlueking/bkpaas/cnb-builder-shim/pkg/utils"

	"github.com/TencentBlueking/bkpaas/smart-app-builder/pkg/builder/handler/runtime"
)

// NewRuntimeHandler return RuntimeHandler by RUNTIME env. The RUNTIME env only support pind and dind, if not set,
// default is pind. pind means podman-in-docker and dind means docker-in-docker.
func NewRuntimeHandler() (RuntimeHandler, error) {
	runtimeName := utils.EnvOrDefault("RUNTIME", "pind")
	switch runtimeName {
	case "pind":
		return runtime.NewPindHandler()
	// case "dind":
	//	return runtime.NewDindHandler()
	default:
		return nil, fmt.Errorf("unsupported runtime %s. only support pind and dind", runtimeName)
	}
}

// RuntimeHandler is build handler
type RuntimeHandler interface {
	// GetSourceDir return the source code directory
	GetSourceDir() string
	// GetDestDir return the s-mart artifact directory
	GetDestDir() string
	// GetTmpDir return the tmp directory. The tmp directory is used to store saas module tgz
	GetTmpDir() string
}
