/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - PaaS 平台 (BlueKing - PaaS System) available.
 * Copyright (C) 2017 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 *     http://opensource.org/licenses/MIT
 *
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package devcontainer

import (
	"io"
	"os"
	"os/exec"
	"path"

	"github.com/pkg/errors"

	"github.com/TencentBlueking/bkpaas/cnb-builder-shim/internal/devcontainer/phase"
)

var (
	ReloadDir    = "/cnb/devcontainer/reload"
	ReloadLogDir = path.Join(ReloadDir, "log")
)

// ReloadStatus is the status of a reload operation.
type ReloadStatus string

const (
	ReloadProcessing ReloadStatus = "Processing"
	ReloadSuccess    ReloadStatus = "Success"
	ReloadFailed     ReloadStatus = "Failed"
	ReloadUnknown    ReloadStatus = "Unknown"
)

// NewHotReloadManager creates a new instance of HotReloadManager.
//
// It initializes a ResultFileRW and returns a pointer to a HotReloadManager and an error.
func NewHotReloadManager() (*HotReloadManager, error) {
	storage := ReloadResultFile{}
	if err := storage.Init(); err != nil {
		return nil, err
	}
	return &HotReloadManager{storage}, nil
}

// HotReloadManager ...
type HotReloadManager struct {
	ReloadResultStorage
}

// Rebuild ...
func (m HotReloadManager) Rebuild(reloadID string) error {
	cmd := phase.MakeBuilderCmd()
	return m.runCmd(reloadID, cmd)
}

// Relaunch ...
func (m HotReloadManager) Relaunch(reloadID string) error {
	cmd := phase.MakeLauncherCmd()
	return m.runCmd(reloadID, cmd)
}

func (m HotReloadManager) runCmd(reloadID string, cmd *exec.Cmd) error {
	w, err := m.ResultLogWriter(reloadID)
	if err != nil {
		return err
	}
	defer w.Close()

	cmd.Stdin = os.Stdin
	multiWriter := io.MultiWriter(os.Stdout, w)
	cmd.Stdout = multiWriter
	cmd.Stderr = multiWriter

	if err = cmd.Run(); err != nil {
		return err
	}
	return nil
}

// ReloadResultStorage is the interface that stores the result of app reload.
type ReloadResultStorage interface {
	// Init initializes some resources before use the storage if necessary.
	Init() error
	// ReadStatus returns the reload status for the given reload ID.
	ReadStatus(reloadID string) (ReloadStatus, error)
	// WriteStatus writes the status of a reload operation.
	WriteStatus(reloadID string, status ReloadStatus) error
	// ReadLog is a function that reads a log file based on the given reloadID.
	ReadLog(reloadID string) (string, error)
	// ResultLogWriter is a function that takes a reloadID as a parameter and returns a writer and an error.
	ResultLogWriter(reloadID string) (io.WriteCloser, error)
}

// ReloadResultFile implements ReloadResultStorage with local file system
type ReloadResultFile struct{}

// Init initializes the ReloadResultFile structure.
//
// It creates the necessary directory structure for the ReloadLogDir.
// Returns an error if there was a problem creating the directory.
func (f ReloadResultFile) Init() error {
	return os.MkdirAll(ReloadLogDir, 0o755)
}

// ReadStatus returns the reload status for the given reload ID.
func (f ReloadResultFile) ReadStatus(reloadID string) (ReloadStatus, error) {
	status, err := os.ReadFile(path.Join(ReloadDir, reloadID))
	if err != nil {
		return ReloadUnknown, errors.Wrap(err, "failed to read status")
	}
	return ReloadStatus(status), nil
}

// WriteStatus writes the status of a reload operation.
func (f ReloadResultFile) WriteStatus(reloadID string, status ReloadStatus) error {
	file, err := os.Create(path.Join(ReloadDir, reloadID))
	if err != nil {
		return errors.Wrap(err, "failed to write status")
	}
	defer file.Close()

	_, err = file.WriteString(string(status))
	return err
}

// ReadLog is a function that reads a log file based on the given reloadID.
func (f ReloadResultFile) ReadLog(reloadID string) (string, error) {
	content, err := os.ReadFile(path.Join(ReloadLogDir, reloadID))
	if err != nil {
		return "", errors.Wrap(err, "failed to read log")
	}
	return string(content), nil
}

// ResultLogWriter is a function that takes a reloadID as a parameter and returns a io.WriteCloser and an error.
// shell command will use this writer to write log
func (f ReloadResultFile) ResultLogWriter(reloadID string) (io.WriteCloser, error) {
	file, err := os.OpenFile(path.Join(ReloadLogDir, reloadID), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open log file")
	}
	return file, nil
}

var _ ReloadResultStorage = (*ReloadResultFile)(nil)
