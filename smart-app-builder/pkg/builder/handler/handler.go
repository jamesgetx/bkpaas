package handler

import (
	"fmt"

	"github.com/TencentBlueking/bkpaas/cnb-builder-shim/pkg/utils"
	"github.com/TencentBlueking/bkpaas/smart-app-builder/pkg/builder/handler/runtime"
)

// NewRuntimeHandler return RuntimeHandler by RUNTIME env. The RUNTIME env only support pind and dind, if not set,
// default is pind. pind means podman in docker and dind means docker in docker.
func NewRuntimeHandler() (RuntimeHandler, error) {
	runtimeName := utils.EnvOrDefault("RUNTIME", "pind")
	switch runtimeName {
	case "pind":
		return &runtime.PindHandler{Workspace: "/podman"}, nil
	case "dind":
		return &runtime.DindHandler{}, nil
	default:
		return nil, fmt.Errorf("unsupported runtime %s. only support pind and dind", runtimeName)
	}
}

type RuntimeHandler interface {
	GetWorkspace() string
}
