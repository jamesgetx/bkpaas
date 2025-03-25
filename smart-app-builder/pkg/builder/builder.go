package builder

import (
	"github.com/TencentBlueking/bkpaas/smart-app-builder/pkg/builder/handler"
)

// buildExecutor execute the process which build the source code to s-mart artifact
type buildExecutor struct {
	sourceURL string
	destURL   string
	handler   handler.RuntimeHandler
}

func (b *buildExecutor) Execute() error {
	return nil
}

// NewBuildExecutor create a new buildExecutor
func NewBuildExecutor(sourceURL, destURL string) (*buildExecutor, error) {
	if h, err := handler.NewRuntimeHandler(); err != nil {
		return nil, err
	} else {
		return &buildExecutor{sourceURL: sourceURL, destURL: destURL, handler: h}, nil
	}
}
