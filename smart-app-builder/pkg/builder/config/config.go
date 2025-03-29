package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var G = struct {
	*viper.Viper

	// The url of the source code, which begins with file:// or http(s)://
	SourceURL string
	// The url of the s-mart artifact to put, which begins with file:// or http(s)://
	DestURL string

	// RuntimeWorkspace is runtime workspace
	RuntimeWorkspace string

	// ScratchImageTarPath is cnb run-scratch image tar path
	ScratchImageTarPath string
	// ScratchImage is cnb run-scratch image
	ScratchImage string
	// CNBBuilderImage is CNB builder image used to build the source code
	CNBBuilderImage string

	// DaemonSockFile is container daemon sock file
	DaemonSockFile string
}{Viper: viper.New()}

// SetGlobalConfig set global config
func SetGlobalConfig() {
	G.BindPFlags(pflag.CommandLine)
	G.AutomaticEnv()

	G.SourceURL = G.GetString("source-url")
	G.DestURL = G.GetString("dest-url")

	G.ScratchImageTarPath = G.GetString("SCRATCH_IMAGE_TAR_PATH")
	G.ScratchImage = G.GetString("SCRATCH_IMAGE")
	G.SetDefault("CNB_BUILDER_IMAGE", "bk-builder-heroku-bionic:v1.0.2")
	G.CNBBuilderImage = G.GetString("CNB_BUILDER_IMAGE")

	G.SetDefault("RUNTIME_WORKSPACE", "/podman/smart-app")
	G.RuntimeWorkspace = G.GetString("RUNTIME_WORKSPACE")

	G.DaemonSockFile = G.GetString("DAEMON_SOCK")
}
