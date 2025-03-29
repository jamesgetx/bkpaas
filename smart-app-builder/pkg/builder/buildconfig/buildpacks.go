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

// BuildpacksEnvs buildpacks 的环境变量
func BuildpacksEnvs(name string) []string {
	var envs []string
	switch name {
	case "bk-buildpack-python":
		envs = append(envs,
			"PIP_VERSION=20.2.3",
			"DISABLE_COLLECTSTATIC=1",
			"BUILDPACK_S3_BASE_URL=https://bkpaas-runtimes-1252002024.file.myqcloud.com/python",
			"PIP_INDEX_URL=https://mirrors.cloud.tencent.com/pypi/simple/",
			"PIP_EXTRA_INDEX_URL=https://mirrors.tencent.com/tencent_pypi/simple/",
		)
		break
	case "bk-buildpack-nodejs":
		envs = append(
			envs,
			"STDLIB_FILE_URL=https://bkpaas-runtimes-1252002024.file.myqcloud.com/common/buildpack-stdlib/bk-node/stdlib.sh",
			"S3_DOMAIN=https://bkpaas-runtimes-1252002024.file.myqcloud.com/nodejs/node/release/linux-x64/",
			"NPM_REGISTRY=https://mirrors.tencent.com/npm/",
		)
		break
	case "bk-buildpack-go":
		envs = append(envs, "GO_BUCKET_URL=https://bkpaas-runtimes-1252002024.file.myqcloud.com/golang")
		break
	}
	return envs
}
