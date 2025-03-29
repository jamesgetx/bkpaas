package runtime

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/pkg/errors"

	"github.com/TencentBlueking/bkpaas/smart-app-builder/pkg/builder/config"
	"github.com/mholt/archives"
)

// startDaemon start container daemon
func startDaemon(cmd *exec.Cmd) error {
	if err := cmd.Start(); err != nil {
		return err
	}

	sockFile := config.G.DaemonSockFile

	// 10 次重试, 探测 daemon 是否就绪
	retryCount := 10

	for i := 0; i < retryCount; i++ {
		if _, err := os.Stat(sockFile); err == nil {
			conn, connErr := net.Dial("unix", sockFile)
			if conn != nil {
				conn.Close()
			}
			if connErr == nil {
				return nil
			}
		}
		time.Sleep(1 * time.Second)
	}

	return errors.New("daemon not ready")
}

// buildSourceTarball 生成 Procfile 文件, 并将模块源码打包成 source.tgz
func buildSourceTarball(srcDir, destDir string, Procfile map[string]string) error {
	procfileContent := ""
	for procName, procCommand := range Procfile {
		procfileContent += fmt.Sprintf("%s: %s\n", procName, procCommand)
	}

	procfilePath := filepath.Join(srcDir, "Procfile")
	if err := os.WriteFile(procfilePath, []byte(procfileContent), 0o755); err != nil {
		return errors.Wrap(err, "failed to write Procfile")
	}

	ctx := context.TODO()
	files, err := archives.FilesFromDisk(ctx, nil, map[string]string{srcDir: ""})
	if err != nil {
		return err
	}

	if err = os.MkdirAll(destDir, 0o744); err != nil {
		return err
	}
	out, err := os.Create(filepath.Join(destDir, "source.tgz"))
	if err != nil {
		return err
	}
	defer out.Close()

	format := archives.CompressedArchive{
		Compression: archives.Gz{},
		Archival:    archives.Tar{},
	}

	err = format.Archive(context.TODO(), out, files)
	if err != nil {
		return errors.Wrap(err, "failed to archive src")
	}
	return nil
}
