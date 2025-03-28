package runtime

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/TencentBlueking/bkpaas/smart-app-builder/pkg/builder/plan"
	"github.com/TencentBlueking/bkpaas/smart-app-builder/pkg/utils"
)

// pindHandler is build handler that runs in podman-in-docker mode
type pindHandler struct {
	// execPath is the path of podman command
	execPath string
}

// GetSourceDir return the source code directory
func (h *pindHandler) GetSourceDir() string {
	return filepath.Join(h.getWorkspace(), "source")
}

// GetDestDir return the s-mart artifact directory
func (h *pindHandler) GetDestDir() string {
	return filepath.Join(h.getWorkspace(), "dest")
}

// GetTmpDir return the tmp directory. The tmp directory is used to store saas module tgz
func (h *pindHandler) GetTmpDir() string {
	return filepath.Join(h.getWorkspace(), "tmp")
}

func (h *pindHandler) Build(buildPlan *plan.BuildPlan) error {
	args := []string{"info"}
	return utils.RunCommand(h.execPath, args...)
}

func (h *pindHandler) getWorkspace() string {
	// return "/podman/smart-app"
	return "/Users/jamesge/Downloads/podman/smart-app"
}

// initWorkspace init workspace
func (h *pindHandler) initWorkspace() error {
	for _, dir := range []string{h.GetSourceDir(), h.GetDestDir()} {
		if err := os.MkdirAll(dir, 0o744); err != nil {
			return err
		}
	}

	return nil
}

func NewPindHandler() (*pindHandler, error) {
	execPath, err := exec.LookPath("podman")
	if err != nil {
		return nil, errors.New("podman command not found")
	}

	h := &pindHandler{execPath: execPath}

	if err := h.initWorkspace(); err != nil {
		return nil, err
	}
	return h, nil
}
