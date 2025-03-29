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

	bcfg "github.com/TencentBlueking/bkpaas/smart-app-builder/pkg/builder/buildconfig"
	"github.com/TencentBlueking/bkpaas/smart-app-builder/pkg/builder/config"
	"github.com/TencentBlueking/bkpaas/smart-app-builder/pkg/builder/plan"
	"github.com/mholt/archives"
	"github.com/samber/lo"
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

// runCNBBuilder run cnb builder
func runCNBBuilder(buildPlan *plan.BuildPlan, execPath, sourceDir, tmpDir string) error {
	for _, step := range buildPlan.Steps {
		moduleSrcDir := filepath.Join(sourceDir, step.SourceDir)
		moduleSrcTgz := filepath.Join(tmpDir, step.BuildModuleName, "source.tgz")
		if err := buildSourceTarball(moduleSrcDir, moduleSrcTgz, buildPlan.Procfile); err != nil {
			return err
		}

		bindTarget := "/tmp/source.tgz"
		args := []string{"run"}

		envSlice := buildEnvArgs(step, bindTarget)
		lo.ForEach(envSlice, func(env string, index int) {
			args = append(args, "-e", env)
		})

		// 挂载源码压缩包
		args = append(args, "--mount", fmt.Sprintf("type=bind,source=%s,target=%s", moduleSrcTgz, bindTarget))

		args = append(
			args,
			"--mount",
			fmt.Sprintf("type=bind,source=%s,target=/var/run/docker.sock", config.G.DaemonSockFile),
		)

		args = append(args, config.G.CNBBuilderImage)

		cmd := exec.Command(execPath, args...)
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

// buildSourceTarball 生成 Procfile 文件, 并将模块源码打包成 source.tgz
func buildSourceTarball(sourceDir, destTgz string, Procfile map[string]string) error {
	procfileContent := ""
	for procName, procCommand := range Procfile {
		procfileContent += fmt.Sprintf("%s: %s\n", procName, procCommand)
	}

	procfilePath := filepath.Join(sourceDir, "Procfile")
	if err := os.WriteFile(procfilePath, []byte(procfileContent), 0o755); err != nil {
		return errors.Wrap(err, "failed to write Procfile")
	}

	ctx := context.TODO()
	files, err := archives.FilesFromDisk(ctx, nil, map[string]string{sourceDir: ""})
	if err != nil {
		return err
	}

	if err = os.MkdirAll(filepath.Dir(destTgz), 0o744); err != nil {
		return err
	}
	out, err := os.Create(destTgz)
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

// buildEnvArgs
func buildEnvArgs(step *plan.ModuleBuildStep, bindTarget string) []string {
	envSlice := lo.MapToSlice(step.Envs, func(k string, v string) string {
		return fmt.Sprintf("%s=%s", k, v)
	})
	lo.ForEach(step.Buildpacks, func(bp bcfg.Buildpack, index int) {
		envSlice = append(envSlice, bcfg.BuildpacksEnvs(bp.Name)...)
	})

	envSlice = append(
		envSlice,
		[]string{
			fmt.Sprintf("REQUIRED_BUILDPACKS=%s", step.RequiredBuildpacks),
			fmt.Sprintf("OUTPUT_IMAGE=%s", step.OutPutImage),
			fmt.Sprintf("CNB_RUN_IMAGE=%s", config.G.ScratchImage),
			fmt.Sprintf("SOURCE_GET_URL=file://%s", bindTarget),
			"USE_DOCKER_DAEMON=true",
		}...)

	return envSlice
}
