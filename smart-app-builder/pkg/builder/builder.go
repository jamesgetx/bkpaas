package builder

import (
	"fmt"
	"net/url"
	"os"

	"github.com/TencentBlueking/bkpaas/cnb-builder-shim/pkg/fetcher/fs"
	"github.com/TencentBlueking/bkpaas/cnb-builder-shim/pkg/fetcher/http"
	"github.com/TencentBlueking/bkpaas/cnb-builder-shim/pkg/utils"
	"github.com/go-logr/logr"

	"github.com/TencentBlueking/bkpaas/smart-app-builder/pkg/builder/handler"
	"github.com/TencentBlueking/bkpaas/smart-app-builder/pkg/builder/plan"
)

// BuildExecutor execute the process which build the source code to s-mart artifact
type BuildExecutor struct {
	sourceURL string
	destURL   string
	handler   handler.RuntimeHandler

	logger logr.Logger
}

// Execute run build process
func (b *BuildExecutor) Execute() error {
	sourceDir := b.handler.GetSourceDir()

	// 获取源码
	if err := b.fetchSource(sourceDir); err != nil {
		return err
	}

	buildPlan, err := plan.PrepareBuildPlan(sourceDir)
	if err != nil {
		return err
	}

	for _, step := range buildPlan.Steps {
		fmt.Printf("step: %v\n", step)
	}

	return b.handler.Build(buildPlan)
}

// fetchSource fetch the source code to destDir
func (b *BuildExecutor) fetchSource(destDir string) error {
	parsedURL, err := url.Parse(b.sourceURL)
	if err != nil {
		return err
	}

	switch parsedURL.Scheme {
	case "file":
		filePath := parsedURL.Path

		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return err
		}

		if !fileInfo.IsDir() {
			if err = fs.NewFetcher(b.logger).Fetch(filePath, destDir); err != nil {
				return err
			} else {
				return nil
			}
		}

		// 如果是文件目录, 目录不同时, 直接将源码拷贝到 destDir 下
		if filePath != destDir {
			return utils.CopyDir(filePath, destDir)
		}
		return nil

	case "http", "https":
		if err := http.NewFetcher(b.logger).Fetch(b.sourceURL, destDir); err != nil {
			return err
		}
	default:
		return fmt.Errorf("not support source-url scheme: %s", parsedURL.Scheme)
	}

	return nil
}

// NewBuildExecutor create a new buildExecutor
func NewBuildExecutor(logger logr.Logger, sourceURL, destURL string) (*BuildExecutor, error) {
	if h, err := handler.NewRuntimeHandler(); err != nil {
		return nil, err
	} else {
		return &BuildExecutor{logger: logger, sourceURL: sourceURL, destURL: destURL, handler: h}, nil
	}
}
