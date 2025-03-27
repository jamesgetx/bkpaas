package v3

import (
	"github.com/TencentBlueking/bkpaas/smart-app-builder/pkg/builder/buildconfig"
)

// AppDescConfig specVersion: 3 版本的 app_desc
type AppDescConfig struct {
	SpecVersion int                 `yaml:"specVersion"`
	AppVersion  string              `yaml:"appVersion"`
	AppInfo     AppInfo             `yaml:"app"`
	Modules     []ModuleDescription `yaml:"modules"`
}

// Validate 验证 app_desc
func (cfg *AppDescConfig) Validate() error {
	return nil
}

// GenerateProcfile 生成 Procfile
func (cfg *AppDescConfig) GenerateProcfile() map[string]string {
	return nil
}

// GenerateModuleBuildConfig 生成 ModuleBuildConfig
func (cfg *AppDescConfig) GenerateModuleBuildConfig() ([]buildconfig.ModuleBuildConfig, error) {
	return nil, nil
}

// AppInfo app 字段
type AppInfo struct {
	AppCode string `yaml:"bkAppCode"`
}

// ModuleDescription 单个 module 字段
type ModuleDescription struct {
	Name      string    `yaml:"name"`
	SourceDir string    `yaml:"sourceDir"`
	Language  string    `yaml:"language"`
	Spec      BkAppSpec `yaml:"spec"`
}

// BkAppSpec bkapp spec
type BkAppSpec struct {
	Processes     []Process               `yaml:"processes"`
	Configuration AppConfig               `yaml:"configuration"`
	Build         buildconfig.BuildConfig `yaml:"build,omitempty"`
}

// Process 进程配置
type Process struct {
	// Name of process
	Name string `yaml:"name"`

	// ProcCommand is the script command to start the process
	ProcCommand string `yaml:"procCommand"`
}

// AppConfig is bkapp related configuration, such as environment variables, etc.
type AppConfig struct {
	// List of environment variables to set in the container.
	Env []AppEnvVar `yaml:"env,omitempty"`
}

// AppEnvVar 单个环境变量
type AppEnvVar struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}
