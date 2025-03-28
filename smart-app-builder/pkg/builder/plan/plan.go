package plan

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/exp/maps"

	"github.com/TencentBlueking/bkpaas/smart-app-builder/pkg/builder/appdesc"
	bcfg "github.com/TencentBlueking/bkpaas/smart-app-builder/pkg/builder/buildconfig"
)

const AppDescFileName = "app_desc.yaml"

// BuildPlan 构建计划
type BuildPlan struct {
	Procfile map[string]string
	Steps    []*ModuleBuildStep
}

// ModuleBuildStep 模块构建步骤配置
type ModuleBuildStep struct {
	SourceDir          string
	RequiredBuildpacks string
	ModuleNames        []string
	Envs               map[string]string
}

// PrepareBuildPlan 解析 app_desc.yaml 文件, 生成 BuildPlan
func PrepareBuildPlan(sourceDir string) (*BuildPlan, error) {
	if err := ensureAppDescYaml(sourceDir); err != nil {
		return nil, err
	}

	desc, err := appdesc.ParseAppDescYAML(filepath.Join(sourceDir, AppDescFileName))
	if err != nil {
		return nil, errors.Wrap(err, "parse app_desc.yaml error")
	}

	if err = desc.Validate(); err != nil {
		return nil, err
	}

	procfile := desc.GenerateProcfile()
	if len(procfile) == 0 {
		return nil, errors.New("no valid processes found in app_desc.yaml")
	}

	steps, err := buildSteps(desc)
	if err != nil {
		return nil, err
	}

	plan := &BuildPlan{Procfile: procfile, Steps: steps}

	return plan, nil
}

// ensureAppDescYaml 确保 app_desc.yaml 存在
func ensureAppDescYaml(sourceDir string) error {
	_, err := os.Stat(filepath.Join(sourceDir, AppDescFileName))
	// 如果 app_desc.yaml 不存在, 则尝试探测 app_desc.yml
	if err != nil && os.IsNotExist(err) {
		if _, err = os.Stat(filepath.Join(sourceDir, "app_desc.yml")); err != nil {
			return err
		}
		// 重命名为 app_desc.yaml
		return os.Rename(filepath.Join(sourceDir, "app_desc.yml"), filepath.Join(sourceDir, AppDescFileName))

	}
	return err
}

// buildSteps 生成构建步骤
func buildSteps(desc appdesc.AppDesc) ([]*ModuleBuildStep, error) {
	configs, err := desc.GenerateModuleBuildConfig()
	if err != nil {
		return nil, err
	}

	stepMap := make(map[string]*ModuleBuildStep)

	for _, cfg := range configs {
		rBuildpacks, err := ToRequiredBuildpacks(cfg.Buildpacks)
		if err != nil {
			return nil, err
		}

		// 合并采用相同构建方案的模块
		k := fmt.Sprintf("%s%s", cfg.SourceDir, rBuildpacks)
		if v, ok := stepMap[k]; !ok {
			stepMap[k] = &ModuleBuildStep{
				SourceDir:          cfg.SourceDir,
				RequiredBuildpacks: rBuildpacks,
				ModuleNames:        []string{cfg.ModuleName},
				Envs:               cfg.Envs,
			}
		} else {
			v.ModuleNames = append(v.ModuleNames, cfg.ModuleName)
			maps.Copy(v.Envs, cfg.Envs)
		}
	}

	return maps.Values(stepMap), nil
}

// ToRequiredBuildpacks 将 buildpacks 从 list 排序后, 转换成 string 结构.
// 格式如: tgz bk-buildpack-apt ... v2;tgz bk-buildpack-python ... v213
// NOTE: 在转换的过程中, 如果 buildpack 的 version 为空版本, 会采用默认版本替换
func ToRequiredBuildpacks(buildpacks []bcfg.Buildpack) (string, error) {
	sort.Slice(buildpacks, func(i, j int) bool {
		return strings.Compare(buildpacks[i].Name, buildpacks[j].Name) > 0
	})

	var parts []string
	var err error

	for _, bp := range buildpacks {
		v := bp.Version
		if v == "" {
			v, err = bcfg.GetDefaultVersionByBPName(bp.Name)
			if err != nil {
				return "", err
			}
		}

		parts = append(parts, strings.Join([]string{"tgz", bp.Name, "...", v}, " "))

	}

	return strings.Join(parts, ";"), nil
}
