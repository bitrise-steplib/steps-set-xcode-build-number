package step

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/bitrise-io/go-steputils/v2/export"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-io/go-xcode/v2/autocodesign/projectmanager"
	"github.com/bitrise-io/go-xcode/xcodeproject/xcodeproj"
)

type Updater struct {
	inputParser stepconf.InputParser
	exporter    export.Exporter
	logger      log.Logger
}

func NewUpdater(inputParser stepconf.InputParser, exporter export.Exporter, logger log.Logger) Updater {
	return Updater{
		inputParser: inputParser,
		exporter:    exporter,
		logger:      logger,
	}
}

func (u Updater) ProcessConfig() (Config, error) {
	var input Input
	err := u.inputParser.Parse(&input)
	if err != nil {
		return Config{}, err
	}

	if input.BuildVersionOffset < 0 {
		return Config{}, fmt.Errorf("build version offset cannot be a negative value (%d)", input.BuildVersionOffset)
	}

	stepconf.Print(input)
	u.logger.Println()

	// golangci told me that I should do this instead of creating a Config struct.
	return Config(input), nil
}

func (u Updater) Run(config Config) (Result, error) {
	helper, err := projectmanager.NewProjectHelper(config.ProjectPath, config.Scheme, config.Configuration)
	if err != nil {
		return Result{}, err
	}

	generated, err := generatesInfoPlist(helper, config.Target, config.Configuration)
	if err != nil {
		return Result{}, err
	}

	buildVersion := config.BuildVersion + config.BuildVersionOffset

	if generated {
		u.logger.Printf("The version numbers are stored in the project file.")

		err := updateVersionNumbersInProject(helper, config.Target, config.Configuration, buildVersion, config.BuildShortVersionString)
		if err != nil {
			return Result{}, err
		}
	} else {
		u.logger.Printf("The version numbers are stored in the plist file.")

		err := updateVersionNumbersInInfoPlist(helper, config.Target, config.Configuration, buildVersion, config.BuildShortVersionString)
		if err != nil {
			return Result{}, err
		}
	}

	u.logger.Donef("Version numbers successfully updated.")

	return Result{BuildVersion: buildVersion}, nil
}

func (u Updater) Export(result Result) error {
	return u.exporter.ExportOutput("XCODE_BUNDLE_VERSION", strconv.FormatInt(result.BuildVersion, 10))
}

func generatesInfoPlist(helper *projectmanager.ProjectHelper, targetName, configuration string) (bool, error) {
	buildConfig, err := buildConfiguration(helper, targetName, configuration)
	if err != nil {
		return false, err
	}

	generatesInfoPlist := buildConfig.BuildSettings["GENERATE_INFOPLIST_FILE"] == "YES"

	return generatesInfoPlist, err
}

func updateVersionNumbersInProject(helper *projectmanager.ProjectHelper, targetName, configuration string, bundleVersion int64, shortVersion string) error {
	if targetName == "" {
		targetName = helper.MainTarget.Name
	}

	for _, target := range helper.XcProj.Proj.Targets {
		if target.Name != targetName {
			continue
		}

		for _, buildConfig := range target.BuildConfigurationList.BuildConfigurations {
			if configuration != "" && buildConfig.Name != configuration {
				continue
			}

			buildConfig.BuildSettings["CURRENT_PROJECT_VERSION"] = bundleVersion

			if shortVersion != "" {
				buildConfig.BuildSettings["MARKETING_VERSION"] = shortVersion
			}
		}
	}

	err := helper.XcProj.Save()
	if err != nil {
		return err
	}

	return nil
}

func updateVersionNumbersInInfoPlist(helper *projectmanager.ProjectHelper, targetName, configuration string, bundleVersion int64, shortVersion string) error {
	buildConfig, err := buildConfiguration(helper, targetName, configuration)
	if err != nil {
		return err
	}

	infoPlistPath, err := buildConfig.BuildSettings.String("INFOPLIST_FILE")
	if err != nil {
		return err
	}

	absoluteInfoPlistPath := filepath.Join(filepath.Dir(helper.XcProj.Path), infoPlistPath)

	infoPlist, format, _ := xcodeproj.ReadPlistFile(absoluteInfoPlistPath)
	infoPlist["CFBundleVersion"] = strconv.FormatInt(bundleVersion, 10)

	if shortVersion != "" {
		infoPlist["CFBundleShortVersionString"] = shortVersion
	}

	err = xcodeproj.WritePlistFile(absoluteInfoPlistPath, infoPlist, format)
	if err != nil {
		return err
	}

	return nil
}

func buildConfiguration(helper *projectmanager.ProjectHelper, targetName, configuration string) (*xcodeproj.BuildConfiguration, error) {
	if targetName == "" {
		targetName = helper.MainTarget.Name
	}

	if configuration == "" {
		configuration = helper.MainTarget.BuildConfigurationList.DefaultConfigurationName
	}

	for _, target := range helper.XcProj.Proj.Targets {
		if target.Name != targetName {
			continue
		}

		for _, buildConfig := range target.BuildConfigurationList.BuildConfigurations {
			if buildConfig.Name == configuration {
				return &buildConfig, nil
			}
		}
	}

	return nil, fmt.Errorf("")
}
