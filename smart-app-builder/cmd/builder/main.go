package main

import (
	"fmt"
	"os"

	"github.com/TencentBlueking/bkpaas/cnb-builder-shim/pkg/logging"
	"github.com/go-logr/logr"
	flag "github.com/spf13/pflag"

	"github.com/TencentBlueking/bkpaas/smart-app-builder/pkg/builder"
)

const (
	SourceURLEnvVarKey = "SOURCE_GET_URL"
	DestURLEnvVarKey   = "DEST_PUT_URL"
)

var (
	sourceURL = flag.String("source-url", os.Getenv(SourceURLEnvVarKey), "The url of the source code.")
	destURL   = flag.String("dest-url", os.Getenv(DestURLEnvVarKey), "The url of the s-mart artifact to put.")
)

func main() {
	logger := logging.Default()

	parseFlags(logger)

	executor, err := builder.NewBuildExecutor(*sourceURL, *destURL)
	if err != nil {
		logger.Error(err, "failed to create build executor")
		os.Exit(1)
	}

	if err := executor.Execute(); err != nil {
		logger.Error(err, "failed to build s-mart package")
		os.Exit(1)
	}

	logger.Info("build s-mart package successfully")
}

func parseFlags(logger logr.Logger) {
	flag.Parse()

	if *sourceURL == "" {
		logger.Error(
			fmt.Errorf("sourceURL is empty"),
			fmt.Sprintf(
				"please provide by setting --source-url option or setting as an environment variable %s",
				SourceURLEnvVarKey,
			),
		)
		os.Exit(1)
	}

	if *destURL == "" {
		logger.Error(
			fmt.Errorf("destURL is empty"),
			fmt.Sprintf(
				"please provide by setting --dest-url option or setting as an environment variable %s",
				DestURLEnvVarKey,
			),
		)
		os.Exit(1)
	}
}
