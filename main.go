package main

import (
	"os"

	"github.com/bitrise-io/go-steputils/v2/export"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-steputils/v2/stepenv"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	. "github.com/bitrise-io/go-utils/v2/exitcode"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-steplib/steps-set-xcode-build-number/step"
)

func main() {
	exitCode := run()
	os.Exit(int(exitCode))
}

func run() ExitCode {
	logger := log.NewLogger()

	updater := createUpdater(logger)
	config, err := updater.ProcessConfig()
	if err != nil {
		logger.Errorf("Process config: %s", err)
		return Failure
	}

	result, err := updater.Run(config)
	if err != nil {
		logger.Errorf("Run: %s", err)
		return Failure
	}

	if err := updater.Export(result); err != nil {
		logger.Errorf("Export outputs: %s", err)
		return Failure
	}

	return Success
}

func createUpdater(logger log.Logger) step.Updater {
	envRepository := stepenv.NewRepository(env.NewRepository())
	inputParser := stepconf.NewInputParser(envRepository)
	exporter := export.NewExporter(command.NewFactory(envRepository))

	return step.NewUpdater(inputParser, exporter, logger)
}
