package step

import (
	"testing"

	"github.com/bitrise-io/go-steputils/v2/export"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-steplib/steps-set-xcode-build-number/step/mocks"
	"github.com/stretchr/testify/assert"
)

func TestExport(t *testing.T) {
	result := Result{BuildVersion: "999"}

	mockFactory := mocks.NewFactory(t)
	arguments := []string{"add", "--key", "XCODE_BUNDLE_VERSION", "--value", result.BuildVersion}
	mockFactory.On("Create", "envman", arguments, (*command.Opts)(nil)).Return(testCommand())

	inputParser := stepconf.NewInputParser(env.NewRepository())
	exporter := export.NewExporter(mockFactory)

	updater := NewUpdater(inputParser, exporter, log.NewLogger())
	err := updater.Export(result)
	assert.NoError(t, err)

	mockFactory.AssertExpectations(t)
}

func testCommand() command.Command {
	factory := command.NewFactory(env.NewRepository())
	return factory.Create("pwd", []string{}, nil)
}
