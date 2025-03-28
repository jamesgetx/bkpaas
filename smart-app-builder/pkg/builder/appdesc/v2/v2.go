package v2

import (
	"fmt"

	bcfg "github.com/TencentBlueking/bkpaas/smart-app-builder/pkg/builder/buildconfig"
	"github.com/TencentBlueking/bkpaas/smart-app-builder/pkg/utils"
)

// AppDescConfig spec_version: 2 版本的 app_desc
type AppDescConfig struct {
	SpecVersion int                          `yaml:"spec_version"`
	AppVersion  string                       `yaml:"app_version"`
	AppInfo     AppInfo                      `yaml:"app"`
	Modules     map[string]ModuleDescription `yaml:"modules"`
}

// Validate 验证 app_desc
func (cfg *AppDescConfig) Validate() error {
	if cfg.SpecVersion != 2 {
		return fmt.Errorf("spec version must be 2")
	}

	if cfg.AppInfo.AppCode == "" {
		return fmt.Errorf("app code is empty")
	}

	if len(cfg.Modules) == 0 {
		return fmt.Errorf("modules is empty")
	}

	return utils.ValidateVersion(cfg.AppVersion)
}

// GenerateProcfile 生成 Procfile
func (cfg *AppDescConfig) GenerateProcfile() map[string]string {
	procfile := make(map[string]string)

	for moduleName, module := range cfg.Modules {
		for processName, process := range module.Processes {
			procfile[moduleName+"-"+processName] = process.ProcCommand
		}
	}

	return procfile
}

// GenerateModuleBuildConfig 生成 ModuleBuildConfig
func (cfg *AppDescConfig) GenerateModuleBuildConfig() ([]bcfg.ModuleBuildConfig, error) {
	config := make([]bcfg.ModuleBuildConfig, 0)

	for moduleName, module := range cfg.Modules {
		envs := make(map[string]string)
		for _, env := range module.EnvVariables {
			envs[env.Key] = env.Value
		}

		// 如果未指定, 表示当前目录
		src := module.SourceDir
		if src == "" {
			src = "."
		}

		buildpacks := module.Build.Buildpacks
		if buildpacks == nil {
			if bp, err := bcfg.GetBuildpackByLanguage(module.Language); err != nil {
				return nil, err
			} else {
				buildpacks = []bcfg.Buildpack{*bp}
			}
		}

		config = append(config, bcfg.ModuleBuildConfig{
			SourceDir:  src,
			ModuleName: moduleName,
			Envs:       envs,
			Buildpacks: buildpacks,
		})
	}

	return config, nil
}

// AppInfo app 字段
type AppInfo struct {
	AppCode string `yaml:"bk_app_code"`
}

// ModuleDescription 单个 module 字段
type ModuleDescription struct {
	SourceDir    string             `yaml:"source_dir"`
	Language     string             `yaml:"language"`
	Processes    map[string]Process `yaml:"processes"`
	EnvVariables []AppEnvVar        `yaml:"env_variables"`
	Build        bcfg.BuildConfig   `yaml:"build,omitempty"`
}

// Process 进程配置
type Process struct {
	ProcCommand string `yaml:"command"`
}

// AppEnvVar 单个环境变量
type AppEnvVar struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}
