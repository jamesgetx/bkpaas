package buildconfig

import (
	"fmt"
	"strings"

	"github.com/TencentBlueking/bkpaas/cnb-builder-shim/pkg/utils"
)

// Buildpack 构建工具
type Buildpack struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version,omitempty"`
}

var defaultBuildpacks = map[string]Buildpack{
	"python": {Name: "bk-buildpack-python", Version: utils.EnvOrDefault("PYTHON_BUILDPACK_VERSION", "v213")},
	"nodejs": {Name: "bk-buildpack-nodejs", Version: utils.EnvOrDefault("NODEJS_BUILDPACK_VERSION", "v163")},
	"go":     {Name: "bk-buildpack-go", Version: utils.EnvOrDefault("GO_BUILDPACK_VERSION", "v191")},
	"apt":    {Name: "bk-buildpack-apt", Version: utils.EnvOrDefault("APT_BUILDPACK_VERSION", "v2")},
}

// GetBuildpackByLanguage 根据语言获取对应的 buildpack
func GetBuildpackByLanguage(language string) (*Buildpack, error) {
	if bp, ok := defaultBuildpacks[strings.ToLower(language)]; ok {
		return &bp, nil
	}
	return nil, fmt.Errorf("no buildpacks match with language: %s", language)
}

// GetDefaultVersionByBPName 根据 buildpack 的名字, 获取它的默认版本
func GetDefaultVersionByBPName(name string) (string, error) {
	// python 的 buildpack 配置 https://github.com/TencentBlueKing/blueking-paas/blob/builder-stack/cloudnative-buildpacks/buildpacks/bk-buildpack-python/cnb-buildpack/buildpack.toml
	for _, v := range defaultBuildpacks {
		if v.Name == name {
			return v.Version, nil
		}
	}
	return "", fmt.Errorf("no buildpacks match with name: %s", name)
}
