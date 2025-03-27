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
)

// buildExecutor execute the process which build the source code to s-mart artifact
type buildExecutor struct {
	sourceURL string
	destURL   string
	handler   handler.RuntimeHandler

	logger logr.Logger
}

// Execute run build process
func (b *buildExecutor) Execute() error {
	// 获取源码
	if err := b.fetchSource(); err != nil {
		return err
	}

	plan, err := PrepareBuildPlan(b.getSourceDir())
	if err != nil {
		return err
	}

	for _, step := range plan.Steps {
		fmt.Printf("step: %+v\n", step)
	}

	return nil
}

// fetchSource fetch the source code to sourceDir
func (b *buildExecutor) fetchSource() error {
	parsedURL, err := url.Parse(b.sourceURL)
	if err != nil {
		return err
	}

	sourceDir := b.getSourceDir()

	switch parsedURL.Scheme {
	case "file":
		filePath := parsedURL.Path

		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return err
		}

		if !fileInfo.IsDir() {
			if err = fs.NewFetcher(b.logger).Fetch(filePath, sourceDir); err != nil {
				return err
			} else {
				return nil
			}
		}

		// 如果是文件目录, 目录不同时, 直接将源码拷贝到 sourceDir 下
		if filePath != sourceDir {
			return utils.CopyDir(filePath, sourceDir)
		}
		return nil

	case "http", "https":
		if err := http.NewFetcher(b.logger).Fetch(b.sourceURL, sourceDir); err != nil {
			return err
		}
	default:
		return fmt.Errorf("not support source-url scheme: %s", parsedURL.Scheme)
	}

	return nil
}

// getSourceDir return the source code directory
func (b *buildExecutor) getSourceDir() string {
	return b.handler.GetSourceDir()
}

// NewBuildExecutor create a new buildExecutor
func NewBuildExecutor(logger logr.Logger, sourceURL, destURL string) (*buildExecutor, error) {
	if h, err := handler.NewRuntimeHandler(); err != nil {
		return nil, err
	} else {
		return &buildExecutor{logger: logger, sourceURL: sourceURL, destURL: destURL, handler: h}, nil
	}
}
