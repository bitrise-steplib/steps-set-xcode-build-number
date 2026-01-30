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
	"github.com/stretchr/testify/require"
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

func Test_incrementBuildVersion(t *testing.T) {
	logger := log.NewLogger()
	tests := []struct {
		name         string
		buildVersion string
		offset       int64
		want         string
		wantErr      bool
	}{
		{
			name:         "simple increment",
			buildVersion: "42",
			offset:       1,
			want:         "43",
			wantErr:      false,
		},
		{
			name:         "skip increment",
			buildVersion: "42",
			offset:       -1,
			want:         "42",
			wantErr:      false,
		},
		{
			name:         "non-numeric build version",
			buildVersion: "1.2.3.4",
			offset:       0,
			want:         "1.2.3.4",
			wantErr:      false,
		},
		{
			name:         "non-numeric build version with non-zero offset",
			buildVersion: "1.2.3.4",
			offset:       3,
			want:         "",
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := incrementBuildVersion(logger, tt.buildVersion, tt.offset)
			if tt.wantErr {
				require.Error(t, gotErr)
			} else {
				require.NoError(t, gotErr)
			}

			require.Equal(t, tt.want, got)
		})
	}
}
