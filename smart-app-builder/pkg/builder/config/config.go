package config

import (
	"github.com/spf13/viper"
)

var G = struct {
	// The url of the source code, which begins with file:// or http(s)://
	SourceURL string
	// The url of the s-mart artifact to put, which begins with file:// or http(s)://
	DestURL string
	// CNBBuilderImage is CNB builder image used to build the source code
	CNBBuilderImage string
}{}

// SetGlobalConfig set global config
func SetGlobalConfig() {
	G.SourceURL = viper.GetString("source-url")
	G.DestURL = viper.GetString("dest-url")

	viper.SetDefault("CNB_BUILDER_IMAGE", "bk-builder-heroku-bionic:v1.0.2")
	G.CNBBuilderImage = viper.GetString("CNB_BUILDER_IMAGE")
}
